package golangmockserver

import (
	"reflect"
	"testing"
)


func TestMockHttpServer_toBytes(t *testing.T) {

	type testStruct struct{
		Foo string `json:"foo"`
	}

	type args struct {
		data interface{}
	}
	tests := []struct {
		name   string
		data interface{}
		want   []byte
	}{
		{
			"bytes",
			[]byte("abcd"),
			[]byte("abcd"),

		},
		{
			"string",
			"abcd",
			[]byte("abcd"),

		},
		{
			"struct",
			testStruct{"bar"},
			[]byte("{\"foo\":\"bar\"}"),

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MockHttpServer{}
			if got := s.toBytes(tt.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}