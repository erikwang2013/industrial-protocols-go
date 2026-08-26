# IO-Link GatewayBridge SDK

IO-LinkプロトコルをTCPゲートウェイブリッジ経由で提供します。

## ハードウェア

| ゲートウェイ | インターフェース | デフォルトIP |
|---------|-----------|-------------|
| ifm AL1332 | IO-Linkマスター（EtherNet/IP） | 192.168.0.40 |
| Balluff BNI00AZ | IO-Linkマスター（PROFINET） | 192.168.0.41 |
| SICK SIG200 | IO-Linkマスター（Ethernet） | 192.168.0.42 |

## 配線

- IO-LinkポートごとにM12コネクタ（4ピン）：L+（茶）、L-（青）、C/Q（黒）、未使用（白）
- マスターとデバイス用の24 VDC電源
- TCPブリッジ用にマスターへEthernetを接続

## ゲートウェイIP設定

```go
driver, codec, err := iolink.NewGatewayDriver("192.168.0.40:2004")
handler, err := iolink.ReadyHandler("192.168.0.40:2004")
```
