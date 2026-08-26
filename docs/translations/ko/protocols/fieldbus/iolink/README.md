# IO-Link GatewayBridge SDK

TCP 게이트웨이 브리지를 통한 IO-Link 프로토콜.

## 하드웨어

| 게이트웨이 | 인터페이스 | 기본 IP |
|---------|-----------|-------------|
| ifm AL1332 | IO-Link Master(EtherNet/IP) | 192.168.0.40 |
| Balluff BNI00AZ | IO-Link Master(PROFINET) | 192.168.0.41 |
| SICK SIG200 | IO-Link Master(이더넷) | 192.168.0.42 |

## 배선

- 각 IO-Link 포트용 M12 커넥터(4핀): L+(갈색), L-(파란색), C/Q(검은색), 미사용(흰색)
- 마스터 및 장치용 24VDC 전원 공급
- TCP 브리지용 마스터의 이더넷

## 게이트웨이 IP 구성

```go
driver, codec, err := iolink.NewGatewayDriver("192.168.0.40:2004")
handler, err := iolink.ReadyHandler("192.168.0.40:2004")
```
