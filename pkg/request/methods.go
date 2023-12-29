package request

import (
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/http/cookiejar"
)

const defaultUserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1"

var client = &http.Client{}

type Config struct {
	Headers map[string]string
	Body    interface{}
	Cookies []*http.Cookie
}

func init() {
	jar, _ := cookiejar.New(nil)
	client.Jar = jar
}

func Get[T interface{}](url string, opts *Config) (T, error) {
	var r T
	if opts == nil {
		opts = defaultConfig()
	}

	req, err := newRequest("GET", url, nil, *opts)
	if err != nil {
		return r, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return r, err
	}

	err = handleResponse(resp, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}

func Post[T interface{}](url string, opts *Config) (T, error) {
	var r T
	if opts == nil {
		opts = defaultConfig()
	}

	serialized, err := json.Marshal(opts.Body)
	if err != nil {
		return r, err
	}

	req, err := newRequest("POST", url, bytes.NewBuffer(serialized), *opts)
	if err != nil {
		return r, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return r, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}

func GetDocument(url string, opts *Config) (*goquery.Document, error) {
	if opts == nil {
		opts = defaultConfig()
	}

	req, err := newRequest("GET", url, nil, *opts)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
