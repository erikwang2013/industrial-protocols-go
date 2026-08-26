# SDK do Protocolo Modbus Plus

O Modbus Plus (MB+) é uma rede industrial de alta velocidade com passagem de token desenvolvida pela Modicon (Schneider Electric).

## Variantes

| Variante | Tipo de Bridge | Descrição |
|---------|-------------|-------------|
| `gateway` | GatewayBridge | Conexão TCP com adaptador SA85/BM85 |
| `cmd` | CmdBridge | Wrapper CLI para o utilitário `sa85_cli` |

## Driver de Gateway

```go
driver, codec, err := modbusplus.NewGatewayDriver("192.168.0.160:2010")
handler, err := modbusplus.ReadyHandler("192.168.0.160:2010")
```

## Driver Cmd

```go
driver, codec, err := modbusplus.NewCmdDriver("/usr/bin/sa85_cli")
```

### Instalação da Ferramenta CLI

```bash
# Install SA85 Modbus Plus driver and tools
# Refer to Schneider Electric documentation for SA85 adapter
```

Verifique: `sa85_cli --help`

## Hardware

| Gateway | Interface | IP Padrão |
|---------|-----------|-------------|
| Schneider SA85 | Adaptador ISA Modbus Plus | atribuído pelo host |
| Schneider BM85 | Ponte/multiplexador Modbus Plus | 192.168.0.A0 |
| ProSoft MVI56-MBP | Módulo ControlLogix MB+ | atribuído pelo host |

## Cabeamento

- Cabo twinaxial (RG-62) com conectores BNC para o tronco MB+
- Resistor de terminação (78 ohms em cada extremidade)
- A ponte BM85 conecta o MB+ ao TCP via Ethernet

## Formato do Quadro do Protocolo

- Magic (2 bytes): `0x4D42`
- Destino (1 byte): endereço do nó
- Comando (1 byte): 0x01=read, 0x02=write
- Comprimento (2 bytes): tamanho do payload (big-endian)
- Payload (N bytes): dados
- CRC (2 bytes): Modbus CRC-16 (little-endian)

## Testes

```bash
go test ./... -v
```
