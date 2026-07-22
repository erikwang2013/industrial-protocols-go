package transport

import "io"

// Transport is the minimal network/stream abstraction.
// Every industrial protocol session reads and writes over a Transport.
type Transport interface {
	io.ReadWriter
	io.Closer
	Addr()  string
	Alive() bool
}
