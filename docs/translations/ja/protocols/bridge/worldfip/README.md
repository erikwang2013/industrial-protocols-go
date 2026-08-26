# WorldFIP GatewayBridge SDK

WorldFIPプロトコルをTCPゲートウェイブリッジ経由で提供します。

## ハードウェア

| ゲートウェイ | インターフェース | デフォルトIP |
|---------|-----------|-------------|
| FIPIO Agent | WorldFIPフィールドバスエージェント | 192.168.0.70 |
| FIP Gateway (Alstom) | WorldFIPからEthernetへ | 192.168.0.71 |
| NI FIP-USB | USB WorldFIPインターフェース | host割り当て |

## 配線

- ゲートウェイの9ピンD-SUBをWorldFIPトランクへ接続（FIP1 = Data+、FIP2 = Data-）
- バス両端にライオンターミネータ（120オーム）
- TCPブリッジ用にEthernetを接続

## ゲートウェイIP設定

```go
driver, codec, err := worldfip.NewGatewayDriver("192.168.0.70:2007")
handler, err := worldfip.ReadyHandler("192.168.0.70:2007")
```
