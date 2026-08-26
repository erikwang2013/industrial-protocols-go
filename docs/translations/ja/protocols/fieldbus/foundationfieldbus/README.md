# Foundation Fieldbus GatewayBridge SDK

Foundation Fieldbus H1/HSEをTCPゲートウェイブリッジ経由で提供します。

## ハードウェア

| ゲートウェイ | インターフェース | デフォルトIP |
|---------|-----------|-------------|
| NI USB-8486 | USB H1インターフェース | host割り当て |
| Softing FFusb | USB H1インターフェース | host割り当て |
| P+F HD2-GTR-4PA | H1からEthernetへのゲートウェイ | 192.168.0.20 |

## 配線

- H1トランク（ツイストペア、シールド付き）で両端にターミネータ
- H1用フィールドバス電源コンディショナ（24 VDC、セグメントあたり350-500 mA）
- ゲートウェイがEthernet経由でH1セグメントをTCPブリッジに接続

## ゲートウェイIP設定

```go
driver, codec, err := foundationfieldbus.NewGatewayDriver("192.168.0.20:2002")
handler, err := foundationfieldbus.ReadyHandler("192.168.0.20:2002")
```
