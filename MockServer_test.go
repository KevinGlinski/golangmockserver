package golangmockserver

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

type testStruct struct {
	Foo string `json:"foo"`
}

func TestMockServer_BasicCall(t *testing.T) {

	//Setup
	mockServer := NewMockServer([]*MockServerRequest{
		{
			URI:    "/foo",
			Method: "GET",
		},
	})
	defer mockServer.Close()

	//Make request to the mock server
	request, _ := http.NewRequest("GET", mockServer.BaseURL()+"/foo", nil)
	client := &http.Client{}

	response, _ := client.Do(request)

	//validate responses
	assert.Equal(t, 200, response.StatusCode)

	assert.Equal(t, 1, mockServer.Requests()[0].InvokeCount())
}


func TestMockServer_PUTWithAnyBody(t *testing.T) {

	//Setup
	mockServer := NewMockServer([]*MockServerRequest{
		{
			URI:    "/foo",
			Method: "PUT",
		},
	})
	defer mockServer.Close()

	//Make request to the mock server
	bodyData := bytes.NewReader([]byte("hello world"))
	request, _ := http.NewRequest("PUT", mockServer.BaseURL()+"/foo", bodyData)
	client := &http.Client{}

	response, _ := client.Do(request)

	//validate responses
	assert.Equal(t, 200, response.StatusCode)

	assert.Equal(t, 1, mockServer.Requests()[0].InvokeCount())
}

func TestMockServer_NotFoundUri(t *testing.T) {

	mockServer := NewMockServer([]*MockServerRequest{
		{
			URI:    "/foo",
			Method: "GET",
		},
	})
	defer mockServer.Close()

	request, _ := http.NewRequest("GET", mockServer.BaseURL()+"/bar", nil)
	client := &http.Client{}

	response, _ := client.Do(request)

	assert.Equal(t, 404, response.StatusCode)
}

func TestMockServer_RequestBodyMatching(t *testing.T) {

	mockServer := NewMockServer([]*MockServerRequest{
		{
			URI:    "/foo",
			Method: "GET",
			Body:   []byte("hello world"),
			Response: &MockServerResponse{
				StatusCode: 201,
			},
		},
	})
	defer mockServer.Close()

	bodyData := bytes.NewReader([]byte("hello world"))
	request, _ := http.NewRequest("GET", mockServer.BaseURL()+"/foo", bodyData)
	client := &http.Client{}

	response, _ := client.Do(request)

	assert.Equal(t, 201, response.StatusCode)
}

func TestMockServer_HeaderMatch(t *testing.T) {

	mockServer := NewMockServer([]*MockServerRequest{
		{
			URI:    "/foo",
			Method: "GET",
			Headers: map[string]string{
				"Authorization": "basic .*",
			},
			Response: &MockServerResponse{
				StatusCode: 200,
			},
		},
	})
	defer mockServer.Close()

	request, _ := http.NewRequest("GET", mockServer.BaseURL()+"/foo", nil)

	client := &http.Client{}

	response, _ := client.Do(request)

	assert.Equal(t, 404, response.StatusCode)

	request, _ = http.NewRequest("GET", mockServer.BaseURL()+"/foo", nil)
	request.Header.Add("Authorization", "basic adfljkasdlfj")

	response, _ = client.Do(request)
	assert.Equal(t, 200, response.StatusCode)

}


func TestMockServer_MaximumMatch(t *testing.T) {

	mockServer := NewMockServer([]*MockServerRequest{
		{
			URI:    "/foo",
			Method: "GET",
			MaxMatchCount: 1,
			Response: &MockServerResponse{
				StatusCode: 500,
			},
		},
		{
			URI:    "/foo",
			Method: "GET",
			Response: &MockServerResponse{
				StatusCode: 200,
			},
		},
	})
	defer mockServer.Close()

	request, _ := http.NewRequest("GET", mockServer.BaseURL()+"/foo", nil)

	client := &http.Client{}

	response, _ := client.Do(request)

	assert.Equal(t, 500, response.StatusCode)

	request, _ = http.NewRequest("GET", mockServer.BaseURL()+"/foo", nil)
	response, _ = client.Do(request)
	assert.Equal(t, 200, response.StatusCode)


	request, _ = http.NewRequest("GET", mockServer.BaseURL()+"/foo", nil)
	response, _ = client.Do(request)
	assert.Equal(t, 200, response.StatusCode)

}



func TestMockServer_QueryStringParamsMatch(t *testing.T) {

	mockServer := NewMockServer([]*MockServerRequest{
		{
			URI:    "/foo",
			Method: "GET",
			QueryParameters: map[string]string{
				"page": "1",
			},
			Response: &MockServerResponse{
				StatusCode: 200,
			},
		},
	})
	defer mockServer.Close()

	request, _ := http.NewRequest("GET", mockServer.BaseURL()+"/foo", nil)

	client := &http.Client{}

	response, _ := client.Do(request)

	assert.Equal(t, 404, response.StatusCode)

	request, _ = http.NewRequest("GET", mockServer.BaseURL()+"/foo?page=1", nil)

	response, _ = client.Do(request)
	assert.Equal(t, 200, response.StatusCode)

}

func TestMockServer_ResponseHeaders(t *testing.T) {

	mockServer := NewMockServer([]*MockServerRequest{
		{
			URI:    "/foo",
			Method: "GET",
			Response: &MockServerResponse{
				Headers: map[string]string{
					"foo": "bar",
				},
			},
		},
	})
	defer mockServer.Close()

	request, _ := http.NewRequest("GET", mockServer.BaseURL()+"/foo", nil)
	client := &http.Client{}

	response, _ := client.Do(request)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "bar", response.Header.Get("foo"))
}

func TestMockServer_toBytes(t *testing.T) {



	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		data interface{}
		want []byte
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
			s := &MockServer{}
			if got := s.toBytes(tt.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
