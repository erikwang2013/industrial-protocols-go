# CompactPCI 프로토콜 SDK

Linux sysfs를 통한 CompactPCI 버스 접근을 위해 `kernel.Protocol`을 구현합니다.

## 프로토콜

- **Name**: `cpci`
- **Variants**: `cpci`
- **Default Port**: 0(메모리 매핑 버스)
- **Transport**: sysfs `/sys/bus/pci/devices/<BDF>/config`

CompactPCI(PICMG 2.0)은 기존 PCI와 동일한 전기적·소프트웨어 인터페이스를
사용합니다. 코덱은 sysfs를 통해 애플리케이션과 PCI 구성 공간 사이의
원시 바이트를 그대로 전달합니다.

## 커널 요구 사항

다음 커널 모듈이 로드되어 있어야 합니다:

- `pcieport` -- PCI Express 포트 드라이버(하이브리드 CPCIe 시스템용)
- `pci_sysfs` -- sysfs PCI 인터페이스(대부분의 커널에 내장)
- `cpci_hotplug` -- CompactPCI 핫플러그 컨트롤러(선택 사항, 핫 스왑용)

sysfs 파일시스템은 `/sys`에 마운트되어 있어야 합니다. 이는 모든 최신
Linux 배포판의 기본 설정입니다.

필요한 권한:
- 구성 공간 접근을 위한 루트 권한 또는 `CAP_SYS_ADMIN`
- 구성 파일은 root:root 소유이며 모드 0600

## 드라이버

```go
d, err := cpci.NewCPCIDriver("0000:02:00.0")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// 원시 CPCI 구성 공간 접근에는 tr.Read/tr.Write 사용
```

## CompactPCI vs PCI

CompactPCI은 표준 PCI 버스 열거 및 구성 공간을 사용합니다.
데스크톱 PCI와의 주요 차이점:

- **3U/6U 폼 팩터**: 핀-소켓 커넥터를 사용하는 유로카드(Eurocard) 구조
- **버스 번호 지정**: 각 CPCI 섀시 세그먼트에 자체 PCI 버스 번호 부여
- **핫 스왑**: PICMG 2.1 핫 스왑은 표준 PCI 핫플러그 모델을 사용
- **시스템 슬롯**: Bus 0, device 0이 시스템 슬롯 컨트롤러

## BDF 형식

버스 주소는 BDF(Bus:Device.Function) 문자열입니다:
- `0000:02:00.0` -- domain 0000, bus 02, device 00, function 0
- `0000:02:08.0` -- domain 0000, bus 02, device 08, function 0
