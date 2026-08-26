# SDK SERCOS III CmdBridge

O SERCOS III (SErial Real-time COmmunication System) é uma interface digital para controle de movimento, usando topologia em anel sobre fibra óptica ou cobre. Este pacote fornece um codec SERCOS III por meio do utilitário de linha de comando `netx_cli`.

## Ferramenta CLI

Usa o utilitário CLI Hilscher netX SERCOS III.

### Instalação

```bash
# Install Hilscher netX driver and tools
# See https://www.hilscher.com/ for netX drivers
sudo apt-get install netx-driver
```

Verifique: `netx_cli --help`

## Uso

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos"
)

func main() {
    p := sercos.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "S-0-51",
        Count:    2,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := sercos.NewCmdDriver("/usr/bin/netx_cli")
```

## Funções Suportadas

| Função | Descrição                 |
|----------|-----------------------------|
| `read`   | Ler parâmetro IDN/S SERCOS |
| `write`  | Gravar parâmetro IDN/S SERCOS |
| `phase`  | Definir fase de comunicação |

## Testes

```bash
go test ./... -v
```
