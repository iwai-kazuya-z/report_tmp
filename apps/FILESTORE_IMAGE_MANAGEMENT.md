# Cloud Filestore 画像管理の仕組み

## 概要

ドラピタでは、ユーザーアップロード画像を**Cloud Filestore（NFS）**に保存し、複数のアプリケーション間で共有しています。

## アーキテクチャ

```
┌─────────────────────────────────────────────────────────────────┐
│                        ユーザー                                  │
│                           │                                       │
│                           ▼                                       │
│  ┌─────────────────────────────────────────────┐                │
│  │  dorapita.com / cadm.dorapita.com (管理)    │                │
│  │  - AssetService::uploadFile()                │                │
│  │  - ファイルをuploadし、DBにメタデータ保存    │                │
│  └──────────────────┬──────────────────────────┘                │
│                     │                                             │
│                     ▼                                             │
│  ┌─────────────────────────────────────────────┐                │
│  │  UPFILE_PATH (resources/upfiles/)            │                │
│  │  = /var/www/dorapita.com/resources/upfiles   │                │
│  │                                               │                │
│  │  ファイル名: {asset_id} (数値のみ)           │                │
│  └──────────────────┬──────────────────────────┘                │
│                     │                                             │
│                     │ マウント（環境による違い）                  │
│                     │                                             │
│        ┌────────────┴────────────┐                               │
│        │                         │                               │
│        ▼                         ▼                               │
│  ┌─────────┐              ┌──────────────┐                      │
│  │ ローカル │              │  Cloud Run   │                      │
│  │開発環境  │              │  本番・開発   │                      │
│  │         │              │               │                      │
│  │ ./upfiles│              │ /mnt (Filestore)                    │
│  │ (Docker) │              │ - 10.227.113.2 (dev)                │
│  │         │              │ - 10.117.161.2 (prod)               │
│  └─────────┘              └──────────────┘                      │
│                                  │                                │
│                                  ▼                                │
│  ┌─────────────────────────────────────────────┐                │
│  │  img.dorapita.com (画像配信)                 │                │
│  │  - AssetsController::index()                 │                │
│  │  - /mnt/upfiles/{asset_id}から読み込み       │                │
│  │  - Vips\Imageでリサイズ・WebP変換            │                │
│  │  - 透かし（dora_sukashi.png）追加            │                │
│  └─────────────────────────────────────────────┘                │
│                           │                                       │
│                           ▼                                       │
│                     CDN/ユーザーへ配信                            │
└─────────────────────────────────────────────────────────────────┘
```

## リポジトリ内の画像ファイル

### .gitignoreで除外

```gitignore
# dorapita_code/.gitignore
upfiles/*
```

**実態**: リポジトリに画像ファイルは**含まれていない**。

**理由**:
- 画像はユーザーアップロードで動的に生成される
- Git管理すると肥大化する
- 環境ごとに異なる（開発/本番で別のデータ）

### ダミーデータの取得方法

開発環境構築時には、本番データのダンプを取得：

```bash
# README.mdに記載されている手順
gsutil -m cp -r gs://dorapita-infra-stg/assets/upfiles .
```

**保存先**: `dorapita_code/upfiles/`（Gitでは管理しない）

## ファイルアップロード処理

### dorapita.com/AssetService.php

```php
// ファイルアップロード
public function uploadFile(UploadedFile $file, ?string $existingAssetId = null): array
{
    // 1. バリデーション（サイズ、MIME、危険なコンテンツ検査）
    $validationResult = $this->validateImage($file);

    // 2. DBにメタデータ保存
    $asset = $this->Assets->newEntity([
        'name' => $file->getClientFilename(),
        'type_name' => $file->getClientMediaType()
    ]);
    $this->Assets->save($asset);

    // 3. 物理ファイル保存
    $save_path = $this->uploadPath . DS . $asset->id;  // UPFILE_PATH/{id}
    move_uploaded_file($tmpPath, $save_path);
}
```

**保存パス**: `resources/upfiles/{asset_id}`
- ファイル名: アセットID（整数）のみ
- 拡張子なし
- 例: `resources/upfiles/12345`

