package httpclient

import (
	"crypto/tls"
	"net/http"
	"time"
)

type config struct {
	request             *RequestBuilder
	response            ResponseProcessor
	transport           http.RoundTripper
	timeout             time.Duration
	keepalive           time.Duration
	credential          *tls.Config
	doRetries           int
	executeRetries      int
	maxConnsPerHost     int
	maxIdleConnsPerHost int
	dialer              DialContext
}

type authConfig struct {
	username string
	password string
}

type bodyConfig struct {
	bodyType   string
	bodyObject interface{}
}

type fileConfig struct {
	Field    string
	FileName string
}

type requestConfig struct {
	Method   string
	URL      string
	Headers  map[string]string
	Queries  map[string]string
	Fragment string
	Cookies  []*http.Cookie
	Auth     *authConfig
	Content  *Body
}

type responseConfig struct {
	StatusCode int
	Content    *Body
}

type transportConfig struct {
	maxIdleConnsPerHost int
	retry               int
	timeout             time.Duration
}
