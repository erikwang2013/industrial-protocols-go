# SDK EtherCAT CmdBridge

O EtherCAT (Ethernet for Control Automation Technology) é um fieldbus Ethernet industrial de alto desempenho. Este pacote fornece um codec EtherCAT por meio do utilitário de linha de comando `ethercat`.

## Ferramenta CLI

Usa a ferramenta de linha de comando do IgH EtherCAT Master.

### Instalação

```bash
# Debian/Ubuntu
sudo apt-get install ethercat-master

# From source
git clone https://gitlab.com/etherlab.org/ethercat.git
cd ethercat
make && sudo make install
```

Verifique: `ethercat slaves`

## Uso

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/ethercat"
)

func main() {
    p := ethercat.New()
    codec, _ := p.NewCodec("cmd")

    // Upload SDO from address 0x1000
    req := &kernel.Request{
        Function: "upload",
        Address:  "0x1000",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "upload 0x1000 4\n"
    _ = raw
}
```

## Driver

```go
b, c, err := ethercat.NewCmdDriver("/usr/bin/ethercat")
```

## Funções Suportadas

| Função   | Descrição            |
|------------|------------------------|
| `upload`   | Ler SDO de um endereço  |
| `download` | Gravar SDO em um endereço |
| `slaves`   | Listar slaves EtherCAT |

## Testes

```bash
go test ./... -v
```
