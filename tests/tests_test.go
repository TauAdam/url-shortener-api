package tests

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
	"github.com/tauadam/url-shortener-api/internal/http_server/handlers/alias/save"
	"github.com/tauadam/url-shortener-api/lib/api"
	"github.com/tauadam/url-shortener-api/lib/random"
	"net/http"
	"net/url"
	"testing"
)

const (
	host = "localhost:5000"
)

func TestUnauthorized(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewString(10),
		}).
		WithBasicAuth("wrong", "wrong").
		Expect().
		Status(http.StatusUnauthorized)
}

func TestAuthorized(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewString(10),
		}).
		WithBasicAuth("admin", "admin").
		Expect().
		Status(http.StatusOK).JSON().Object().ContainsKey("alias")
}

func TestCreateRedirectDelete(t *testing.T) {
	testCases := []struct {
		name        string
		url         string
		alias       string
		responseErr string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:        "Invalid URL",
			url:         "some_invalid_url",
			alias:       gofakeit.Word(),
			responseErr: "URL is not a valid URL",
		},
		{
			name:        "URL with whitespace",
			url:         "whitespace http://example.com/with",
			alias:       gofakeit.Word(),
			responseErr: "URL is not a valid URL",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}
			e := httpexpect.Default(t, u.String())

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("admin", "admin").
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.responseErr != "" {
				resp.NotContainsKey("alias")
				resp.Value("error").String().IsEqual(tc.responseErr)
				return
			}

			alias := tc.alias
			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}
			if tc.alias != "" {
				testRedirect(t, alias, tc.url)
				return
			}
			testInvalidRedirect(t, alias)

			//respDelete := e.DELETE(fmt.Sprintf("/url/%s", alias)).
			//	WithBasicAuth("admin", "admin").Expect().
			//	Status(http.StatusOK).JSON().Object()
			//respDelete.Value("status").String().IsEqual("success")
			//testInvalidRedirect(t, alias)
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedURL, err := api.ProvokeRedirect(u.String())
	require.NoError(t, err)
	require.Equal(t, urlToRedirect, redirectedURL)
}

func testInvalidRedirect(t *testing.T, alias string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}
	_, err := api.ProvokeRedirect(u.String())
	require.ErrorIs(t, err, api.ErrAliasNotFound)
}
