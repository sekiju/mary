package speed_binb

import (
	json "github.com/bytedance/sonic"
	"image"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DescrambleView struct {
	Width, Height int
	Transfers     []DescrambleTransfer
}

type DescrambleTransfer struct {
	Index  int
	Coords []DescrambleCord
}

// Piece is one scrambled tile: its grid position (x, y) and span (w, h) in
// tile units, used by the SpeedbinbA prototype's tile-shuffle scheme.
type Piece struct {
	x, y, w, h int
}

// Yt is a decoded tile-shuffle layout for SpeedbinbA: grid dimensions
// (ndx x ndy) plus the per-tile placement order.
type Yt struct {
	ndx, ndy int
	piece    *[]Piece
}

// SpeedBinbDecoder is implemented by each descramble prototype variant
// (SpeedbinbF, SpeedbinbA, SpeedbinbH). vt reports whether the prototype
// parsed successfully; bt reports whether an image is large enough to be
// scrambled under this prototype; dt returns the descrambled image's
// bounds; gt returns the source->destination copy regions for
// descrambleImage.
type SpeedBinbDecoder interface {
	vt() bool
	bt(t image.Rectangle) bool
	dt(t image.Rectangle) image.Rectangle
	gt(t image.Rectangle) []DescrambleCord
}

type DescrambleCord struct {
	XSrc, YSrc, Width, Height, XDest, YDest int
}

// SpeedbinbF is the grid-descramble prototype variant identified by a
// "=C-I(+|-)Jt-<data>" table entry: a C x I grid of tiles with Jt-pixel
// gutters, permuted per Mt (built from the encoder/decoder tile orderings
// Xt/Et and It/St).
type SpeedbinbF struct {
	C, I, Jt           int
	Xt, Et, It, St, Mt []int
}

// Tt is a base64url-alphabet reverse lookup table (char code -> 0-63, or -1
// if the char isn't in base64chars), used by SpeedbinbF.yt to decode the
// packed tile-ordering data.
var Tt = []int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 62, -1, -1, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, -1, -1, -1, -1, -1, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, -1, -1, -1, -1, 63, -1, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, -1, -1, -1, -1, -1}

func padStart(input string, length int, padChar string) string {
	if len(input) >= length {
		return input
	}

	padding := strings.Repeat(padChar, length-len(input))
	return padding + input
}

const base64chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

// generateSharingKey derives the "k" query-param sharing key from the
// content ID (t), mixing a hex timestamp with an XOR checksum of t.
func generateSharingKey(t string) string {
	timestampMilliseconds := time.Now().UnixNano() / int64(time.Millisecond)
	n := strconv.FormatInt(timestampMilliseconds, 16)
	nHex := padStart(n, 16, "x")
	repeatedString := strings.Repeat(t, (16+len(t)-1)/len(t))
	r := repeatedString[:16]
	e := repeatedString[len(repeatedString)-16:]

	s, u, h := 0, 0, 0
	result := ""
	for i, char := range nHex {
		s ^= int(char)
		u ^= int(r[i])
		h ^= int(e[i])

		result += string(char)
		result += string(base64chars[(s+u+h)&63])
	}

	return result
}

// decodeTable XOR/rotation-decodes the ctbl/ptbl prototype table (n) using
// the content ID (t) and sharing key (i) as a keystream seed, returning the
// resulting JSON array of prototype descriptors.
func decodeTable(t, i, n string) []string {
	r := t + ":" + i
	e := 0
	for s := 0; s < len(r); s++ {
		e += int(r[s]) << (s % 16)
	}
	e &= 2147483647
	if e == 0 {
		e = 305419896
	}

	u := ""
	h := e
	for s := 0; s < len(n); s++ {
		h = h>>1 ^ 1210056708&-(1&h)
		o := (int(n[s])-32+h)%94 + 32
		u += string(rune(o))
	}

	var result []string
	err := json.Unmarshal([]byte(u), &result)
	if err != nil {
		return nil
	}
	return result
}

