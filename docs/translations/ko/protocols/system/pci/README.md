# PCI / PCIe 프로토콜 SDK

Linux sysfs를 통한 PCI 및 PCI Express 버스 접근을 위해 `kernel.Protocol`을 구현합니다.

## 프로토콜

- **Name**: `pci`
- **Variants**: `pci`
- **Default Port**: 0(메모리 매핑 버스)
- **Transport**: sysfs `/sys/bus/pci/devices/<BDF>/config`

코덱은 애플리케이션과 PCI 구성 공간 사이의 원시 바이트를 그대로
전달합니다. 읽기와 쓰기는 sysfs 구성 파일을 통해 장치의 구성
레지스터에 직접 수행됩니다.

## 커널 요구 사항

다음 커널 모듈이 로드되어 있어야 합니다:

- `pcieport` -- PCI Express 포트 드라이버
- `pci_sysfs` -- sysfs PCI 인터페이스(대부분의 커널에 내장)

sysfs 파일시스템은 `/sys`에 마운트되어 있어야 합니다. 이는 모든 최신
Linux 배포판의 기본 설정입니다.

필요한 권한:
- 구성 공간 접근을 위한 루트 권한 또는 `CAP_SYS_ADMIN`
- 구성 파일은 root:root 소유이며 모드 0600

## 드라이버

```go
d, err := pci.NewPCIDriver("0000:00:1f.3")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// 원시 PCI 구성 공간 접근에는 tr.Read/tr.Write 사용
```

## BDF 형식

버스 주소는 BDF(Bus:Device.Function) 문자열입니다:
- `0000:00:1f.3` -- domain 0000, bus 00, device 1f, function 3
- `0000:01:00.0` -- domain 0000, bus 01, device 00, function 0
