package golangmockserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
)


type MockHttpServer struct {
	requests []*MockHttpServerRequest
	server  *httptest.Server

}

func NewMockHttpServer(mockedRequests []*MockHttpServerRequest) *MockHttpServer {
	server := &MockHttpServer{
		requests: mockedRequests,
	}

	handler := http.NewServeMux()

	registeredUris := make([]string, 0)

	for _, request := range mockedRequests {

		alreadyRegistered := false

		for _, uri := range registeredUris {
			if uri == request.Uri {
				alreadyRegistered = true
			}
		}

		if !alreadyRegistered{
			handler.HandleFunc(request.Uri, server.handleRequest)
			registeredUris = append(registeredUris, request.Uri)
		}
	}

	server.server = httptest.NewServer(handler)

	return server
}


func (s *MockHttpServer) Close() {
	s.server.Close()
}


func (s *MockHttpServer) BaseUrl() string{
	return s.server.URL
}


func (s *MockHttpServer) toBytes(data interface{}) []byte {
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

func (s *MockHttpServer) handleRequest(w http.ResponseWriter, r *http.Request) {

	respdata, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()


	for _, request := range s.requests {
		if request.Uri == r.URL.Path &&
			request.Method == r.Method &&
			(request.Body != nil && reflect.DeepEqual(s.toBytes(request.Body), respdata)) {

			if request.Response == nil {
				_, _ = w.Write(nil)
				return
			}else{
				if request.Response.StatusCode != 0 {
					w.WriteHeader(request.Response.StatusCode)
				}

				if request.Response.Headers != nil {
					for k, v := range request.Response.Headers {
						w.Header().Add(k,v)
					}
				}

				_, _ = w.Write(s.toBytes(request.Response.Body))
			}

			return
		}
	}
}