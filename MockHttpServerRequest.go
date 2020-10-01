package golangmockserver

//Mocked request to the server
type MockHttpServerRequest struct{
	//The URI to match
	Uri string
	//The HTTP method to match (GET, PUT, POST, etc)
	Method string
	//Optional: The Body to match against
	Body interface{}
	//Optional: Any headers to match against.  Header values can be a regex
	Headers map[string]string

	//Optional: Response to send from the matched request
	Response *MockHttpServerResponse
}

type MockHttpServerResponse struct {
	//Optional: HTTP status code to return, default is 200
	StatusCode int
	//Optional: body to return, default is nil
	Body       interface{}
	//Optional: Additional headers to return, default is none
	Headers map[string]string
}
