package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	devopscmd "github.com/ranjith-ka/Devops/cmd"
)

type jokeInput struct{}

type jokeOutput struct {
	Joke string `json:"joke" jsonschema:"the joke returned to the caller"`
}

func tellJoke(ctx context.Context, _ *mcp.CallToolRequest, _ jokeInput) (*mcp.CallToolResult, jokeOutput, error) {
	joke, err := devopscmd.GetRandomJoke(ctx)
	if err != nil {
		return nil, jokeOutput{}, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: joke},
		},
	}, jokeOutput{Joke: joke}, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "Here is my first http program")
	if err != nil {
		fmt.Printf("%v", err)
	}
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			_, err := fmt.Fprintf(w, "%v: %v\n", name, h)
			if err != nil {
				fmt.Printf("%v", err)
			}
		}
	}
}

func main() {
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "canary-joke-server",
		Version: "1.0.0",
	}, nil)
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "tell_joke",
		Description: "Fetches and returns a random dad joke.",
	}, tellJoke)

	mcpHandler := mcp.NewStreamableHTTPHandler(
		func(_ *http.Request) *mcp.Server { return mcpServer },
		&mcp.StreamableHTTPOptions{
			Stateless:    true,
			JSONResponse: true,
		},
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Welcome to my website!")
		if err != nil {
			fmt.Printf("%v", err)
		}
	})
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/headers", headers)
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/mcp", mcpHandler)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("MCP endpoint: http://localhost:8080/mcp")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
