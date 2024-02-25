package request

import (
	"encoding/json"
	"fmt"
	"golang.org/x/image/webp"
	"image"
	"io"
	"net/http"
)

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

func handleResponse(res *http.Response, r interface{}) error {
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("request failed with status: %d", res.StatusCode)
	}

	switch res.Header.Get("content-type") {
	case "application/json", "application/json; charset=utf-8", "application/json; charset=UTF-8":
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &r)
		if err != nil {
			return fmt.Errorf("request: %v", err)
		}
	case "image/png", "image/jpeg":
		switch v := r.(type) {
		case *image.Image:
			img, _, err := image.Decode(res.Body)
			if err != nil {
				return err
			}

			*v = img
		case *[]byte:
			body, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}

			*v = body
		default:
			return fmt.Errorf("unexpected type for image decoding")
		}
	case "image/webp":
		switch v := r.(type) {
		case *image.Image:
			img, err := webp.Decode(res.Body)
			if err != nil {
				return err
			}

			*v = img
		case *[]byte:
			body, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}

			*v = body
		default:
			return fmt.Errorf("unexpected type for image decoding")
		}
	case "text/html", "text/html; charset=utf-8", "text/html; charset=UTF-8":
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		switch v := r.(type) {
		case *string:
			*v = string(body)
		default:
			return fmt.Errorf("unexpected type for HTML response body")
		}
	default:
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		switch v := r.(type) {
		case *[]byte:
			*v = body
		default:
			return fmt.Errorf("unexpected type for response body, response: %d", res.StatusCode)
		}
	}

	return nil
}
