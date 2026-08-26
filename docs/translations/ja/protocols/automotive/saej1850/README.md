# SAE J1850 CAN プロトコルSDK

SAE J1850は、オンボード診断（OBD-II）に使用される車両通信規格です。このパッケージは、CANバス（ISO 15765-4 / CAN TP）上のJ1850コーデックを提供します。

## プロトコル概要

CAN上のSAE J1850 OBD-IIは、次の構造を持つ29ビット拡張CAN識別子を使用します：

### CAN ID形式（29ビット）

| ビット       | フィールド    | 説明                        |
|-----------|----------|------------------------------------|
| 28-26     | Priority | メッセージ優先度（0-7、デフォルト6）  |
| 25        | Ext ID   | 拡張フレームでは常に1       |
| 24-16     | PF       | Parameter Format（ヘッダー）          |
| 15-8      | PS       | Parameter Specific（送信先/送信元） |
| 7-0       | SA       | 送信元アドレス                     |

### 標準OBD-II CAN ID

| タイプ                | CAN ID (hex)    | 説明                    |
|--------------------|-----------------|--------------------------------|
| Physical Request   | 0x18DAxxF1      | 特定ECUへのリクエスト（xx=アドレス） |
| Physical Response  | 0x18DAF1xx      | ECUからのレスポンス（xx=アドレス）   |
| Functional Request | 0x18DB33F1      | 全ECUへのブロードキャスト         |

### ISO 15765-2 フレーム形式

単一フレーム：バイト0の上位ニブル = データ長（0-7）、下位ニブル + 残りのバイト = 診断データ。

## ハードウェア要件

- SocketCAN対応の**Linux**システム
- J1850-CAN OBD-IIアダプタ（例：ELM327互換USB-CAN、OBDLink SX、Macchina M2）
- OBD-IIコネクタ付き車両（1996年以降のほとんどの車両）

### CANインターフェースのセットアップ

```bash
sudo modprobe can
sudo modprobe can_raw
sudo ip link set can0 type can bitrate 500000
sudo ip link set up can0
```

## 使用方法

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/saej1850"
)

func main() {
    p := saej1850.New()
    codec, _ := p.NewCodec("can")

    // Mode $01 PID $0C: Engine RPM
    req := &kernel.Request{
        Function: "mode01",
        Metadata: map[string]any{"pid": float64(0x0C)},
    }
    raw, _ := codec.Encode(req)
    _ = raw

    // Mode $03: Request emission-related DTCs
    req2 := &kernel.Request{Function: "mode03"}
    raw2, _ := codec.Encode(req2)
    _ = raw2
}
```

## サポートされる機能

| 機能        | OBD-IIモード | 説明                         |
|----------------|-------------|-------------------------------------|
| `mode01`       | $01         | 現在のパワートレインデータをリクエスト     |
| `mode03`       | $03         | 排出関連DTCをリクエスト       |
| `mode0A`       | $0A         | パーマネントDTCをリクエスト              |
| `diag_request` | カスタム      | 汎用診断リクエスト          |
| `diag_response`| -           | 診断レスポンスフレーム           |
| `broadcast`    | -           | 全ECUへの機能ブロードキャスト    |

## テスト

```bash
go test ./... -v
```
