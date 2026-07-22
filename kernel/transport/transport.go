// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package transport

import "io"

// Transport 抽象底层通信信道，不关心上层协议。
// 内置实现包括 TCPTransport、UDPTransport、SerialTransport、PipeTransport。
type Transport interface {
	io.ReadWriter
	io.Closer
	Addr()  string
	Alive() bool
}
