package request

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest_Parse(t *testing.T) {
	defaultIP := net.ParseIP("192.168.0.1")
	defaultHTTPRequest := func(url string) *http.Request {
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req.RemoteAddr = defaultIP.String() + ":80"

		return req
	}

	testCases := []struct {
		name           string
		request        func() *http.Request
		expectedResult *Request
		expectedError  error
	}{
		{
			name:           "Count not passed",
			request:        func() *http.Request { return defaultHTTPRequest("/") },
			expectedResult: NewRequest(defaultCount, 0, defaultIP),
		},
		{
			name:           "Count valid",
			request:        func() *http.Request { return defaultHTTPRequest("/?count=10") },
			expectedResult: NewRequest(10, 0, defaultIP),
		},
		{
			name:          "Count over maximum",
			request:       func() *http.Request { return defaultHTTPRequest("/?count=101") },
			expectedError: ErrInvalidParameterValue,
		},
		{
			name:          "Count invalid",
			request:       func() *http.Request { return defaultHTTPRequest("/?count=test") },
			expectedError: ErrInvalidParameterValue,
		},
		{
			name:           "Offset not passed",
			request:        func() *http.Request { return defaultHTTPRequest("/") },
			expectedResult: NewRequest(defaultCount, 0, defaultIP),
		},
		{
			name:           "Offset valid",
			request:        func() *http.Request { return defaultHTTPRequest("/?offset=10") },
			expectedResult: NewRequest(defaultCount, 10, defaultIP),
		},
		{
			name:          "Offset over maximum",
			request:       func() *http.Request { return defaultHTTPRequest("/?offset=10001") },
			expectedError: ErrInvalidParameterValue,
		},
		{
			name:          "Offset invalid",
			request:       func() *http.Request { return defaultHTTPRequest("/?offset=test") },
			expectedError: ErrInvalidParameterValue,
		},
		{
			name:           "Valid IP from remote address field",
			request:        func() *http.Request { return defaultHTTPRequest("/") },
			expectedResult: NewRequest(defaultCount, 0, defaultIP),
		},
		{
			name: "Valid IP from X-FORWARDED-FOR header",
			request: func() *http.Request {
				req := defaultHTTPRequest("/")
				req.Header.Add(forwardedForHeaderName, "8.8.8.8")

				return req
			},
			expectedResult: NewRequest(defaultCount, 0, net.ParseIP("8.8.8.8")),
		},
		{
			name: "Valid IP from X-FORWARDED-FOR header: multiple IPs",
			request: func() *http.Request {
				req := defaultHTTPRequest("/")
				req.Header.Add(forwardedForHeaderName, "1.1.1.1, 8.8.8.8")

				return req
			},
			expectedResult: NewRequest(defaultCount, 0, net.ParseIP("1.1.1.1")),
		},
		{
			name: "Invalid X-FORWARDED-FOR header",
			request: func() *http.Request {
				req := defaultHTTPRequest("/")
				req.Header.Add(forwardedForHeaderName, "test")

				return req
			},
			expectedResult: NewRequest(defaultCount, 0, defaultIP),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			request := Request{}
			err := request.Parse(test.request())

			// check error returned
			if test.expectedError != nil && test.expectedError != err {
				t.Fatalf("err check failed: expected to get '%v', but got '%v'", test.expectedError, err)
			}

			if test.expectedResult != nil {
				// check count
				if request.Count != test.expectedResult.Count {
					t.Fatalf("count check failed: expected to get '%v', but got '%v'", test.expectedResult.Count, request.Count)
				}

				// check offset
				if request.Offset != test.expectedResult.Offset {
					t.Fatalf("offset check failed: expected to get '%v', but got '%v'", test.expectedResult.Offset, request.Offset)
				}

				// check user IP
				if request.UserIP.String() != test.expectedResult.UserIP.String() {
					t.Fatalf("user IP check failed: expected to get '%v', but got '%v'", test.expectedResult.UserIP.String(), request.UserIP.String())
				}
			}
		})
	}
}
