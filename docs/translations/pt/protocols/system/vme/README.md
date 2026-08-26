# SDK do Protocolo VME / VPX

Implementa `kernel.Protocol` para acesso aos barramentos VMEbus e VPX via procfs do Linux.

## Protocolo

- **Nome**: `vme`
- **Variantes**: `vme`
- **Porta padrão**: 0 (barramento mapeado em memória)
- **Transporte**: procfs `/proc/vme/<slot>`

O codec repassa bytes brutos entre a aplicação e o espaço de endereçamento
VME. Leituras e gravações vão diretamente ao barramento por meio da
interface procfs fornecida pelo driver VME do kernel.

## Requisitos do Kernel

O seguinte módulo do kernel deve estar carregado:

- `vme_tsi148` -- driver de ponte VME Tundra TSI148 (mais comum)
  - Também suportados: `vme_ca91cx42` (Universe II), `vme_user`

O sistema de arquivos procfs deve estar montado em `/proc`.

Permissões necessárias:
- Acesso root é necessário para abrir `/proc/vme/<slot>` para leitura-gravação
- Os nós de dispositivo pertencem a root:root

## Driver

```go
d, err := vme.NewVMEDriver(0) // slot 0
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Use tr.Read/tr.Write for raw VME bus access
```

## Modos de Endereçamento VME

O codec repassa todos os bytes de modificador de endereço e endereço de
forma transparente. As aplicações devem antepor as informações de
endereçamento ao payload de dados:

- **A16**: espaço de endereço de I/O curto de 16 bits
- **A24**: espaço de endereço padrão de 24 bits
- **A32**: espaço de endereço estendido de 32 bits

## Compatibilidade VPX

Sistemas VPX (VITA 46) que expõem uma interface procfs compatível com VME
podem usar este driver. A numeração de slots é a mesma.
