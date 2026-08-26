# Foundation Fieldbus GatewayBridge SDK

TCP 게이트웨이 브리지를 통한 Foundation Fieldbus H1/HSE.

## 하드웨어

| 게이트웨이 | 인터페이스 | 기본 IP |
|---------|-----------|-------------|
| NI USB-8486 | USB H1 인터페이스 | 호스트 할당 |
| Softing FFusb | USB H1 인터페이스 | 호스트 할당 |
| P+F HD2-GTR-4PA | H1 to Ethernet 게이트웨이 | 192.168.0.20 |

## 배선

- 양쪽 끝에 종단 저항이 있는 H1 트렁크(차폐 트위스트 페어)
- H1용 필드버스 전원 컨디셔너(24VDC, 세그먼트당 350-500mA)
- 게이트웨이가 이더넷을 통해 H1 세그먼트를 TCP 브리지에 연결

## 게이트웨이 IP 구성

```go
driver, codec, err := foundationfieldbus.NewGatewayDriver("192.168.0.20:2002")
handler, err := foundationfieldbus.ReadyHandler("192.168.0.20:2002")
```
