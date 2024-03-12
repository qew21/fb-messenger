package assistant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
)

type InputMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Input struct {
	Messages []InputMessage `json:"messages"`
}

type RequestData struct {
	Model string `json:"model"`
	Input Input  `json:"input"`
}

var conversationHistory map[string][]InputMessage

const QianWenUrl = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"

func init() {
	conversationHistory = make(map[string][]InputMessage)
}

type ResponseData struct {
	Output struct {
		Text string `json:"text"`
	} `json:"output"`

	RequestID string `json:"request_id"`
}

func QianWen(userID string, newMessage string, key string) (string, error) {
	updateConversationHistory(userID, newMessage)

	var messages []InputMessage
	if hist, ok := conversationHistory[userID]; ok {
		if len(hist) > 0 && hist[len(hist)-1].Role == "assistant" {
			conversationHistory[userID] = append(hist, InputMessage{Role: "user", Content: newMessage})
			messages = conversationHistory[userID]
		} else {
			return "", fmt.Errorf("failed to find assistant message in conversation history")
		}
	} else {
		messages = []InputMessage{{Role: "user", Content: newMessage}}
	}
	input := Input{Messages: messages}

	requestData := RequestData{
		Model: "qwen-max",
		Input: input,
	}

	jsonPayload, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request data: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, QianWenUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var responseData ResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	output := responseData.Output.Text
	if output != "" {
		latestReply := InputMessage{Role: "assistant", Content: output}
		log.Info().Str("userID", userID).Str("message", newMessage).Msg(output)
		conversationHistory[userID] = append(messages, latestReply)
	}

	return output, nil
}

func updateConversationHistory(userID string, newMessage string) {
	hist, ok := conversationHistory[userID]
	if !ok || len(hist) == 0 {
		hist = []InputMessage{{Role: "system", Content: "You are a helpful assistant."}}
	}
	if len(hist) > 20 {
		hist = hist[len(hist)-20:]
	}

}
