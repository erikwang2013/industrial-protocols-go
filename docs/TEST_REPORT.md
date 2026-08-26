# 测试报告 — industrial-protocols-go

日期：2026-08-27
范围：Go 多模块工作区（go.work，42 个模块：kernel / protocols / examples）
命令：`go test ./...` 与 `go test -cover ./...`（因工作区根目录无 go.mod，逐模块等价执行，结果一致）

## 1. 总体结论

- **全绿**：42 个模块、52 个含测试的包全部通过（多次全量重跑验证，含并行测试写入期间）。
- 测试文件 106 个，测试函数 659 个。
- `go vet ./...` 全模块无告警。
- 平均语句覆盖率 **83.3%**；kernel 核心包覆盖率 95%~100%。

## 2. 各模块测试统计与覆盖率

### kernel（11 个包，除 examples 外覆盖率最高）

| 包 | 覆盖率 |
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

### protocols 重点模块

| 模块 | 覆盖率 | 模块 | 覆盖率 |
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

> examples/modbus_basic 无测试文件（示例程序，符合预期）。

## 3. 修复清单

### 3.1 源码 bug 修复（tester-kernel / tester-protocols 发现并验证）

| 文件:行号 | 问题 | 修复 |
|---|---|---|
| kernel/gateway/engine.go:27 | `Transform` 对 `Src`/`Dst`/`Map` 为 nil 的规则在 `rule.Src.Name()` / `rule.Dst.Name()` / `rule.Map(...)` 处 nil 解引用 panic | 循环开头跳过未配置完整的规则（与既有"跳过坏规则"语义一致） |
| protocols/ethernet/modbus/modbus.go:215 | `decodePDU` 对 1 字节异常 PDU（如 `0x81`）读取 `pdu[1]` 越界 panic | 增加 `len(pdu) < 2` 检查，返回解析错误 |
| protocols/ethernet/modbus/modbus.go:230 | `decodePDU` 字节计数（`pdu[1]`）超过实际数据时切片越界 panic（恶意/损坏帧） | 增加 `2+n > len(pdu)` 边界检查，返回解析错误 |
| protocols/ethernet/opcua/opcua.go:63,74 | `encodeHello` 将消息长度回填到偏移 8（协议版本槽），规范要求偏移 4 | 提前记录 `sizePos`（偏移 4），回填到正确位置 |
| protocols/ethernet/opcua/opcua.go:82,89 | `encodeOpenSecureChannel` 的长度字段从未回填，恒为 0 | 末尾按实际帧长回填 |
| protocols/ethernet/ethernetip/ethernetip.go:111,117 | `encodeReadTag` 帧分配 16+len(cipReq) 但 CIP 请求写在 [28:]，多出的 12 字节导致 `copy` 静默丢弃整个 CIP 请求；长度字段同步少 12 | 分配改为 28+len(cipReq)，length 字段改为 4+len(cipReq) |
| protocols/ethernet/profinet/profinet.go:72 | `Decode` 对超界 `blockLen` 切片越界 panic | 增加 `len(data) < 10+blockLen` 检查 |
| protocols/fieldbus/canopen/canopen.go:102 | `Decode` 最小长度检查为 4，但 `unmarshalCAN` 需要 `data[4:8]`，4~7 字节帧 panic | 最小长度改为 8（线格式固定 8 字节） |
| protocols/fieldbus/dnp3/dnp3.go:121 | `decodeTPDU` 对 `length<5`（`apduEnd < apduStart`）或超长 `length`（`apduEnd > len(data)`）切片越界 panic | 统一校验 `length < 5 || apduEnd > len(data)` 后返回错误 |
| protocols/fieldbus/dnp3/dnp3.go:157 | `crc16DNP` 工作变量未按规范掩码 `& 0xFF`，已知向量 `crc16DNP("123456789")` 得 0x69FF（应为 0xEA82） | `temp := (crc ^ b) & 0xFF` |
| protocols/automotive/kline/kline.go:74 | `Decode` 的 `len<5` 检查先于 fast-init 同步帧判断，1 字节同步帧（0x55）被误拒 | 先判空帧与同步帧，再判长度 |

