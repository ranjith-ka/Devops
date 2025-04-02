package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
)

type RandomJoke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Random joke from the package Devops",
	Long:  `Get a Random joke from the package Devops`,
	Run: func(cmd *cobra.Command, args []string) {
		// getRandomJoke()
		getRandomJokeWithLLMStudio()
	},
}

const url = "https://icanhazdadjoke.com/"

// getRandomJoke func use the method getJokeData to get JSON and unmarshal into string to print in the screen
func getRandomJoke() string {
	fmt.Println("Here is your Joke")
	_, err := emoji.Println(":beer::beer:Beer!!!")
	if err != nil {
		return "joksu"
	}

	responseBytes := getJokeData(url)

	joke := RandomJoke{}

	if err := json.Unmarshal(responseBytes, &joke); err != nil {
		fmt.Printf("Could not unmarshal reponseBytes. %v", err)
	}

	fmt.Println(string(joke.Joke))
	return string(joke.Joke)
}

// getJokeData connect the external URL and get a JSON response from the API
func getJokeData(baseAPI string) []byte {
	request, err := http.NewRequest(
		http.MethodGet, // method
		baseAPI,        // url
		nil,            // body
	)
	if err != nil {
		log.Printf("Could not request a dadjoke. %v", err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "Ranjith KA (https://github.com/ranjith-ka/Docker)")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Could not make a request. %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error here")
		}
	}(response.Body)
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Could not read response body. %v", err)
	}

	return responseBytes
}

// Removed Copilot request and response logic and replaced it with LM Studio local server API implementation.
func getLocalServerResponse(prompt string) string {
	requestBody, err := createLocalServerRequestBody(prompt)
	if err != nil {
		return fmt.Sprintf("Error creating request body: %v", err)
	}

	responseBytes, err := sendLocalServerRequest(requestBody)
	if err != nil {
		return fmt.Sprintf("Error sending request: %v", err)
	}

	response := parseLocalServerResponse(responseBytes)
	if response == "" {
		return "Error: Empty response from local server."
	}

	return response
}

func createLocalServerRequestBody(prompt string) ([]byte, error) {
	messages := []map[string]string{
		{"role": "system", "content": "You are a helpful assistant."},
		{"role": "user", "content": prompt},
	}
	requestBody := map[string]interface{}{
		"model":       "meta-llama-3.1-8b-instruct",
		"messages":    messages,
		"temperature": 0.7,
		"max_tokens":  -1,
		"stream":      false,
	}
	return json.Marshal(requestBody)
}

func sendLocalServerRequest(requestBody []byte) ([]byte, error) {
	url := "http://127.0.0.1:1234/v1/chat/completions"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("api_key", "lm-studio") // Added API key header

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	return responseBytes, nil
}

func parseLocalServerResponse(responseBytes []byte) string {
	var localServerResponse map[string]interface{}
	if err := json.Unmarshal(responseBytes, &localServerResponse); err != nil {
		log.Printf("Could not unmarshal local server response. %v", err)
		return "Error: Failed to parse local server response."
	}

	choices, ok := localServerResponse["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		log.Println("Error: Local server response does not contain valid choices.")
		return ""
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		log.Println("Error: Invalid choice format in local server response.")
		return ""
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		log.Println("Error: Invalid message format in local server response.")
		return ""
	}

	content, ok := message["content"].(string)
	if !ok {
		log.Println("Error: Content missing or invalid in local server response.")
		return ""
	}

	return content
}

// Example usage of the Copilot agent
func getRandomJokeWithLLMStudio() string {
	fmt.Println("Fetching a joke using LLM Studio...")
	prompt := "Tell me a funny programming joke."
	copilotResponse := getLocalServerResponse(prompt)

	if copilotResponse == "" {
		fmt.Println("Failed to fetch a joke. Please try again later.")
		return ""
	}

	fmt.Println("Here's a joke for you:")
	fmt.Println(copilotResponse)
	return copilotResponse
}
