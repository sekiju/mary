package corocoro

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/knst0/mdl/extractor/registry"
	"github.com/knst0/mdl/sdk/manga"
	"github.com/knst0/mdl/sdk/manga/pluginutil"

	"resty.dev/v3"
)

type Extractor struct {
	pluginutil.Base
}

func (e *Extractor) FindChapters(ctx context.Context, URL string) ([]*manga.Chapter, error) {
	return nil, manga.ErrChapterListingUnsupported
}

func (e *Extractor) SupportsChapterListing() bool {
	return false
}

var re = regexp.MustCompile(`https://www.corocoro.jp/chapter/(\d*)/viewer`)

func (e *Extractor) FindChapter(ctx context.Context, URL string) (*manga.Chapter, error) {
	matches := re.FindStringSubmatch(URL)
	if len(matches) != 2 {
		return nil, manga.ErrInvalidChapterURL
	}

	chapterID := matches[1]

	req := resty.New().R().SetContext(ctx)
	e.ApplyCookie(req)

	res, err := req.Get(URL)
	if err != nil {
		return nil, err
	}

	html := res.String()

	titleStart := strings.Index(html, "<title>")
	titleEnd := strings.Index(html, "</title>")
	if titleStart == -1 || titleEnd == -1 {
		return nil, fmt.Errorf("title not found in HTML")
	}
	fullTitle := html[titleStart+7 : titleEnd]
	// Format: "{manga_title} {chapter_title} {author} | 週刊コロコロコミック"
	// Strip trailing " | 週刊コロコロコミック"
	chapterMainName := strings.TrimSuffix(fullTitle, " | 週刊コロコロコミック")

	return &manga.Chapter{
		ID:      chapterID,
		Number:  "",
		Title:   chapterMainName,
		Index:   0,
		URL:     URL,
		MangaID: "",
	}, nil
}

func (e *Extractor) FindChapterPages(ctx context.Context, chapter *manga.Chapter) ([]*manga.Page, error) {
	apiURL := fmt.Sprintf("https://www.corocoro.jp/api/csr?rq=chapter/viewer&chapter_id=%s&use_ticket=0&event_point=0&paid_point=0", chapter.ID)

	req := resty.New().R().
		SetContext(ctx).
		SetHeader("RSC", "1").
		SetHeader("Content-Type", "text/plain;charset=UTF-8")
	e.ApplyCookie(req)

	res, err := req.Put(apiURL)
	if err != nil {
		return nil, err
	}

	data := res.Bytes()

	pb := &protobufReader{buf: data}
	viewer := pb.parseViewer()

	return pluginutil.BuildPages(len(viewer.pages), ".webp", func(index int, filename string) (*manga.Page, error) {
		page := viewer.pages[index]
		return &manga.Page{
			Index:    uint(index),
			URL:      page.src,
			Filename: filename,
			Decode: func(b []byte) ([]byte, error) {
				if viewer.aesKey == "" || viewer.aesIv == "" {
					return b, nil
				}

				key, err := hex.DecodeString(viewer.aesKey)
				if err != nil {
					return nil, fmt.Errorf("invalid key hex: %w: %w", err, manga.ErrMalformedChapterData)
				}

				iv, err := hex.DecodeString(viewer.aesIv)
				if err != nil {
					return nil, fmt.Errorf("invalid IV hex: %w: %w", err, manga.ErrMalformedChapterData)
				}

				block, err := aes.NewCipher(key)
				if err != nil {
					return nil, err
				}

				mode := cipher.NewCBCDecrypter(block, iv)
				mode.CryptBlocks(b, b)

				// PKCS7 unpad
				padLen := int(b[len(b)-1])
				if padLen > 0 && padLen <= aes.BlockSize {
					return b[:len(b)-padLen], nil
				}
				return b, nil
			},
		}, nil
	})
}

// protobufReader is a minimal protobuf parser for the corocoro viewer response.
type protobufReader struct {
	buf []byte
	pos int
}

