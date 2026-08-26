# AS-Interface GatewayBridge SDK

AS-Interface（ASi）プロトコルをTCPゲートウェイブリッジ経由で提供します。

## ハードウェア

| ゲートウェイ | インターフェース | デフォルトIP |
|---------|-----------|-------------|
| Bihl+Wiedemann BWU2044 | ASiマスター（Ethernet） | 192.168.0.30 |
| P+F VBG-ENX-K20-DMD-EV | ASiゲートウェイ | 192.168.0.31 |
| ifm AC1375 | ASi ControllerE | 192.168.0.32 |

## 配線

- ゲートウェイからスレーブへ黄色のASiケーブル（電源 + データ）
- 黒色の補助電源ケーブル（アクチュエータ用24 VDC）はオプション
- TCPブリッジ用にゲートウェイへEthernetを接続

## ゲートウェイIP設定

```go
driver, codec, err := asinterface.NewGatewayDriver("192.168.0.30:2003")
handler, err := asinterface.ReadyHandler("192.168.0.30:2003")
```
