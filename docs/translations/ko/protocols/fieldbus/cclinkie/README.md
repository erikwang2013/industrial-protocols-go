# CC-Link IE Field GatewayBridge SDK

TCP 게이트웨이 브리지를 통한 CC-Link IE Field 프로토콜.

## 하드웨어

| 게이트웨이 | 인터페이스 | 기본 IP |
|---------|-----------|-------------|
| Mitsubishi MELSEC IQ-R | CC-Link IE Field 마스터 | 192.168.3.40 |
| Mitsubishi RJ71GF11-T2 | CC-Link IE Field 모듈 | 호스트 할당 |
| HMS Anybus CC-Link IE | 임베디드 게이트웨이 | 192.168.0.52 |

## 배선

- CC-Link IE Field용 RJ45 이더넷(1Gbps 링 또는 스타 토폴로지)
- 관리 포트는 별도 네트워크에 구성
- 게이트웨이가 필드 네트워크를 TCP 브리지에 연결

## 게이트웨이 IP 구성

```go
driver, codec, err := cclinkie.NewGatewayDriver("192.168.3.40:2005")
handler, err := cclinkie.ReadyHandler("192.168.3.40:2005")
```
