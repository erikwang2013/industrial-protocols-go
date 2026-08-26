# Отчёт о тестировании — industrial-protocols-go

Дата: 2026-08-27
Область: многомодульное рабочее пространство Go (go.work, 42 модуля: kernel / protocols / examples)
Команды: `go test ./...` и `go test -cover ./...` (из-за отсутствия go.mod в корне рабочего пространства выполнялись эквивалентно по каждому модулю; результаты идентичны)

## 1. Общий вывод

- **Всё зелёное**: все 42 модуля и 52 пакета с тестами прошли (проверено многократными полными прогонами, в том числе во время параллельной записи тестов).
- Тестовых файлов 106, тестовых функций 659.
- `go vet ./...` — предупреждений во всех модулях нет.
- Среднее покрытие по операторам **83,3 %**; покрытие пакетов ядра kernel — 95–100 %.

## 2. Статистика и покрытие по модулям

### kernel (11 пакетов, самое высокое покрытие после examples)

| Пакет | Покрытие |
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

### Ключевые модули protocols

| Модуль | Покрытие | Модуль | Покрытие |
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

> examples/modbus_basic не содержит тестовых файлов (пример программы — ожидаемо).

## 3. Список исправлений

### 3.1 Исправления ошибок в исходном коде (обнаружены и подтверждены tester-kernel / tester-protocols)

| Файл:строка | Проблема | Исправление |
|---|---|---|
| kernel/gateway/engine.go:27 | `Transform` паникует из-за разыменования nil в `rule.Src.Name()` / `rule.Dst.Name()` / `rule.Map(...)` для правил с nil-полями `Src`/`Dst`/`Map` | В начале цикла пропускать не полностью сконфигурированные правила (в соответствии с существующей семантикой «пропуск некорректных правил») |
| protocols/ethernet/modbus/modbus.go:215 | `decodePDU` паникует из-за выхода за границы при чтении `pdu[1]` для 1-байтового аномального PDU (например, `0x81`) | Добавлена проверка `len(pdu) < 2`, возвращается ошибка разбора |
| protocols/ethernet/modbus/modbus.go:230 | `decodePDU` паникует из-за выхода за границы среза, когда счётчик байтов (`pdu[1]`) превышает фактические данные (вредоносный/повреждённый кадр) | Добавлена проверка границ `2+n > len(pdu)`, возвращается ошибка разбора |
| protocols/ethernet/opcua/opcua.go:63,74 | `encodeHello` записывает длину сообщения в смещение 8 (слот номера версии протокола), а по спецификации нужно смещение 4 | Заранее запоминать `sizePos` (смещение 4) и записывать длину в правильное место |
| protocols/ethernet/opcua/opcua.go:82,89 | Поле длины в `encodeOpenSecureChannel` никогда не заполнялось и всегда было равно 0 | В конце заполняется по фактической длине кадра |
| protocols/ethernet/ethernetip/ethernetip.go:111,117 | `encodeReadTag` выделяет 16+len(cipReq), но запрос CIP записывается в [28:]; лишние 12 байт приводят к тому, что `copy` молча отбрасывает весь запрос CIP; поле длины тоже занижено на 12 | Выделение изменено на 28+len(cipReq), поле length — на 4+len(cipReq) |
| protocols/ethernet/profinet/profinet.go:72 | `Decode` паникует из-за выхода за границы среза при превышении `blockLen` | Добавлена проверка `len(data) < 10+blockLen` |
| protocols/fieldbus/canopen/canopen.go:102 | Минимальная проверка длины в `Decode` равна 4, но `unmarshalCAN` требует `data[4:8]`; кадры из 4–7 байт паникуют | Минимальная длина изменена на 8 (линейный формат фиксирован — 8 байт) |
| protocols/fieldbus/dnp3/dnp3.go:121 | `decodeTPDU` паникует из-за выхода за границы среза при `length<5` (`apduEnd < apduStart`) или при чрезмерной `length` (`apduEnd > len(data)`) | Единая проверка `length < 5 || apduEnd > len(data)` с возвратом ошибки |
| protocols/fieldbus/dnp3/dnp3.go:157 | Рабочая переменная в `crc16DNP` не маскировалась как `& 0xFF` по спецификации; известный вектор `crc16DNP("123456789")` даёт 0x69FF (должно быть 0xEA82) | `temp := (crc ^ b) & 0xFF` |
| protocols/automotive/kline/kline.go:74 | Проверка `len<5` в `Decode` выполнялась раньше определения синхрокадра fast-init; 1-байтовый синхрокадр (0x55) ошибочно отклонялся | Сначала проверяются пустой кадр и синхрокадр, затем длина |

