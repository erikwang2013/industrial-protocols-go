// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package transport

import "io"

// PipeTransport wraps io.ReadCloser and io.WriteCloser as a Transport.
// It is used for in-process connections, e.g. bridging to sub-process stdin/stdout.
type PipeTransport struct {
	r    io.ReadCloser
	w    io.WriteCloser
	addr string
}

// NewPipeTransport creates a PipeTransport from a reader, writer, and address label.
func NewPipeTransport(r io.ReadCloser, w io.WriteCloser, addr string) *PipeTransport {
	return &PipeTransport{r: r, w: w, addr: addr}
}

func (t *PipeTransport) Read(p []byte) (int, error)  { return t.r.Read(p) }
func (t *PipeTransport) Write(p []byte) (int, error) { return t.w.Write(p) }
func (t *PipeTransport) Close() error {
	t.r.Close()
	_ = t.w.Close()
	return nil
}
func (t *PipeTransport) Addr() string { return t.addr }
func (t *PipeTransport) Alive() bool  { return true }
