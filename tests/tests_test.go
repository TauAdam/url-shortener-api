package tests

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
	"vigilant-octo-spoon/internal/http_server/handlers/alias/save"
	"vigilant-octo-spoon/lib/api"
	"vigilant-octo-spoon/lib/random"
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
		Status(401)
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
		Status(200).JSON().Object().ContainsKey("alias")
}

func TestCreateRedirectDelete(t *testing.T) {
	testCases := []struct {
		name        string
		url         string
		alias       string
		responseErr string
	}{
		{
			name:  "Without alias",
			url:   gofakeit.URL(),
			alias: "",
		},
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
				resp.Value("responseErr").String().IsEqual(tc.responseErr)
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

	redirectedURL, err := api.ProvokeRedirect(u.String())
	require.NoError(t, err)
	require.Equal(t, urlToRedirect, redirectedURL)
}
