# Relatório de Testes — industrial-protocols-go

Data: 2026-08-27
Escopo: workspace Go multi-módulo (go.work, 42 módulos: kernel / protocols / examples)
Comando: `go test ./...` e `go test -cover ./...` (executados por módulo, de forma equivalente, pois a raiz do workspace não possui go.mod; os resultados são idênticos)

## 1. Conclusão Geral

- **Tudo verde**: os 42 módulos e os 52 pacotes com testes passaram integralmente (verificado com múltiplas reexecuções completas, inclusive durante a escrita de testes em paralelo).
- 106 arquivos de teste, 659 funções de teste.
- `go vet ./...` sem avisos em todos os módulos.
- Cobertura média de statements de **83,3%**; cobertura de 95%~100% nos pacotes principais do kernel.

## 2. Estatísticas de Teste e Cobertura por Módulo

### kernel (11 pacotes, maior cobertura excluindo examples)

| Pacote | Cobertura |
|---|---|
| kernel | 100% |
| kernel/bridge | 80.0% |
| kernel/config | 100% |
| kernel/connection | 95.3% |
| kernel/event | 100% |
| kernel/gateway | 100% |
| kernel/metrics | 100% |
| kernel/pipeline | 100% |
| kernel/security | 100% |
| kernel/transport | 100% |
| kernel/vendor | 100% |

### Módulos de destaque em protocols

| Módulo | Cobertura | Módulo | Cobertura |
|---|---|---|---|
| automotive/flexray | 79.4% | fieldbus/canopen | 73.9% |
| automotive/kline | 93.3% | fieldbus/cclink | 93.8% |
| automotive/lin | 92.5% | fieldbus/cclinkie | 73.5% |
| automotive/most | 100% | fieldbus/devicenet | 73.2% |
| automotive/saej1850 | 84.2% | fieldbus/dnp3 | 88.9% |
| bridge/controlnet | 100% | fieldbus/hart | 96.6% |
| bridge/ethercat | 100% | fieldbus/iec61850 | 91.8% |
| bridge/interbus | 70.0% | iot/hartip | 88.6% |
| bridge/isa100 | 95.2% | iot/mqtt | 89.1% |
| bridge/powerlink | 97.8% | ethernet/bacnet | 95.2% |
| bridge/sercos | 97.9% | ethernet/ethernetip | 96.9% |
| bridge/sercos1 | 95.6% | ethernet/modbus | 67.7% |
| bridge/safej1850 | 100% | ethernet/opcua | 95.8% |
| bridge/wirelesshart | 95.2% | ethernet/profinet | 95.4% |
| building/dali | 92.9% | system/cpci | 78.6% |
| building/lonworks | **37.9%** | system/pci | 78.6% |
| fieldbus/asinterface | **37.9%** | system/vme | 78.6% |
| fieldbus/foundationfieldbus | **37.9%** | fieldbus/profibus | **35.7%** |
| fieldbus/iolink | **37.9%** | bridge/lightbus / modbusplus / worldfip | 72~73.5% |

> examples/modbus_basic não possui arquivos de teste (programa de exemplo, conforme o esperado).

## 3. Lista de Correções

### 3.1 Correções de bugs no código-fonte (descobertos e verificados por tester-kernel / tester-protocols)

| Arquivo:linha | Problema | Correção |
|---|---|---|
| kernel/gateway/engine.go:27 | `Transform` causava panic por dereferência de nil em regras cujo `Src`/`Dst`/`Map` é nil, nos pontos `rule.Src.Name()` / `rule.Dst.Name()` / `rule.Map(...)` | No início do laço, pular regras configuradas de forma incompleta (consistente com a semântica existente de "pular regras inválidas") |
| protocols/ethernet/modbus/modbus.go:215 | `decodePDU` causava panic por acesso fora dos limites ao ler `pdu[1]` em PDU de exceção de 1 byte (ex.: `0x81`) | Adicionada verificação `len(pdu) < 2`, retornando erro de parse |
| protocols/ethernet/modbus/modbus.go:230 | `decodePDU` causava panic por slice fora dos limites quando a contagem de bytes (`pdu[1]`) excedia os dados reais (quadro malicioso/corrompido) | Adicionada verificação de limites `2+n > len(pdu)`, retornando erro de parse |
| protocols/ethernet/opcua/opcua.go:63,74 | `encodeHello` gravava o tamanho da mensagem no offset 8 (slot do número de versão do protocolo); a especificação exige o offset 4 | Registrar `sizePos` antecipadamente (offset 4) e gravar no local correto |
| protocols/ethernet/opcua/opcua.go:82,89 | O campo de tamanho de `encodeOpenSecureChannel` nunca era preenchido, permanecendo sempre 0 | Preencher ao final com o tamanho real do quadro |
| protocols/ethernet/ethernetip/ethernetip.go:111,117 | `encodeReadTag` alocava 16+len(cipReq), mas a requisição CIP era gravada em [28:]; os 12 bytes excedentes faziam o `copy` descartar silenciosamente toda a requisição CIP; o campo de tamanho ficava 12 bytes menor | Alocação alterada para 28+len(cipReq) e campo length alterado para 4+len(cipReq) |
| protocols/ethernet/profinet/profinet.go:72 | `Decode` causava panic por slice fora dos limites com `blockLen` fora dos limites | Adicionada verificação `len(data) < 10+blockLen` |
| protocols/fieldbus/canopen/canopen.go:102 | O tamanho mínimo verificado em `Decode` era 4, mas `unmarshalCAN` precisa de `data[4:8]`; quadros de 4~7 bytes causavam panic | Tamanho mínimo alterado para 8 (formato de linha fixado em 8 bytes) |
| protocols/fieldbus/dnp3/dnp3.go:121 | `decodeTPDU` causava panic por slice fora dos limites com `length<5` (`apduEnd < apduStart`) ou `length` excedente (`apduEnd > len(data)`) | Validação unificada `length < 5 || apduEnd > len(data)`, retornando erro |
| protocols/fieldbus/dnp3/dnp3.go:157 | A variável de trabalho de `crc16DNP` não era mascarada com `& 0xFF` conforme a especificação; o vetor conhecido `crc16DNP("123456789")` retornava 0x69FF (deveria ser 0xEA82) | `temp := (crc ^ b) & 0xFF` |
| protocols/automotive/kline/kline.go:74 | A verificação `len<5` em `Decode` ocorria antes da detecção do quadro de sincronização fast-init; o quadro de sincronização de 1 byte (0x55) era rejeitado indevidamente | Verificar primeiro quadro vazio e quadro de sincronização, depois o tamanho |

