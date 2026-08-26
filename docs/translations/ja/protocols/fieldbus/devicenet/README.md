# DeviceNet プロトコルSDK

DeviceNetは、ファクトリーオートメーション向けのCANベースの産業ネットワークプロトコルです。このパッケージは、SocketCANとTCPゲートウェイの両方のドライバを持つDeviceNetコーデックを提供します。

## プロトコル概要

DeviceNetはCAN上でCommon Industrial Protocol (CIP) を使用します。定義済みのMaster/Slave接続セットは、ポーリングI/Oと明示的メッセージングにGroup 2メッセージ（CAN ID 0x400 + NodeID）を使用します。

## ハードウェア要件

### CANモード（SocketCAN）
- SocketCAN対応の**Linux**システム（`CONFIG_CAN`有効）
- CANインターフェース（例：`can0`、`vcan0`）
- DeviceNet対応CANハードウェア（例：Anybus Communicator、HMS IXXAT）

### ゲートウェイモード
- DeviceNetゲートウェイへのTCP/IP接続
- プレーンテキストコマンドプロトコルをサポートするゲートウェイ（例：HMS Anybus、Hilscher netX）

### 仮想CANインターフェースのセットアップ（テスト用）

```bash
sudo modprobe can
sudo modprobe can_raw
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

## 使用方法

### CANモード

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/devicenet"
)

func main() {
    p := devicenet.New()
    codec, _ := p.NewCodec("can")

    // Poll request
    req := &kernel.Request{
        Function: "poll",
        Data:     []byte{0x01, 0x00},
    }
    raw, _ := codec.Encode(req)

    // Open connection
    req2 := &kernel.Request{Function: "open"}
    raw2, _ := codec.Encode(req2)
    _ = raw
    _ = raw2
}
```

### ゲートウェイモード

```go
p := devicenet.New()
codec, _ := p.NewCodec("gateway")

req := &kernel.Request{
    Function: "poll",
    Data:     []byte{0xAB, 0xCD},
}
raw, _ := codec.Encode(req)
// raw will be: "poll abcd\n"
```

## サポートされる機能

| 機能 | 説明 |
|----------|--------------------------------|
| `open`   | 明示的接続を開く |
| `poll`   | I/Oデータをポーリング（Group 2） |

## テスト

```bash
go test ./... -v
```

注：SocketCANドライバのテストはLinuxが必要です。ゲートウェイドライバのテストは任意のOSで実行できます。
