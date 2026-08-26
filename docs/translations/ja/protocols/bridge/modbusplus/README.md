# Modbus Plus プロトコルSDK

Modbus Plus（MB+）は、Modicon（Schneider Electric）によって開発された高速トークンパッシング産業ネットワークです。

## バリアント

| バリアント | ブリッジタイプ | 説明 |
|---------|-------------|-------------|
| `gateway` | GatewayBridge | SA85/BM85アダプタへのTCP接続 |
| `cmd` | CmdBridge | `sa85_cli` ユーティリティのCLIラッパー |

## ゲートウェイドライバ

```go
driver, codec, err := modbusplus.NewGatewayDriver("192.168.0.160:2010")
handler, err := modbusplus.ReadyHandler("192.168.0.160:2010")
```

## Cmdドライバ

```go
driver, codec, err := modbusplus.NewCmdDriver("/usr/bin/sa85_cli")
```

### CLIツールのインストール

```bash
# Install SA85 Modbus Plus driver and tools
# Refer to Schneider Electric documentation for SA85 adapter
```

確認：`sa85_cli --help`

## ハードウェア

| ゲートウェイ | インターフェース | デフォルトIP |
|---------|-----------|-------------|
| Schneider SA85 | ISA Modbus Plusアダプタ | host割り当て |
| Schneider BM85 | Modbus Plusブリッジ/マルチプレクサ | 192.168.0.A0 |
| ProSoft MVI56-MBP | ControlLogix MB+モジュール | host割り当て |

## 配線

- MB+トランク用のツイナックスケーブル（RG-62）+ BNCコネクタ
- 終端抵抗（各端に78オーム）
- BM85ブリッジがEthernet経由でMB+をTCPに接続

## プロトコルフレーム形式

- マジック（2バイト）：`0x4D42`
- 宛先（1バイト）：ノードアドレス
- コマンド（1バイト）：0x01=読み取り、0x02=書き込み
- 長さ（2バイト）：ペイロードサイズ（ビッグエンディアン）
- ペイロード（Nバイト）：データ
- CRC（2バイト）：Modbus CRC-16（リトルエンディアン）

## テスト

```bash
go test ./... -v
```
