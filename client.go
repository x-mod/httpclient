package httpclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
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
	//DefaultTLSHandhakeTimeout default client tls hands hake timeout
	DefaultTLSHandhakeTimeout = 10 * time.Second
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

//Proxy opt
func Proxy(host string) Opt {
	return func(cf *config) {
		cf.proxy = host
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

//Transport opt, When this option is SET, Timeout/MaxConnsPerHost/MaxIdleConnsPerHost option will be IGNORED!!!
func Transport(transport http.RoundTripper) Opt {
	return func(cf *config) {
		cf.transport = transport
	}
}

//Debug opt
func Debug(flag bool) Opt {
	return func(cf *config) {
		cf.debug = flag
	}
}

//HTTPClient opt, when this option is SET, All above option will ignore
func HTTPClient(client *http.Client) Opt {
	return func(cf *config) {
		cf.client = client
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
		tlsHandsHakeTimeout: DefaultTLSHandhakeTimeout,
	}
	for _, opt := range opts {
		opt(cf)
	}
	client := getClient(cf)
	return &Client{config: cf, Client: client}
}

//get client from config
func getClient(cf *config) *http.Client {
	if cf.client != nil {
		return cf.client
	}
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   cf.timeout,   //must less then config.timeout
			KeepAlive: cf.keepalive, //zero, keep-alives are enabled
			DualStack: true,
		}).DialContext,
		TLSClientConfig:     cf.credential,
		TLSHandshakeTimeout: cf.tlsHandsHakeTimeout,
		MaxConnsPerHost:     cf.maxConnsPerHost,
		MaxIdleConnsPerHost: cf.maxIdleConnsPerHost,
	}
	if cf.dialer != nil {
		tr.DialContext = cf.dialer
	}
	if len(cf.proxy) > 0 {
		if u, err := url.Parse(cf.proxy); err == nil {
			tr.Proxy = http.ProxyURL(u)
		}
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

//GetTransport get standard http.RoundTripper Transport
func (c *Client) GetTransport() http.RoundTripper {
	return c.Client.Transport
}

//GetClient get standard http.Client
func (c *Client) GetClient() *http.Client {
	return c.Client
}

//Execute client
func (c *Client) Execute(ctx context.Context, req *http.Request, processor ResponseProcessor) error {
	if req == nil {
		return errors.New("request required")
	}
	if c.config.debug {
		if b, err := httputil.DumpRequest(req, true); err == nil {
			fmt.Println("------------------------------")
			fmt.Println(string(b))
			fmt.Println("------------------------------")
		}
	}
	rsp, err := c.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	if c.config.debug {
		if b, err := httputil.DumpResponse(rsp, true); err == nil {
			fmt.Println(string(b))
			fmt.Println("------------------------------")
		}
	}

	if processor != nil {
		return processor.Process(ctx, rsp)
	}
	return nil
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
