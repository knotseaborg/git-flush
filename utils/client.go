package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Joker interface {
	makeJoke()
	makeArt()
}

type GitJoker struct {
	apiKey   string
	model    string
	endpoint string
}

// RequestPayload defines the structure of the JSON payload
type RequestPayload struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

// ResponsePayload defines the structure of the expected API response
type ResponsePayload struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
	// Add other fields as needed based on the API response
}

func (gj *GitJoker) setup() {
	if err := godotenv.Load(); err != nil {
		logger.Panic("No .env file found, or the env file is of incorrect format")
	}
	gj.apiKey = os.Getenv("API_KEY")
	if gj.apiKey == "" {
		logger.Panic(os.Stderr, "Error: API_KEY environment variable not set")
		os.Exit(1)
	}
	gj.model = os.Getenv("MODEL")
	if gj.model == "" {
		logger.Warn("Warning: LLM Model not found, defaulting to gpt-4-mini")
	}
	gj.endpoint = os.Getenv("ENDPOINT")
	if gj.endpoint == "" {
		fmt.Fprintln(os.Stderr, "Error: ENDPOINT environment variable not set")
	}
}

func (gj *GitJoker) competion(input string) {
	// Define the payload
	payload := RequestPayload{
		Model: gj.model,
		Input: input,
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	// Create HTTP RequestPayload
	req, err := http.NewRequest("POST", gj.endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
		os.Exit(1)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+gj.apiKey)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response: %v\n", err)
		os.Exit(1)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Unexpected status code: %d\nResponse: %s\n", resp.StatusCode, string(body))
		os.Exit(1)
	}

	// Parse response (optional, depending on your needs)
	var responseData ResponsePayload
	if err := json.Unmarshal(body, &responseData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing response JSON: %v\n", err)
		os.Exit(1)
	}

	// Print the response (modify based on actual response structure)
	fmt.Println("Response:", string(body))
}
