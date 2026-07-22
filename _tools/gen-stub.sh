#!/bin/bash
set -e

CATEGORY="$1"
DISPLAY="$2"
VARIANTS="$3"
PORT="$4"
NAME=$(basename "$CATEGORY")

DIR="protocols/${CATEGORY}"
MOD_PATH="github.com/erikwang2013/industrial-protocols-go/protocols/${CATEGORY}"
KERNEL_VERSION="v0.0.0"

mkdir -p "$DIR"
cd "$DIR"
go mod init "$MOD_PATH"
go mod edit -require "github.com/erikwang2013/industrial-protocols-go/kernel@${KERNEL_VERSION}"
go mod edit -replace "github.com/erikwang2013/industrial-protocols-go/kernel=../../../kernel"

cat > "${NAME}.go" << GOFILE
package ${NAME}

import "github.com/erikwang2013/industrial-protocols-go/kernel"

type ${DISPLAY}Protocol struct{}

func New() *${DISPLAY}Protocol { return &${DISPLAY}Protocol{} }

func (p *${DISPLAY}Protocol) Name() string        { return "${NAME}" }
func (p *${DISPLAY}Protocol) Variants() []string   { return []string{${VARIANTS}} }
func (p *${DISPLAY}Protocol) DefaultPort() int     { return ${PORT} }
func (p *${DISPLAY}Protocol) NewCodec(variant string) (kernel.Codec, error) {
	return nil, kernel.ErrInvalidAddress
}
GOFILE

cat > "${NAME}_test.go" << TESTFILE
package ${NAME}

import "testing"

func TestProtocol_ImplementsInterface(t *testing.T) {
	p := New()
	if p.Name() != "${NAME}" {
		t.Errorf("expected ${NAME}, got %s", p.Name())
	}
	if len(p.Variants()) == 0 {
		t.Error("should have at least one variant")
	}
	if p.DefaultPort() < 0 {
		t.Error("default port must not be negative")
	}
	_, err := p.NewCodec("invalid")
	if err == nil {
		t.Error("stub NewCodec should return error")
	}
}
TESTFILE

echo "Stub created: $DIR"
