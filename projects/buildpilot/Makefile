HTTP_ADDRESS ?= :8090
REPOSITORY_ROOT ?= ../..
OLLAMA_URL ?= http://localhost:11434
OLLAMA_MODEL ?= qwen2.5-coder:7b

.PHONY: test build start run agent

test:
	go test ./...

build:
	go build -o buildpilot-control-plane ./cmd/control-plane
	go build -o buildpilot-cluster-agent ./cmd/cluster-agent

start:
	HTTP_ADDRESS="$(HTTP_ADDRESS)" \
	REPOSITORY_ROOT="$(REPOSITORY_ROOT)" \
	OLLAMA_URL="$(OLLAMA_URL)" \
	OLLAMA_MODEL="$(OLLAMA_MODEL)" \
	go run ./cmd/control-plane

run: start

agent:
	go run ./cmd/cluster-agent
