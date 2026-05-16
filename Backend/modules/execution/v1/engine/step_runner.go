package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type StepResult struct {
	Output  map[string]any
	Skipped bool
}

func RunStep(ctx context.Context, step Step, outputs map[string]map[string]any) (StepResult, error) {
	switch step.Type {
	case StepTypeHTTPCall:
		return runHTTP(ctx, step)
	case StepTypeDelay:
		return runDelay(ctx, step)
	case StepTypeBranch:
		return runBranch(step, outputs)
	case StepTypeScript:
		return StepResult{}, errors.New("script execution is disabled until Docker sandbox is available")
	default:
		return StepResult{}, fmt.Errorf("unsupported step type %s", step.Type)
	}
}
func runHTTP(ctx context.Context, step Step) (StepResult, error) {
	method := str(step.Config["method"], "GET")
	url := str(step.Config["url"], "")
	if url == "" {
		return StepResult{}, errors.New("http url is required")
	}
	var body io.Reader
	if b, ok := step.Config["body"]; ok {
		raw, _ := json.Marshal(b)
		body = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return StepResult{}, err
	}
	if headers, ok := step.Config["headers"].(map[string]any); ok {
		for k, v := range headers {
			req.Header.Set(k, fmt.Sprint(v))
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return StepResult{}, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	expected := intVal(step.Config["expected_status"], 0)
	if expected == 0 {
		expected = resp.StatusCode
	}
	if resp.StatusCode != expected {
		return StepResult{Output: map[string]any{"status": resp.StatusCode, "body": string(data)}}, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	return StepResult{Output: map[string]any{"status": resp.StatusCode, "body": string(data)}}, nil
}
func runDelay(ctx context.Context, step Step) (StepResult, error) {
	ms := intVal(step.Config["duration_ms"], 1000)
	select {
	case <-time.After(time.Duration(ms) * time.Millisecond):
		return StepResult{Output: map[string]any{"delayed_ms": ms}}, nil
	case <-ctx.Done():
		return StepResult{}, ctx.Err()
	}
}
func runBranch(step Step, outputs map[string]map[string]any) (StepResult, error) {
	source := str(step.Config["source_step"], "")
	key := str(step.Config["key"], "")
	equals := fmt.Sprint(step.Config["equals"])
	if source == "" || key == "" {
		return StepResult{}, errors.New("branch source_step and key are required")
	}
	actual := fmt.Sprint(outputs[source][key])
	matched := actual == equals
	return StepResult{Output: map[string]any{"matched": matched, "actual": actual, "expected": equals}}, nil
}
func str(v any, def string) string {
	if v == nil {
		return def
	}
	return fmt.Sprint(v)
}
func intVal(v any, def int) int {
	switch t := v.(type) {
	case int:
		return t
	case float64:
		return int(t)
	case string:
		n, err := strconv.Atoi(t)
		if err == nil {
			return n
		}
	}
	return def
}
