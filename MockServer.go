// golangmockserver

//Golang MockServer is a wrapper around httptest.Server and provides helpers to mock out HTTP request/responses

package golangmockserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
)

//NewMockServer creates a new Mock server with mocked requests
func NewMockServer(mockedRequests []*MockServerRequest) *MockServer {
	server := &MockServer{
		requests: mockedRequests,
	}

	handler := http.NewServeMux()

	registeredUris := make([]string, 0)

	for _, request := range mockedRequests {

		alreadyRegistered := false

		for _, uri := range registeredUris {
			if uri == request.URI {
				alreadyRegistered = true
			}
		}

		if !alreadyRegistered {
			handler.HandleFunc(request.URI, server.handleRequest)
			registeredUris = append(registeredUris, request.URI)
		}
	}

	server.server = httptest.NewServer(handler)

	return server
}

// MockServer is a mock http server
type MockServer struct {
	requests []*MockServerRequest
	server   *httptest.Server
}

//Close shuts down the mock http server
func (s *MockServer) Close() {
	s.server.Close()
}

//BaseURL gets the localhost base url to call
func (s *MockServer) BaseURL() string {
	return s.server.URL
}

//Requests gets the configured requests
func (s *MockServer) Requests() []*MockServerRequest {
	return s.requests
}

func (s *MockServer) toBytes(data interface{}) []byte {
	switch data.(type) {
	case string:
		return []byte(data.(string))
	case []byte:
		return data.([]byte)
	default:
		marshaled, _ := json.Marshal(data)
		return marshaled
	}
}

func (s *MockServer) doHeadersMatch(request *MockServerRequest, r *http.Request) bool {
	//match headers
	if request.Headers != nil {
		for k, v := range request.Headers {
			headerValue := r.Header.Get(k)
			if !regexp.MustCompile(v).MatchString(headerValue) {
				return false
			}
		}
	}

	return true
}


func (s *MockServer) doQueryParametersMatch(request *MockServerRequest, r *http.Request) bool {
	//match query parameters
	if request.QueryParameters != nil {
		for k, v := range request.QueryParameters {
			paramValue := r.FormValue(k)
			if !regexp.MustCompile(v).MatchString(paramValue) {
				return false
			}
		}
	}

	return true
}

func (s *MockServer) handleRequest(w http.ResponseWriter, r *http.Request) {

	requestdata, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	for _, request := range s.requests {
		if request.URI == r.URL.Path &&
			request.Method == r.Method &&
			(request.Body == nil || (request.Body != nil && reflect.DeepEqual(s.toBytes(request.Body), requestdata))) {

			if !s.doHeadersMatch(request, r) {
				continue
			}

			if !s.doQueryParametersMatch(request, r) {
				continue
			}

			if request.MaxMatchCount > 0 && request.invokeCount >= request.MaxMatchCount {
				continue
			}

			request.invokeCount++

			if request.Response == nil {
				_, _ = w.Write(nil)
				return
			}

			if request.Response.Headers != nil {
				for k, v := range request.Response.Headers {
					w.Header().Add(k, v)
				}
			}

			if request.Response.StatusCode != 0 {
				w.WriteHeader(request.Response.StatusCode)
			}

			_, _ = w.Write(s.toBytes(request.Response.Body))

			return

		}
	}

	w.WriteHeader(404)
}
