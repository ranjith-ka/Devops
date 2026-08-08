package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/ranjith-ka/buildpilot/internal/domain"
)

const version = "0.1.0-dev"

type agent struct {
	id           string
	clusterName  string
	controlPlane string
	http         *http.Client
	logger       *slog.Logger
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	agentID := requiredEnv(logger, "AGENT_ID")
	clusterName := requiredEnv(logger, "CLUSTER_NAME")
	controlPlane := strings.TrimRight(env("CONTROL_PLANE_URL", "http://localhost:8090"), "/")
	runner := &agent{
		id: agentID, clusterName: clusterName, controlPlane: controlPlane,
		http: &http.Client{Timeout: 20 * time.Second}, logger: logger,
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := runner.run(ctx); err != nil && ctx.Err() == nil {
		logger.Error("agent stopped", "error", err)
		os.Exit(1)
	}
}

func (a *agent) run(ctx context.Context) error {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for {
		if err := a.heartbeat(ctx); err != nil {
			a.logger.Warn("heartbeat failed", "error", err)
		} else if job, found, err := a.nextJob(ctx); err != nil {
			a.logger.Warn("job poll failed", "error", err)
		} else if found {
			a.logger.Info("job received; executor not implemented in phase 0", "job_id", job.ID, "image", job.Image)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (a *agent) heartbeat(ctx context.Context) error {
	payload := domain.Agent{ID: a.id, ClusterName: a.clusterName, Architecture: runtime.GOARCH, Version: version}
	return a.jsonRequest(ctx, http.MethodPost, "/api/v1/agents/"+a.id+"/heartbeat", payload, &domain.Agent{})
}

func (a *agent) nextJob(ctx context.Context) (domain.Job, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.controlPlane+"/api/v1/agents/"+a.id+"/jobs/next", nil)
	if err != nil {
		return domain.Job{}, false, err
	}
	resp, err := a.http.Do(req)
	if err != nil {
		return domain.Job{}, false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return domain.Job{}, false, nil
	}
	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		return domain.Job{}, false, fmt.Errorf("poll jobs: %s: %s", resp.Status, strings.TrimSpace(string(data)))
	}
	var job domain.Job
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		return domain.Job{}, false, err
	}
	return job, true, nil
}

func (a *agent) jsonRequest(ctx context.Context, method, path string, input, output any) error {
	body, err := json.Marshal(input)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, method, a.controlPlane+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		return fmt.Errorf("control plane returned %s: %s", resp.Status, strings.TrimSpace(string(data)))
	}
	return json.NewDecoder(resp.Body).Decode(output)
}

func env(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func requiredEnv(logger *slog.Logger, name string) string {
	value := os.Getenv(name)
	if value == "" {
		logger.Error("required environment variable is missing", "name", name)
		os.Exit(2)
	}
	return value
}
