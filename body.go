package httpclient

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/url"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	json "github.com/json-iterator/go"
	"github.com/x-mod/errors"
)

var types = map[string]string{
	"html":       "text/html",
	"json":       "application/json",
	"pb":         "application/octet-stream",
	"pbjson":     "application/json",
	"xml":        "application/xml",
	"text":       "text/plain",
	"binary":     "application/octet-stream",
	"urlencoded": "application/x-www-form-urlencoded",
	"form":       "application/x-www-form-urlencoded",
	"form-data":  "application/x-www-form-urlencoded",
	"multipart":  "multipart/form-data",
}

//Body struct
type Body struct {
	config *bodyConfig
}

//BodyOpt type
type BodyOpt func(*bodyConfig)

//Text opt
func Text(txt string) BodyOpt {
	return func(cf *bodyConfig) {
		cf.bodyType = "text"
		cf.bodyObject = txt
	}
}

//Binary opt
func Binary(bytes []byte) BodyOpt {
	return func(cf *bodyConfig) {
		cf.bodyType = "binary"
		cf.bodyObject = bytes
	}
}

//JSON opt
func JSON(obj interface{}) BodyOpt {
	return func(cf *bodyConfig) {
		cf.bodyType = "json"
		cf.bodyObject = obj
	}
}

//PB opt
func PB(obj proto.Message) BodyOpt {
	return func(cf *bodyConfig) {
		cf.bodyType = "pb"
		cf.bodyObject = obj
	}
}

//PBJSON opt
func PBJSON(obj proto.Message) BodyOpt {
	return func(cf *bodyConfig) {
		cf.bodyType = "pbjson"
		cf.bodyObject = obj
	}
}

//XML opt
func XML(obj map[string]interface{}) BodyOpt {
	return func(cf *bodyConfig) {
		cf.bodyType = "xml"
		cf.bodyObject = obj
	}
}

//Form opt
func Form(obj url.Values) BodyOpt {
	return func(cf *bodyConfig) {
		cf.bodyType = "form"
		cf.bodyObject = obj
	}
}

//Reader opt
func Reader(rd io.Reader) BodyOpt {
	return func(cf *bodyConfig) {
		cf.bodyType = "reader"
		cf.bodyObject = rd
	}
}

//Get Body io.Reader
func (b *Body) Get() (io.Reader, error) {
	if b.config != nil {
		switch strings.ToLower(b.config.bodyType) {
		case "text":
			return bytes.NewBufferString(b.config.bodyObject.(string)), nil
		case "binary":
			return bytes.NewBuffer(b.config.bodyObject.([]byte)), nil
		case "json":
			byts, err := json.Marshal(b.config.bodyObject.(map[string]interface{}))
			if err != nil {
				return nil, errors.Annotate(err, "json marshal failed")
			}
			return bytes.NewBuffer(byts), nil
		case "pb":
			byts, err := proto.Marshal(b.config.bodyObject.(proto.Message))
			if err != nil {
				return nil, errors.Annotate(err, "pb marshal failed")
			}
			return bytes.NewBuffer(byts), nil
		case "pbjson":
			wr := bytes.NewBuffer([]byte{})
			marshaler := &jsonpb.Marshaler{EmitDefaults: true}
			if err := marshaler.Marshal(wr, b.config.bodyObject.(proto.Message)); err != nil {
				return nil, errors.Annotate(err, "pbjson marshal failed")
			}
			return wr, nil
		case "xml":
			byts, err := xml.Marshal(b.config.bodyObject)
			if err != nil {
				return nil, errors.Annotate(err, "xml marshal failed")
			}
			return bytes.NewBuffer(byts), nil
		case "form":
			data := b.config.bodyObject.(url.Values).Encode()
			return strings.NewReader(data), nil
		case "reader":
			return b.config.bodyObject.(io.Reader), nil
		}
	}
	return bytes.NewBuffer([]byte{}), nil
}

//ContentType Body Content-Type
func (b *Body) ContentType() string {
	if b.config != nil {
		if v, ok := types[strings.ToLower(b.config.bodyType)]; ok {
			return v
		}
	}
	return types["html"]
}
