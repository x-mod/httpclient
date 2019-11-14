package grpc

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	spb "google.golang.org/genproto/googleapis/rpc/status"
)

type statusError struct {
	*spb.Status
}

func (st *statusError) Error() string {
	return st.Message
}
func (st *statusError) Value() int32 {
	return st.Code
}
func (st *statusError) String() string {
	return st.Message
}

type PBJSONProcessor struct {
	data proto.Message
}

func PBJSONResponse(data proto.Message) *PBJSONProcessor {
	return &PBJSONProcessor{data: data}
}

func (pbp *PBJSONProcessor) Process(ctx context.Context, rsp *http.Response) error {
	defer rsp.Body.Close()
	if rsp.StatusCode == 200 {
		if err := jsonpb.Unmarshal(rsp.Body, pbp.data); err != nil {
			return err
		}
		return nil
	}

	st := spb.Status{}
	if err := jsonpb.Unmarshal(rsp.Body, &st); err != nil {
		return err
	}
	return &statusError{Status: &st}
}

type PBProcessor struct {
	data proto.Message
}

func PBResponse(data proto.Message) *PBProcessor {
	return &PBProcessor{data: data}
}

func (pbp *PBProcessor) Process(ctx context.Context, rsp *http.Response) error {
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if rsp.StatusCode == 200 {
		if err := proto.Unmarshal(b, pbp.data); err != nil {
			return err
		}
		return nil
	}

	st := spb.Status{}
	if err := proto.Unmarshal(b, &st); err != nil {
		return err
	}
	return &statusError{Status: &st}
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
