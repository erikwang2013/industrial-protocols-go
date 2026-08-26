// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package gateway

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// mockTransport is an in-memory transport.Transport for exercising Transform
// without any real I/O.
type mockTransport struct {
	readData []byte
	readErr  error
	writeErr error
	written  []byte
}

func (m *mockTransport) Read(p []byte) (int, error) {
	n := copy(p, m.readData)
	return n, m.readErr
}

func (m *mockTransport) Write(p []byte) (int, error) {
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	m.written = append(m.written, p...)
	return len(p), nil
}

func (m *mockTransport) Close() error        { return nil }
func (m *mockTransport) Addr() string        { return "mock" }
func (m *mockTransport) Alive() bool         { return true }

// mockProto is a kernel.Protocol whose NewCodec can be forced to fail,
// and which can return a custom codec for error injection.
type mockProto struct {
	name     string
	codecErr error
	codec    kernel.Codec
}

func (p *mockProto) Name() string       { return p.name }
func (p *mockProto) Variants() []string { return []string{"tcp"} }
func (p *mockProto) DefaultPort() int   { return 0 }

func (p *mockProto) NewCodec(variant string) (kernel.Codec, error) {
	if p.codecErr != nil {
		return nil, p.codecErr
	}
	if p.codec != nil {
		return p.codec, nil
	}
	return &mockCodec{}, nil
}

// failableCodec fails on Encode or Decode for error-path tests.
type failableCodec struct {
	encodeErr error
	decodeErr error
}

func (c *failableCodec) Encode(req *kernel.Request) ([]byte, error) {
	if c.encodeErr != nil {
		return nil, c.encodeErr
	}
	return []byte(req.Function + ":" + req.Address), nil
}

func (c *failableCodec) Decode(data []byte) (*kernel.Response, error) {
	if c.decodeErr != nil {
		return nil, c.decodeErr
	}
	return &kernel.Response{Data: data}, nil
}

func mapToPublish(r *kernel.Response) *kernel.Request {
	return &kernel.Request{Function: "publish", Address: "topic", Data: r.Data}
}

func TestTransform_Success(t *testing.T) {
	engine := New()
	src := &mockTransport{readData: []byte("payload-1")}
	dst := &mockTransport{readData: []byte("ack")}
	mapCalls := 0
	engine.Add(Rule{
		Src: &mockProto{name: "modbus"},
		Dst: &mockProto{name: "mqtt"},
		Map: func(r *kernel.Response) *kernel.Request {
			mapCalls++
			return mapToPublish(r)
		},
	})

	resp, err := engine.Transform(context.Background(), src, dst,
		&kernel.Request{Function: "read", Address: "40001"})
	if err != nil {
		t.Fatal(err)
	}
	if string(src.written) != "read:40001" {
		t.Errorf("src write: expected 'read:40001', got %q", src.written)
	}
	if string(dst.written) != "publish:topic" {
		t.Errorf("dst write: expected 'publish:topic', got %q", dst.written)
	}
	if string(resp.Data) != "ack" {
		t.Errorf("expected response data 'ack', got %q", resp.Data)
	}
	if mapCalls != 1 {
		t.Errorf("expected 1 map call, got %d", mapCalls)
	}
}

