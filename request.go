package httpclient

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/x-mod/errors"
)

//RequestBuilder struct
type RequestBuilder struct {
	config *requestConfig
}

//ReqOpt opt
type ReqOpt func(*requestConfig)

//Method opt
func Method(method string) ReqOpt {
	return func(cf *requestConfig) {
		cf.Method = strings.ToUpper(method)
	}
}

//URL opt
func SetURL(Url string) ReqOpt {
	return func(cf *requestConfig) {
		u, err := url.Parse(Url)
		if err != nil {
			panic(err)
		}
		cf.URL = u
	}
}

//URL opt
func URL(opts ...URLOpt) ReqOpt {
	return func(cf *requestConfig) {
		if cf.URL == nil {
			cf.URL = &url.URL{}
		}
		for _, opt := range opts {
			opt(cf.URL)
		}
	}
}

type URLOpt func(*url.URL)

//URI opt
func URI(uri string) URLOpt {
	return func(u *url.URL) {
		u.Path = uri
	}
}

//Scheme opt
func Scheme(scheme string) URLOpt {
	return func(u *url.URL) {
		u.Scheme = scheme
	}
}

//User opt
func User(user string) URLOpt {
	return func(u *url.URL) {
		u.User = url.User(user)
	}
}

//UserPassword opt
func UserPassword(username string, password string) URLOpt {
	return func(u *url.URL) {
		u.User = url.UserPassword(username, password)
	}
}

//Host opt [ip:port]
func Host(host string) URLOpt {
	return func(u *url.URL) {
		u.Host = host
	}
}

//Fragment opt
func Fragment(name string) URLOpt {
	return func(u *url.URL) {
		u.Fragment = name
	}
}

//Query opt
func Query(name string, value string) ReqOpt {
	return func(cf *requestConfig) {
		cf.Queries[name] = value
	}
}

//Header opt
func Header(name string, value string) ReqOpt {
	return func(cf *requestConfig) {
		cf.Header.Add(name, value)
	}
}

//Cookie opt
func Cookie(cookie *http.Cookie) ReqOpt {
	return func(cf *requestConfig) {
		if cookie != nil {
			cf.Cookies = append(cf.Cookies, cookie)
		}
	}
}

//BasicAuth opt
func BasicAuth(username string, password string) ReqOpt {
	return func(cf *requestConfig) {
		cf.Auth = &authConfig{
			username: username,
			password: password,
		}
	}
}

//BearerAuth opt
func BearerAuth(token string) ReqOpt {
	return func(cf *requestConfig) {
		cf.Token = &tokenConfig{
			token: token,
		}
	}
}

//BearerAuthFunc opt
func BearerAuthFunc(token TokenFunc) ReqOpt {
	return func(cf *requestConfig) {
		cf.Token = &tokenConfig{
			tokenGet: token,
		}
	}
}

//Content opt
func Content(opts ...BodyOpt) ReqOpt {
	return func(cf *requestConfig) {
		body := &bodyConfig{}
		for _, opt := range opts {
			opt(body)
		}
		cf.Content = &Body{config: body}
	}
}

//NewRequestBuilder new
func NewRequestBuilder(opts ...ReqOpt) *RequestBuilder {
	config := &requestConfig{
		Queries: make(map[string]string),
		Cookies: []*http.Cookie{},
	}
	for _, opt := range opts {
		opt(config)
	}
	return &RequestBuilder{config: config}
}

//MakeRequest make a http.Request
func MakeRequest(opts ...ReqOpt) (*http.Request, error) {
	builder := NewRequestBuilder(opts...)
	return builder.makeRequest()
}

func (req *RequestBuilder) makeRequest() (*http.Request, error) {
	if req.config.URL == nil {
		return nil, errors.New("url required")
	}
	//url
	if len(req.config.Queries) > 0 {
		q := req.config.URL.Query()
		for k, v := range req.config.Queries {
			q.Add(k, v)
		}
		req.config.URL.RawQuery = q.Encode()
	}

	//body
	var body io.Reader
	if req.config.Content != nil {
		rd, err := req.config.Content.Get()
		if err != nil {
			return nil, err
		}
		body = rd
	}

	//new request
	rr, err := http.NewRequest(req.config.Method, req.config.URL.String(), body)
	if err != nil {
		return nil, err
	}
	rr.Header = req.config.Header

	// content-type
	if req.config.Content != nil {
		rr.Header.Set("Content-Type", req.config.Content.ContentType())
	}
	// cookies
	for _, v := range req.config.Cookies {
		rr.AddCookie(v)
	}
	// basic auth
	if req.config.Auth != nil {
		rr.SetBasicAuth(req.config.Auth.username, req.config.Auth.password)
	}
	// bearer auth
	if req.config.Token != nil {
		t := req.config.Token.token
		if req.config.Token.tokenGet != nil {
			t = req.config.Token.tokenGet()
		}
		rr.Header.Set("Authorization", strings.Join([]string{"Bearer", t}, " "))
	}
	return rr, nil
}
