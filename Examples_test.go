package golangmockserver

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
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
	defer mockServer.Close()

	//Make request to the mock server
	request, _ := http.NewRequest("GET", mockServer.BaseUrl() + "/foo", nil)
	client := &http.Client {}

	response, _ := client.Do(request)

	//validate responses
	assert.Equal(t, 200, response.StatusCode)
}


func TestMockHttpServer_NotFoundUri(t *testing.T) {

	mockServer := NewMockHttpServer([]*MockHttpServerRequest{
		{
			Uri:      "/foo",
			Method:   "GET",
		},
	})
	defer mockServer.Close()

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
	defer mockServer.Close()

	bodyData := bytes.NewReader([]byte("hello world"))
	request, _ := http.NewRequest("GET", mockServer.BaseUrl() + "/foo", bodyData)
	client := &http.Client {}

	response, _ := client.Do(request)

	assert.Equal(t, 201, response.StatusCode)
}



func TestMockHttpServer_HeaderMatch(t *testing.T) {

	mockServer := NewMockHttpServer([]*MockHttpServerRequest{
		{
			Uri:      "/foo",
			Method:   "GET",
			Headers: map[string]string{
				"Authorization": "basic .*",
			},
			Response: &MockHttpServerResponse{
				StatusCode: 200,
			},
		},
	})
	defer mockServer.Close()

	request, _ := http.NewRequest("GET", mockServer.BaseUrl() + "/foo", nil)

	client := &http.Client {}

	response, _ := client.Do(request)

	assert.Equal(t, 404, response.StatusCode)


	request, _ = http.NewRequest("GET", mockServer.BaseUrl() + "/foo", nil)
	request.Header.Add("Authorization", "basic adfljkasdlfj")

	response, _ = client.Do(request)
	assert.Equal(t, 200, response.StatusCode)

}

func TestMockHttpServer_ResponseHeaders(t *testing.T) {

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
	defer mockServer.Close()

	request, _ := http.NewRequest("GET", mockServer.BaseUrl() + "/foo", nil)
	client := &http.Client {}

	response, _ := client.Do(request)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "bar", response.Header.Get("foo"))
}