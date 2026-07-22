package connection

import (
	"net"
	"time"
)

type ConnStrategy interface {
	Dial(network, addr string, timeout time.Duration) (net.Conn, error)
}

type LazyStrategy struct{}

func (s *LazyStrategy) Dial(network, addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

type EagerStrategy struct{}

func (s *EagerStrategy) Dial(network, addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

type PooledStrategy struct {
	DefaultPoolSize int
}

func (s *PooledStrategy) Dial(network, addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

type DeviceConfig struct {
	Name    string
	Network string
	Addr    string
	Timeout time.Duration
	Pool    *PoolConfig
}

type PoolConfig struct {
	MaxSize     int
	IdleTimeout time.Duration
}

type HealthStatus struct {
	Healthy   bool
	LatencyMs float64
}
