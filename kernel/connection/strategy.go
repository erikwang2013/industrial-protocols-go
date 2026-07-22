// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package connection

import (
	"net"
	"time"
)

// ConnStrategy 定义了连接创建策略的接口。
type ConnStrategy interface {
	Dial(network, addr string, timeout time.Duration) (net.Conn, error)
}

// LazyStrategy 延迟连接——仅在首次操作时才 Dial。
type LazyStrategy struct{}

func (s *LazyStrategy) Dial(network, addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

// EagerStrategy 急切连接——Register 时立即 Dial 验证连通性。
type EagerStrategy struct{}

func (s *EagerStrategy) Dial(network, addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

// PooledStrategy 池化连接——维护连接池，每次操作从池中获取，用完归还。
type PooledStrategy struct {
	DefaultPoolSize int
}

func (s *PooledStrategy) Dial(network, addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

// DeviceConfig 设备连接参数配置。
type DeviceConfig struct {
	Name    string
	Network string
	Addr    string
	Timeout time.Duration
	Pool    *PoolConfig
}

// PoolConfig 连接池配置。
type PoolConfig struct {
	MaxSize     int
	IdleTimeout time.Duration
}

// HealthStatus 设备健康状态。
type HealthStatus struct {
	Healthy   bool
	LatencyMs float64
}
