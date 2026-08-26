# AS-Interface GatewayBridge SDK

TCP 게이트웨이 브리지를 통한 AS-Interface(ASi) 프로토콜.

## 하드웨어

| 게이트웨이 | 인터페이스 | 기본 IP |
|---------|-----------|-------------|
| Bihl+Wiedemann BWU2044 | ASi Master(이더넷) | 192.168.0.30 |
| P+F VBG-ENX-K20-DMD-EV | ASi 게이트웨이 | 192.168.0.31 |
| ifm AC1375 | ASi ControllerE | 192.168.0.32 |

## 배선

- 게이트웨이에서 슬레이브까지 노란색 ASi 케이블(전원 + 데이터)
- 검은색 보조 전원 케이블(액추에이터용 24VDC) 선택 사항
- TCP 브리지용 게이트웨이의 이더넷

## 게이트웨이 IP 구성

```go
driver, codec, err := asinterface.NewGatewayDriver("192.168.0.30:2003")
handler, err := asinterface.ReadyHandler("192.168.0.30:2003")
```
