// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package metrics

// Recorder provides a simple interface for recording metrics.
// It supports incrementing counters and observing numeric values
// with optional labels.
type Recorder interface {
	Inc(name string, labels map[string]string)
	Observe(name string, value float64, labels map[string]string)
}

type noopRecorder struct{}

func (n noopRecorder) Inc(name string, labels map[string]string)           {}
func (n noopRecorder) Observe(name string, v float64, labels map[string]string) {}

// NoopRecorder returns a Recorder implementation that discards all metrics.
// It can be used as a default when no metrics backend is configured.
func NoopRecorder() Recorder { return noopRecorder{} }
