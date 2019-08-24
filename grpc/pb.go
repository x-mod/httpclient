package grpc

import (
	"context"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/x-mod/httpclient"
	spb "google.golang.org/genproto/googleapis/rpc/status"
)

type PBResponseFunc func(proto.Message) httpclient.ResponseProcessor

var PBResponse PBResponseFunc

func defaultPBResponse(out proto.Message) httpclient.ResponseProcessor {
	return &PBResponsor{out: out}
}

type PBResponsor struct {
	out proto.Message
}

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

func (pr PBResponsor) Process(ctx context.Context, rsp *http.Response) error {
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusExpectationFailed {
		st := spb.Status{}
		if err := jsonpb.Unmarshal(rsp.Body, &st); err != nil {
			return err
		}
		return &statusError{Status: &st}
	}
	if rsp.StatusCode != 200 {
		st := spb.Status{}
		st.Code = int32(rsp.StatusCode)
		st.Message = rsp.Status
		return &statusError{Status: &st}
	}
	if err := jsonpb.Unmarshal(rsp.Body, pr.out); err != nil {
		return err
	}
	return nil
}

func init() {
	PBResponse = defaultPBResponse
}
