# POWERLINK CmdBridge SDK

POWERLINK（Ethernet POWERLINK）は、産業オートメーション向けのリアルタイムEthernetプロトコルです。このパッケージは、`openPOWERLINK_demo` コマンドラインユーティリティ経由のPOWERLINKコーデックを提供します。

## CLIツール

openPOWERLINKスタックデモアプリケーションを使用します。

### インストール

```bash
git clone https://github.com/OpenAutomationTechnologies/openPOWERLINK.git
cd openPOWERLINK
cmake -B build && cmake --build build
sudo cmake --install build
```

確認：`openPOWERLINK_demo --help`

## 使用方法

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/powerlink"
)

func main() {
    p := powerlink.New()
    codec, _ := p.NewCodec("cmd")

    req := &kernel.Request{
        Function: "read",
        Address:  "0x2000",
        Count:    8,
    }
    raw, _ := codec.Encode(req)
    _ = raw
}
```

## ドライバ

```go
b, c, err := powerlink.NewCmdDriver("/usr/bin/openPOWERLINK_demo")
```

## サポートされる機能

| 機能 | 説明                   |
|----------|-------------------------------|
| `read`   | オブジェクトディクショナリエントリを読み取る  |
| `write`  | オブジェクトディクショナリエントリを書き込む |
| `status` | ノード/NMT状態を照会          |

## テスト

```bash
go test ./... -v
```
