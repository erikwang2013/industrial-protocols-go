# SDK POWERLINK CmdBridge

O POWERLINK (Ethernet POWERLINK) é um protocolo Ethernet em tempo real para automação industrial. Este pacote fornece um codec POWERLINK por meio do utilitário de linha de comando `openPOWERLINK_demo`.

## Ferramenta CLI

Usa o aplicativo de demonstração da stack openPOWERLINK.

### Instalação

```bash
git clone https://github.com/OpenAutomationTechnologies/openPOWERLINK.git
cd openPOWERLINK
cmake -B build && cmake --build build
sudo cmake --install build
```

Verifique: `openPOWERLINK_demo --help`

## Uso

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/powerlink"
)

func main() {
    p := powerlink.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x2000",
        Count:    8,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## Driver

```go
b, c, err := powerlink.NewCmdDriver("/usr/bin/openPOWERLINK_demo")
```

## Funções Suportadas

| Função | Descrição                   |
|----------|-------------------------------|
| `read`   | Ler entrada do dicionário de objetos |
| `write`  | Gravar entrada do dicionário de objetos |
| `status` | Consultar estado do nó/NMT   |

## Testes

```bash
go test ./... -v
```
