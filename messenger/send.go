package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/qew21/fb-messenger/config"
)

type Payload struct {
	Recipient     Recipient `json:"recipient"`
	Message       Message   `json:"message"`
	MessagingType string    `json:"messaging_type"`
}

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	Text string `json:"text"`
}

func SendMessage(psid string, messageText string, appConfig *config.AppConfig) error {
	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/messages", appConfig.APIVersion, appConfig.PageID)
	accessToken := appConfig.PageAccesToken
	recipientData := Recipient{ID: psid}
	messageData := Message{Text: messageText}
	payload := Payload{
		Recipient:     recipientData,
		Message:       messageData,
		MessagingType: "RESPONSE",
	}
	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("Failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("Failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response body: %w", err)
	}

	return fmt.Errorf("Failed to send message: %s", body)
}
