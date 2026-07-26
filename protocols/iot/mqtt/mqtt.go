// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package mqtt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

type MqttProtocol struct{}

func New() *MqttProtocol                             { return &MqttProtocol{} }
func (p *MqttProtocol) Name() string                 { return "mqtt" }
func (p *MqttProtocol) Variants() []string           { return []string{"tcp", "ws"} }
func (p *MqttProtocol) DefaultPort() int             { return 1883 }

func (p *MqttProtocol) NewCodec(v string) (kernel.Codec, error) {
	if v != "tcp" && v != "ws" {
		return nil, fmt.Errorf("mqtt: unsupported variant %q", v)
	}
	return &mqttCodec{}, nil
}

type mqttCodec struct{}

func (c *mqttCodec) Encode(req *kernel.Request) ([]byte, error) {
	switch req.Function {
	case "connect":
		return c.encodeConnect(req)
	case "publish":
		return c.encodePublish(req)
	case "subscribe":
		return c.encodeSubscribe(req)
	case "pingreq":
		return []byte{0xC0, 0x00}, nil
	case "disconnect":
		return []byte{0xE0, 0x00}, nil
	default:
		return nil, fmt.Errorf("mqtt: unknown function %q", req.Function)
	}
}

func (c *mqttCodec) Decode(data []byte) (*kernel.Response, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("mqtt: frame too short")
	}
	msgType := data[0] >> 4
	switch msgType {
	case 2:
		return &kernel.Response{Function: "connack", Data: data}, nil
	case 3:
		return c.decodePublish(data)
	case 9:
		return &kernel.Response{Function: "suback", Data: data}, nil
	case 13:
		return &kernel.Response{Function: "pingresp"}, nil
	default:
		return &kernel.Response{Data: data}, nil
	}
}

func (c *mqttCodec) encodeConnect(req *kernel.Request) ([]byte, error) {
	clientID := "goclient"
	if v, ok := req.Metadata["client_id"].(string); ok {
		clientID = v
	}
	keepAlive := uint16(60)
	if v, ok := req.Metadata["keep_alive"].(float64); ok {
		keepAlive = uint16(v)
	}

	buf := &bytes.Buffer{}
	buf.WriteByte(0x10)
	payload := &bytes.Buffer{}
	if err := writeString(payload, "MQTT"); err != nil {
		return nil, err
	}
	payload.WriteByte(4)
	payload.WriteByte(2)
	binary.Write(payload, binary.BigEndian, keepAlive)
	if err := writeString(payload, clientID); err != nil {
		return nil, err
	}
	encodeRemainingLength(buf, payload.Len())
	buf.Write(payload.Bytes())
	return buf.Bytes(), nil
}

func (c *mqttCodec) encodePublish(req *kernel.Request) ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteByte(0x30)
	payload := &bytes.Buffer{}
	if err := writeString(payload, req.Address); err != nil {
		return nil, err
	}
	payload.Write(req.Data)
	encodeRemainingLength(buf, payload.Len())
	buf.Write(payload.Bytes())
	return buf.Bytes(), nil
}

func (c *mqttCodec) encodeSubscribe(req *kernel.Request) ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteByte(0x82)
	payload := &bytes.Buffer{}
	binary.Write(payload, binary.BigEndian, uint16(1))
	if err := writeString(payload, req.Address); err != nil {
		return nil, err
	}
	payload.WriteByte(0)
	encodeRemainingLength(buf, payload.Len())
	buf.Write(payload.Bytes())
	return buf.Bytes(), nil
}

func (c *mqttCodec) decodePublish(data []byte) (*kernel.Response, error) {
	offset := 1
	_, n := decodeRemainingLength(data[offset:])
	offset += n
	if offset+2 > len(data) {
		return nil, fmt.Errorf("mqtt: publish frame too short for topic length")
	}
	topicLen := int(data[offset])<<8 | int(data[offset+1])
	offset += 2
	if offset+topicLen > len(data) {
		return nil, fmt.Errorf("mqtt: publish frame truncated at topic")
	}
	topic := string(data[offset : offset+topicLen])
	offset += topicLen
	return &kernel.Response{Function: "publish", Address: topic, Data: data[offset:]}, nil
}

func encodeRemainingLength(w io.Writer, n int) {
	for {
		b := byte(n % 128)
		n /= 128
		if n > 0 {
			b |= 0x80
		}
		w.(*bytes.Buffer).WriteByte(b)
		if n == 0 {
			break
		}
	}
}

func decodeRemainingLength(data []byte) (int, int) {
	var length, multiplier int
	for i, b := range data {
		length += int(b&0x7F) << multiplier
		if b&0x80 == 0 {
			return length, i + 1
		}
		multiplier += 7
	}
	return 0, 0
}

func writeString(w io.Writer, s string) error {
	if err := binary.Write(w, binary.BigEndian, uint16(len(s))); err != nil {
		return err
	}
	_, err := w.Write([]byte(s))
	return err
}
