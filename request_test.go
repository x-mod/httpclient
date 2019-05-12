package httpclient

import (
	"testing"
)

func TestBody_Get(t *testing.T) {
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Body{
				config: tt.fields.config,
			}
			_, err := b.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("Body.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewRequestBuilder(t *testing.T) {
	type args struct {
		opts []ReqOpt
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"options",
			args{
				[]ReqOpt{},
			},
			true,
		},
		{
			"options",
			args{
				[]ReqOpt{
					URL("http://aa"),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRequestBuilder(tt.args.opts...).Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestBuilder.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