// selectPrototype picks the SpeedBinbDecoder variant (F/A/H) matching the
// image path t, indexing into ctbl/ptbl by a checksum of t's basename.
func selectPrototype(t string, ctbl, ptbl []string) SpeedBinbDecoder {
	i := [2]int{0, 0}

	if t != "" {
		n := strings.LastIndex(t, "/") + 1
		r := len(t) - n
		for e := 0; e < r; e++ {
			i[e%2] += int(t[e+n])
		}

		i[0] %= 8
		i[1] %= 8
	}

	s, u := ptbl[i[0]], ctbl[i[1]]

	if s[0] == '=' && u[0] == '=' {
		return NewSpeedbinbF(u, s)
	} else if regexp.MustCompile("^[0-9]").MatchString(u) && regexp.MustCompile("^[0-9]").MatchString(s) {
		return NewSpeedbinbA(u, s)
	} else if "" == u && "" == s {
		return NewSpeedbinbH()
	}

	return nil
}

func NewSpeedbinbF(t, i string) *SpeedbinbF {
	sb := SpeedbinbF{Mt: []int{}}

	reg := regexp.MustCompile("^=([0-9]+)-([0-9]+)([-+])([0-9]+)-([-_0-9A-Za-z]+)$")

	n, r := reg.FindStringSubmatch(t), reg.FindStringSubmatch(i)

	if len(n) > 0 && len(r) > 0 && n[1] == r[1] && n[2] == r[2] && n[4] == r[4] && "+" == n[3] && "-" == r[3] {
		sb.C, _ = strconv.Atoi(n[1])
		sb.I, _ = strconv.Atoi(n[2])
		sb.Jt, _ = strconv.Atoi(n[4])

		if !(8 < sb.C || 8 < sb.I || 64 < sb.C*sb.I) {
			e := sb.C + sb.I + sb.C*sb.I
			if len(n[5]) == e && len(r[5]) == e {
				s, u := sb.yt(n[5]), sb.yt(r[5])

				sb.Xt = s["n"]
				sb.Et = s["t"]
				sb.It = u["n"]
				sb.St = u["t"]

				for h := 0; h < sb.C*sb.I; h++ {
					sb.Mt = append(sb.Mt, s["p"][u["p"][h]])
				}
			}
		}
	}

	return &sb
}

// vt reports whether the grid params parsed successfully (Mt is set).
func (s *SpeedbinbF) vt() bool {
	return s.Mt != nil
}

// bt reports whether image t is large enough to actually be grid-scrambled
// (accounting for the C x I grid and Jt gutters); small images are served
// unscrambled.
func (s *SpeedbinbF) bt(t image.Rectangle) bool {
	i := 2 * s.C * s.Jt
	n := 2 * s.I * s.Jt
	return t.Dx() >= 64+i && t.Dy() >= 64+n && t.Dx()*t.Dy() >= (320+i)*(320+n)
}

// dt returns the descrambled image bounds: t minus the grid gutters, or t
// unchanged if bt reports the image isn't scrambled.
func (s *SpeedbinbF) dt(t image.Rectangle) image.Rectangle {
	if s.bt(t) {
		return image.Rect(0, 0, t.Dx()-2*s.C*s.Jt, t.Dy()-2*s.I*s.Jt)
	}

	return t
}

