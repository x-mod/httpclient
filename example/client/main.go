package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/x-mod/httpclient"
)

type _id struct{}

func print(ctx context.Context, resp *http.Response) error {
	buf, err := ioutil.ReadAll(resp.Body)
	log.Printf("%d: %s -- %v\n", ctx.Value(_id{}), string(buf), err)
	return resp.Body.Close()
}

func main() {
	client := httpclient.New(
		httpclient.Retry(3),
		httpclient.Keepalive(500*time.Millisecond),
		httpclient.Dialer(httpclient.DebugDialer),
		httpclient.MaxConnsPerHost(4),
		httpclient.MaxIdleConnsPerHost(2),
	)

	req, _ := httpclient.MakeRequest(
		httpclient.SetURL("http://localhost:12345/hello"),
		httpclient.Method("POST"),
		httpclient.Content(
			httpclient.JSON(map[string]interface{}{
				"Hello": "JayL",
			}),
		),
	)

	ctx := context.TODO()
	wg := &sync.WaitGroup{}
	fn := func(id int) {
		defer wg.Done()
		if err := client.Execute(context.WithValue(ctx, _id{}, id), req, httpclient.ResponseProcessorFunc(print)); err != nil {
			log.Println("client execute failed:", err)
		}
	}
	//concurrency to test MaxConnsPerHost
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go fn(i)
		go fn(i)
	}
	wg.Wait()

}

func init() {
	httpclient.DefaultMaxConnsPerHost = 4
	httpclient.DefaultMaxIdleConnsPerHost = 4
}
