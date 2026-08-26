# Interbus GatewayBridge SDK

InterbusプロトコルをTCPゲートウェイブリッジ経由で提供します。

## ハードウェア

| ゲートウェイ | インターフェース | デフォルトIP |
|---------|-----------|-------------|
| Phoenix Contact IBS S5 DSC/I-T | Interbusマスターコントローラ | 192.168.0.60 |
| Phoenix Contact IBS PCI SC/I-T | PCI Interbusマスター | host割り当て |
| HMS Anybus Interbus | 組み込みゲートウェイ | 192.168.0.61 |

## 配線

- ゲートウェイの9ピンD-SUBをInterbusリモートバス（入力/出力）へ接続
- 両端でシールドをFEに接続
- TCPブリッジ用にゲートウェイへEthernetを接続

## ゲートウェイIP設定

```go
driver, codec, err := interbus.NewGatewayDriver("192.168.0.60:2006")
handler, err := interbus.ReadyHandler("192.168.0.60:2006")
```
