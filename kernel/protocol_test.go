// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package kernel

import "testing"

type testProtocol struct{}

func (p *testProtocol) Name() string                 { return "test" }
func (p *testProtocol) Variants() []string           { return []string{"v1", "v2"} }
func (p *testProtocol) DefaultPort() int             { return 9000 }
func (p *testProtocol) NewCodec(v string) (Codec, error) {
	return nil, nil
}

func TestProtocolInterface(t *testing.T) {
	var p Protocol = &testProtocol{}
	if p.Name() != "test" {
		t.Errorf("expected test, got %s", p.Name())
	}
	if len(p.Variants()) != 2 {
		t.Errorf("expected 2 variants, got %d", len(p.Variants()))
	}
	if p.DefaultPort() != 9000 {
		t.Errorf("expected 9000, got %d", p.DefaultPort())
	}
}