// gt computes the per-tile source->destination copy regions that undo the
// C x I grid permutation encoded in Mt.
func (s *SpeedbinbF) gt(t image.Rectangle) []DescrambleCord {
	if !s.vt() {
		return nil
	}
	if !s.bt(t) {
		return []DescrambleCord{
			{XSrc: 0, YSrc: 0, Width: t.Dx(), Height: t.Dy(), XDest: 0, YDest: 0},
		}
	}

	i := t.Dx() - 2*s.C*s.Jt
	n := t.Dy() - 2*s.I*s.Jt
	r := (i + s.C - 1) / s.C
	e := i - (s.C-1)*r
	sy := (n + s.I - 1) / s.I
	u := n - (s.I-1)*sy
	var h []DescrambleCord

	for o := 0; o < s.C*s.I; o++ {
		a := o % s.C
		f := o / s.C
		_c := 0
		if s.It[f] < a {
			_c = e - r
		}
		c := s.Jt + a*(r+2*s.Jt) + _c
		_l := 0
		if s.St[a] < f {
			_l = u - sy
		}
		l := s.Jt + f*(sy+2*s.Jt) + _l
		v := s.Mt[o] % s.C
		d := s.Mt[o] / s.C
		_g := 0
		if s.Xt[d] < v {
			_g = e - r
		}
		_p := 0
		if s.Et[v] < d {
			_p = u - sy
		}

		g := v*r + _g
		p := d*sy + _p
		b := r
		if s.It[f] == a {
			b = e
		}
		m := sy
		if s.St[a] == f {
			m = u
		}

		if 0 < i && 0 < n {
			h = append(h, DescrambleCord{
				XSrc:   c,
				YSrc:   l,
				Width:  b,
				Height: m,
				XDest:  g,
				YDest:  p,
			})
		}
	}

	return h
}

// yt decodes a packed tile-ordering string t (via the Tt lookup table) into
// its "n" (encoder order), "t" (decoder order), and "p" (permutation) slices.
func (s *SpeedbinbF) yt(t string) map[string][]int {
	var n, r, e []int

	for i := 0; i < s.C; i++ {
		n = append(n, Tt[int(t[i])])
	}

	for i := 0; i < s.I; i++ {
		r = append(r, Tt[int(t[s.C+i])])
	}

	for i := 0; i < s.C*s.I; i++ {
		e = append(e, Tt[int(t[s.C+s.I+i])])
	}

	return map[string][]int{
		"t": n,
		"n": r,
		"p": e,
	}
}

// SpeedbinbA is the tile-shuffle descramble prototype variant identified by
// a numeric "<ndx>-<ndy>-<data>" table entry: mt is the encoder's tile
// layout, wt is the decoder's (matching) layout.
type SpeedbinbA struct {
	mt, wt *Yt
}

func NewSpeedbinbA(t, i string) *SpeedbinbA {
	sb := &SpeedbinbA{
		mt: nil,
		wt: nil,
	}

	n, r := sb.yt(t), sb.yt(i)

	if n != nil && r != nil && n.ndx == r.ndx && n.ndy == r.ndy {
		sb.mt = n
		sb.wt = r
	}

	return sb
}

// vt reports whether both tile layouts parsed and are compatible.
func (sb *SpeedbinbA) vt() bool {
	return sb.mt != nil && sb.wt != nil
}

// bt reports whether image t is large enough to be tile-shuffled.
func (sb *SpeedbinbA) bt(t image.Rectangle) bool {
	return 64 <= t.Dx() && 64 <= t.Dy() && 102400 <= t.Dx()*t.Dy()
}

// dt: SpeedbinbA never trims bounds, unlike SpeedbinbF.
func (sb *SpeedbinbA) dt(t image.Rectangle) image.Rectangle {
	return t
}

