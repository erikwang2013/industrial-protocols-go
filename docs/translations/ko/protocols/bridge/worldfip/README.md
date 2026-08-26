# WorldFIP GatewayBridge SDK

TCP 게이트웨이 브리지를 통한 WorldFIP 프로토콜.

## 하드웨어

| 게이트웨이 | 인터페이스 | 기본 IP |
|---------|-----------|-------------|
| FIPIO Agent | WorldFIP 필드버스 에이전트 | 192.168.0.70 |
| FIP Gateway (Alstom) | WorldFIP to Ethernet | 192.168.0.71 |
| NI FIP-USB | USB WorldFIP 인터페이스 | 호스트 할당 |

## 배선

- 게이트웨이의 9핀 D-SUB를 WorldFIP 트렁크에 연결(FIP1 = Data+, FIP2 = Data-)
- 버스 양쪽 끝에 라인 종단 저항(120옴)
- TCP 브리지용 이더넷

## 게이트웨이 IP 구성

```go
driver, codec, err := worldfip.NewGatewayDriver("192.168.0.70:2007")
handler, err := worldfip.ReadyHandler("192.168.0.70:2007")
```
