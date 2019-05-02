package httpclient

import (
	"context"
	"log"
	"testing"
)

func TestClient_Do(t *testing.T) {
	c := New(
		Request(
			NewRequestBuilder(
				URL("https://baidu.com"),
				Method("GET"),
			),
		),
		Response(
			NewDumpResponse(),
		),
	)
	err := c.Execute(context.TODO())
	log.Println("client execute error:", err)
}