// gt computes the per-tile source->destination copy regions that undo the
// mt/wt tile shuffle, plus any leftover right/bottom strip outside the
// tile grid.
func (sb *SpeedbinbA) gt(t image.Rectangle) []DescrambleCord {
	if !sb.vt() {
		return nil
	}

	if !sb.bt(t) {
		return []DescrambleCord{{
			XSrc:   0,
			YSrc:   0,
			Width:  t.Dx(),
			Height: t.Dy(),
			XDest:  0,
			YDest:  0,
		}}
	}

	var rects []DescrambleCord

	n := t.Dx() - t.Dx()%8   // stabilize w
	r := (n-1)/7 - (n-1)/7%8 // block w
	e := n - 7*r             // padding w
	s := t.Dy() - t.Dy()%8   // same...
	u := (s-1)/7 - (s-1)/7%8
	h := s - 7*u

	for a := 0; a < len(*sb.mt.piece); a++ {
		f := (*sb.mt.piece)[a]
		c := (*sb.wt.piece)[a]

		rects = append(rects, DescrambleCord{
			f.x/2*r + f.x%2*e,
			f.y/2*u + f.y%2*h,
			f.w/2*r + f.w%2*e,
			f.h/2*u + f.h%2*h,
			c.x/2*r + c.x%2*e,
			c.y/2*u + c.y%2*h,
		})
	}

	l := r*(sb.mt.ndx-1) + e
	v := u*(sb.mt.ndy-1) + h

	if l < t.Dx() {
		rects = append(rects, DescrambleCord{
			XSrc:   l,
			YSrc:   0,
			Width:  t.Dx() - l,
			Height: v,
			XDest:  l,
			YDest:  0,
		})
	}

	if v < t.Dy() {
		rects = append(rects, DescrambleCord{
			XSrc:   0,
			YSrc:   v,
			Width:  t.Dx(),
			Height: t.Dy() - v,
			XDest:  0,
			YDest:  v,
		})
	}

	return rects
}

// yt parses a "<ndx>-<ndy>-<data>" tile layout string into a Yt, where each
// pair of chars in data (decoded via Ot) is one tile's (x, y) grid position.
func (sb *SpeedbinbA) yt(t string) *Yt {
	if t == "" {
		return nil
	}

	i := strings.Split(t, "-")
	if len(i) != 3 {
		return nil
	}

	n, err := strconv.Atoi(i[0])
	if err != nil {
		return nil
	}

	r, err := strconv.Atoi(i[1])
	if err != nil {
		return nil
	}

	e := i[2]

	if len(e) != n*r*2 {
		return nil
	}

	var v []Piece
	a := (n-1)*(r-1) - 1
	f := a + (n - 1)
	c := f + (r - 1)
	l := c + 1

	for d := 0; d < n*r; d++ {
		s := sb.Ot(e[2*d])
		u := sb.Ot(e[2*d+1])

		var o, h int
		if d <= a {
			h, o = 2, 2
		} else if d <= f {
			h, o = 2, 1
		} else if d <= c {
			h, o = 1, 2
		} else if d <= l {
			h, o = 1, 1
		}

		v = append(v, Piece{
			x: s,
			y: u,
			w: h,
			h: o,
		})
	}

	return &Yt{
		ndx:   n,
		ndy:   r,
		piece: &v,
	}
}

// Ot decodes a single base62-ish digit char (A-Z, a-z) into 0..51.
func (sb *SpeedbinbA) Ot(t byte) int {
	i := 0
	n := strings.Index("ABCDEFGHIJKLMNOPQRSTUVWXYZ", string(t))
	if n < 0 {
		n = strings.Index("abcdefghijklmnopqrstuvwxyz", string(t))
	} else {
		i = 1
	}
	return i + 2*n
}

// SpeedbinbH is the no-op prototype variant used when both the ctbl and
// ptbl entries for an image are empty, meaning the image isn't scrambled at
// all — gt returns a single identity copy region.
type SpeedbinbH struct {
	mt, wt *Yt
}

func NewSpeedbinbH() *SpeedbinbH {
	return &SpeedbinbH{
		mt: nil,
		wt: nil,
	}
}

// vt: SpeedbinbH always "parses" successfully — there's nothing to parse.
func (s SpeedbinbH) vt() bool {
	return true
}

// bt: SpeedbinbH images are never treated as scrambled.
func (s SpeedbinbH) bt(_ image.Rectangle) bool {
	return false
}

// dt: no bounds change for an unscrambled image.
func (s SpeedbinbH) dt(t image.Rectangle) image.Rectangle {
	return t
}

// gt returns a single identity (no-op) copy region covering the whole image.
func (s SpeedbinbH) gt(t image.Rectangle) []DescrambleCord {
	return []DescrambleCord{{
		XSrc:   0,
		YSrc:   0,
		Width:  t.Dx(),
		Height: t.Dy(),
		XDest:  0,
		YDest:  0,
	}}
}
