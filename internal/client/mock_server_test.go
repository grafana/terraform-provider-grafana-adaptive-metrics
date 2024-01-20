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
	t.Helper()

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

func (s *mockServer) addExpected(
	method, path string, reqHeader http.Header, reqBody []byte,
	respHeader http.Header, respBody []byte,
) {

	s.responses = append(s.responses, &mockServerResponse{
		method:    method,
		path:      path,
		params:    url.Values{},
		reqHeader: reqHeader,
		reqBody:   reqBody,

		statusCode: 200,
		respHeader: respHeader,
		respBody:   respBody,
	})
}