func TestTransform_ErrorPaths(t *testing.T) {
	encodeErr := errors.New("encode boom")
	decodeErr := errors.New("decode boom")
	writeErr := errors.New("write boom")
	readErr := errors.New("read boom")
	codecErr := errors.New("codec boom")

	tests := []struct {
		name      string
		src       *mockProto
		dst       *mockProto
		srcTr     *mockTransport
		dstTr     *mockTransport
		wantError string
	}{
		{
			name: "src encode", src: &mockProto{name: "a", codec: &failableCodec{encodeErr: encodeErr}},
			dst: &mockProto{name: "b"}, srcTr: &mockTransport{}, dstTr: &mockTransport{},
			wantError: "src encode",
		},
		{
			name: "src write", src: &mockProto{name: "a"},
			dst: &mockProto{name: "b"}, srcTr: &mockTransport{writeErr: writeErr}, dstTr: &mockTransport{},
			wantError: "src write",
		},
		{
			name: "src read", src: &mockProto{name: "a"},
			dst: &mockProto{name: "b"}, srcTr: &mockTransport{readErr: readErr}, dstTr: &mockTransport{},
			wantError: "src read",
		},
		{
			name: "src decode", src: &mockProto{name: "a", codec: &failableCodec{decodeErr: decodeErr}},
			dst: &mockProto{name: "b"}, srcTr: &mockTransport{readData: []byte("x")}, dstTr: &mockTransport{},
			wantError: "src decode",
		},
		{
			name: "dst codec", src: &mockProto{name: "a"},
			dst: &mockProto{name: "b", codecErr: codecErr}, srcTr: &mockTransport{readData: []byte("x")}, dstTr: &mockTransport{},
			wantError: "dst codec",
		},
		{
			name: "dst encode", src: &mockProto{name: "a"},
			dst: &mockProto{name: "b", codec: &failableCodec{encodeErr: encodeErr}},
			srcTr: &mockTransport{readData: []byte("x")}, dstTr: &mockTransport{},
			wantError: "dst encode",
		},
		{
			name: "dst write", src: &mockProto{name: "a"},
			dst: &mockProto{name: "b"}, srcTr: &mockTransport{readData: []byte("x")},
			dstTr: &mockTransport{writeErr: writeErr},
			wantError: "dst write",
		},
		{
			name: "dst read", src: &mockProto{name: "a"},
			dst: &mockProto{name: "b"}, srcTr: &mockTransport{readData: []byte("x")},
			dstTr: &mockTransport{readErr: readErr},
			wantError: "dst read",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := New()
			engine.Add(Rule{Src: tt.src, Dst: tt.dst, Map: mapToPublish})
			_, err := engine.Transform(context.Background(), tt.srcTr, tt.dstTr, &kernel.Request{})
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.wantError) {
				t.Errorf("expected error containing %q, got %v", tt.wantError, err)
			}
		})
	}
}

func TestTransform_SkipsRuleWithBadSrcCodec(t *testing.T) {
	engine := New()
	engine.Add(Rule{
		Src: &mockProto{name: "bad", codecErr: errors.New("no codec")},
		Dst: &mockProto{name: "b"}, Map: mapToPublish,
	})
	engine.Add(Rule{
		Src: &mockProto{name: "good"}, Dst: &mockProto{name: "b"}, Map: mapToPublish,
	})

	dst := &mockTransport{readData: []byte("done")}
	resp, err := engine.Transform(context.Background(), &mockTransport{readData: []byte("x")}, dst, &kernel.Request{})
	if err != nil {
		t.Fatal(err)
	}
	if string(resp.Data) != "done" {
		t.Errorf("expected 'done', got %q", resp.Data)
	}
}

func TestTransform_AllSrcCodecsFail(t *testing.T) {
	engine := New()
	engine.Add(Rule{
		Src: &mockProto{name: "a", codecErr: errors.New("x")},
		Dst: &mockProto{name: "b"}, Map: mapToPublish,
	})
	_, err := engine.Transform(context.Background(), &mockTransport{}, &mockTransport{}, &kernel.Request{})
	if err == nil || !strings.Contains(err.Error(), "no matching rule") {
		t.Errorf("expected 'no matching rule' error, got %v", err)
	}
}

// TestTransform_NilMapDoesNotPanic is a prove-it test: currently Rule without
// a Map function panics (rule.Map(srcResp) on nil func). It should return an
// error instead.
func TestTransform_NilMapDoesNotPanic(t *testing.T) {
	engine := New()
	engine.Add(Rule{Src: &mockProto{name: "a"}, Dst: &mockProto{name: "b"}})

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Transform panicked on rule with nil Map (bug): %v", r)
		}
	}()
	_, err := engine.Transform(context.Background(),
		&mockTransport{readData: []byte("x")}, &mockTransport{readData: []byte("y")}, &kernel.Request{})
	if err == nil {
		t.Error("expected error for rule without Map")
	}
}

// TestTransform_NilDstDoesNotPanic is a prove-it test: currently a rule with a
// nil Dst panics at rule.Dst.Name(). It should return an error instead.
func TestTransform_NilDstDoesNotPanic(t *testing.T) {
	engine := New()
	engine.Add(Rule{Src: &mockProto{name: "a"}, Dst: nil, Map: mapToPublish})

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Transform panicked on rule with nil Dst (bug): %v", r)
		}
	}()
	_, err := engine.Transform(context.Background(),
		&mockTransport{readData: []byte("x")}, &mockTransport{}, &kernel.Request{})
	if err == nil {
		t.Error("expected error for rule without Dst")
	}
}
