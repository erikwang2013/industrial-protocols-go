# Interbus GatewayBridge SDK

TCP 게이트웨이 브리지를 통한 Interbus 프로토콜.

## 하드웨어

| 게이트웨이 | 인터페이스 | 기본 IP |
|---------|-----------|-------------|
| Phoenix Contact IBS S5 DSC/I-T | Interbus 마스터 컨트롤러 | 192.168.0.60 |
| Phoenix Contact IBS PCI SC/I-T | PCI Interbus 마스터 | 호스트 할당 |
| HMS Anybus Interbus | 임베디드 게이트웨이 | 192.168.0.61 |

## 배선

- 게이트웨이의 9핀 D-SUB를 Interbus 리모트 버스(입력/출력)에 연결
- 양쪽 끝의 실드를 FE에 연결
- TCP 브리지용 이더넷을 게이트웨이에 연결

## 게이트웨이 IP 구성

```go
driver, codec, err := interbus.NewGatewayDriver("192.168.0.60:2006")
handler, err := interbus.ReadyHandler("192.168.0.60:2006")
```
