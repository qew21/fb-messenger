package analysis

import (
	"testing"

	"github.com/qew21/fb-messenger/config"
	"github.com/stretchr/testify/assert"
)

func TestSentimentAnalysis(t *testing.T) {
	testCases := []struct {
		Text     string
		Expected string
	}{
		{"This film is terrible!", "negative"},
		{"This film is great!", "positive"},
		{"This film is not terrible!", "positive"},
		{"This film is not great!", "negative"},
		{"Where can I find this film?", "neutral"},
	}

	appConfig, _ := config.LoadConfig("../config.yaml")
	for _, tc := range testCases {
		sentiment, err := Sentiment(appConfig.PredictUrl, tc.Text)
		if err != nil {
			t.Errorf("Failed to analyze sentiment for text '%s': %s", tc.Text, err)
			continue
		}

		assert.Equal(t, tc.Expected, sentiment, "Unexpected sentiment for text '%s'", tc.Text)
	}
}
