package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_Execute(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.RequestURI {
			case "/head":
				if r.Header.Get("X-HEAD") != "x-head-value" {
					http.Error(w, "head not equal", http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `ok`)
				return
			case "/auth":
				if user, pass, ok := r.BasicAuth(); ok {
					if user == "jay" && pass == "123" {
						w.WriteHeader(http.StatusOK)
						io.WriteString(w, `ok`)
						return
					}
				}
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			case "/error":
				http.Error(w, "error", http.StatusBadRequest)
				return
			case "/sleep":
				time.Sleep(time.Second)
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `sleeped`)
				return
			case "/ping":
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, `ok`)
				return
			default:
				http.NotFound(w, r)
			}
			if r.Header.Get("X-HEAD") != "x-head-value" {
				http.Error(w, "head not equal", http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `ok`)
		}),
	)
	defer ts.Close()

	client1 := New(
		Retry(3),
		ExecuteRetry(3),
		MaxConnsPerHost(2),
		MaxIdleConnsPerHost(2),
	)
	assert.NotNil(t, client1)

	//head
	req1, err := MakeRequest(
		SetURL(ts.URL+"/head"),
		Header("X-HEAD", "xxx"),
	)
	assert.NotNil(t, req1)
	assert.Nil(t, err)

	fn := func(ctx context.Context, rsp *http.Response) error {
		if rsp.StatusCode != http.StatusOK {
			return fmt.Errorf("status != 200")
		}
		return nil
	}
	assert.NotNil(t, client1.Execute(context.TODO(), req1, ResponseProcessorFunc(fn)))

	req2, err := MakeRequest(
		SetURL(ts.URL+"/head"),
		Header("X-HEAD", "x-head-value"),
	)
	assert.NotNil(t, req2)
	assert.Nil(t, err)
	assert.Nil(t, client1.Execute(context.TODO(), req2, nil))

	//auth
	req3, err := MakeRequest(
		SetURL(ts.URL+"/auth"),
		BasicAuth("jay", "123"),
	)
	assert.NotNil(t, req3)
	assert.Nil(t, err)
	assert.Nil(t, client1.Execute(context.TODO(), req3, nil))

	//sleep
	client2 := New(
		Timeout(500 * time.Millisecond),
	)
	assert.NotNil(t, client2)

	//timeout
	req4, err := MakeRequest(
		SetURL(ts.URL + "/sleep"),
	)
	assert.NotNil(t, req4)
	assert.Nil(t, err)
	assert.NotNil(t, client2.Execute(context.TODO(), req4, nil))

	//retry execute
	client3 := New(
		Timeout(500*time.Millisecond),
		ExecuteRetry(3),
	)
	assert.NotNil(t, client3)
	req5, err := MakeRequest(
		SetURL(ts.URL + "/sleep"),
	)
	assert.NotNil(t, req5)
	assert.Nil(t, err)
	assert.NotNil(t, client2.Execute(context.TODO(), req5, nil))

	//ping
	client4 := New()
	assert.NotNil(t, client4)
	req6, err := MakeRequest(
		SetURL(ts.URL + "/ping"),
	)
	assert.NotNil(t, req6)
	assert.Nil(t, err)
	rsp, err := client4.DoRequest(context.TODO(), req6)
	assert.NotNil(t, rsp)
	assert.Nil(t, err)

	dump := NewDumpResponse()
	assert.NotNil(t, dump)
	assert.Nil(t, dump.Process(context.TODO(), rsp))
}
