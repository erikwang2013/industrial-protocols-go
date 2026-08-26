# PCI / PCIe プロトコルSDK

Linux sysfs経由のPCIおよびPCI Expressバスアクセス用に `kernel.Protocol` を実装します。

## プロトコル

- **名前**: `pci`
- **バリアント**: `pci`
- **デフォルトポート**: 0（メモリマップドバス）
- **トランスポート**: sysfs `/sys/bus/pci/devices/<BDF>/config`

コーデックは、アプリケーションとPCIコンフィグ空間の間で生のバイトを透過的に受け渡します。読み書きはsysfsのコンフィグファイル経由でデバイスの設定レジスタに直接アクセスします。

## カーネル要件

以下のカーネルモジュールがロードされている必要があります：

- `pcieport` -- PCI Expressポートドライバ
- `pci_sysfs` -- sysfs PCIインターフェース（ほとんどのカーネルに組み込み）

sysfsファイルシステムは `/sys` にマウントされている必要があります。これはすべての近代的なLinuxディストリビューションでデフォルトです。

必要な権限：
- コンフィグ空間へのアクセスにはroot権限または `CAP_SYS_ADMIN`
- コンフィグファイルはroot:root所有、モード0600

## ドライバ

```go
d, err := pci.NewPCIDriver("0000:00:1f.3")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Use tr.Read/tr.Write for raw PCI config space access
```

## BDF形式

バスアドレスはBDF（Bus:Device.Function）文字列です：
- `0000:00:1f.3` -- ドメイン0000、バス00、デバイス1f、ファンクション3
- `0000:01:00.0` -- ドメイン0000、バス01、デバイス00、ファンクション0