### 3.2 Correções em testes (erros de compilação / asserções incorretas / nomes duplicados / asserções desatualizadas)

| Arquivo:linha | Problema | Correção |
|---|---|---|
| protocols/automotive/kline/kline_extra_test.go:43 | A constante `0x68+0x10+0xF1+0x01` (362) atribuída a byte falhava na compilação | Expectativa alterada para o byte baixo do checksum `0x6A` |
| protocols/automotive/saej1850/saej1850_extra_test.go | Chamava métodos/campos não exportados na interface `kernel.Codec`, falhando na compilação | Novo helper `newCodec` com asserção de tipo `*j1850Codec`; a asserção de offset em `TestEncodeMode01PID` foi corrigida para raw[4..6] (ID de 4 bytes + primeiros 4 bytes de dados) |
| protocols/bridge/sercos1/sercos1_extra_test.go:45 | `TestEncodeReadDefaults` duplicava o nome de um teste existente, falhando na compilação | Renomeado para `TestEncodeReadDefaultsFiber` (mantendo a cobertura da variante fiber) |
| protocols/building/dali/dali_extra_test.go:68 | `TestEncodeDirectArc` duplicava o nome de um teste existente, falhando na compilação | Renomeado para `TestEncodeDirectArcCmd` |
| protocols/bridge/powerlink/powerlink_extra_test.go:51 | Expectativa `write 0x1000 A`, mas o resultado real era `0A` (`%X` fixa dois dígitos por byte, consistente com a convenção `DEADBEEF` dos testes existentes) | Expectativa alterada para `0A` |
| protocols/bridge/sercos/sercos_extra_test.go:62 | Idem ao anterior; a expectativa `1` deveria ser `01` | Expectativa alterada para `01` |
| protocols/ethernet/opcua/opcua_extra_test.go:35 | Asserção do comportamento antigo do bug (tamanho no offset 8), invalidada após a correção no código-fonte | Alterado para asserção do campo de versão do protocolo igual a 0 (o tamanho é coberto por TestExtraHELMessageSizeAtOffset4) |
| protocols/fieldbus/canopen/canopen_extra_test.go:135 | Teste "prove-it" com asserção invertida (falhava quando `err != nil`, mas o erro de parse é exatamente o resultado esperado) | Alterado para falhar quando `err == nil` |
| protocols/ethernet/profinet/profinet_extra_test.go:23 | Idem ao anterior, asserção invertida | Alterado para falhar quando `err == nil` |
| protocols/fieldbus/dnp3/dnp3_extra_test.go:35,53 | Idem ao anterior ×2 | Alterado para falhar quando `err == nil` |

> Nota: o problema de implementação da interface `stubTransport` em `kernel/security/tls_test.go` e os valores esperados de CRC de `cclink/cclinkie` foram corrigidos pelo próprio engenheiro de testes durante o trabalho em paralelo; não foram alterados nesta rodada.

## 4. Riscos Remanescentes

1. **Módulos com baixa cobertura**: profibus (35,7%), lonworks / asinterface / foundationfieldbus / iolink (37,9%) — os testes cobrem apenas poucos caminhos; recomenda-se complementar os ramos de Decode/Encode e os caminhos de exceção no futuro.
2. **CRC do Modbus RTU não verificado**: `decodeRTU` não valida o CRC (o teste `TestExtraRTUDecodesCorruptCRC` documenta explicitamente essa lacuna e passa como está). Recomenda-se implementar a verificação para interoperabilidade com dispositivos reais.
3. **Limitação do formato de linha do CANopen**: `marshalCAN` transporta apenas os primeiros 4 bytes de dados (documentado em `TestExtraSDOWritePayloadLostOnWire`); gravações SDO com payload de múltiplos bytes perdem dados.
4. **Dependências de hardware ignoradas**: cpci / pci / vme usam `t.Skip` em ambientes sem hardware (razoável).
5. **Arquivos temporários remanescentes no repositório**: `main.go` não rastreado na raiz (programa de sondagem que referencia o módulo `probe/`, inexistente), `scripts/`, `docs/*.png` — não fazem parte desta entrega; recomenda-se limpá-los.
6. **Testes ainda sendo adicionados em paralelo**: este relatório é baseado no último snapshot verde da execução completa; se o engenheiro de testes continuar enviando `*_extra_test.go`, será necessário reexecutar a regressão completa.
