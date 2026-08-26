# SDK do Protocolo PCI / PCIe

Implementa `kernel.Protocol` para acesso aos barramentos PCI e PCI Express via sysfs do Linux.

## Protocolo

- **Nome**: `pci`
- **Variantes**: `pci`
- **Porta padrão**: 0 (barramento mapeado em memória)
- **Transporte**: sysfs `/sys/bus/pci/devices/<BDF>/config`

O codec repassa bytes brutos entre a aplicação e o espaço de configuração
PCI. Leituras e gravações vão diretamente aos registradores de configuração
do dispositivo por meio do arquivo de configuração sysfs.

## Requisitos do Kernel

Os seguintes módulos do kernel devem estar carregados:

- `pcieport` -- driver de porta PCI Express
- `pci_sysfs` -- interface sysfs PCI (integrada na maioria dos kernels)

O sistema de arquivos sysfs deve estar montado em `/sys`. Este é o padrão
em todas as distribuições Linux modernas.

Permissões necessárias:
- Acesso root ou `CAP_SYS_ADMIN` para acesso ao espaço de configuração
- O arquivo de configuração pertence a root:root com modo 0600

## Driver

```go
d, err := pci.NewPCIDriver("0000:00:1f.3")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Use tr.Read/tr.Write for raw PCI config space access
```

## Formato BDF

O endereço do barramento é uma string BDF (Bus:Device.Function):
- `0000:00:1f.3` -- domínio 0000, barramento 00, dispositivo 1f, função 3
- `0000:01:00.0` -- domínio 0000, barramento 01, dispositivo 00, função 0
