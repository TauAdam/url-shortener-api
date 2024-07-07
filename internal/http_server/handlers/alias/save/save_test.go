package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"vigilant-octo-spoon/internal/http_server/handlers/alias/save"
	"vigilant-octo-spoon/internal/http_server/handlers/alias/save/mocks"
	"vigilant-octo-spoon/lib/logger/handler/slogdiscard"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name        string
		alias       string
		url         string
		responseErr string
		mockError   error
	}{
		{
			name:  "Should Save Shortcut Successfully",
			alias: "new_test_alias",
			url:   "https://example.com",
		},
		{
			name:  "Should Save Shortcut Successfully with new alias",
			alias: "",
			url:   "https://example.com",
		},
		{
			name:        "Should Return Error For Empty URL",
			url:         "",
			alias:       "new_alias",
			responseErr: "URL is required",
		},
		{
			name:        "Should Return Error For Invalid URL",
			url:         "some invalid URL",
			alias:       "new_alias",
			responseErr: "URL is not a valid URL",
		},
		{
			name:        "Should Return Error For Already Existing url",
			alias:       "new_test_alias",
			url:         "https://example.com",
			responseErr: "failed to save url",
			mockError:   errors.New("URL already exists"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewShortcutSaver(t)

			if tc.responseErr == "" || tc.mockError != nil {
				urlSaverMock.On("SaveShortcut", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"alias": "%s", "url": "%s"}`, tc.alias, tc.url)
			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()
			var resp save.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.responseErr, resp.Error)
		})
	}
}
