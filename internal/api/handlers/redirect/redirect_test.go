package redirect

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/api/handlers/redirect/mocks"
	"url-shortener/internal/api/response"
	"url-shortener/internal/lib/logger/handlers"
	"url-shortener/internal/storage"
)

func TestRedirectHandler(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		alias     string
		mockError error
		respError string
	}{
		{
			name:  "Valid alias",
			url:   "https://example.com",
			alias: "validAlias",
		},
		{
			name:      "Alias not found",
			url:       "",
			alias:     "notFoundAlias",
			mockError: storage.ErrURLNotFound,
			respError: "not found",
		},
		{
			name:      "Internal error",
			url:       "",
			alias:     "internalErrorAlias",
			mockError: errors.New("internal error"),
			respError: "internal error",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlGetterMock := mocks.NewURLGetter(t)

			if tc.alias != "" {
				urlGetterMock.On("GetURL", tc.alias).Return(tc.url, tc.mockError)
			}

			handler := New(handlers.NewDiscardLogger(), urlGetterMock)

			router := chi.NewRouter()
			router.Get("/{alias}", handler)

			req, err := http.NewRequest(http.MethodGet, "/"+tc.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			body := rr.Body.String()

			var resp response.Response

			if tc.url == "" {
				require.NoError(t, json.Unmarshal([]byte(body), &resp))

				if tc.respError != "" {
					require.Equal(t, tc.respError, resp.Error)
				} else {
					require.Empty(t, resp.Error)
				}
			} else {
				require.Equal(t, tc.url, rr.Header().Get("Location"))
			}
		})
	}
}
