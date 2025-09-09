package imgboard

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchHtml(t *testing.T) {
	testCases := []struct {
		name        string
		reqPath     string
		respBody    string
		respStatus  int
		expectedRes string
		expectedErr error
	}{
		{
			name:        "200 ok",
			reqPath:     "/image.png",
			respBody:    "html",
			respStatus:  200,
			expectedRes: "html",
			expectedErr: nil,
		},
		{
			name:        "404 not found",
			reqPath:     "/missing.png",
			respBody:    "Not found",
			respStatus:  404,
			expectedRes: "",
			expectedErr: ErrFetchFailed,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, testCase.reqPath, r.URL.Path)
				assert.Equal(t, "GET", r.Method)
				w.WriteHeader(testCase.respStatus)
				_, _ = w.Write([]byte(testCase.respBody))
			}))
			defer server.Close()
			requestUrl := server.URL + testCase.reqPath
			res, err := FetchHtml(requestUrl)
			assert.Equal(t, res, testCase.expectedRes)
			assert.ErrorIs(t, err, testCase.expectedErr)
		})
	}
}
