package golangmockserver

type MockHttpServerRequest struct{
	Uri string
	Method string
	Body interface{}

	Response *MockHttpServerResponse
}

type MockHttpServerResponse struct {
	StatusCode int
	Body       interface{}
	Headers map[string]string
}
