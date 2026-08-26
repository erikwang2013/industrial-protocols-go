# ControlNet CmdBridge SDK

ControlNetは、高速で時間厳守のデータ交換のためにAllen-Bradley（Rockwell Automation）が開発したリアルタイム産業ネットワークプロトコルです。このパッケージは、`1784-pcic-cli` コマンドラインユーティリティ経由のControlNetコーデックを提供します。

## CLIツール

1784-PCIC ControlNetインターフェースカードCLIユーティリティを使用します。

### インストール

```bash
# Install Rockwell 1784-PCIC driver and tools
# Refer to Rockwell Automation documentation for RSLinx Classic SDK
```

確認：`1784-pcic-cli --help`

## 使用方法

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/controlnet"
)

func main() {
    p := controlnet.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x10",
        Count:    8,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## ドライバ

```go
b, c, err := controlnet.NewCmdDriver("/usr/bin/1784-pcic-cli")
```

## サポートされる機能

| 機能 | 説明               |
|----------|---------------------------|
| `read`   | ControlNetノードから読み取る |
| `write`  | ControlNetノードに書き込む  |
| `status` | PCICカードの状態を照会    |

## テスト

```bash
go test ./... -v
```
