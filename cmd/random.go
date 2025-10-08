package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type RandomJoke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

var tracer trace.Tracer

// initTracing initializes OpenTelemetry tracing
func initTracing() func(context.Context) error {
	ctx := context.Background()

	// Get Tempo endpoint from environment or use default
	tempoEndpoint := os.Getenv("TEMPO_ENDPOINT")
	if tempoEndpoint == "" {
		tempoEndpoint = "tempo-ingest.example.com:80" // Use your ingress endpoint (hostname:port format)
	}

	// Create resource
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("devops-cli"),
			semconv.ServiceVersion("1.0.0"),
			semconv.DeploymentEnvironment("development"),
		),
	)
	if err != nil {
		log.Printf("Failed to create resource: %v", err)
		return func(context.Context) error { return nil }
	}

	// Create OTLP HTTP exporter
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(tempoEndpoint),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithHeaders(map[string]string{
			"Content-Type": "application/json",
		}),
	)
	if err != nil {
		log.Printf("Failed to create OTLP exporter: %v", err)
		return func(context.Context) error { return nil }
	}

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithMaxExportBatchSize(512),
			sdktrace.WithBatchTimeout(5*time.Second),
		),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Create tracer
	tracer = otel.Tracer("devops-cli")

	log.Printf("✅ Tracing initialized - sending to: %s", tempoEndpoint)

	return tp.Shutdown
}

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Random joke from the package Devops",
	Long:  `Get a Random joke from the package Devops`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize tracing
		shutdown := initTracing()
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := shutdown(ctx); err != nil {
				log.Printf("Failed to shutdown tracing: %v", err)
			}
		}()

		// Create root span for the command
		ctx, span := tracer.Start(context.Background(), "random_command",
			trace.WithAttributes(
				attribute.String("command", "random"),
				attribute.StringSlice("args", args),
			),
		)
		defer span.End()

		// Add event
		span.AddEvent("command_started")

		// Execute the joke fetching functions with tracing context
		getRandomJoke(ctx)
		getRandomJokeWithLLMStudio(ctx)

		span.AddEvent("command_completed")
	},
}

const url = "https://icanhazdadjoke.com/"

// getRandomJoke func use the method getJokeData to get JSON and unmarshal into string to print in the screen
func getRandomJoke(ctx context.Context) string {
	// Create span for this function
	ctx, span := tracer.Start(ctx, "get_random_joke",
		trace.WithAttributes(
			attribute.String("joke.source", "icanhazdadjoke.com"),
		),
	)
	defer span.End()

	fmt.Println("Here is your Joke")
	span.AddEvent("joke_request_started")

	_, err := emoji.Println(":beer::beer:Beer!!!")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(500, "Failed to print emoji")
		return "joksu"
	}

	responseBytes := getJokeData(ctx, url)

	joke := RandomJoke{}

	if err := json.Unmarshal(responseBytes, &joke); err != nil {
		fmt.Printf("Could not unmarshal reponseBytes. %v", err)
		span.RecordError(err)
		span.SetStatus(500, "Failed to unmarshal joke response")
		return ""
	}

	fmt.Println(string(joke.Joke))

	// Add attributes to span
	span.SetAttributes(
		attribute.String("joke.id", joke.ID),
		attribute.String("joke.content", joke.Joke),
		attribute.Int("joke.status", joke.Status),
	)

	span.AddEvent("joke_fetched_successfully")

	return string(joke.Joke)
}

