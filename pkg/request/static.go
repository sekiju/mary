package request

import "net/http"

type Config struct {
	Headers map[string]string
	Body    interface{}
	Cookies []*http.Cookie
}

type OptsFn func(*Config)

type Response[T interface{}] struct {
	Status int
	Body   T
}
