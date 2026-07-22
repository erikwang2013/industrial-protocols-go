// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package kernel

import (
	"errors"
	"fmt"
)

// 预定义的哨兵错误，供中间件和协议实现统一使用。
var (
	// ErrTimeout 表示操作超时
	ErrTimeout = errors.New("protocol: timeout")
	// ErrCircuitOpen 表示熔断器已打开，请求被拒绝
	ErrCircuitOpen = errors.New("protocol: circuit breaker open")
	// ErrRetryExhausted 表示所有重试次数已用尽
	ErrRetryExhausted = errors.New("protocol: all retries exhausted")
	// ErrTransportClosed 表示传输层已关闭
	ErrTransportClosed = errors.New("protocol: transport closed")
	// ErrInvalidAddress 表示协议地址格式无效
	ErrInvalidAddress = errors.New("protocol: invalid address")
)

// ProtocolError 携带协议级别的错误信息，包括原生错误码和原始响应帧。
type ProtocolError struct {
	Code    string
	Message string
	Raw     []byte
}

func (e *ProtocolError) Error() string {
	return fmt.Sprintf("protocol error [%s]: %s", e.Code, e.Message)
}
