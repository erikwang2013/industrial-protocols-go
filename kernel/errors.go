package kernel

import (
	"errors"
	"fmt"
)

var (
	ErrTimeout         = errors.New("protocol: timeout")
	ErrCircuitOpen     = errors.New("protocol: circuit breaker open")
	ErrRetryExhausted  = errors.New("protocol: all retries exhausted")
	ErrTransportClosed = errors.New("protocol: transport closed")
	ErrInvalidAddress  = errors.New("protocol: invalid address")
)

type ProtocolError struct {
	Code    string
	Message string
	Raw     []byte
}

func (e *ProtocolError) Error() string {
	return fmt.Sprintf("protocol error [%s]: %s", e.Code, e.Message)
}
