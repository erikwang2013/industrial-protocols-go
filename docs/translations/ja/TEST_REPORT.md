# テストレポート — industrial-protocols-go

日付：2026-08-27
範囲：Goマルチモジュールワークスペース（go.work、42モジュール：kernel / protocols / examples）
コマンド：`go test ./...` と `go test -cover ./...`（ワークスペースルートにgo.modがないため、モジュール単位で同等に実行、結果は一致）

## 1. 総合結論

- **オールグリーン**：42モジュール、テストを含む52パッケージがすべて合格（並行テスト書き込み期間中を含め、複数回の全量再実行で検証済み）。
- テストファイル106個、テスト関数659個。
- `go vet ./...` は全モジュールで警告なし。
- 平均ステートメントカバレッジ **83.3%**；kernelコアパッケージは95%~100%。

## 2. 各モジュールのテスト統計とカバレッジ

### kernel（11パッケージ、examplesを除きカバレッジ最高）

| パッケージ | カバレッジ |
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

### protocols 主要モジュール

| モジュール | カバレッジ | モジュール | カバレッジ |
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

> examples/modbus_basic にはテストファイルがない（サンプルプログラムのため、想定どおり）。

## 3. 修正リスト

### 3.1 ソースコードのバグ修正（tester-kernel / tester-protocols が発見・検証）

| ファイル:行 | 問題 | 修正 |
|---|---|---|
| kernel/gateway/engine.go:27 | `Transform` が `Src`/`Dst`/`Map` がnilのルールに対して `rule.Src.Name()` / `rule.Dst.Name()` / `rule.Map(...)` でnil参照によりpanic | ループ先頭で未設定のルールをスキップ（既存の「不正ルールはスキップ」という意味論と一致） |
| protocols/ethernet/modbus/modbus.go:215 | `decodePDU` が1バイトの異常PDU（例：`0x81`）に対して `pdu[1]` を読み取り範囲外でpanic | `len(pdu) < 2` のチェックを追加し、パースエラーを返す |
| protocols/ethernet/modbus/modbus.go:230 | `decodePDU` のバイトカウント（`pdu[1]`）が実データを超える場合にスライス範囲外でpanic（悪意のある/破損フレーム） | `2+n > len(pdu)` の境界チェックを追加し、パースエラーを返す |
| protocols/ethernet/opcua/opcua.go:63,74 | `encodeHello` がメッセージ長をオフセット8（プロトコルバージョンスロット）に書き戻しており、仕様ではオフセット4 | `sizePos`（オフセット4）を事前に記録し、正しい位置に書き戻す |
| protocols/ethernet/opcua/opcua.go:82,89 | `encodeOpenSecureChannel` の長さフィールドが一度も書き戻されず、常に0 | 末尾で実際のフレーム長に基づき書き戻す |
| protocols/ethernet/ethernetip/ethernetip.go:111,117 | `encodeReadTag` が16+len(cipReq) を確保するが、CIPリクエストは [28:] に書き込まれており、余分な12バイトにより `copy` がCIPリクエスト全体を黙って破棄。長さフィールドも12バイト不足 | 確保を 28+len(cipReq) に変更、lengthフィールドを 4+len(cipReq) に変更 |
| protocols/ethernet/profinet/profinet.go:72 | `Decode` が範囲外の `blockLen` でスライス範囲外panic | `len(data) < 10+blockLen` のチェックを追加 |
| protocols/fieldbus/canopen/canopen.go:102 | `Decode` の最小長チェックが4だが、`unmarshalCAN` は `data[4:8]` を必要とし、4〜7バイトのフレームでpanic | 最小長を8に変更（線形式は固定8バイト） |
| protocols/fieldbus/dnp3/dnp3.go:121 | `decodeTPDU` が `length<5`（`apduEnd < apduStart`）または過大な `length`（`apduEnd > len(data)`）でスライス範囲外panic | `length < 5 || apduEnd > len(data)` を一括検証してからエラーを返す |
| protocols/fieldbus/dnp3/dnp3.go:157 | `crc16DNP` の作業変数が仕様どおり `& 0xFF` でマスクされておらず、既知ベクトル `crc16DNP("123456789")` が 0x69FF（正しくは 0xEA82） | `temp := (crc ^ b) & 0xFF` |
| protocols/automotive/kline/kline.go:74 | `Decode` の `len<5` チェックがfast-init同期フレーム判定より先にあり、1バイトの同期フレーム（0x55）が誤って拒否される | 空フレームと同期フレームを先に判定し、その後で長さを判定 |

