package golangmockserver

import (
	"bytes"
	"net/http"
)

func ExampleF__BasicCall() {

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


func ExampleF__NotFoundUri() {

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
}

func ExampleF_RequestBodyMatching() {

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



func ExampleF_HeaderMatch() {

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

func ExampleF_ResponseHeaders() {

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