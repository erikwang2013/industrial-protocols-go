package cclinkie

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type CCLinkIEProtocol struct{}

func New() *CCLinkIEProtocol { return &CCLinkIEProtocol{} }

func (p *CCLinkIEProtocol) Name() string        { return "cclinkie" }
func (p *CCLinkIEProtocol) Variants() []string   { return []string{"ethernet"} }
func (p *CCLinkIEProtocol) DefaultPort() int     { return 0 }
func (p *CCLinkIEProtocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
