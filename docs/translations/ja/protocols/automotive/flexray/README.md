# FlexRay CAN プロトコルSDK

FlexRayは、高速で決定論的な自動車通信プロトコルです。このパッケージは、CRC-16/XMODEM整合性チェックによるサイクルベースのフレーム処理をサポートする、CANバス上のFlexRayコーデックを提供します。

## プロトコル概要

FlexRayは、繰り返し通信サイクルによる時分割多重アクセス（TDMA）方式を使用します。各サイクルは静的セグメントと動的セグメントで構成されます。このコーデックはFlexRayフレームを拡張29ビットCANフレームにマッピングします。

### ワイヤ形式

FlexRayサイクルペイロード：
- **ヘッダー**（2バイト）：サイクル番号（リトルエンディアン）
- **ステータス**（1バイト）：ビット7=PPI（Payload Preamble Indicator）、ビット6=NFI、ビット5=SYF、ビット4=SUF
- **データ**（Nバイト）：ペイロード（最大254バイト）
- **CRC**（2バイト）：ヘッダー+ステータス+データに対するCRC-16/XMODEM（リトルエンディアン）

CAN IDエンコーディング（29ビット拡張）：
- ビット28-24：メッセージタイプ（0x01=フレーム、0x02=ステータス）
- ビット23-10：予約
- ビット15-10：スロットID（6ビット）
- ビット9-0：サイクル番号（10ビット）

## ハードウェア要件

- SocketCAN対応の**Linux**システム
- Vector VN7600/VN7640またはBosch FlexRay-CANアダプタ
- 適切な終端（2.5Vバイアス）を持つFlexRayネットワーク

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
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/flexray"
)

func main() {
    p := flexray.New()
    codec, _ := p.NewCodec("can")

    // Send a FlexRay frame in cycle 5 with PPI set
    req := &kernel.Request{
        Function: "frame",
        Data:     []byte{0x42, 0x01},
        Metadata: map[string]any{
            "cycle": float64(5),
            "ppi":   true,
        },
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## サポートされる機能

| 機能 | 説明                              |
|----------|------------------------------------------|
| `frame`  | FlexRayフレームペイロードを送信               |
| `status` | スロット構成を照会（スロット、サイクル）   |

## テスト

```bash
go test ./... -v
```
