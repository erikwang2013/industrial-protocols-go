# CC-Link IE Field GatewayBridge SDK

CC-Link IE FieldプロトコルをTCPゲートウェイブリッジ経由で提供します。

## ハードウェア

| ゲートウェイ | インターフェース | デフォルトIP |
|---------|-----------|-------------|
| Mitsubishi MELSEC IQ-R | CC-Link IE Field マスター | 192.168.3.40 |
| Mitsubishi RJ71GF11-T2 | CC-Link IE Field モジュール | host割り当て |
| HMS Anybus CC-Link IE | 組み込みゲートウェイ | 192.168.0.52 |

## 配線

- CC-Link IE Field用のRJ45 Ethernet（1Gbpsリングまたはスター構成）
- 管理ポートは別ネットワークに接続
- ゲートウェイがフィールドネットワークをTCPブリッジに接続

## ゲートウェイIP設定

```go
driver, codec, err := cclinkie.NewGatewayDriver("192.168.3.40:2005")
handler, err := cclinkie.ReadyHandler("192.168.3.40:2005")
```
