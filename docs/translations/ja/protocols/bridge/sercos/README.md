# SERCOS III CmdBridge SDK

SERCOS III（SErial Real-time COmmunication System）は、モーションコントロール用のデジタルインターフェースで、光ファイバまたは銅線でリングトポロジを使用します。このパッケージは、`netx_cli` コマンドラインユーティリティ経由のSERCOS IIIコーデックを提供します。

## CLIツール

Hilscher netX SERCOS III CLIユーティリティを使用します。

### インストール

```bash
# Install Hilscher netX driver and tools
# See https://www.hilscher.com/ for netX drivers
sudo apt-get install netx-driver
```

確認：`netx_cli --help`

## 使用方法

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/sercos"
)

func main() {
    p := sercos.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "S-0-51",
        Count:    2,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## ドライバ

```go
b, c, err := sercos.NewCmdDriver("/usr/bin/netx_cli")
```

## サポートされる機能

| 機能 | 説明                 |
|----------|-----------------------------|
| `read`   | SERCOS IDN/Sパラメータを読み取る |
| `write`  | SERCOS IDN/Sパラメータを書き込む|
| `phase`  | 通信フェーズを設定     |

## テスト

```bash
go test ./... -v
```
