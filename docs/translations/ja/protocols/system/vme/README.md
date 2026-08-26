# VME / VPX プロトコルSDK

Linux procfs経由のVMEbusおよびVPXバスアクセス用に `kernel.Protocol` を実装します。

## プロトコル

- **名前**: `vme`
- **バリアント**: `vme`
- **デフォルトポート**: 0（メモリマップドバス）
- **トランスポート**: procfs `/proc/vme/<slot>`

コーデックは、アプリケーションとVMEアドレス空間の間で生のバイトを透過的に受け渡します。読み書きはVMEカーネルドライバが提供するprocfsインターフェース経由でバスに直接アクセスします。

## カーネル要件

以下のカーネルモジュールがロードされている必要があります：

- `vme_tsi148` -- Tundra TSI148 VMEブリッジドライバ（最も一般的）
  - その他対応：`vme_ca91cx42`（Universe II）、`vme_user`

procfsファイルシステムは `/proc` にマウントされている必要があります。

必要な権限：
- 読み書きで `/proc/vme/<slot>` を開くにはroot権限が必要
- デバイスノードはroot:root所有

## ドライバ

```go
d, err := vme.NewVMEDriver(0) // slot 0
if err != nil {
    log.Fatal(err)
}
defer d.Close()

tr := d.Transport()
// Use tr.Read/tr.Write for raw VME bus access
```

## VMEアドレッシングモード

コーデックはすべてのアドレス修飾子とアドレスバイトを透過的に受け渡します。アプリケーションはデータペイロードの前にアドレッシング情報を付加する必要があります：

- **A16**: 16ビット短I/Oアドレス空間
- **A24**: 24ビット標準アドレス空間
- **A32**: 32ビット拡張アドレス空間

## VPX互換性

VME互換のprocfsインターフェースを公開するVPX（VITA 46）システムは、このドライバを使用できます。スロット番号の付け方も同じです。
