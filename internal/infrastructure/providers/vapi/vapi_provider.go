package vapi

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/CallPilotReceptionist/internal/domain/errors"
	"github.com/CallPilotReceptionist/internal/domain/providers"
)

type VapiProvider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	webhookSecret string
}

func NewVapiProvider(apiKey, baseURL, webhookSecret string) *VapiProvider {
	return &VapiProvider{
		apiKey:     apiKey,
		baseURL:    baseURL,
		webhookSecret: webhookSecret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (v *VapiProvider) InitiateCall(ctx context.Context, req providers.CallRequest) (*providers.CallSession, error) {
	payload := map[string]interface{}{
		"phoneNumber": req.PhoneNumber,
	}

	if req.AssistantID != "" {
		payload["assistantId"] = req.AssistantID
	}

	if req.AssistantConfig != nil {
		payload["assistant"] = v.convertAssistantConfig(req.AssistantConfig)
	}

	if req.Metadata != nil {
		payload["metadata"] = req.Metadata
	}

	respData, err := v.makeRequest(ctx, "POST", "/call/phone", payload)
	if err != nil {
		return nil, errors.NewProviderError(err, "failed to initiate call")
	}

	session := &providers.CallSession{
		ID:          getString(respData, "id"),
		Status:      getString(respData, "status"),
		PhoneNumber: req.PhoneNumber,
	}

	if startedAt := getString(respData, "startedAt"); startedAt != "" {
		if t, err := time.Parse(time.RFC3339, startedAt); err == nil {
			session.StartedAt = &t
		}
	}

	return session, nil
}

func (v *VapiProvider) HandleWebhook(ctx context.Context, payload []byte, signature string) (*providers.CallEvent, error) {
	if !v.ValidateWebhookSignature(payload, signature) {
		return nil, errors.NewUnauthorizedError("invalid webhook signature")
	}

	var webhookData map[string]interface{}
	if err := json.Unmarshal(payload, &webhookData); err != nil {
		return nil, errors.NewInvalidInputError("invalid webhook payload")
	}

	event := &providers.CallEvent{
		Type:      getString(webhookData, "type"),
		CallID:    getString(webhookData, "callId"),
		Status:    getString(webhookData, "status"),
		Timestamp: time.Now(),
		Data:      webhookData,
	}

	if timestamp := getString(webhookData, "timestamp"); timestamp != "" {
		if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
			event.Timestamp = t
		}
	}

	return event, nil
}

func (v *VapiProvider) GetCallDetails(ctx context.Context, callID string) (*providers.CallDetails, error) {
	respData, err := v.makeRequest(ctx, "GET", fmt.Sprintf("/call/%s", callID), nil)
	if err != nil {
		return nil, errors.NewProviderError(err, "failed to get call details")
	}

	details := &providers.CallDetails{
		ID:          getString(respData, "id"),
		Status:      getString(respData, "status"),
		PhoneNumber: getString(respData, "customer", "number"),
		Duration:    getInt(respData, "duration"),
		Cost:        getFloat(respData, "cost"),
	}

	if startedAt := getString(respData, "startedAt"); startedAt != "" {
		if t, err := time.Parse(time.RFC3339, startedAt); err == nil {
			details.StartedAt = &t
		}
	}

	if endedAt := getString(respData, "endedAt"); endedAt != "" {
		if t, err := time.Parse(time.RFC3339, endedAt); err == nil {
			details.EndedAt = &t
		}
	}

	if metadata, ok := respData["metadata"].(map[string]interface{}); ok {
		details.Metadata = metadata
	}

	return details, nil
}

func (v *VapiProvider) GetTranscript(ctx context.Context, callID string) (*providers.Transcript, error) {
	respData, err := v.makeRequest(ctx, "GET", fmt.Sprintf("/call/%s", callID), nil)
	if err != nil {
		return nil, errors.NewProviderError(err, "failed to get transcript")
	}

	transcript := &providers.Transcript{
		CallID:   callID,
		Messages: []providers.TranscriptMessage{},
	}

	if messages, ok := respData["messages"].([]interface{}); ok {
		for _, msg := range messages {
			if msgMap, ok := msg.(map[string]interface{}); ok {
				message := providers.TranscriptMessage{
					Role:    getString(msgMap, "role"),
					Message: getString(msgMap, "message"),
				}

				if timestamp := getString(msgMap, "timestamp"); timestamp != "" {
					if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
						message.Timestamp = t
					}
				}

				transcript.Messages = append(transcript.Messages, message)
			}
		}
	}

	return transcript, nil
}

