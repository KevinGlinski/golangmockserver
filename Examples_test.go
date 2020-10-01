package golangmockserver

import (
	"bytes"
	"net/http"
)

//Basic example matches the route GET /foo and returns a 200
func ExampleNewMockServer() {

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

	_, _ = client.Do(request)

}

//A 404 is returned when there aren't any matching URIs
func ExampleNewMockServer_uriNotFound() {

	mockServer := NewMockServer([]*MockServerRequest{
		{
			URI:    "/foo",
			Method: "GET",
		},
	})
	defer mockServer.Close()

	request, _ := http.NewRequest("GET", mockServer.BaseURL()+"/bar", nil)
	client := &http.Client{}

	client.Do(request)
	//Returns a 404
}

func ExampleNewMockServer_requestBodyMatching() {

	mockServer := NewMockServer([]*MockServerRequest{
		{
			URI:    "/foo",
			Method: "GET",
			Response: &MockServerResponse{
				StatusCode: 500,
			},
		},
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

	client.Do(request)
}

//This example show how a match based on a header only returns 200 when the header matches.  The first request responds with a 404 because the header isn't set.
func ExampleNewMockServer_headerMatch() {

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

	client.Do(request)
	//returns 404

	request, _ = http.NewRequest("GET", mockServer.BaseURL()+"/foo", nil)
	request.Header.Add("Authorization", "basic adfljkasdlfj")
	//returns 200

}

//In this example the foo response header is returned.
func ExampleNewMockServer_responseHeaders() {

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

	client.Do(request)

}
