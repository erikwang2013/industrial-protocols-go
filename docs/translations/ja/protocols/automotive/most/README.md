# MOST シリアルプロトコルSDK

MOST（Media Oriented Systems Transport）は、主に自動車のインフォテインメントシステムで使用される高速マルチメディアネットワーク技術です。このパッケージは、ATコマンドインターフェースを使用するシリアルアダプタ経由のMOSTコーデックを提供します。

## プロトコル概要

MOSTは光ファイバ物理層上で同期シリアル通信を使用します。この実装は、115200ボーでATコマンドインターフェースを公開するシリアルアダプタ経由で接続します。

### ATコマンド

| コマンド            | 説明                  |
|--------------------|------------------------------|
| `AT+READ=<addr>,<len>`  | アドレスからバイトを読み取る |
| `AT+WRITE=<addr>,<hex>` | アドレスに16進バイトを書き込む |
| `AT+STATUS`             | リング/ネットワーク状態を照会 |

### レスポンス形式

- `+OK:<hex_data>` -- 成功レスポンス
- `+ERR:<code>` -- エラーレスポンス

## ハードウェア要件

- 適切な終端を持つMOST光ファイバネットワーク
- MOST-シリアルアダプタ（例：MOST150 USBアダプタ）
- 115200ボー、8N1のシリアルポート

## 使用方法

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/automotive/most"
)

func main() {
    p := most.New()
    codec, _ := p.NewCodec("serial")

    // Read 4 bytes from address 0x0100
    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "AT+READ=0x0100,4\r\n"
    _ = raw

    // Write data to address 0x0200
    req2 := &kernel.Request{
        Function: "write",
        Address:  "0x0200",
        Data:     []byte{0xDE, 0xAD, 0xBE, 0xEF},
    }
    raw2, _ := codec.Encode(req2)
    // raw2 = "AT+WRITE=0x0200,DEADBEEF\r\n"
    _ = raw2
}
```

## ドライバ

```go
b, c, err := most.NewSerialDriver("/dev/ttyUSB0")
```

## サポートされる機能

| 機能 | 説明                 |
|----------|-----------------------------|
| `read`   | MOSTアドレスから読み取る    |
| `write`  | MOSTアドレスにデータを書き込む |
| `status` | リング/ネットワーク状態を照会   |

## テスト

```bash
go test ./... -v
```
