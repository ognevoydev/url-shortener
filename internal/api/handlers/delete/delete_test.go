package delete

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/api/handlers/delete/mocks"
	"url-shortener/internal/api/response"
	"url-shortener/internal/lib/logger/handlers"
)

func TestDeleteHandler(t *testing.T) {
	tests := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
		},
		{
			name:      "Failed to delete",
			alias:     "test_alias",
			respError: "delete url error",
			mockError: errors.New("internal error"),
		},
		{
			name:      "Alias not found",
			alias:     "non_existing_alias",
			respError: "",
			mockError: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlRemoverMock := mocks.NewURLRemover(t)

			if tc.respError == "" || tc.mockError != nil {
				urlRemoverMock.On("DeleteURL", tc.alias).
					Return(tc.mockError).
					Once()
			}

			handler := New(handlers.NewDiscardLogger(), urlRemoverMock)

			router := chi.NewRouter()
			router.Delete("/{alias}", handler)

			req, err := http.NewRequest(http.MethodDelete, "/"+tc.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			body := rr.Body.String()

			var resp response.Response

			if tc.alias != "" {
				require.Equal(t, rr.Code, http.StatusOK)
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
			}

			if tc.respError != "" {
				require.Equal(t, tc.respError, resp.Error)
			} else {
				require.Empty(t, resp.Error)
			}
		})
	}
}
