# SDK ControlNet CmdBridge

O ControlNet é um protocolo de rede industrial em tempo real desenvolvido pela Allen-Bradley (Rockwell Automation) para troca de dados de alta velocidade e sensível ao tempo. Este pacote fornece um codec ControlNet por meio do utilitário de linha de comando `1784-pcic-cli`.

## Ferramenta CLI

Usa o utilitário CLI da placa de interface ControlNet 1784-PCIC.

### Instalação

```bash
# Install Rockwell 1784-PCIC driver and tools
# Refer to Rockwell Automation documentation for RSLinx Classic SDK
```

Verifique: `1784-pcic-cli --help`

## Uso

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/controlnet"
)

func main() {
    p := controlnet.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x10",
        Count:    8,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := controlnet.NewCmdDriver("/usr/bin/1784-pcic-cli")
```

## Funções Suportadas

| Função | Descrição               |
|----------|---------------------------|
| `read`   | Ler de um nó ControlNet  |
| `write`  | Gravar em um nó ControlNet |
| `status` | Consultar status da placa PCIC |

## Testes

```bash
go test ./... -v
```
