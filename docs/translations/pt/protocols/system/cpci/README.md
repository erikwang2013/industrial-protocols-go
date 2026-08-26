# SDK do Protocolo CompactPCI

Implementa `kernel.Protocol` para acesso ao barramento CompactPCI via sysfs do Linux.

## Protocolo

- **Nome**: `cpci`
- **Variantes**: `cpci`
- **Porta padrão**: 0 (barramento mapeado em memória)
- **Transporte**: sysfs `/sys/bus/pci/devices/<BDF>/config`

O CompactPCI (PICMG 2.0) usa a mesma interface elétrica e de software
do PCI convencional. O codec repassa bytes brutos entre a aplicação e o
espaço de configuração PCI via sysfs.

## Requisitos do Kernel

Os seguintes módulos do kernel devem estar carregados:

- `pcieport` -- driver de porta PCI Express (para sistemas CPCIe híbridos)
- `pci_sysfs` -- interface sysfs PCI (integrada na maioria dos kernels)
- `cpci_hotplug` -- controlador de hotplug CompactPCI (opcional, para hot-swap)

O sistema de arquivos sysfs deve estar montado em `/sys`. Este é o padrão
em todas as distribuições Linux modernas.

Permissões necessárias:
- Acesso root ou `CAP_SYS_ADMIN` para acesso ao espaço de configuração
- O arquivo de configuração pertence a root:root com modo 0600

## Driver

```go
d, err := cpci.NewCPCIDriver("0000:02:00.0")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Use tr.Read/tr.Write for raw CPCI config space access
```

## CompactPCI vs PCI

O CompactPCI usa a enumeração e o espaço de configuração padrão do barramento PCI.
Principais diferenças em relação ao PCI de desktop:

- **Fator de forma 3U/6U**: mecânica Eurocard com conectores pino-soquete
- **Numeração do barramento**: cada segmento de chassis CPCI recebe seu próprio número de barramento PCI
- **Hot swap**: o hot swap PICMG 2.1 usa o modelo padrão de hotplug do PCI
- **Slot de sistema**: Barramento 0, dispositivo 0 é o controlador do slot de sistema

## Formato BDF

O endereço do barramento é uma string BDF (Bus:Device.Function):
- `0000:02:00.0` -- domínio 0000, barramento 02, dispositivo 00, função 0
- `0000:02:08.0` -- domínio 0000, barramento 02, dispositivo 08, função 0
