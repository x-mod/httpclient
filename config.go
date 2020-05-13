package httpclient

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

type config struct {
	proxy               string
	timeout             time.Duration
	keepalive           time.Duration
	tlsHandsHakeTimeout time.Duration
	credential          *tls.Config
	doRetries           int
	executeRetries      int
	maxConnsPerHost     int
	maxIdleConnsPerHost int
	debug               bool
	dialer              DialContext
	transport           http.RoundTripper
	client              *http.Client
}

type authConfig struct {
	username string
	password string
}

type TokenFunc func() string

type tokenConfig struct {
	token    string
	tokenGet TokenFunc
}

type bodyConfig struct {
	bodyType   string
	bodyObject interface{}
}

type requestConfig struct {
	Method  string
	URL     *url.URL
	Header  http.Header
	Queries map[string]string
	Cookies []*http.Cookie
	Auth    *authConfig
	Token   *tokenConfig
	Content *Body
}
