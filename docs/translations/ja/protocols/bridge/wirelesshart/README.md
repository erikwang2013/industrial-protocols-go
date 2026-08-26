# WirelessHART CmdBridge SDK

WirelessHART（IEC 62591）は、HARTプロトコルに基づくワイヤレス工業ネットワーク規格です。このパッケージは、`emerson_1410_cli`（Emerson 1410/1420 Wireless Gateway）コマンドラインユーティリティ経由のWirelessHARTコーデックを提供します。

## CLIツール

Emerson 1410/1420 Wireless Gateway CLIユーティリティを使用します。

### インストール

```bash
# Install Emerson Wireless Gateway software and tools
# Refer to Emerson documentation for 1410/1420 Gateway setup
```

確認：`emerson_1410_cli --help`

## 使用方法

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/wirelesshart"
)

func main() {
    p := wirelesshart.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "TT101",
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## ドライバ

```go
b, c, err := wirelesshart.NewCmdDriver("/usr/bin/emerson_1410_cli")
```

## サポートされる機能

| 機能 | 説明                |
|----------|----------------------------|
| `read`   | デバイスパラメータを読み取る      |
| `write`  | デバイスパラメータを書き込む     |
| `scan`   | ワイヤレスデバイスをスキャン  |

## テスト

```bash
go test ./... -v
```
