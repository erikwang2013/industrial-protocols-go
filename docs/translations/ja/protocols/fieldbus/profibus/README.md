# PROFIBUS GatewayBridge SDK

PROFIBUS DP/PAプロトコルをTCPゲートウェイブリッジ経由で提供します。

## ハードウェア

| ゲートウェイ | インターフェース | デフォルトIP |
|---------|-----------|-------------|
| Anybus Communicator | PROFIBUS DP-V1 スレーブ | 192.168.0.50 |
| Siemens CP 5611 proxy | PCI/PCIe PROFIBUS マスター | host割り当て |
| HMS Fieldbus Gateway | Anybus NP40 | 192.168.0.51 |

## 配線

- ゲートウェイのDB9メスをPROFIBUSネットワークへ接続（A線は緑、B線は赤）
- ネットワーク両端で終端抵抗をONにする
- ゲートウェイの管理ポートにEthernetを接続

## ゲートウェイIP設定

```go
driver, codec, err := profibus.NewGatewayDriver("192.168.0.50:2000")
handler, err := profibus.ReadyHandler("192.168.0.50:2000")
```
