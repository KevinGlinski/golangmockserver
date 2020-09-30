package golangmockserver

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestMockHttpServer_BasicCall(t *testing.T) {

	//Setup
	mockServer := NewMockHttpServer([]*MockHttpServerRequest{
		{
			Uri:      "/foo",
			Method:   "GET",
		},
	})

	request, _ := http.NewRequest("GET", mockServer.BaseUrl() + "/foo", nil)
	client := &http.Client {}

	response, _ := client.Do(request)

	assert.Equal(t, 200, response.StatusCode)
}


func TestMockHttpServer_NotFoundUri(t *testing.T) {

	//Setup
	mockServer := NewMockHttpServer([]*MockHttpServerRequest{
		{
			Uri:      "/foo",
			Method:   "GET",
		},
	})

	request, _ := http.NewRequest("GET", mockServer.BaseUrl() + "/bar", nil)
	client := &http.Client {}

	response, _ := client.Do(request)

	assert.Equal(t, 404, response.StatusCode)
}

func TestMockHttpServer_RequestBodyMatching(t *testing.T) {

	mockServer := NewMockHttpServer([]*MockHttpServerRequest{
		{
			Uri:      "/foo",
			Method:   "GET",
			Response: &MockHttpServerResponse{
				StatusCode: 500,
			},
		},
		{
			Uri:      "/foo",
			Method:   "GET",
			Body: []byte("hello world"),
			Response: &MockHttpServerResponse{
				StatusCode: 201,
			},
		},
	})

	bodyData := bytes.NewReader([]byte("hello world"))
	request, _ := http.NewRequest("GET", mockServer.BaseUrl() + "/foo", bodyData)
	client := &http.Client {}

	response, _ := client.Do(request)

	assert.Equal(t, 201, response.StatusCode)
}


func TestMockHttpServer_ResponseHeaders(t *testing.T) {

	//Setup
	mockServer := NewMockHttpServer([]*MockHttpServerRequest{
		{
			Uri:      "/foo",
			Method:   "GET",
			Response: &MockHttpServerResponse{
				Headers: map[string]string{
					"foo": "bar",
				},
			},
		},
	})

	request, _ := http.NewRequest("GET", mockServer.BaseUrl() + "/foo", nil)
	client := &http.Client {}

	response, _ := client.Do(request)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "bar", response.Header.Get("foo"))
}

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