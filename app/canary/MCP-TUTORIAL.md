# Tutorial: Build a Joke MCP Server in Go

This tutorial explains the MCP server implemented in `hello-world.go`.

You need Go 1.25 or newer because the official MCP SDK version used by this
example requires it. Check with `go version`.

## What you are building

The existing Go web application still serves its normal HTTP endpoints:

```text
GET /          Welcome page
GET /hello     Hello message
GET /headers   Request headers
GET /metrics   Prometheus metrics
```

It now also serves MCP at:

```text
POST /mcp      MCP Streamable HTTP endpoint
```

An MCP-compatible AI client can discover and call the `tell_joke` tool. The
tool delegates the actual work to `cmd.GetRandomJoke`.

## MCP in simple terms

MCP stands for Model Context Protocol. It provides a standard connection
between an AI application and capabilities supplied by another application.

```text
AI application -> MCP request -> tellJoke -> cmd.GetRandomJoke -> joke API
AI application <- MCP result  <- tellJoke <- joke text        <- joke API
```

The AI model does not directly execute Go code. It asks the MCP server to call
a named tool with structured arguments.

## Step 1: Define the tool input and output

The joke tool takes no arguments, so its input is an empty struct:

```go
type jokeInput struct{}
```

Its structured output contains a joke:

```go
type jokeOutput struct {
    Joke string `json:"joke" jsonschema:"the joke returned to the caller"`
}
```

The official Go SDK uses these Go types to generate JSON schemas. MCP clients
use the schemas to understand valid tool input and output.

## Step 2: Implement the tool

The handler is an ordinary Go function with the signature expected by the SDK:

```go
func tellJoke(
    _ context.Context,
    _ *mcp.CallToolRequest,
    _ jokeInput,
) (*mcp.CallToolResult, jokeOutput, error) {
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
```

The MCP handler is deliberately thin: `cmd.GetRandomJoke` owns the HTTP request
to the joke provider, while `tellJoke` converts its result to MCP content. It
returns both text content for conversational clients and typed structured
output for clients that want to process the result as data.

## Step 3: Create the MCP server

`mcp.NewServer` supplies the server identity:

```go
mcpServer := mcp.NewServer(&mcp.Implementation{
    Name:    "canary-joke-server",
    Version: "1.0.0",
}, nil)
```

The identity is shown to MCP clients when they connect.

## Step 4: Register the tool

```go
mcp.AddTool(mcpServer, &mcp.Tool{
    Name:        "tell_joke",
    Description: "Fetches and returns a random dad joke.",
}, tellJoke)
```

The name is the identifier clients use to call the tool. The description helps
an AI model decide when the tool is relevant.

## Step 5: Add the HTTP transport

```go
mcpHandler := mcp.NewStreamableHTTPHandler(
    func(_ *http.Request) *mcp.Server { return mcpServer },
    &mcp.StreamableHTTPOptions{
        Stateless:    true,
        JSONResponse: true,
    },
)
```

- Streamable HTTP lets remote clients communicate with MCP over HTTP.
- Stateless mode avoids storing a client session between requests.
- JSON responses make this small example easy to inspect with `curl`.

The handler is mounted alongside the existing routes:

```go
mux.Handle("/mcp", mcpHandler)
```

## Run the server

From the repository root:

```bash
go run ./app/canary/hello-world.go
```

The server listens at `http://localhost:8080`, and the MCP endpoint is
`http://localhost:8080/mcp`.

Check the original web application:

```bash
curl http://localhost:8080/hello
```

## Discover the tool

Send the MCP `tools/list` JSON-RPC request:

```bash
curl --request POST http://localhost:8080/mcp \
  --header 'Content-Type: application/json' \
  --header 'Accept: application/json, text/event-stream' \
  --data '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'
```

Look for a tool named `tell_joke` in the response. Tool discovery allows a
client to learn the server's capabilities without hard-coding them.

## Call the joke tool

```bash
curl --request POST http://localhost:8080/mcp \
  --header 'Content-Type: application/json' \
  --header 'Accept: application/json, text/event-stream' \
  --data '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/call",
    "params": {
      "name": "tell_joke",
      "arguments": {}
    }
  }'
```

The response contains:

```json
{
  "content": [
    {
      "type": "text",
      "text": "A random joke returned by icanhazdadjoke.com"
    }
  ]
}
```

## Learning exercise

Change the tool so it accepts a joke category:

```go
type jokeInput struct {
    Category string `json:"category" jsonschema:"joke category, for example programming or devops"`
}
```

Then use `input.Category` inside `tellJoke` to select a joke. After restarting
the server, call `tools/list` again and inspect how the generated input schema
changed.

## Production considerations

This tutorial tool is intentionally harmless and public. Before adding tools
that read private information or change infrastructure:

- Add authentication and authorization.
- Validate every argument.
- Give tools narrow permissions.
- Avoid returning secrets in errors or logs.
- Add timeouts and cancellation for external calls.
- Record tool calls with structured logging or tracing.

The official SDK documentation is available at
https://github.com/modelcontextprotocol/go-sdk.
