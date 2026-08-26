# LonWorks GatewayBridge SDK

TCP 게이트웨이 브리지를 통한 LonWorks(ANSI/CEA-709.1) 프로토콜.

## 하드웨어

| 게이트웨이 | 인터페이스 | 기본 IP |
|---------|-----------|-------------|
| Echelon U60 | USB FT-10 네트워크 인터페이스 | 호스트 할당 |
| Echelon U70 | USB TP/XF-1250 인터페이스 | 호스트 할당 |
| Loytec L-IP | LonWorks/IP 라우터 | 192.168.0.90 |

## 배선

- FT-10(Free Topology): 극성 무관 트위스트 페어, 프리 토폴로지 최대 500m
- TP/XF-1250: 105옴 종단 저항이 있는 버스 토폴로지
- TCP 브리지용 L-IP 라우터의 이더넷

## 게이트웨이 IP 구성

```go
driver, codec, err := lonworks.NewGatewayDriver("192.168.0.90:2009")
handler, err := lonworks.ReadyHandler("192.168.0.90:2009")
```
