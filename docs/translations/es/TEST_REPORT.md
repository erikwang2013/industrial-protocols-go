# Informe de pruebas — industrial-protocols-go

Fecha: 2026-08-27
Alcance: workspace multi-módulo de Go (go.work, 42 módulos: kernel / protocols / examples)
Comando: `go test ./...` y `go test -cover ./...` (al no existir go.mod en la raíz del workspace, se ejecutaron equivalentemente por módulo, con resultados idénticos)

## 1. Conclusión general

- **Todo en verde**: los 42 módulos y los 52 paquetes con pruebas pasaron todos (verificado con múltiples ejecuciones completas, incluso durante la escritura de pruebas en paralelo).
- 106 archivos de prueba, 659 funciones de prueba.
- `go vet ./...` sin avisos en ningún módulo.
- Cobertura media por sentencias **83,3 %**; los paquetes centrales del kernel entre 95 % y 100 %.

## 2. Estadísticas y cobertura de pruebas por módulo

### kernel (11 paquetes, la cobertura más alta salvo examples)

| Paquete | Cobertura |
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

### Módulos destacados de protocols

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

> examples/modbus_basic no tiene archivos de prueba (es un programa de ejemplo, lo esperado).

## 3. Lista de correcciones

### 3.1 Correcciones de bugs en el código fuente (detectados y verificados por tester-kernel / tester-protocols)

| Archivo:línea | Problema | Corrección |
|---|---|---|
| kernel/gateway/engine.go:27 | `Transform` provoca panic por desreferencia de nil en `rule.Src.Name()` / `rule.Dst.Name()` / `rule.Map(...)` para reglas con `Src`/`Dst`/`Map` nil | Saltar las reglas mal configuradas al inicio del bucle (coherente con la semántica existente de "saltar reglas incorrectas") |
| protocols/ethernet/modbus/modbus.go:215 | `decodePDU` lee `pdu[1]` fuera de límites con una PDU anómala de 1 byte (p. ej. `0x81`) | Añadir comprobación de `len(pdu) < 2` y devolver error de análisis |
| protocols/ethernet/modbus/modbus.go:230 | `decodePDU` provoca panic por corte fuera de límites cuando el conteo de bytes (`pdu[1]`) excede los datos reales (tramas maliciosas/corruptas) | Añadir comprobación de límites `2+n > len(pdu)` y devolver error de análisis |
| protocols/ethernet/opcua/opcua.go:63,74 | `encodeHello` rellenaba la longitud del mensaje en el offset 8 (ranura de versión de protocolo); la especificación exige offset 4 | Registrar `sizePos` (offset 4) por adelantado y rellenar en la posición correcta |
| protocols/ethernet/opcua/opcua.go:82,89 | El campo de longitud de `encodeOpenSecureChannel` nunca se rellenaba, quedaba siempre en 0 | Rellenar al final con la longitud real de la trama |
| protocols/ethernet/ethernetip/ethernetip.go:111,117 | `encodeReadTag` asignaba 16+len(cipReq) para la trama, pero la solicitud CIP se escribía en [28:]; los 12 bytes sobrantes hacían que `copy` descartara silenciosamente toda la solicitud CIP; el campo de longitud iba 12 bytes corto | Asignar 28+len(cipReq) y poner el campo length en 4+len(cipReq) |
| protocols/ethernet/profinet/profinet.go:72 | `Decode` provocaba panic por corte fuera de límites con `blockLen` desbordado | Añadir comprobación de `len(data) < 10+blockLen` |
| protocols/fieldbus/canopen/canopen.go:102 | La comprobación de longitud mínima de `Decode` era 4, pero `unmarshalCAN` necesita `data[4:8]`; las tramas de 4~7 bytes provocaban panic | Longitud mínima cambiada a 8 (el formato de línea fija 8 bytes) |
| protocols/fieldbus/dnp3/dnp3.go:121 | `decodeTPDU` provocaba panic por corte fuera de límites con `length<5` (`apduEnd < apduStart`) o `length` excesivo (`apduEnd > len(data)`) | Validar unificadamente `length < 5 || apduEnd > len(data)` y devolver error |
| protocols/fieldbus/dnp3/dnp3.go:157 | La variable de trabajo de `crc16DNP` no se enmascaraba con `& 0xFF` según la especificación; el vector conocido `crc16DNP("123456789")` daba 0x69FF (debería ser 0xEA82) | `temp := (crc ^ b) & 0xFF` |
| protocols/automotive/kline/kline.go:74 | La comprobación `len<5` de `Decode` se ejecutaba antes de la detección de tramas de sincronización fast-init; una trama de sincronización de 1 byte (0x55) se rechazaba por error | Primero comprobar tramas vacías y de sincronización, luego la longitud |

