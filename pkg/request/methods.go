package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"net/http/cookiejar"
)

var client = &http.Client{}

func init() {
	jar, _ := cookiejar.New(nil)
	client.Jar = jar
}

func Get[T interface{}](url string, opts ...OptsFn) (Response[T], error) {
	cfg := defaultOpts()
	for _, fn := range opts {
		fn(&cfg)
	}

	var responseWithBody Response[T]
	req, err := newRequest("GET", url, nil, cfg)
	if err != nil {
		return responseWithBody, err
	}

	res, err := client.Do(req)
	responseWithBody.Status = res.StatusCode
	if err != nil {
		return responseWithBody, err
	}

	defer res.Body.Close()

	err = handleResponse(res, &responseWithBody.Body)
	if err != nil {
		return responseWithBody, err
	}

	return responseWithBody, nil
}

func Post[T interface{}](url string, opts ...OptsFn) (Response[T], error) {
	cfg := defaultOpts()
	for _, fn := range opts {
		fn(&cfg)
	}

	var responseWithBody Response[T]

	serialized, err := json.Marshal(cfg.Body)
	if err != nil {
		return responseWithBody, err
	}

	req, err := newRequest("POST", url, bytes.NewBuffer(serialized), cfg)
	if err != nil {
		return responseWithBody, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	responseWithBody.Status = res.StatusCode
	if err != nil {
		return responseWithBody, err
	}

	defer res.Body.Close()

	err = handleResponse(res, &responseWithBody.Body)
	if err != nil {
		return responseWithBody, err
	}

	return responseWithBody, nil
}

func Document(url string, opts ...OptsFn) (*goquery.Document, error) {
	cfg := defaultOpts()
	for _, fn := range opts {
		fn(&cfg)
	}

	req, err := newRequest("GET", url, nil, cfg)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
