package golangmockserver

//MockServerRequest is a mocked request to the server
type MockServerRequest struct {
	//The URI to match
	URI string
	//The HTTP method to match (GET, PUT, POST, etc)
	Method string
	//Optional: The Body to match against
	Body interface{}
	//Optional: Any headers to match against.  Header values can be a regex
	Headers map[string]string

	//Optional: Any querystring parameters to match against.  Parameter values can be a regex
	QueryParameters map[string]string

	//Optional: Response to send from the matched request
	Response *MockServerResponse

	invokeCount int
}

// Invoke count returns the number of times this request was invoked
func (m *MockServerRequest) InvokeCount() int {
	return m.invokeCount
}

//MockServerResponse is the response from a mocked request.
type MockServerResponse struct {
	//Optional: HTTP status code to return, default is 200
	StatusCode int
	//Optional: body to return, default is nil
	Body interface{}
	//Optional: Additional headers to return, default is none
	Headers map[string]string
}
