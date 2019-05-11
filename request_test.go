package httpclient

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestBody_Get(t *testing.T) {
	obj := map[string]interface{}{
		"a": "xxx",
		"b": 2,
		"c": false,
	}
	byts, err := json.Marshal(obj)
	if err != nil {
		t.Errorf("Body.Get() prepare obj error = %v, wantErr %v", err, nil)
		return
	}
	type fields struct {
		config *bodyConfig
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "json",
			fields: fields{
				&bodyConfig{
					bodyType: "json",
					bodyObject: map[string]interface{}{
						"a": "xxx",
						"b": 2,
						"c": false,
					},
				},
			},
			want:    byts,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Body{
				config: tt.fields.config,
			}
			got, err := b.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("Body.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotByts, err := ioutil.ReadAll(got)
			if !reflect.DeepEqual(gotByts, tt.want) {
				t.Errorf("Body.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
