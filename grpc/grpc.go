package grpc

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/x-mod/httpclient"
	"github.com/x-mod/tlsconfig"
)

type HTTPClient struct {
	*httpclient.Client
	cfg *HTTPClientCfg
}

type HTTPClientCfg struct {
	version     string
	packageName string
	serviceName string
	host        string
	schema      string
	tls         *tls.Config
	debug       bool
}

type HTTPClientOpt func(*HTTPClient)

func Version(version string) HTTPClientOpt {
	return func(c *HTTPClient) {
		c.cfg.version = version
	}
}
func PackageName(pkg string) HTTPClientOpt {
	return func(c *HTTPClient) {
		c.cfg.packageName = pkg
	}
}
func ServiceName(svc string) HTTPClientOpt {
	return func(c *HTTPClient) {
		c.cfg.serviceName = svc
	}
}
func Schema(schema string) HTTPClientOpt {
	return func(c *HTTPClient) {
		c.cfg.schema = schema
	}
}
func Host(host string) HTTPClientOpt {
	return func(c *HTTPClient) {
		c.cfg.host = host
	}
}
func TLSConfig(opts ...tlsconfig.Option) HTTPClientOpt {
	return func(c *HTTPClient) {
		if len(opts) > 0 {
			c.cfg.tls = tlsconfig.New(opts...)
		}
	}
}
func Debug(v bool) HTTPClientOpt {
	return func(c *HTTPClient) {
		c.cfg.debug = v
	}
}

func NewHTTPClient(opts ...HTTPClientOpt) *HTTPClient {
	c := &HTTPClient{
		cfg: &HTTPClientCfg{
			schema: "http",
			host:   "127.0.0.1",
		},
	}
	for _, o := range opts {
		o(c)
	}
	copts := []httpclient.Opt{
		httpclient.Keepalive(30 * time.Second),
		httpclient.MaxConnsPerHost(32),
		httpclient.MaxIdleConnsPerHost(16),
	}
	if c.cfg.debug {
		copts = append(copts, httpclient.Debug(true))
	}
	client := httpclient.New(copts...)
	c.Client = client
	return c
}

func (c *HTTPClient) MakeRequest(methodName string, opts ...httpclient.ReqOpt) (*http.Request, error) {
	copts := []httpclient.ReqOpt{
		httpclient.Method("post"),
		httpclient.URL(
			httpclient.Host(c.cfg.host),
			httpclient.Schema(c.cfg.schema),
			httpclient.URI(URIFormat(c.cfg.version, c.cfg.packageName, c.cfg.serviceName, methodName)),
		),
	}
	copts = append(copts, opts...)
	return httpclient.MakeRequest(copts...)
}

//default URIFormat: /v1/pkg.Service/Method
func defaultURIFormat(version string, pkg string, service string, method string) string {
	return fmt.Sprintf("/%s/%s.%s/%s", version, pkg, service, method)
}

type URIFormatFunc func(version string, pkg string, service string, method string) string

var URIFormat URIFormatFunc

func init() {
	URIFormat = defaultURIFormat
}
