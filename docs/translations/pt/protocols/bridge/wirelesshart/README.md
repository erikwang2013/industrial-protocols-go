# SDK WirelessHART CmdBridge

O WirelessHART (IEC 62591) é um padrão de rede industrial sem fio baseado no protocolo HART. Este pacote fornece um codec WirelessHART por meio do utilitário de linha de comando `emerson_1410_cli` (Emerson 1410/1420 Wireless Gateway).

## Ferramenta CLI

Usa o utilitário CLI do Emerson 1410/1420 Wireless Gateway.

### Instalação

```bash
# Install Emerson Wireless Gateway software and tools
# Refer to Emerson documentation for 1410/1420 Gateway setup
```

Verifique: `emerson_1410_cli --help`

## Uso

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/wirelesshart"
)

func main() {
    p := wirelesshart.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "TT101",
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := wirelesshart.NewCmdDriver("/usr/bin/emerson_1410_cli")
```

## Funções Suportadas

| Função | Descrição                |
|----------|----------------------------|
| `read`   | Ler parâmetro de dispositivo |
| `write`  | Gravar parâmetro de dispositivo |
| `scan`   | Procurar dispositivos sem fio |

## Testes

```bash
go test ./... -v
```