**セキュリティ**:
- ファイルサイズ制限: 2MB
- MIME Type検証: `image/jpeg`, `image/png`のみ
- 危険なコンテンツ検査: `<?php`, `<script>`, `eval(`等をチェック

## ファイル配信処理

### img.dorapita.com/AssetsController.php

```php
public function index($id = '', int $image_size = 720, $format = 'webp'): ?Response
{
    // 1. Filestoreから画像読み込み
    $url = $id ? "/mnt/upfiles/$id" : self::NO_IMAGE_FILE;
    $blob = @file_get_contents($url);

    // 2. Vipsで画像処理
    $image = Vips\Image::newFromBuffer($blob, '', ['access' => 'sequential']);

    // 3. 透かし追加（300px以上の画像）
    if ($image->width >= 300 || $image->height >= 300) {
        $watermak = Vips\Image::newFromFile(WWW_ROOT . DS .'img' . DS . 'dora_sukashi.png');
        $image = $image->composite($watermak, 'over', [...]);
    }

    // 4. リサイズ
    $image = $image->thumbnail_image($image_size, ['height' => $image_size]);

    // 5. WebP変換して返却
    return $this->response->withType($format)->withStringBody($image->writeToBuffer(".{$format}"));
}
```

**機能**:
- 動的リサイズ（デフォルト720px）
- WebP/JPEG変換
- 透かし自動追加
- CDNキャッシュ対応（Cache-Control: 1日）

## 環境別マウント設定

### ローカル開発環境（docker-compose.yml）

```yaml
volumes:
  - ./upfiles:/var/www/legacy.dorapita.com/app/upfiles
  - ./upfiles:/var/www/cadm.dorapita.com/upfiles
  - ./upfiles:/var/www/kanri.dorapita.com/resources/upfiles
```

**特徴**:
- ホストの`./upfiles`ディレクトリをマウント
- 複数サービスで共有
- ホットリロード有効

### 開発環境 Cloud Run（dorapita-core-dev）

#### dorapita-com, img-dorapita-com

```yaml
volumes:
- name: filestore
  nfs:
    path: /dorapita
    server: 10.227.113.2  # Filestore IP (dev)

volumeMounts:
- mountPath: /mnt
  name: filestore
```

**マウント先**: `/mnt/upfiles/`（NFSサブディレクトリ）

### 本番環境 Cloud Run（dorapita-core）

#### img-dorapita-com

```yaml
volumes:
- name: filestore
  nfs:
    path: /dorapita
    server: 10.117.161.2  # Filestore IP (prod)

volumeMounts:
- mountPath: /mnt
  name: filestore
```

**マウント先**: `/mnt/upfiles/`

### 本番環境 Compute Engine VM

VMに直接NFSマウント（推測）:

```bash
# /etc/fstab (推測)
10.117.161.2:/dorapita /mnt/upfiles nfs defaults 0 0
```

## Filestore構成

| 環境 | IP | パス | マウント先 |
|------|-----|------|-----------|
| 開発 | 10.227.113.2 | /dorapita | /mnt |
| 本番 | 10.117.161.2 | /dorapita | /mnt |

**共有構造**:

```
Filestore: /dorapita/
├── upfiles/           ← 画像ファイル
│   ├── 1
│   ├── 2
│   ├── 3
│   └── ...
└── (その他のファイル？)
```

## パス定義まとめ

| アプリ | 定義ファイル | UPFILE_PATH |
|--------|-------------|-------------|
| dorapita.com | config/paths.php:82 | `RESOURCES.'upfiles'` → `resources/upfiles/` |
| cadm.dorapita.com | (同上) | (同上) |
| kanri.dorapita.com | config/paths.php:99 | `RESOURCES.'upfiles'` → `resources/upfiles/` |
| legacy.dorapita.com | app/config/settings.php.default:3 | `ROOT.DS.APP_DIR.DS.'upfiles'` → `app/upfiles/` |
| img.dorapita.com | ハードコード | `/mnt/upfiles/` |

**注意**: img.dorapita.comのみ`/mnt/upfiles`を直接参照。他はアプリ内パスを使用し、マウントで吸収。

## データフロー

### アップロード時

