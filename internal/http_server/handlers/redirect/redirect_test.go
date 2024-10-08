package redirect_test

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tauadam/url-shortener-api/internal/http_server/handlers/redirect"
	"github.com/tauadam/url-shortener-api/internal/http_server/handlers/redirect/mocks"
	"github.com/tauadam/url-shortener-api/lib/api"
	"github.com/tauadam/url-shortener-api/lib/logger/handler/slogdiscard"
	"net/http/httptest"
	"testing"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name        string
		alias       string
		url         string
		responseErr string
		mockErr     error
	}{{
		name:  "Should get url successfully",
		alias: "new_test_alias",
		url:   "https://example.com",
	},
		{
			name:  "Should make redirect successfully",
			alias: "pkg",
			url:   "https://pkg.go.dev",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if len(c.responseErr) == 0 || c.mockErr != nil {
				urlGetterMock.On("GetURL", c.alias, mock.AnythingOfType("string")).
					Return(c.url, c.mockErr).
					Once()
			}
			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedURL, err := api.ProvokeRedirect(fmt.Sprintf("%s/%s", ts.URL, c.alias))
			require.NoError(t, err)

			require.Equal(t, c.url, redirectedURL)
		})
	}
}
