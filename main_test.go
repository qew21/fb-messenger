package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/qew21/fb-messenger/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestConfiguration() {
	if appConfig == nil {
		appConfig, _ = config.LoadConfig("config.yaml")
	}
}

func TestHandlePostWebHook(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
	}{
		{"Message", "message.json"},
		{"Feed", "feed.json"},
		{"Ratings", "ratings.json"},
	}
	setupTestConfiguration()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonFile, err := os.Open(filepath.Join("mock", tc.filename))
			require.NoError(t, err)
			defer jsonFile.Close()

			byteValue, err := ioutil.ReadAll(jsonFile)
			require.NoError(t, err)

			var payload map[string]interface{}
			err = json.Unmarshal(byteValue, &payload)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/test", bytes.NewReader(byteValue))
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			handlePostWebHook(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Contains(t, recorder.Body.String(), "Webhook processed successfully.")
		})
	}
}