### 3.2 测试修复（修正编译错误 / 断言错误 / 重名 / 过期断言）

| 文件:行号 | 问题 | 修复 |
|---|---|---|
| protocols/automotive/kline/kline_extra_test.go:43 | 常量 `0x68+0x10+0xF1+0x01`(362) 赋值给 byte 编译失败 | 期望改为和校验和的低 8 位 `0x6A` |
| protocols/automotive/saej1850/saej1850_extra_test.go | 在 `kernel.Codec` 接口上调用未导出方法/字段，编译失败 | 新增 `newCodec` helper 断言为 `*j1850Codec`；`TestEncodeMode01PID` 断言偏移修正为 raw[4..6]（ID 4 字节 + 数据前 4 字节） |
| protocols/bridge/sercos1/sercos1_extra_test.go:45 | `TestEncodeReadDefaults` 与既有测试重名，编译失败 | 更名 `TestEncodeReadDefaultsFiber`（保留 fiber 变体覆盖） |
| protocols/building/dali/dali_extra_test.go:68 | `TestEncodeDirectArc` 与既有测试重名，编译失败 | 更名 `TestEncodeDirectArcCmd` |
| protocols/bridge/powerlink/powerlink_extra_test.go:51 | 期望 `write 0x1000 A`，实际为 `0A`（`%X` 按字节固定两位，与既有测试 `DEADBEEF` 约定一致） | 期望改为 `0A` |
| protocols/bridge/sercos/sercos_extra_test.go:62 | 同上，期望 `1` 应为 `01` | 期望改为 `01` |
| protocols/ethernet/opcua/opcua_extra_test.go:35 | 断言旧的 bug 行为（长度在偏移 8），源码修复后失效 | 改为断言协议版本字段为 0（长度由 TestExtraHELMessageSizeAtOffset4 覆盖） |
| protocols/fieldbus/canopen/canopen_extra_test.go:135 | "prove-it" 测试断言反转（`err != nil` 时失败，但解析错误正是期望结果） | 改为 `err == nil` 失败 |
| protocols/ethernet/profinet/profinet_extra_test.go:23 | 同上，断言反转 | 改为 `err == nil` 失败 |
| protocols/fieldbus/dnp3/dnp3_extra_test.go:35,53 | 同上 ×2 | 改为 `err == nil` 失败 |

> 注：`kernel/security/tls_test.go` 的 `stubTransport` 接口实现问题与 `cclink/cclinkie` 的 CRC 期望值，由测试工程师在并行工作中自行修正，本次未改动。

## 4. 遗留风险

1. **低覆盖率模块**：profibus（35.7%）、lonworks / asinterface / foundationfieldbus / iolink（37.9%）——测试仅覆盖少量路径，建议后续补充 Decode/Encode 分支与异常路径。
2. **modbus RTU CRC 未校验**：`decodeRTU` 不校验 CRC（测试 `TestExtraRTUDecodesCorruptCRC` 明确文档化此 GAP 并按现状通过）。与真实设备互通时建议补上。
3. **canopen 线格式限制**：`marshalCAN` 只承载前 4 字节数据（`TestExtraSDOWritePayloadLostOnWire` 文档化），SDO 写多字节负载会丢失。
4. **硬件依赖跳过**：cpci / pci / vme 在无硬件环境下 `t.Skip`（合理）。
5. **仓库残留临时文件**：根目录未跟踪的 `main.go`（探测程序，引用不存在的 `probe/` 模块）、`scripts/`、`docs/*.png`，不属于本次交付，建议清理。
6. **测试仍在并行落地**：本报告基于最后一次全量绿测快照；若测试工程师后续继续提交 `*_extra_test.go`，需重新跑全量回归。
