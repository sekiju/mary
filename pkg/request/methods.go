package request

import (
	"bytes"
	"encoding/json"
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

	resp, err := client.Do(req)
	responseWithBody.Status = resp.StatusCode
	if err != nil {
		return responseWithBody, err
	}

	defer resp.Body.Close()

	err = handleResponse(resp, &responseWithBody.Body)
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

	resp, err := client.Do(req)
	responseWithBody.Status = resp.StatusCode
	if err != nil {
		return responseWithBody, err
	}

	defer resp.Body.Close()

	err = handleResponse(resp, &responseWithBody.Body)
	if err != nil {
		return responseWithBody, err
	}

	return responseWithBody, nil
}