### 3.2 Correcciones en las pruebas (errores de compilación / aserciones erróneas / nombres duplicados / aserciones obsoletas)

| Archivo:línea | Problema | Corrección |
|---|---|---|
| protocols/automotive/kline/kline_extra_test.go:43 | La constante `0x68+0x10+0xF1+0x01`(362) no compilaba al asignarla a un byte | Cambiar el valor esperado al byte bajo del checksum, `0x6A` |
| protocols/automotive/saej1850/saej1850_extra_test.go | Se llamaban métodos/campos no exportados sobre la interfaz `kernel.Codec`; no compilaba | Añadir helper `newCodec` que afirme el tipo `*j1850Codec`; la aserción de offset de `TestEncodeMode01PID` se corrige a raw[4..6] (ID de 4 bytes + 4 primeros bytes de datos) |
| protocols/bridge/sercos1/sercos1_extra_test.go:45 | `TestEncodeReadDefaults` duplicaba el nombre de una prueba existente; no compilaba | Renombrar a `TestEncodeReadDefaultsFiber` (conserva la cobertura de la variante fiber) |
| protocols/building/dali/dali_extra_test.go:68 | `TestEncodeDirectArc` duplicaba el nombre de una prueba existente; no compilaba | Renombrar a `TestEncodeDirectArcCmd` |
| protocols/bridge/powerlink/powerlink_extra_test.go:51 | Se esperaba `write 0x1000 A`, el valor real es `0A` (`%X` fija dos dígitos por byte, coherente con la convención `DEADBEEF` de las pruebas existentes) | Cambiar el valor esperado a `0A` |
| protocols/bridge/sercos/sercos_extra_test.go:62 | Igual que arriba, se esperaba `1` y debe ser `01` | Cambiar el valor esperado a `01` |
| protocols/ethernet/opcua/opcua_extra_test.go:35 | La aserción validaba el comportamiento del bug antiguo (longitud en offset 8), dejó de ser válida tras la corrección del código fuente | Cambiar la aserción para verificar que el campo de versión de protocolo es 0 (la longitud queda cubierta por TestExtraHELMessageSizeAtOffset4) |
| protocols/fieldbus/canopen/canopen_extra_test.go:135 | La prueba "prove-it" tenía la aserción invertida (fallaba con `err != nil`, pero el error de análisis es justo el resultado esperado) | Cambiar a que falle con `err == nil` |
| protocols/ethernet/profinet/profinet_extra_test.go:23 | Igual que arriba, aserción invertida | Cambiar a que falle con `err == nil` |
| protocols/fieldbus/dnp3/dnp3_extra_test.go:35,53 | Igual que arriba ×2 | Cambiar a que falle con `err == nil` |

> Nota: el problema de implementación de la interfaz `stubTransport` en `kernel/security/tls_test.go` y los valores esperados de CRC de `cclink/cclinkie` fueron corregidos por el ingeniero de pruebas en su trabajo paralelo; no se modificaron en esta ocasión.

## 4. Riesgos pendientes

1. **Módulos con cobertura baja**: profibus (35,7 %), lonworks / asinterface / foundationfieldbus / iolink (37,9 %) — las pruebas solo cubren unas pocas rutas; se recomienda añadir después ramas de Decode/Encode y rutas de error.
2. **CRC de Modbus RTU sin validar**: `decodeRTU` no valida el CRC (la prueba `TestExtraRTUDecodesCorruptCRC` documenta explícitamente esta laguna y pasa tal cual). Se recomienda añadirla para la interoperabilidad con dispositivos reales.
3. **Limitación del formato de línea de canopen**: `marshalCAN` solo transporta los 4 primeros bytes de datos (`TestExtraSDOWritePayloadLostOnWire` lo documenta); las escrituras SDO con carga útil de varios bytes pierden datos.
4. **Omisión por dependencia de hardware**: cpci / pci / vme usan `t.Skip` en entornos sin hardware (razonable).
5. **Archivos temporales residuales en el repositorio**: en la raíz hay `main.go` sin rastrear (programa de sondeo que referencia el módulo `probe/` inexistente), `scripts/` y `docs/*.png`, que no pertenecen a esta entrega; se recomienda limpiarlos.
6. **Las pruebas aún se están integrando en paralelo**: este informe se basa en la última instantánea completa en verde; si el ingeniero de pruebas continúa enviando `*_extra_test.go`, habrá que volver a ejecutar la regresión completa.