// getJokeData connect the external URL and get a JSON response from the API
func getJokeData(ctx context.Context, baseAPI string) []byte {
	// Create span for HTTP request
	ctx, span := tracer.Start(ctx, "http_get_joke",
		trace.WithAttributes(
			attribute.String("http.method", "GET"),
			attribute.String("http.url", baseAPI),
		),
	)
	defer span.End()

	request, err := http.NewRequestWithContext(ctx,
		http.MethodGet, // method
		baseAPI,        // url
		nil,            // body
	)
	if err != nil {
		log.Printf("Could not request a dadjoke. %v", err)
		span.RecordError(err)
		span.SetStatus(500, "Failed to create request")
		return nil
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "Ranjith KA (https://github.com/ranjith-ka/Docker)")

	// Use instrumented HTTP client
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	span.AddEvent("http_request_started")
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Could not make a request. %v", err)
		span.RecordError(err)
		span.SetStatus(500, "HTTP request failed")
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error here")
		}
	}(response.Body)

	// Add response attributes
	span.SetAttributes(
		attribute.Int("http.status_code", response.StatusCode),
		attribute.String("http.status", response.Status),
	)

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Could not read response body. %v", err)
		span.RecordError(err)
		span.SetStatus(500, "Failed to read response body")
		return nil
	}

	span.SetAttributes(
		attribute.Int("response.size_bytes", len(responseBytes)),
	)
	span.AddEvent("http_response_received")

	return responseBytes
}

// Removed Copilot request and response logic and replaced it with LM Studio local server API implementation.
func getLocalServerResponse(ctx context.Context, prompt string) string {
	// Create span for LLM API call
	ctx, span := tracer.Start(ctx, "llm_api_call",
		trace.WithAttributes(
			attribute.String("llm.api", "lm-studio"),
			attribute.String("llm.endpoint", "http://127.0.0.1:1234/v1/chat/completions"),
		),
	)
	defer span.End()

	requestBody, err := createLocalServerRequestBody(prompt)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(500, "Failed to create request body")
		return fmt.Sprintf("Error creating request body: %v", err)
	}

	span.SetAttributes(
		attribute.Int("request.body_size", len(requestBody)),
	)

	responseBytes, err := sendLocalServerRequest(ctx, requestBody)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(500, "Failed to send request")
		return fmt.Sprintf("Error sending request: %v", err)
	}

	response := parseLocalServerResponse(responseBytes)
	if response == "" {
		span.SetStatus(500, "Empty response from server")
		return "Error: Empty response from local server."
	}

	span.SetAttributes(
		attribute.String("llm.response_preview", response[:min(100, len(response))]),
	)
	span.AddEvent("llm_response_parsed")

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

func sendLocalServerRequest(ctx context.Context, requestBody []byte) ([]byte, error) {
	// Create span for HTTP request to LM Studio
	ctx, span := tracer.Start(ctx, "http_request_lm_studio",
		trace.WithAttributes(
			attribute.String("http.method", "POST"),
			attribute.String("http.url", "http://127.0.0.1:1234/v1/chat/completions"),
		),
	)
	defer span.End()

	url := "http://127.0.0.1:1234/v1/chat/completions"
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(500, "Failed to create HTTP request")
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("api_key", "lm-studio") // Added API key header

	// Use the traced HTTP client
	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	response, err := client.Do(request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(500, "HTTP request failed")
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer response.Body.Close()

	span.SetAttributes(
		attribute.Int("http.status_code", response.StatusCode),
	)

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(500, "Failed to read response body")
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	span.SetAttributes(
		attribute.Int("http.response_size", len(responseBytes)),
	)

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

// Example usage of the LLM Studio agent
func getRandomJokeWithLLMStudio(ctx context.Context) string {
	// Create span for LLM request
	ctx, span := tracer.Start(ctx, "get_llm_joke",
		trace.WithAttributes(
			attribute.String("llm.provider", "lm-studio"),
			attribute.String("llm.model", "meta-llama-3.1-8b-instruct"),
		),
	)
	defer span.End()

	fmt.Println("Fetching a joke using LLM Studio...")
	prompt := "Tell me a funny astronaut joke."
	
	span.SetAttributes(
		attribute.String("llm.prompt", prompt),
	)
	span.AddEvent("llm_request_started")

	copilotResponse := getLocalServerResponse(ctx, prompt)

	if copilotResponse == "" {
		fmt.Println("Failed to fetch a joke. Please try again later.")
		span.SetStatus(500, "Empty response from LLM")
		return ""
	}

	fmt.Println("Here's a joke for you:")
	fmt.Println(copilotResponse)
	
	span.SetAttributes(
		attribute.String("llm.response", copilotResponse),
		attribute.Int("llm.response_length", len(copilotResponse)),
	)
	span.AddEvent("llm_joke_received")

	return copilotResponse
}
