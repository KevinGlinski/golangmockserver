[![Documentation](https://godoc.org/github.com/yangwenmai/how-to-add-badge-in-github-readme?status.svg)](https://pkg.go.dev/github.com/KevinGlinski/golangmockserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/KevinGlinski/golangmockserver)](https://goreportcard.com/report/github.com/KevinGlinski/golangmockserver)
[![license](https://img.shields.io/github/license/KevinGlinski/golangmockserver.svg?maxAge=2592000)](https://github.com/KevinGlinski/golangmockserver/LICENSE)
[![Release](https://img.shields.io/github/release/KevinGlinski/golangmockserver.svg?label=Release)](https://github.com/KevinGlinski/golangmockserver/releases)

Golang MockServer is a wrapper around httptest.Server and provides helpers to mock out HTTP request/responses

see Examples_test.go for example usage



```
// Create a new mock http server and pass in the methods you want it to match
mockServer := NewMockServer([]*MockServerRequest{
    {
        Uri:      "/foo",
        Method:   "GET",
    },
})
defer mockServer.Close()

//Make request to the mock server
// uses mockServer.BaseUrl() to make the call to the localhost server
request, _ := http.NewRequest("GET", mockServer.BaseUrl() + "/foo", nil)
client := &http.Client {}

response, _ := client.Do(request)

//validate responses
assert.Equal(t, 200, response.StatusCode)
```

Inspired by: https://www.mock-server.com