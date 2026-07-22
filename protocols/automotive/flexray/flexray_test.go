package flexray

import "testing"

func TestProtocol_ImplementsInterface(t *testing.T) {
	p := New()
	if p.Name() != "flexray" {
		t.Errorf("expected flexray, got %s", p.Name())
	}
	if len(p.Variants()) == 0 {
		t.Error("should have at least one variant")
	}
	if p.DefaultPort() < 0 {
		t.Error("default port must not be negative")
	}
	_, err := p.NewCodec("invalid")
	if err == nil {
		t.Error("stub NewCodec should return error")
	}
}
