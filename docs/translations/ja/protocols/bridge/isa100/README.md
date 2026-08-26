# ISA100.11a CmdBridge SDK

ISA100.11aは、プロセスオートメーション向けのワイヤレス工業ネットワーク規格です。このパッケージは、`yfgw410_cli`（Yokogawa YFGW410フィールドワイヤレスゲートウェイ）コマンドラインユーティリティ経由のISA100.11aコーデックを提供します。

## CLIツール

Yokogawa YFGW410 Field Wireless Gateway CLIを使用します。

### インストール

```bash
# Install Yokogawa YFGW410 gateway software and tools
# Refer to Yokogawa documentation for Field Wireless Gateway setup
```

確認：`yfgw410_cli --help`

## 使用方法

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/isa100"
)

func main() {
    p := isa100.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "DEV001",
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## ドライバ

```go
b, c, err := isa100.NewCmdDriver("/usr/bin/yfgw410_cli")
```

## サポートされる機能

| 機能 | 説明                 |
|----------|-----------------------------|
| `read`   | デバイス属性を読み取る       |
| `write`  | デバイス属性を書き込む      |
| `list`   | プロビジョニング済みデバイスを一覧表示    |

## テスト

```bash
go test ./... -v
```
