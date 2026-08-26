# SDK del protocolo Modbus Plus

Modbus Plus (MB+) es una red industrial de alta velocidad por paso de testigo desarrollada por Modicon (Schneider Electric).

## Variantes

| Variante | Tipo de puente | Descripción |
|---------|-------------|-------------|
| `gateway` | GatewayBridge | Conexión TCP al adaptador SA85/BM85 |
| `cmd` | CmdBridge | Envoltorio CLI para la utilidad `sa85_cli` |

## Driver de gateway

```go
driver, codec, err := modbusplus.NewGatewayDriver("192.168.0.160:2010")
handler, err := modbusplus.ReadyHandler("192.168.0.160:2010")
```

## Driver Cmd

```go
driver, codec, err := modbusplus.NewCmdDriver("/usr/bin/sa85_cli")
```

### Instalación de la herramienta CLI

```bash
# Instalar el driver y las herramientas de SA85 Modbus Plus
# Consulte la documentación de Schneider Electric sobre el adaptador SA85
```

Verificar: `sa85_cli --help`

## Hardware

| Gateway | Interfaz | IP predeterminada |
|---------|-----------|-------------|
| Schneider SA85 | Adaptador ISA Modbus Plus | asignada por el host |
| Schneider BM85 | Puente/multiplexor Modbus Plus | 192.168.0.A0 |
| ProSoft MVI56-MBP | Módulo ControlLogix MB+ | asignada por el host |

## Cableado

- Cable twinaxial (RG-62) con conectores BNC para el troncal MB+
- Resistencia de terminación (78 ohmios en cada extremo)
- El puente BM85 conecta MB+ a TCP mediante Ethernet

## Formato de trama del protocolo

- Magic (2 bytes): `0x4D42`
- Destino (1 byte): dirección de nodo
- Comando (1 byte): 0x01=lectura, 0x02=escritura
- Longitud (2 bytes): tamaño de la carga útil (big-endian)
- Carga útil (N bytes): datos
- CRC (2 bytes): CRC-16 Modbus (little-endian)

## Pruebas

```bash
go test ./... -v
```
