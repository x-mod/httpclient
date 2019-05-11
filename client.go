package httpclient

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/x-mod/errors"
)

var (
	//DefaultMaxConnsPerHost default max connections for per host
	DefaultMaxConnsPerHost = 32
	//DefaultMaxIdleConnsPerHost default max idle connections for per host
	DefaultMaxIdleConnsPerHost = 8
	//DefaultClientTimeout default client timeout for each do request
	DefaultClientTimeout = 30 * time.Second
)

//DefaultTLSConfig default tls.config is nil
var DefaultTLSConfig *tls.Config

//Client struct
type Client struct {
	*http.Client
	config *config
}

//Opt for client
type Opt func(*config)

//Request opt
func Request(builder *RequestBuilder) Opt {
	return func(cf *config) {
		cf.request = builder
	}
}

//Transport opt, When this option is SET, Timeout/MaxConnsPerHost/MaxIdleConnsPerHost option will be IGNORED!!!
func Transport(transport http.RoundTripper) Opt {
	return func(cf *config) {
		cf.transport = transport
	}
}

//Timeout opt
func Timeout(duration time.Duration) Opt {
	return func(cf *config) {
		cf.timeout = duration
	}
}

//Keepalive opt
func Keepalive(keepalive time.Duration) Opt {
	return func(cf *config) {
		cf.keepalive = keepalive
	}
}

//Retry opt for client.Do, only > 1
func Retry(retry int) Opt {
	return func(cf *config) {
		if retry > 1 {
			cf.doRetries = retry
		}
	}
}

//ExecuteRetry opt for client.Execute, only > 1
func ExecuteRetry(retry int) Opt {
	return func(cf *config) {
		if retry > 1 {
			cf.executeRetries = retry
		}
	}
}

//Credential for TLSConfig
func Credential(cred *tls.Config) Opt {
	return func(cf *config) {
		cf.credential = cred
	}
}

//DialContext dialer function
type DialContext func(ctx context.Context, network, addr string) (net.Conn, error)

//Dialer opt
func Dialer(dialer DialContext) Opt {
	return func(cf *config) {
		cf.dialer = dialer
	}
}

//DebugDialer debug dialer
func DebugDialer(ctx context.Context, network, addr string) (net.Conn, error) {
	dial := net.Dialer{
		Timeout: 30 * time.Second,
	}
	conn, err := dial.DialContext(ctx, network, addr)
	if err != nil {
		return conn, err
	}

	log.Println("dailed connection at", conn.LocalAddr().String())
	return conn, err
}

//MaxConnsPerHost opt
func MaxConnsPerHost(max int) Opt {
	return func(cf *config) {
		cf.maxConnsPerHost = max
	}
}

//MaxIdleConnsPerHost opt
func MaxIdleConnsPerHost(max int) Opt {
	return func(cf *config) {
		cf.maxIdleConnsPerHost = max
	}
}

//Response opt
func Response(processor ResponseProcessor) Opt {
	return func(cf *config) {
		cf.response = processor
	}
}

//New client
func New(opts ...Opt) *Client {
	cf := &config{
		doRetries:           1,
		executeRetries:      1,
		timeout:             DefaultClientTimeout,
		maxConnsPerHost:     DefaultMaxConnsPerHost,
		maxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
	}
	for _, opt := range opts {
		opt(cf)
	}
	client := getClient(cf)
	return &Client{config: cf, Client: client}
}

//get client from config
func getClient(cf *config) *http.Client {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   cf.timeout,   //must less then config.timeout
			KeepAlive: cf.keepalive, //zero, keep-alives are enabled
			DualStack: true,
		}).DialContext,
		TLSClientConfig:     cf.credential,
		MaxConnsPerHost:     cf.maxConnsPerHost,
		MaxIdleConnsPerHost: cf.maxIdleConnsPerHost,
	}
	if cf.dialer != nil {
		tr.DialContext = cf.dialer
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   cf.timeout,
	}
	if cf.transport != nil {
		client.Transport = cf.transport
	}
	return client
}

//Close Client release connection resource
func (c *Client) Close() {
	c.CloseIdleConnections()
}

//Execute client
func (c *Client) Execute(ctx context.Context) error {
	if c.config.request == nil {
		return errors.New("request required")
	}

	req, err := c.config.request.Get()
	if err != nil {
		return err
	}

	return c.ExecuteRequest(ctx, req)
}

//ExecuteRequest do custom request with response processor
func (c *Client) ExecuteRequest(ctx context.Context, req *http.Request) (err error) {
	fn := func() error {
		rsp, err := c.DoRequest(ctx, req)
		if err != nil {
			return err
		}

		if c.config.response != nil {
			return c.config.response.Process(ctx, rsp)
		}
		return defaultProcess(ctx, rsp)
	}
	//retries for executing
	for i := 0; i < c.config.executeRetries; i++ {
		if err = fn(); err == nil {
			return
		}
	}
	return
}

//DoRequest do request with context
func (c *Client) DoRequest(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	return c.Do(req.WithContext(ctx))
}

//Do reimpl
func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	//retries for do
	for i := 0; i < c.config.doRetries; i++ {
		if resp, err = c.Client.Do(req); err == nil {
			return
		}
	}
	return
}