func (v *VapiProvider) UpdateAssistantConfig(ctx context.Context, config providers.AssistantConfig) (string, error) {
	payload := v.convertAssistantConfig(&config)

	respData, err := v.makeRequest(ctx, "POST", "/assistant", payload)
	if err != nil {
		return "", errors.NewProviderError(err, "failed to update assistant config")
	}

	return getString(respData, "id"), nil
}

func (v *VapiProvider) GetAssistantConfig(ctx context.Context, assistantID string) (*providers.AssistantConfig, error) {
	respData, err := v.makeRequest(ctx, "GET", fmt.Sprintf("/assistant/%s", assistantID), nil)
	if err != nil {
		return nil, errors.NewProviderError(err, "failed to get assistant config")
	}

	config := &providers.AssistantConfig{
		Name:         getString(respData, "name"),
		Voice:        getString(respData, "voice"),
		Language:     getString(respData, "language"),
		Prompt:       getString(respData, "prompt"),
		FirstMessage: getString(respData, "firstMessage"),
		Model:        getString(respData, "model"),
	}

	if metadata, ok := respData["metadata"].(map[string]interface{}); ok {
		config.Metadata = metadata
	}

	return config, nil
}

func (v *VapiProvider) DeleteAssistantConfig(ctx context.Context, assistantID string) error {
	_, err := v.makeRequest(ctx, "DELETE", fmt.Sprintf("/assistant/%s", assistantID), nil)
	if err != nil {
		return errors.NewProviderError(err, "failed to delete assistant config")
	}
	return nil
}

func (v *VapiProvider) ValidateWebhookSignature(payload []byte, signature string) bool {
	if v.webhookSecret == "" {
		return true // Skip validation if no secret configured
	}

	mac := hmac.New(sha256.New, []byte(v.webhookSecret))
	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

func (v *VapiProvider) makeRequest(ctx context.Context, method, path string, payload interface{}) (map[string]interface{}, error) {
	url := v.baseURL + path

	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+v.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

func (v *VapiProvider) convertAssistantConfig(config *providers.AssistantConfig) map[string]interface{} {
	result := map[string]interface{}{
		"name": config.Name,
	}

	if config.Voice != "" {
		result["voice"] = config.Voice
	}
	if config.Language != "" {
		result["language"] = config.Language
	}
	if config.Prompt != "" {
		result["prompt"] = config.Prompt
	}
	if config.FirstMessage != "" {
		result["firstMessage"] = config.FirstMessage
	}
	if config.Model != "" {
		result["model"] = config.Model
	}
	if config.Metadata != nil {
		result["metadata"] = config.Metadata
	}

	return result
}

// Helper functions
func getString(data map[string]interface{}, keys ...string) string {
	current := data
	for i, key := range keys {
		if i == len(keys)-1 {
			if val, ok := current[key].(string); ok {
				return val
			}
			return ""
		}
		if next, ok := current[key].(map[string]interface{}); ok {
			current = next
		} else {
			return ""
		}
	}
	return ""
}

func getInt(data map[string]interface{}, keys ...string) int {
	current := data
	for i, key := range keys {
		if i == len(keys)-1 {
			if val, ok := current[key].(float64); ok {
				return int(val)
			}
			return 0
		}
		if next, ok := current[key].(map[string]interface{}); ok {
			current = next
		} else {
			return 0
		}
	}
	return 0
}

func getFloat(data map[string]interface{}, keys ...string) float64 {
	current := data
	for i, key := range keys {
		if i == len(keys)-1 {
			if val, ok := current[key].(float64); ok {
				return val
			}
			return 0
		}
		if next, ok := current[key].(map[string]interface{}); ok {
			current = next
		} else {
			return 0
		}
	}
	return 0
}
