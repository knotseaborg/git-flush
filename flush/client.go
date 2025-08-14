package flush

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type LLMClient struct {
	apiKey   string
	model    string
	endpoint string
}

type RequestPayload struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type ResponsePayload struct {
	Output []struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

func (c *LLMClient) Complete(input string) (string, int, error) {
	req, err := c.MakeRequestPayload(input)
	if err != nil {
		logger.Error("Error: Failed to make request payload; ", err)
		return "", 0, err
	}

	respData, err := c.MakeRequest(req)
	if err != nil {
		logger.Error("Error: Failed to make http request; ", err)
		return "", 0, err
	}

	return respData.Output[0].Content[0].Text,
		respData.Usage.TotalTokens,
		nil
}

func (c *LLMClient) MakeRequest(req *http.Request) (*ResponsePayload, error) {
	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error making request: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error reading response: ", err)
		return nil, err
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		logger.Error(
			"Unexpected status code: ", resp.StatusCode,
			"\nResponse: ", string(body),
		)
		return nil, err
	}

	// Parse response (optional, depending on your needs)
	var responseData ResponsePayload
	if err := json.Unmarshal(body, &responseData); err != nil {
		logger.Error("Failed to parse response JSON: ", err)
		return nil, err
	}
	return &responseData, nil
}

func (c *LLMClient) MakeRequestPayload(input string) (*http.Request, error) {
	// Convert payload to JSON
	payloadBytes, err := json.Marshal(
		RequestPayload{
			Model: c.model, Input: input,
		},
	)
	if err != nil {
		logger.Error("Failed to marshal JSON: ", err)
		return nil, err
	}
	// Create HTTP RequestPayload
	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		logger.Error("Failed to create request: ", err)
		return nil, err
	}
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	return req, nil
}

func InitLLMClient() *LLMClient {
	if config.APIKey == "" {
		logger.Error("API Key not set in config file. Use `git-flush --config` to set your API Key")
	}
	return &LLMClient{
		config.APIKey,
		config.Model,
		config.EndPoint,
	}
}
