package delete_test

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	delete2 "github.com/tauadam/url-shortener-api/internal/http_server/handlers/alias/delete"
	"github.com/tauadam/url-shortener-api/internal/http_server/handlers/alias/delete/mocks"
	"github.com/tauadam/url-shortener-api/lib/logger/handler/slogdiscard"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name        string
		alias       string
		url         string
		responseErr string
		mockErr     error
	}{
		//{
		//	name:        "Should get url successfully",
		//	alias:       "not_existing_alias",
		//	url:         "https://example.com",
		//	responseErr: "not found",
		//	mockErr:     storage.ErrNotFound,
		//},
		{
			name:        "Should make redirect successfully",
			alias:       "",
			url:         "https://pkg.go.dev",
			responseErr: "redirect: invalid alias",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ShortcutDeleterMock := mocks.NewShortcutDeleter(t)

			if len(tc.responseErr) == 0 || tc.mockErr != nil {
				ShortcutDeleterMock.On("DeleteURL", tc.alias, mock.AnythingOfType("string")).
					Return(tc.mockErr).
					Once()
			}
			handler := delete2.New(slogdiscard.NewDiscardLogger(), ShortcutDeleterMock)

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tc.alias), nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			require.Equal(t, http.StatusOK, rr.Code)
		})
	}
}