### 3.2 テスト修正（コンパイルエラー / アサーション誤り / 同名重複 / 陳腐化アサーションの修正）

| ファイル:行 | 問題 | 修正 |
|---|---|---|
| protocols/automotive/kline/kline_extra_test.go:43 | 定数 `0x68+0x10+0xF1+0x01`(362) をbyteに代入するとコンパイルエラー | 期待値をチェックサムの下位8ビット `0x6A` に変更 |
| protocols/automotive/saej1850/saej1850_extra_test.go | `kernel.Codec` インターフェース上で非公開メソッド/フィールドを呼び出しコンパイルエラー | `newCodec` ヘルパーを追加し `*j1850Codec` であることをアサート；`TestEncodeMode01PID` のアサーションオフセットを raw[4..6]（ID 4バイト + データ先頭4バイト）に修正 |
| protocols/bridge/sercos1/sercos1_extra_test.go:45 | `TestEncodeReadDefaults` が既存テストと同名でコンパイルエラー | `TestEncodeReadDefaultsFiber` に改名（fiberバリアントのカバレッジを維持） |
| protocols/building/dali/dali_extra_test.go:68 | `TestEncodeDirectArc` が既存テストと同名でコンパイルエラー | `TestEncodeDirectArcCmd` に改名 |
| protocols/bridge/powerlink/powerlink_extra_test.go:51 | 期待値 `write 0x1000 A`、実際は `0A`（`%X` はバイト単位で固定2桁、既存テストの `DEADBEEF` 規約と一致） | 期待値を `0A` に変更 |
| protocols/bridge/sercos/sercos_extra_test.go:62 | 同上、期待値 `1` は `01` であるべき | 期待値を `01` に変更 |
| protocols/ethernet/opcua/opcua_extra_test.go:35 | 旧バグ挙動（長さがオフセット8）をアサートしており、ソース修正後に失效 | プロトコルバージョンフィールドが0であることをアサートするよう変更（長さは TestExtraHELMessageSizeAtOffset4 がカバー） |
| protocols/fieldbus/canopen/canopen_extra_test.go:135 | "prove-it" テストのアサーションが逆（`err != nil` で失敗するが、パースエラーこそ期待結果） | `err == nil` で失敗するよう変更 |
| protocols/ethernet/profinet/profinet_extra_test.go:23 | 同上、アサーションが逆 | `err == nil` で失敗するよう変更 |
| protocols/fieldbus/dnp3/dnp3_extra_test.go:35,53 | 同上 ×2 | `err == nil` で失敗するよう変更 |

> 注：`kernel/security/tls_test.go` の `stubTransport` インターフェース実装の問題と `cclink/cclinkie` のCRC期待値は、テストエンジニアが並行作業で自己修正したため、今回は変更していない。

## 4. 残存リスク

1. **カバレッジ低のモジュール**：profibus（35.7%）、lonworks / asinterface / foundationfieldbus / iolink（37.9%）——テストは一部のパスのみをカバーしており、今後Decode/Encodeの分岐と異常パスを追加することを推奨。
2. **modbus RTU CRC未検証**：`decodeRTU` はCRCを検証しない（テスト `TestExtraRTUDecodesCorruptCRC` がこのGAPを明示的に文書化し、現状のまま合格）。実機との相互接続時は追加を推奨。
3. **canopen 線形式の制限**：`marshalCAN` は先頭4バイトのデータのみを搬送（`TestExtraSDOWritePayloadLostOnWire` で文書化）。SDO書き込みのマルチバイトペイロードは失われる。
4. **ハードウェア依存のスキップ**：cpci / pci / vme はハードウェアのない環境では `t.Skip`（妥当）。
5. **リポジトリ内の残留一時ファイル**：ルートにある未追跡の `main.go`（プローブプログラム。存在しない `probe/` モジュールを参照）、`scripts/`、`docs/*.png` は今回の成果物に含まれず、削除を推奨。
6. **テストの並行追加が継続中**：本レポートは最後の全量グリーン実行スナップショットに基づく。テストエンジニアが今後も `*_extra_test.go` を追加し続ける場合は、全量リグレッションを再実行する必要がある。
