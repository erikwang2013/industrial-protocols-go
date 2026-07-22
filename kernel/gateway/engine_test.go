package gateway

import (
	"context"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

type mockProtocol struct {
	name     string
	variants []string
	port     int
}

func (p *mockProtocol) Name() string              { return p.name }
func (p *mockProtocol) Variants() []string        { return p.variants }
func (p *mockProtocol) DefaultPort() int          { return p.port }
func (p *mockProtocol) NewCodec(v string) (kernel.Codec, error) {
	return &mockCodec{}, nil
}

type mockCodec struct{}

func (c *mockCodec) Encode(req *kernel.Request) ([]byte, error) {
	return []byte(req.Function + ":" + req.Address), nil
}
func (c *mockCodec) Decode(data []byte) (*kernel.Response, error) {
	return &kernel.Response{Data: data}, nil
}

func TestGatewayAddRule(t *testing.T) {
	engine := New()
	src := &mockProtocol{name: "modbus", variants: []string{"tcp"}, port: 502}
	dst := &mockProtocol{name: "mqtt", variants: []string{"tcp"}, port: 1883}
	engine.Add(Rule{Src: src, Dst: dst, Map: func(r *kernel.Response) *kernel.Request {
		return &kernel.Request{Function: "publish", Address: "topic", Data: r.Data}
	}})
	if len(engine.rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(engine.rules))
	}
}

func TestTransformNoRule(t *testing.T) {
	engine := New()
	_, err := engine.Transform(context.Background(), nil, nil, &kernel.Request{})
	if err == nil {
		t.Error("expected error for no matching rule")
	}
}
