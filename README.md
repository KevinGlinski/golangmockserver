Golang MockServer is a wrapper around httptest.Server and provides helpers to mock out HTTP request/responses

see Examples_test.go for example usage



```
//Setup
// Create a new mock http server and pass in the methods you want it to match
mockServer := NewMockHttpServer([]*MockHttpServerRequest{
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