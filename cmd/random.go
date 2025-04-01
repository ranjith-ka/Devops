package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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
        getRandomJoke()
        getRandomJokeWithCopilot()
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

// getCopilotResponse connects to the GitHub Copilot chat completions API and retrieves a response based on the input prompt
func getCopilotResponse(prompt string) string {
    githubToken := os.Getenv("GITHUB_TOKEN")
    if githubToken == "" {
        log.Println("GITHUB_TOKEN not found. Cannot connect to Copilot service.")
        return "Error: Missing GITHUB_TOKEN for authentication."
    }

    requestBody, err := createCopilotRequestBody(prompt)
    if err != nil {
        return err.Error()
    }

    responseBytes, err := sendCopilotRequest(requestBody, githubToken)
    if err != nil {
        return err.Error()
    }

    return parseCopilotResponse(responseBytes)
}

func createCopilotRequestBody(prompt string) ([]byte, error) {
    messages := []map[string]string{
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": prompt},
    }
    requestBody := map[string]interface{}{
        "messages": messages,
        "stream":   false,
    }
    return json.Marshal(requestBody)
}

func sendCopilotRequest(requestBody []byte, githubToken string) ([]byte, error) {
    copilotAPI := "https://api.githubcopilot.com/chat/completions"
    request, err := http.NewRequest(http.MethodPost, copilotAPI, bytes.NewBuffer(requestBody))
    if err != nil {
        log.Printf("Could not create Copilot request. %v", err)
        return nil, fmt.Errorf("error: failed to create Copilot request")
    }
    request.Header.Add("Content-Type", "application/json")
    request.Header.Add("Authorization", "Bearer "+githubToken)

    log.Printf("Copilot Request Body: %s", string(requestBody))

    response, err := http.DefaultClient.Do(request)
    if err != nil {
        log.Printf("Could not make a request to the Copilot service. %v", err)
        return nil, fmt.Errorf("error: failed to connect to Copilot service")
    }
    defer response.Body.Close()

    responseBody, err := io.ReadAll(response.Body)
    if err != nil {
        log.Printf("Could not read Copilot response body. %v", err)
        return nil, fmt.Errorf("error: failed to read Copilot response body")
    }

    log.Printf("Copilot Response Body: %s", string(responseBody))

    if response.StatusCode != http.StatusOK {
        log.Printf("Copilot request failed with status code: %d", response.StatusCode)
        return nil, fmt.Errorf("error: copilot service returned status code %d", response.StatusCode)
    }

    return responseBody, nil
}

func parseCopilotResponse(responseBytes []byte) string {
    var copilotResponse map[string]interface{}
    if err := json.Unmarshal(responseBytes, &copilotResponse); err != nil {
        log.Printf("Could not unmarshal Copilot response. %v", err)
        return "Error: Failed to parse Copilot response."
    }

    if choices, ok := copilotResponse["choices"].([]interface{}); ok && len(choices) > 0 {
        if choice, ok := choices[0].(map[string]interface{}); ok {
            if message, ok := choice["message"].(map[string]interface{}); ok {
                if content, ok := message["content"].(string); ok {
                    return content
                }
            }
        }
    }

    return "Error: Copilot response did not contain valid text."
}

// Example usage of the Copilot agent
func getRandomJokeWithCopilot() string {
    fmt.Println("Fetching a joke using Copilot...")
    prompt := "Tell me a funny programming joke."
    copilotResponse := getCopilotResponse(prompt)
    fmt.Println("Copilot Response:", copilotResponse)
    return copilotResponse
}