package hart

import (
	"fmt"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

type HartProtocol struct{}

func New() *HartProtocol                { return &HartProtocol{} }
func (p *HartProtocol) Name() string    { return "hart" }
func (p *HartProtocol) Variants() []string { return []string{"fsk"} }
func (p *HartProtocol) DefaultPort() int   { return 0 }

func (p *HartProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "fsk" {
		return nil, fmt.Errorf("hart: unsupported variant %q", v)
	}
	return &hartCodec{}, nil
}

type hartCodec struct{}

func (c *hartCodec) Encode(req *kernel.Request) ([]byte, error) {
	cmd := byte(0)
	switch req.Function {
	case "read_unique_id":
		cmd = 0
	case "read_dynamic_variables":
		cmd = 3
	default:
		cmd = 0
	}

	addr := byte(0x80) // polling address 0 (short frame, master to slave)
	if v, ok := req.Metadata["polling_addr"].(byte); ok {
		addr = v | 0x80
	}

	frame := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x02, addr, cmd, 0}
	frame = append(frame, xorChecksum(frame[5:]))

	return frame, nil
}

func (c *hartCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 9 {
		return nil, fmt.Errorf("hart: frame too short (%d bytes)", len(data))
	}

	// Skip preamble, find delimiter
	delimIdx := 0
	for i, b := range data {
		if b != 0xFF {
			delimIdx = i
			break
		}
	}

	payload := data[delimIdx:]
	if len(payload) < 4 {
		return nil, fmt.Errorf("hart: payload too short")
	}

	delim := payload[0]
	isLong := delim&0x80 != 0

	addrSize := 1
	if isLong {
		addrSize = 5
	}

	cmdOff := 1 + addrSize
	bcOff := cmdOff + 1

	if len(payload) < bcOff+1 {
		return nil, fmt.Errorf("hart: payload too short for frame type")
	}

	cmd := payload[cmdOff]
	byteCount := payload[bcOff]
	dataOff := bcOff + 1
	csOff := dataOff + int(byteCount)

	if csOff >= len(payload) {
		return nil, fmt.Errorf("hart: byte count mismatch")
	}

	// Verify checksum
	expectedCS := xorChecksum(payload[:csOff])
	actualCS := payload[csOff]
	if expectedCS != actualCS {
		return nil, fmt.Errorf("hart: checksum mismatch (expected %02X, got %02X)", expectedCS, actualCS)
	}

	resp := &kernel.Response{
		Data: payload[dataOff:csOff],
		Metadata: map[string]any{
			"command":    int(cmd),
			"long_frame": isLong,
			"delimiter":  delim,
		},
	}

	if cmd == 0 && byteCount >= 2 {
		resp.Metadata["manufacturer"] = payload[dataOff]
		resp.Metadata["device_type"] = payload[dataOff+1]
	}
	if cmd == 3 && byteCount >= 1 {
		resp.Metadata["pv_units"] = payload[dataOff]
	}

	return resp, nil
}

func xorChecksum(data []byte) byte {
	var cs byte
	for _, b := range data {
		cs ^= b
	}
	return cs
}
