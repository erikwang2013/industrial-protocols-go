# CANopen プロトコルSDK

CANopenは、組み込み制御システム向けのCANベースの上位層プロトコルです。このパッケージは、CANopenコーデックとSocketCANドライバを提供します。

## プロトコル概要

CANopenは標準の11ビットCAN識別子を使用し、以下の定義済み接続セットを持ちます：

| 機能    | CAN ID              | 説明                   |
|------------|---------------------|-------------------------------|
| NMT        | 0x000               | ネットワーク管理            |
| SYNC       | 0x080               | 同期メッセージ       |
| SDO (tx)   | 0x580 + NodeID      | Service Data Object (サーバー)  |
| SDO (rx)   | 0x600 + NodeID      | Service Data Object (クライアント)  |
| PDO1 (tx)  | 0x180 + NodeID      | Process Data Object 1         |
| Heartbeat  | 0x700 + NodeID      | Heartbeat / Bootup            |

## ハードウェア要件

- SocketCAN対応の**Linux**システム（`CONFIG_CAN`有効）
- CANインターフェース（例：`can0`、仮想CAN用の`vcan0`）
- CAN対応トランシーバーハードウェア（例：MCP2515、SJA1000、またはUSB-CANアダプタ）

### 仮想CANインターフェースのセットアップ（テスト用）

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## 使用方法

```go
package main

import (
    "fmt"
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/canopen"
)

func main() {
    p := canopen.New()
    codec, _ := p.NewCodec("can")

    // Read SDO object dictionary entry
    req := &kernel.Request{
        Function: "sdo_read",
        Metadata: map[string]any{
            "index": float64(0x1000), // Device type
            "sub":   float64(0),
        },
    }
    raw, _ := codec.Encode(req)
    fmt.Printf("SDO read frame: %X\n", raw)

    // NMT start remote node
    req2 := &kernel.Request{Function: "nmt_start"}
    raw2, _ := codec.Encode(req2)
    fmt.Printf("NMT start frame: %X\n", raw2)
}
```

## サポートされる機能

| 機能     | 説明                          |
|-------------|--------------------------------------|
| `sdo_read`  | オブジェクトディクショナリエントリを読み取る         |
| `sdo_write` | オブジェクトディクショナリエントリを書き込む        |
| `nmt_start` | リモートノードを起動（NMT）              |
| `nmt_stop`  | リモートノードを停止（NMT）               |
| `nmt_reset` | リモートノードをリセット（NMT）              |
| `heartbeat` | Heartbeat / bootupメッセージを送信      |

## テスト

```bash
go test ./... -v
```

注：SocketCANドライバのテストには、CANハードウェアまたは仮想CANインターフェースを持つLinuxシステムが必要です。
