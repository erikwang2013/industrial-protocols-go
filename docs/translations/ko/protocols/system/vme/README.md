# VME / VPX 프로토콜 SDK

Linux procfs를 통한 VMEbus 및 VPX 버스 접근을 위해 `kernel.Protocol`을 구현합니다.

## 프로토콜

- **Name**: `vme`
- **Variants**: `vme`
- **Default Port**: 0(메모리 매핑 버스)
- **Transport**: procfs `/proc/vme/<slot>`

코덱은 애플리케이션과 VME 주소 공간 사이의 원시 바이트를 그대로
전달합니다. 읽기와 쓰기는 VME 커널 드라이버가 제공하는 procfs
인터페이스를 통해 버스에 직접 수행됩니다.

## 커널 요구 사항

다음 커널 모듈이 로드되어 있어야 합니다:

- `vme_tsi148` -- Tundra TSI148 VME 브리지 드라이버(가장 일반적)
  - 추가 지원: `vme_ca91cx42`(Universe II), `vme_user`

procfs 파일시스템은 `/proc`에 마운트되어 있어야 합니다.

필요한 권한:
- `/proc/vme/<slot>`을 읽기-쓰기로 열려면 루트 권한 필요
- 장치 노드는 root:root 소유

## 드라이버

```go
d, err := vme.NewVMEDriver(0) // slot 0
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// 원시 VME 버스 접근에는 tr.Read/tr.Write 사용
```

## VME 어드레싱 모드

코덱은 모든 주소 수정자(Address Modifier) 및 주소 바이트를 투명하게
전달합니다. 애플리케이션은 데이터 페이로드 앞에 어드레싱 정보를
붙여야 합니다:

- **A16**: 16비트 단거리 I/O 주소 공간
- **A24**: 24비트 표준 주소 공간
- **A32**: 32비트 확장 주소 공간

## VPX 호환성

VME 호환 procfs 인터페이스를 노출하는 VPX(VITA 46) 시스템은
이 드라이버를 사용할 수 있습니다. 슬롯 번호 지정 방식은 동일합니다.
