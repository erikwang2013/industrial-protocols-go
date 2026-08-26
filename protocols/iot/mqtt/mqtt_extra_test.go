// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package mqtt

import (
	"bytes"
	"testing"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

func TestExtraRemainingLengthEncoding(t *testing.T) {
	tests := []struct {
		n    int
		want []byte
	}{
		{0, []byte{0x00}},
		{127, []byte{0x7F}},
		{128, []byte{0x80, 0x01}},
		{16383, []byte{0xFF, 0x7F}},
		{16384, []byte{0x80, 0x80, 0x01}},
		{2097151, []byte{0xFF, 0xFF, 0x7F}},
	}
	for _, tt := range tests {
		var buf bytes.Buffer
		encodeRemainingLength(&buf, tt.n)
		if !bytes.Equal(buf.Bytes(), tt.want) {
			t.Errorf("encodeRemainingLength(%d) = % X, want % X", tt.n, buf.Bytes(), tt.want)
		}
		got, n := decodeRemainingLength(buf.Bytes())
		if got != tt.n || n != len(tt.want) {
			t.Errorf("decodeRemainingLength(% X) = (%d, %d), want (%d, %d)", buf.Bytes(), got, n, tt.n, len(tt.want))
		}
	}
}

func TestExtraDecodePublishTruncatedTopic(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	// Remaining length 7 declares topic length 5 but only 2 bytes follow.
	frame := []byte{0x30, 0x07, 0x00, 0x05, 't', 'o'}
	if _, err := c.Decode(frame); err == nil {
		t.Fatal("expected error for truncated topic")
	}
}

func TestExtraDecodePublishTooShortForTopicLength(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	frame := []byte{0x30, 0x03, 0x00}
	if _, err := c.Decode(frame); err == nil {
		t.Fatal("expected error for frame shorter than topic length field")
	}
}

func TestExtraDecodePublishMultiByteLength(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	// Remaining length 128 (multi-byte), topic "ab", payload [01].
	frame := []byte{0x30, 0x80, 0x01, 0x00, 0x02, 'a', 'b', 0x01}
	resp, err := c.Decode(frame)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Address != "ab" {
		t.Errorf("topic = %q, want ab", resp.Address)
	}
	if len(resp.Data) != 1 || resp.Data[0] != 0x01 {
		t.Errorf("data = %v, want [01]", resp.Data)
	}
}

func TestExtraEncodeConnectExact(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	raw, err := c.Encode(&kernel.Request{
		Function: "connect",
		Metadata: map[string]any{"client_id": "test1", "keep_alive": float64(30)},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x10, 0x11, 0x00, 0x04, 'M', 'Q', 'T', 'T', 0x04, 0x02, 0x00, 0x1E, 0x00, 0x05, 't', 'e', 's', 't', '1'}
	if !bytes.Equal(raw, want) {
		t.Errorf("connect = % X, want % X", raw, want)
	}
}

func TestExtraEncodePublishExact(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	raw, err := c.Encode(&kernel.Request{Function: "publish", Address: "a/b", Data: []byte{0x01, 0x02}})
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x30, 0x07, 0x00, 0x03, 'a', '/', 'b', 0x01, 0x02}
	if !bytes.Equal(raw, want) {
		t.Errorf("publish = % X, want % X", raw, want)
	}
}

func TestExtraUnknownFunction(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	if _, err := c.Encode(&kernel.Request{Function: "bogus"}); err == nil {
		t.Fatal("expected error for unknown function")
	}
}

func TestExtraDecodeTooShort(t *testing.T) {
	c, _ := New().NewCodec("tcp")
	if _, err := c.Decode([]byte{0x30}); err == nil {
		t.Fatal("expected error for 1-byte frame")
	}
}
