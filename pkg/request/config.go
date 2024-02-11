package request

import "net/http"

const defaultUserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1"

func defaultOpts() Config {
	return Config{
		Headers: map[string]string{
			"user-agent": defaultUserAgent,
		},
		Cookies: make([]*http.Cookie, 0),
	}
}

func Body(body interface{}) OptsFn {
	return func(c *Config) {
		c.Body = body
	}
}

func SetHeader(key, value string) OptsFn {
	return func(c *Config) {
		c.Headers[key] = value
	}
}

func AddCookies(cookies ...*http.Cookie) OptsFn {
	return func(c *Config) {
		c.Cookies = append(c.Cookies, cookies...)
	}
}
