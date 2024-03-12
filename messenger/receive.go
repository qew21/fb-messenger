package messenger

import (
	"fmt"

	"github.com/qew21/fb-messenger/analysis"
	"github.com/qew21/fb-messenger/assistant"
	"github.com/qew21/fb-messenger/config"
)

func getMessageText(messageMap map[string]interface{}) string {
	messagePart := messageMap["message"].(map[string]interface{})
	return messagePart["text"].(string)
}

func getMessageSender(messageMap map[string]interface{}) string {
	messagePart := messageMap["sender"].(map[string]interface{})
	return messagePart["id"].(string)
}

func analyzeSentimentBasedOnFieldType(field string, predictUrl string, value map[string]interface{}) (string, string, error) {
	var sentiment string
	var senderID string

	switch field {
	case "feed":
		senderID = value["post_id"].(string)
		sentiment, err := analysis.Sentiment(predictUrl, value["message"].(string))
		if err != nil {
			return sentiment, senderID, fmt.Errorf("failed to analyze feed sentiment: %w", err)
		}
	case "ratings":
		senderID = value["comment_id"].(string)
		recommendationType := value["recommendation_type"].(string)
		if recommendationType == "POSITIVE" {
			sentiment = "positive"
		} else {
			sentiment = "negative"
		}
	default:
		return "", "", nil
	}

	return sentiment, senderID, nil
}

func ProcessMessage(update map[string]interface{}, appConfig *config.AppConfig, testMode bool) error {
	switch update["object"].(string) {
	case "page":
		pageEntries := update["entry"].([]interface{})

		for _, entry := range pageEntries {
			entryMap := entry.(map[string]interface{})
			changesOrMessaging, ok := entryMap["changes"]

			if ok {
				changes := changesOrMessaging.([]interface{})
				for _, change := range changes {
					changeMap := change.(map[string]interface{})
					field := changeMap["field"].(string)
					value := changeMap["value"].(map[string]interface{})

					sentiment, senderID, err := analyzeSentimentBasedOnFieldType(field, appConfig.PredictUrl, value)
					if err != nil {
						return fmt.Errorf("failed to analyze %s sentiment: %w from %s", sentiment, err, senderID)
					}
					if !testMode {
						if sentiment == "positive" {
							SendMessage(senderID, "We're so glad to hear that! Could you share more about what you enjoyed?", appConfig)
						} else if sentiment == "negative" {
							SendMessage(senderID, "We're sorry to hear that. Could you share more about what went wrong?", appConfig)
						}
					}
				}
			} else if messaging, ok := entryMap["messaging"]; ok {
				messages := messaging.([]interface{})
				messageMap := messages[0].(map[string]interface{})
				textValue := getMessageText(messageMap)
				senderID := getMessageSender(messageMap)
				if !testMode {
					reply, _ := assistant.QianWen(senderID, textValue, appConfig.QianwenKey)
					if reply != "" {
						SendMessage(senderID, reply, appConfig)
					}
				}
			}
		}
	default:
		return fmt.Errorf("unknown object type in payload")
	}

	return nil
}
