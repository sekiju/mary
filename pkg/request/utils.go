package request

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
	"path"
	"strings"
)

func JoinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}

func defaultConfig() *Config {
	return &Config{
		Headers: map[string]string{
			"user-agent": defaultUserAgent,
		},
		Cookies: make([]*http.Cookie, 0),
	}
}

func newRequest(method, url string, body io.Reader, c Config) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return req, err
	}

	req.Header.Add("User-Agent", defaultUserAgent)
	if c.Headers != nil {
		for key, value := range c.Headers {
			req.Header.Add(key, value)
		}
	}

	for _, cookie := range c.Cookies {
		req.AddCookie(cookie)
	}

	return req, nil
}

func handleResponse(resp *http.Response, r interface{}) error {
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	switch resp.Header.Get("content-type") {
	case "application/json", "application/json; charset=utf-8":
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &r)
		if err != nil {
			return err
		}
	case "image/png", "image/jpeg":
		img, _, err := image.Decode(resp.Body)
		if err != nil {
			return err
		}

		switch v := r.(type) {
		case *image.Image:
			*v = img
		default:
			return fmt.Errorf("unexpected type for image decoding")
		}
	default:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		switch v := r.(type) {
		case *[]byte:
			*v = body
		default:
			return fmt.Errorf("unexpected type for response body")
		}
	}

	return nil
}
