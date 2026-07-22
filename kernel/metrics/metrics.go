// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package metrics

// Recorder 指标采集接口。内置 Prometheus 和 noop 两种实现。
type Recorder interface {
	Inc(name string, labels map[string]string)
	Observe(name string, value float64, labels map[string]string)
}

type noopRecorder struct{}

func (n noopRecorder) Inc(name string, labels map[string]string)           {}
func (n noopRecorder) Observe(name string, v float64, labels map[string]string) {}

// NoopRecorder 返回一个不执行任何操作的 Recorder。
func NoopRecorder() Recorder { return noopRecorder{} }
