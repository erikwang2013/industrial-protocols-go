package kernel

type Request struct {
	Function string
	Address  string
	Count    int
	Data     []byte
	Metadata map[string]any
}

type Response struct {
	Function string
	Address  string
	Data     []byte
	Metadata map[string]any
}

type Codec interface {
	Encode(req *Request) ([]byte, error)
	Decode(data []byte) (*Response, error)
}
