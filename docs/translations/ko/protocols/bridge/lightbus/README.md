# Lightbus GatewayBridge SDK

TCP 게이트웨이 브리지를 통한 Beckhoff Lightbus 광섬유 프로토콜.

## 하드웨어

| 게이트웨이 | 인터페이스 | 기본 IP |
|---------|-----------|-------------|
| Beckhoff FC2001 | Lightbus PCI 카드 | 호스트 할당 |
| Beckhoff BK2000 | Lightbus 버스 커플러 | 192.168.0.80 |
| Beckhoff FC9001 | Lightbus 이더넷 어댑터 | 192.168.0.81 |

## 배선

- 플라스틱 광섬유(POF) 링 토폴로지
- FC2001/FC9001이 이더넷을 통해 링을 TCP 브리지에 연결
- 각 장치에는 TX 및 RX 광섬유 커넥터가 있음

## 게이트웨이 IP 구성

```go
driver, codec, err := lightbus.NewGatewayDriver("192.168.0.80:2008")
handler, err := lightbus.ReadyHandler("192.168.0.80:2008")
```