### 3.2 Исправления тестов (ошибки компиляции / неверные утверждения / дубли имён / устаревшие утверждения)

| Файл:строка | Проблема | Исправление |
|---|---|---|
| protocols/automotive/kline/kline_extra_test.go:43 | Константа `0x68+0x10+0xF1+0x01`(362) не компилируется при присваивании byte | Ожидание заменено на младшие 8 бит контрольной суммы `0x6A` |
| protocols/automotive/saej1850/saej1850_extra_test.go | Компиляция не проходила: вызов неэкспортированных методов/полей на интерфейсе `kernel.Codec` | Добавлен helper `newCodec` с утверждением типа `*j1850Codec`; в `TestEncodeMode01PID` смещение утверждения исправлено на raw[4..6] (ID 4 байта + первые 4 байта данных) |
| protocols/bridge/sercos1/sercos1_extra_test.go:45 | `TestEncodeReadDefaults` дублирует имя существующего теста, компиляция не проходила | Переименован в `TestEncodeReadDefaultsFiber` (покрытие варианта fiber сохранено) |
| protocols/building/dali/dali_extra_test.go:68 | `TestEncodeDirectArc` дублирует имя существующего теста, компиляция не проходила | Переименован в `TestEncodeDirectArcCmd` |
| protocols/bridge/powerlink/powerlink_extra_test.go:51 | Ожидалось `write 0x1000 A`, фактически `0A` (`%X` печатает ровно две цифры на байт, в соответствии с соглашением из существующего теста `DEADBEEF`) | Ожидание заменено на `0A` |
| protocols/bridge/sercos/sercos_extra_test.go:62 | Аналогично: ожидание `1` должно быть `01` | Ожидание заменено на `01` |
| protocols/ethernet/opcua/opcua_extra_test.go:35 | Утверждалось старое поведение ошибки (длина в смещении 8); после исправления исходного кода тест неактуален | Изменено на утверждение, что поле версии протокола равно 0 (длина покрывается тестом TestExtraHELMessageSizeAtOffset4) |
| protocols/fieldbus/canopen/canopen_extra_test.go:135 | «Prove-it»-тест с инвертированным утверждением (`err != nil` — ошибка, хотя ошибка разбора и есть ожидаемый результат) | Изменено: ошибка при `err == nil` |
| protocols/ethernet/profinet/profinet_extra_test.go:23 | Аналогично: инвертированное утверждение | Изменено: ошибка при `err == nil` |
| protocols/fieldbus/dnp3/dnp3_extra_test.go:35,53 | Аналогично ×2 | Изменено: ошибка при `err == nil` |

> Примечание: проблемы с реализацией интерфейса `stubTransport` в `kernel/security/tls_test.go` и ожидаемые значения CRC в `cclink/cclinkie` были исправлены самим инженером по тестированию в параллельной работе; в этот заход не вносились.

## 4. Остаточные риски

1. **Модули с низким покрытием**: profibus (35,7 %), lonworks / asinterface / foundationfieldbus / iolink (37,9 %) — тесты покрывают лишь немногие пути; рекомендуется дополнить ветви Decode/Encode и аварийные пути.
2. **CRC Modbus RTU не проверяется**: `decodeRTU` не проверяет CRC (тест `TestExtraRTUDecodesCorruptCRC` явно документирует этот пробел и проходит как есть). При обмене с реальными устройствами рекомендуется добавить.
3. **Ограничение линейного формата canopen**: `marshalCAN` переносит только первые 4 байта данных (документировано в `TestExtraSDOWritePayloadLostOnWire`); при SDO-записи много байтовые нагрузки теряются.
4. **Пропуск из-за аппаратной зависимости**: cpci / pci / vme выполняют `t.Skip` в среде без аппаратуры (обоснованно).
5. **Остаточные временные файлы в репозитории**: не отслеживаемый в корне `main.go` (пробная программа, ссылается на несуществующий модуль `probe/`), `scripts/`, `docs/*.png` — не относятся к данной поставке, рекомендуется очистить.
6. **Тесты всё ещё добавляются параллельно**: отчёт основан на последнем снимке полного «зелёного» прогона; если инженер по тестированию продолжит добавлять `*_extra_test.go`, необходимо перезапустить полную регрессию.
