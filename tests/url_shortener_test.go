package tests

import (
	"net/http"
	"net/url"
	"testing"
	"url-shortener/internal/http-server/api"
	"url-shortener/pkg/random"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/url/save"
)

const (
	host = "localhost:8081"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.GetRandomAlias(10),
		}).
		WithBasicAuth("admin", "admin").
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		ContainsKey("alias")
}

func TestURLShortener_SaveRedirect(t *testing.T) {
	testCases := []struct {
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
			error: "field URL is not a valid URL",
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())
			var resp *httpexpect.Object
			if tc.error == "" {
				resp = e.POST("/url").
					WithJSON(save.Request{
						URL:   tc.url,
						Alias: tc.alias,
					}).
					WithBasicAuth("admin", "admin").
					Expect().Status(http.StatusCreated).
					JSON().Object()
			} else {
				resp = e.POST("/url").
					WithJSON(save.Request{
						URL:   tc.url,
						Alias: tc.alias,
					}).
					WithBasicAuth("admin", "admin").
					Expect().Status(http.StatusBadRequest).
					JSON().Object()
			}

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().NotEmpty()

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			testRedirect(t, alias, tc.url)
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}

func TestURLSaveAndDelete(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())
	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: "test",
		}).
		WithBasicAuth("admin", "admin").
		Expect()

	e.DELETE("/url/test").
		WithBasicAuth("admin", "admin").
		Expect().
		Status(http.StatusOK)
}

func TestURLDeleteNotFound(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())
	e.DELETE("/url/test").
		WithBasicAuth("admin", "admin").
		Expect().
		Status(http.StatusNotFound)
}

func TestURLRedirectNotFound(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "test",
	}
	e := httpexpect.Default(t, u.String())
	e.GET("/test23").
		Expect().
		Status(http.StatusNotFound)
}
