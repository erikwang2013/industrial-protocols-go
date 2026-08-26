# CompactPCI プロトコルSDK

Linux sysfs経由のCompactPCIバスアクセス用に `kernel.Protocol` を実装します。

## プロトコル

- **名前**: `cpci`
- **バリアント**: `cpci`
- **デフォルトポート**: 0（メモリマップドバス）
- **トランスポート**: sysfs `/sys/bus/pci/devices/<BDF>/config`

CompactPCI（PICMG 2.0）は、従来のPCIと同じ電気的・ソフトウェアインターフェースを使用します。コーデックは、アプリケーションとPCIコンフィグ空間の間で生のバイトをsysfs経由で透過的に受け渡します。

## カーネル要件

以下のカーネルモジュールがロードされている必要があります：

- `pcieport` -- PCI Expressポートドライバ（ハイブリッドCPCIeシステム用）
- `pci_sysfs` -- sysfs PCIインターフェース（ほとんどのカーネルに組み込み）
- `cpci_hotplug` -- CompactPCIホットプラグコントローラ（オプション、ホットスワップ用）

sysfsファイルシステムは `/sys` にマウントされている必要があります。これはすべての近代的なLinuxディストリビューションでデフォルトです。

必要な権限：
- コンフィグ空間へのアクセスにはroot権限または `CAP_SYS_ADMIN`
- コンフィグファイルはroot:root所有、モード0600

## ドライバ

```go
d, err := cpci.NewCPCIDriver("0000:02:00.0")
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Use tr.Read/tr.Write for raw CPCI config space access
```

## CompactPCI と PCI の違い

CompactPCIは標準のPCIバス列挙とコンフィグ空間を使用します。デスクトップPCIとの主な違い：

- **3U/6Uフォームファクタ**: ピン&ソケットコネクタを備えたEurocardメカニクス
- **バス番号付け**: 各CPCIシャーシセグメントが独自のPCIバス番号を取得
- **ホットスワップ**: PICMG 2.1ホットスワップは標準のPCIホットプラグモデルを使用
- **システムスロット**: バス0、デバイス0がシステムスロットコントローラ

## BDF形式

バスアドレスはBDF（Bus:Device.Function）文字列です：
- `0000:02:00.0` -- ドメイン0000、バス02、デバイス00、ファンクション0
- `0000:02:08.0` -- ドメイン0000、バス02、デバイス08、ファンクション0
