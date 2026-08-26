# Lightbus GatewayBridge SDK

Beckhoff Lightbus光ファイバプロトコルをTCPゲートウェイブリッジ経由で提供します。

## ハードウェア

| ゲートウェイ | インターフェース | デフォルトIP |
|---------|-----------|-------------|
| Beckhoff FC2001 | Lightbus PCIカード | host割り当て |
| Beckhoff BK2000 | Lightbusバスカプラ | 192.168.0.80 |
| Beckhoff FC9001 | Lightbus Ethernetアダプタ | 192.168.0.81 |

## 配線

- プラスチック光ファイバ（POF）リングトポロジ
- FC2001/FC9001がEthernet経由でリングをTCPブリッジに接続
- 各デバイスにTXおよびRX光ファイバコネクタ

## ゲートウェイIP設定

```go
driver, codec, err := lightbus.NewGatewayDriver("192.168.0.80:2008")
handler, err := lightbus.ReadyHandler("192.168.0.80:2008")
```
