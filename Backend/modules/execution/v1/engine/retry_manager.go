package engine

import "time"

func RetryDelay(cfg *RetryConfig, attempt int) time.Duration {
	if cfg == nil {
		return 0
	}
	initial := cfg.InitialDelayMs
	if initial <= 0 {
		initial = 500
	}
	max := cfg.MaxDelayMs
	if max <= 0 {
		max = 5000
	}
	factor := cfg.BackoffFactor
	if factor <= 1 {
		factor = 2
	}
	delay := float64(initial)
	for i := 1; i < attempt; i++ {
		delay *= factor
	}
	if delay > float64(max) {
		delay = float64(max)
	}
	return time.Duration(delay) * time.Millisecond
}
func MaxAttempts(cfg *RetryConfig) int {
	if cfg == nil || cfg.MaxAttempts < 1 {
		return 1
	}
	return cfg.MaxAttempts
}
