package engine

import (
	"context"
	"time"
)

type StepEventFunc func(step Step, status string, attempt int, output map[string]any, err error)

type Executor struct{}

func NewExecutor() *Executor { return &Executor{} }
func (e *Executor) Execute(ctx context.Context, dag *DAGDefinition, onEvent StepEventFunc) error {
	if dag.TimeoutMs > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(dag.TimeoutMs)*time.Millisecond)
		defer cancel()
	}
	waves, err := TopologicalWaves(dag)
	if err != nil {
		return err
	}
	outputs := map[string]map[string]any{}
	for _, wave := range waves {
		for _, step := range wave {
			attempts := MaxAttempts(step.RetryConfig)
			var lastErr error
			var out map[string]any
			for attempt := 1; attempt <= attempts; attempt++ {
				if onEvent != nil {
					onEvent(step, "running", attempt, nil, nil)
				}
				res, err := RunStep(ctx, step, outputs)
				out = res.Output
				lastErr = err
				if err == nil {
					outputs[step.ID] = out
					if onEvent != nil {
						onEvent(step, "completed", attempt, out, nil)
					}
					break
				}
				if attempt < attempts {
					time.Sleep(RetryDelay(step.RetryConfig, attempt))
				}
			}
			if lastErr != nil {
				if onEvent != nil {
					onEvent(step, "failed", attempts, out, lastErr)
				}
				return lastErr
			}
		}
	}
	return nil
}
