// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package kernel

import "testing"

func TestRequestStruct(t *testing.T) {
	r := &Request{
		Function: "read_coils",
		Address:  "40001",
		Count:    8,
		Metadata: map[string]any{"unit_id": 1},
	}
	if r.Function != "read_coils" {
		t.Errorf("expected read_coils, got %s", r.Function)
	}
}
