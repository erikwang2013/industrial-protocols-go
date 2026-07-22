// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package gateway

import (
	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

// Rule 定义了一个协议到另一个协议的转换规则。
type Rule struct {
	Src kernel.Protocol
	Dst kernel.Protocol
	Map func(srcResp *kernel.Response) *kernel.Request
}
