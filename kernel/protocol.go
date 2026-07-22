// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package kernel

// Protocol 定义了每个协议模块必须实现的元信息接口。
// 每个协议实现为一个独立的 Go module，通过实现此接口注册到系统中。
type Protocol interface {
	// Name 返回协议名称，如 "modbus"、"mqtt"
	Name() string
	// Variants 返回支持的通信变体，如 ["tcp", "rtu"]
	Variants() []string
	// DefaultPort 返回协议的默认端口号
	DefaultPort() int
	// NewCodec 根据 variant 创建对应的编解码器
	NewCodec(variant string) (Codec, error)
}
