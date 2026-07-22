// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package metrics

import "testing"

func TestNoopRecorder(t *testing.T) {
	r := NoopRecorder()
	r.Inc("test", map[string]string{"k": "v"})
	r.Observe("test", 1.0, nil)
}

func TestRecorderInterface(t *testing.T) {
	var r Recorder = NoopRecorder()
	if r == nil {
		t.Error("NoopRecorder returned nil")
	}
}