func (r *protobufReader) readVarint() (uint64, bool) {
	var result uint64
	var shift uint
	for r.pos < len(r.buf) {
		b := r.buf[r.pos]
		r.pos++
		result |= uint64(b&0x7F) << shift
		shift += 7
		if b < 0x80 {
			return result, true
		}
	}
	return 0, false
}

func (r *protobufReader) readFixed32() (uint32, bool) {
	if r.pos+4 > len(r.buf) {
		return 0, false
	}
	v := binary.LittleEndian.Uint32(r.buf[r.pos:])
	r.pos += 4
	return v, true
}

func (r *protobufReader) readFixed64() (uint64, bool) {
	if r.pos+8 > len(r.buf) {
		return 0, false
	}
	v := binary.LittleEndian.Uint64(r.buf[r.pos:])
	r.pos += 8
	return v, true
}

func (r *protobufReader) readBytes(length int) ([]byte, bool) {
	if r.pos+length > len(r.buf) {
		return nil, false
	}
	b := r.buf[r.pos : r.pos+length]
	r.pos += length
	return b, true
}

func (r *protobufReader) skipWireType(wireType uint64) bool {
	switch wireType {
	case 0: // varint
		_, ok := r.readVarint()
		return ok
	case 1: // 64-bit
		_, ok := r.readFixed64()
		return ok
	case 2: // length-delimited
		length, ok := r.readVarint()
		if !ok {
			return false
		}
		if r.pos+int(length) > len(r.buf) {
			return false
		}
		r.pos += int(length)
		return true
	case 5: // 32-bit
		_, ok := r.readFixed32()
		return ok
	default:
		return false
	}
}

func (r *protobufReader) readString() (string, bool) {
	length, ok := r.readVarint()
	if !ok {
		return "", false
	}
	data, ok := r.readBytes(int(length))
	if !ok {
		return "", false
	}
	return string(data), true
}

type viewerResponse struct {
	pages  []pageInfo
	aesKey string
	aesIv  string
}

type pageInfo struct {
	src string
}

// parseViewer parses the ViewerView protobuf message.
// Schema (from reverse engineering):
//
//	Field 2: repeated PageImage {
//	  1: src (string)
//	  2: width (varint)
//	  3: height (varint)
//	  4: alt (string)
//	}
//	Field 19: aesKey (string, hex-encoded)
//	Field 20: aesIv (string, hex-encoded)
func (r *protobufReader) parseViewer() viewerResponse {
	var resp viewerResponse

	for r.pos < len(r.buf) {
		tag, ok := r.readVarint()
		if !ok {
			break
		}
		fieldNum := tag >> 3
		wireType := tag & 7

		switch fieldNum {
		case 2: // pages
			if wireType == 2 {
				length, ok := r.readVarint()
				if !ok {
					return resp
				}
				end := r.pos + int(length)
				pi := r.parsePageImage()
				if pi.src != "" {
					resp.pages = append(resp.pages, pi)
				}
				// Ensure we're at the right position
				r.pos = end
			} else {
				r.skipWireType(wireType)
			}

		case 19: // aesKey
			if s, ok := r.readString(); ok {
				resp.aesKey = s
			}

		case 20: // aesIv
			if s, ok := r.readString(); ok {
				resp.aesIv = s
			}

		default:
			if !r.skipWireType(wireType) {
				return resp
			}
		}
	}

	return resp
}

func (r *protobufReader) parsePageImage() pageInfo {
	var pi pageInfo

	for r.pos < len(r.buf) {
		tag, ok := r.readVarint()
		if !ok {
			break
		}
		fieldNum := tag >> 3
		wireType := tag & 7

		switch fieldNum {
		case 1: // src
			if s, ok := r.readString(); ok {
				pi.src = s
			}

		default:
			if !r.skipWireType(wireType) {
				return pi
			}
		}
	}

	return pi
}

func init() {
	registry.Register("www.corocoro.jp", registry.WithSession(New))
}

func New() (manga.Extractor, error) {
	return &Extractor{Base: pluginutil.Base{Settings: &manga.Settings{}}}, nil
}
