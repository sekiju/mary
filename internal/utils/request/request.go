package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
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

func Get[T interface{}](url string, args *Config) (T, error) {
	var r T
	if args == nil {
		args = &Config{}
	}

	req, err := newRequest("GET", url, nil, *args)
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

func Post[T interface{}](url string, args *Config) (T, error) {
	var r T
	if args == nil {
		args = &Config{}
	}

	serialized, err := json.Marshal(args.Body)
	if err != nil {
		return r, err
	}

	req, err := newRequest("POST", url, bytes.NewBuffer(serialized), *args)
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
