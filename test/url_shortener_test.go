package test

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
	"url-shortener/internal/api/handlers/save"
	"url-shortener/internal/lib/random"
)

const (
	host = "localhost:8080"
)

func TestURLShortener_Save(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.POST("/save").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.String(6),
		}).
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("alias")
}

func TestURLShortener_SaveRedirect(t *testing.T) {
	tests := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: gofakeit.Word(),
			error: "failed to validate request URL",
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// Save
			resp := e.POST("/save").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")
				resp.Value("error").String().IsEqual(tc.error)
				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()
				alias = resp.Value("alias").String().Raw()
			}

			// Redirect
			u = url.URL{
				Scheme: "http",
				Host:   host,
				Path:   alias,
			}

			client := new(http.Client)
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}

			response, err := client.Get(u.String())
			if err != nil {
				t.Fail()
			}

			redirectUrl := response.Header.Get("Location")
			require.Equal(t, tc.url, redirectUrl)
		})
	}
}
