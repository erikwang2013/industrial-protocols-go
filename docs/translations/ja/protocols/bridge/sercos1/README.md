# SERCOS I/II CmdBridge SDK

SERCOS I/IIは、デジタルモーションコントロール用SERCOSインターフェースの旧型シリアル光ファイバ版です。このパッケージは、`sercos_cli` コマンドラインユーティリティ経由のSERCOS I/IIコーデックを提供します。

## CLIツール

SERCOS光ファイバインターフェースCLIユーティリティを使用します。

### インストール

```bash
# SERCOS interface card driver and tools
# Refer to vendor documentation for the specific SERCOS master card
```

確認：`sercos_cli --help`

## 使用方法

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos1"
)

func main() {
    p := sercos1.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x0100",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## ドライバ

```go
b, c, err := sercos1.NewCmdDriver("/usr/bin/sercos_cli")
```

## サポートされる機能

| 機能 | 説明              |
|----------|--------------------------|
| `read`   | SERCOS IDNを読み取る          |
| `write`  | SERCOS IDNを書き込む         |
| `status` | ドライブ状態を照会       |

## テスト

```bash
go test ./... -v
```
