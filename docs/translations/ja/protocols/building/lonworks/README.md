# LonWorks GatewayBridge SDK

LonWorks（ANSI/CEA-709.1）プロトコルをTCPゲートウェイブリッジ経由で提供します。

## ハードウェア

| ゲートウェイ | インターフェース | デフォルトIP |
|---------|-----------|-------------|
| Echelon U60 | USB FT-10ネットワークインターフェース | host割り当て |
| Echelon U70 | USB TP/XF-1250インターフェース | host割り当て |
| Loytec L-IP | LonWorks/IPルータ | 192.168.0.90 |

## 配線

- FT-10（Free Topology）：極性不要のツイストペア、フリートポロジ最大500 m
- TP/XF-1250：105オーム終端のバストポロジ
- TCPブリッジ用にL-IPルータへEthernetを接続

## ゲートウェイIP設定

```go
driver, codec, err := lonworks.NewGatewayDriver("192.168.0.90:2009")
handler, err := lonworks.ReadyHandler("192.168.0.90:2009")
```
