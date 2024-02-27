package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockServerResponse struct {
	method    string
	path      string
	params    url.Values
	reqHeader http.Header
	reqBody   []byte

	statusCode int
	respHeader http.Header
	respBody   []byte
}

type mockServer struct {
	server    *httptest.Server
	responses []*mockServerResponse
}

func newMockServer(t *testing.T) *mockServer {
	s := &mockServer{}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(s.responses) == 0 {
			t.Fatalf("unexpected request: %+v", r)
		}

		next := s.responses[0]
		s.responses = s.responses[1:]

		require.Equal(t, next.method, r.Method)
		require.Equal(t, next.path, r.URL.Path)
		require.Equal(t, next.params, r.URL.Query())

		for k, v := range next.reqHeader {
			require.Equal(t, v, r.Header.Values(k))
		}

		if next.reqBody == nil {
			next.reqBody = []byte{}
		}
		actualReqBody, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, next.reqBody, actualReqBody)

		for k, vs := range next.respHeader {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(next.statusCode)
		if next.respBody != nil {
			_, err = w.Write(next.respBody)
			require.NoError(t, err)
		}
	}))

	return s
}

func (s *mockServer) close() {
	s.server.Close()
}

func (s *mockServer) addExpected(method, path string, options ...mockRequestOption) {
	resp := &mockServerResponse{
		method:     method,
		path:       path,
		params:     url.Values{},
		statusCode: 200,
	}

	for _, opt := range options {
		opt(resp)
	}

	s.responses = append(s.responses, resp)
}

type mockRequestOption func(*mockServerResponse)

func withReqHeader(h http.Header) mockRequestOption {
	return func(r *mockServerResponse) {
		r.reqHeader = h
	}
}

func withReqBody(b []byte) mockRequestOption {
	return func(r *mockServerResponse) {
		r.reqBody = b
	}
}

func withRespHeader(h http.Header) mockRequestOption {
	return func(r *mockServerResponse) {
		r.respHeader = h
	}
}

func withRespBody(b []byte) mockRequestOption {
	return func(r *mockServerResponse) {
		r.respBody = b
	}
}

func withParams(p url.Values) mockRequestOption {
	return func(r *mockServerResponse) {
		r.params = p
	}
}
