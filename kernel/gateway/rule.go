package gateway

import (
	"github.com/erikwang2013/industrial-protocols-go/kernel"
)

type Rule struct {
	Src kernel.Protocol
	Dst kernel.Protocol
	Map func(srcResp *kernel.Response) *kernel.Request
}
