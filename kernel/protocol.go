package kernel

type Protocol interface {
	Name()        string
	Variants()    []string
	DefaultPort() int
	NewCodec(variant string) (Codec, error)
}