```
1. ユーザー → dorapita.com/cadm/kanri (フォーム送信)
2. AssetService::uploadFile()
   - バリデーション
   - DBにメタデータ保存（Assetsテーブル）
   - ファイル保存: resources/upfiles/{asset_id}
3. ファイルはFilestoreに永続化（NFSマウント経由）
```

### 配信時

```
1. ユーザー → img.dorapita.com/{asset_id}/{size}/{format}
2. AssetsController::index()
   - /mnt/upfiles/{asset_id}から読み込み
   - Vipsで画像処理（リサイズ、透かし、WebP変換）
   - レスポンス返却（Cache-Control: 1日）
3. CDN/ブラウザキャッシュ
```

## Database Schema（推測）

### Assetsテーブル

```sql
CREATE TABLE assets (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),        -- 元のファイル名
    type_name VARCHAR(50),    -- MIME Type (image/jpeg等)
    created DATETIME,
    modified DATETIME
);
```

**特徴**:
- IDがそのままファイル名になる
- 拡張子情報なし（MIMEから判定）
- メタデータのみ保存（実ファイルはFilestore）

## セキュリティ考慮事項

### アップロード時

✅ **実装済み**:
- ファイルサイズ制限（2MB）
- MIME Type検証
- 危険なコンテンツ検査（PHPコード、スクリプト等）
- getimagesize()による画像検証

⚠️ **追加推奨**:
- ファイル名のサニタイズ（実装済み：IDのみ使用）
- ウイルススキャン
- レート制限（DoS対策）

### 配信時

✅ **実装済み**:
- Access-Control-Allow-Origin設定
- Cache-Control適切
- 存在しないファイルは404 + noimg.png

⚠️ **追加推奨**:
- CDN経由配信（オリジンシールド）
- 画像最適化自動化（既にVips使用中）

## 運用上の注意点

### ディスク使用量監視

```bash
# Filestoreディスク使用量確認
gcloud filestore instances describe <INSTANCE_NAME> \
  --project=dorapita-core \
  --location=asia-northeast1
```

### バックアップ

**現状**: 不明（Filestoreのバックアップポリシー未確認）

**推奨**:
```bash
# 定期バックアップ設定
gcloud filestore backups create <BACKUP_NAME> \
  --instance=<INSTANCE_NAME> \
  --location=asia-northeast1
```

### 画像削除

```php
// AssetService::deleteFile()
public function deleteFile($id): bool
{
    $file_path = $this->uploadPath . DS . $id;
    if (file_exists($file_path)) {
        unlink($file_path);  // 物理ファイル削除
    }
    return $this->Assets->delete($asset);  // DBから削除
}
```

**注意**: 論理削除ではなく物理削除。復元不可。

## トラブルシューティング

### 画像が表示されない

```bash
# 1. Filestoreマウント確認
kubectl exec -it <POD> -- df -h | grep /mnt

# 2. ファイル存在確認
kubectl exec -it <POD> -- ls -la /mnt/upfiles/{asset_id}

# 3. img.dorapita.comログ確認
gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=img-dorapita-com"
```

### アップロード失敗

```bash
# 1. ディスク容量確認
df -h /var/www/dorapita.com/resources/upfiles

# 2. パーミッション確認
ls -la /var/www/dorapita.com/resources/

# 3. アプリケーションログ
tail -f /var/www/dorapita.com/logs/error.log
```

## まとめ

| 項目 | 内容 |
|------|------|
| **リポジトリ内** | ❌ 画像ファイルなし（.gitignore除外） |
| **ローカル開発** | ./upfiles をDockerマウント |
| **本番・開発環境** | Cloud Filestore (NFS) を/mntにマウント |
| **アップロード先** | resources/upfiles/{asset_id} |
| **配信元** | /mnt/upfiles/{asset_id} (Filestore) |
| **画像処理** | Vips (リサイズ、WebP変換、透かし) |
| **メタデータ** | PostgreSQL Assetsテーブル |

**結論**: リポジトリに画像ファイルは含まれず、Filestoreが単一の真実の源（Single Source of Truth）として機能。

## 関連ドキュメント

- [アプリケーション構造](README.md)
- [ローカル開発環境](LOCAL_DEVELOPMENT.md)
- [GCPインフラ詳細](../gcp_inventory_reports/dorapita-inventory.md)
