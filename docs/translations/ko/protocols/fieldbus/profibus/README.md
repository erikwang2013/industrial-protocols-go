# PROFIBUS GatewayBridge SDK

TCP 게이트웨이 브리지를 통한 PROFIBUS DP/PA 프로토콜.

## 하드웨어

| 게이트웨이 | 인터페이스 | 기본 IP |
|---------|-----------|-------------|
| Anybus Communicator | PROFIBUS DP-V1 슬레이브 | 192.168.0.50 |
| Siemens CP 5611 proxy | PCI/PCIe PROFIBUS 마스터 | 호스트 할당 |
| HMS Fieldbus Gateway | Anybus NP40 | 192.168.0.51 |

## 배선

- 게이트웨이의 DB9 암 커넥터를 PROFIBUS 네트워크에 연결(A-라인 녹색, B-라인 빨간색)
- 네트워크 양쪽 끝에 종단 저항 ON
- 게이트웨이 관리 포트에 이더넷 연결

## 게이트웨이 IP 구성

```go
driver, codec, err := profibus.NewGatewayDriver("192.168.0.50:2000")
handler, err := profibus.ReadyHandler("192.168.0.50:2000")
```
