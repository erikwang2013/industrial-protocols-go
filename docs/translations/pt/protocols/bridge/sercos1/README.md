# SDK SERCOS I/II CmdBridge

O SERCOS I/II é a versão legada serial por fibra óptica da interface SERCOS para controle de movimento digital. Este pacote fornece um codec SERCOS I/II por meio do utilitário de linha de comando `sercos_cli`.

## Ferramenta CLI

Usa um utilitário CLI de interface SERCOS de fibra óptica.

### Instalação

```bash
# SERCOS interface card driver and tools
# Refer to vendor documentation for the specific SERCOS master card
```

Verifique: `sercos_cli --help`

## Uso

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos1"
)

func main() {
    p := sercos1.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := sercos1.NewCmdDriver("/usr/bin/sercos_cli")
```

## Funções Suportadas

| Função | Descrição              |
|----------|--------------------------|
| `read`   | Ler IDN SERCOS          |
| `write`  | Gravar IDN SERCOS       |
| `status` | Consultar status do drive |

## Testes

```bash
go test ./... -v
```
