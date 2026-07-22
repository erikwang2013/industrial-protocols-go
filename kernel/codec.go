// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package kernel

// Request 是协议无关的通用请求结构体。
// Function 描述操作类型（如 "read_coils"、"publish"），
// Address 是协议自定义的寻址方式（如 Modbus 的 "40001"），
// Metadata 携带协议特有字段（如 Modbus 的 unit_id）。
type Request struct {
	Function string
	Address  string
	Count    int
	Data     []byte
	Metadata map[string]any
}

// Response 是协议无关的通用响应结构体。
type Response struct {
	Function string
	Address  string
	Data     []byte
	Metadata map[string]any
}

// Codec 负责协议的编解码——把应用层数据结构编码为二进制帧，或从二进制帧解码。
type Codec interface {
	// Encode 将请求编码为 wire format 字节流
	Encode(req *Request) ([]byte, error)
	// Decode 将 wire format 字节流解码为响应
	Decode(data []byte) (*Response, error)
}
