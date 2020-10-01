package golangmockserver

import (
	"bytes"
	"net/http"
)

//Basic example matches the route GET /foo and returns a 200
func ExampleNewMockHttpServer() {

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

	_, _ = client.Do(request)

}

//A 404 is returned when there aren't any matching URIs
func ExampleNewMockHttpServer_uriNotFound() {

	mockServer := NewMockHttpServer([]*MockHttpServerRequest{
		{
			Uri:      "/foo",
			Method:   "GET",
		},
	})
	defer mockServer.Close()

	request, _ := http.NewRequest("GET", mockServer.BaseUrl() + "/bar", nil)
	client := &http.Client {}

	client.Do(request)
	//Returns a 404
}


func ExampleNewMockHttpServer_requestBodyMatching() {

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

	client.Do(request)
}


//This example show how a match based on a header only returns 200 when the header matches.  The first request responds with a 404 because the header isn't set.
func ExampleNewMockHttpServer_headerMatch() {

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

	client.Do(request)
	//returns 404

	request, _ = http.NewRequest("GET", mockServer.BaseUrl() + "/foo", nil)
	request.Header.Add("Authorization", "basic adfljkasdlfj")
	//returns 200

}

//In this example the foo response header is returned.
func ExampleNewMockHttpServer_responseHeaders() {

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

	client.Do(request)

}