# EtherCAT CmdBridge SDK

EtherCAT（Ethernet for Control Automation Technology）は、高性能な産業用イーサネットフィールドバスです。このパッケージは、`ethercat` コマンドラインユーティリティ経由のEtherCATコーデックを提供します。

## CLIツール

IgH EtherCAT Masterコマンドラインツールを使用します。

### インストール

```bash
# Debian/Ubuntu
sudo apt-get install ethercat-master

# From source
git clone https://gitlab.com/etherlab.org/ethercat.git
cd ethercat
make && sudo make install
```

確認：`ethercat slaves`

## 使用方法

```go
package main

import (
    "github.com/erikwang2013/industrial-protocols-go/kernel"
    "github.com/erikwang2013/industrial-protocols-go/protocols/bridge/ethercat"
)

func main() {
    p := ethercat.New()
    codec, _ := p.NewCodec("cmd")

    // Upload SDO from address 0x1000
    req := &kernel.Request{
        Function: "upload",
        Address:  "0x1000",
        Count:    4,
    }
    raw, _ := codec.Encode(req)
    // raw = "upload 0x1000 4\n"
    _ = raw
}
```

## ドライバ

```go
b, c, err := ethercat.NewCmdDriver("/usr/bin/ethercat")
```

## サポートされる機能

| 機能   | 説明            |
|------------|------------------------|
| `upload`   | アドレスからSDOを読み取る  |
| `download` | アドレスにSDOを書き込む   |
| `slaves`   | EtherCATスレーブを一覧表示   |

## テスト

```bash
go test ./... -v
```
