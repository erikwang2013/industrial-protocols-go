# SDK ISA100.11a CmdBridge

O ISA100.11a é um padrão de rede industrial sem fio para automação de processos. Este pacote fornece um codec ISA100.11a por meio do utilitário de linha de comando `yfgw410_cli` (gateway sem fio de campo Yokogawa YFGW410).

## Ferramenta CLI

Usa o CLI do Yokogawa YFGW410 Field Wireless Gateway.

### Instalação

```bash
# Install Yokogawa YFGW410 gateway software and tools
# Refer to Yokogawa documentation for Field Wireless Gateway setup
```

Verifique: `yfgw410_cli --help`

## Uso

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/isa100"
)

func main() {
    p := isa100.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "DEV001",
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := isa100.NewCmdDriver("/usr/bin/yfgw410_cli")
```

## Funções Suportadas

| Função | Descrição                 |
|----------|-----------------------------|
| `read`   | Ler atributo de dispositivo |
| `write`  | Gravar atributo de dispositivo |
| `list`   | Listar dispositivos provisionados |

## Testes

```bash
go test ./... -v
```
