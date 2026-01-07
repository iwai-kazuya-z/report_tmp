# Dorapita GCPインフラ インベントリ概要

**調査日時**: 2025-12-21

---

## サービス構成・接続関係図

### 本番環境 (dorapita-core)

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              【 Cloud DNS 】                                     │
│   dorapita.com (公開)  ──┬── dorapita.com → 34.111.187.230                      │
│                          ├── kanri.dorapita.com → 35.186.245.94                 │
│                          ├── img.dorapita.com → 35.190.92.252                   │
│                          ├── help.dorapita.com → 34.36.124.133                  │
│                          └── cadm.dorapita.com → 34.49.139.134                  │
│   gcp (private)         ── 内部名前解決用                                        │
│   dorapita.com (private) ── VPC内名前解決用                                      │
└────────────────────────────────────┬────────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        【 GCLB (Global HTTPS LB) 】                              │
│   + Cloud Armor (WAF: allow-dorapita-system等)                                  │
│   + SSLポリシー: 未設定 ⚠️                                                       │
│                                                                                  │
│   バックエンドサービス:                                                           │
│   ├── web-dorapita-backend ────────────────────┐                                │
│   ├── web-dorapita-kanri-backend ──────────────┼──→ Compute Engine (web-*)      │
│   ├── web-dorapita-cadm-backend ───────────────┘                                │
│   ├── web-help-dorapita-backend ───────────────────→ Cloud Run (NEG)            │
│   ├── web-img-dorapita-backend ────────────────────→ Cloud Run (NEG)            │
│   └── web-maintenance-dorapita-com ────────────────→ Cloud Run (NEG)            │
└────────────────────────────────────┬────────────────────────────────────────────┘
                                     │
              ┌──────────────────────┴──────────────────────┐
              ▼                                             ▼
┌──────────────────────────────┐              ┌──────────────────────────────┐
│     【 Compute Engine 】      │              │       【 Cloud Run 】         │
│                              │              │                              │
│  web-120011 (n2-standard-2)  │              │  help-dorapita-com           │
│  web-120021 (n2d-standard-2) │              │    └─ Ingress: internal+LB   │
│  web-120022 (n2d-standard-2) │              │    └─ VPC: なし              │
│  web-120023 (n2d-standard-2) │              │                              │
│  web-120031 (t2d-standard-2) │              │  img-dorapita-com            │
│  gw-120011 (e2-small)        │              │    └─ Ingress: internal+LB   │
│                              │              │    └─ VPC: あり              │
│  VPC: dorapita-vpc           │              │                              │
│  Subnet: 10.120.101.0/24     │              │  dorapita-maintenance        │
│                              │              │    └─ Ingress: internal+LB   │
└──────────────┬───────────────┘              │    └─ VPC: なし              │
               │                              └──────────────────────────────┘
               │  プライベートIP接続                        │ DBアクセス不要
               ▼                                           │
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           【 データベース層 】                                   │
│                                                                                  │
│  ┌─────────────────────────┐  ┌─────────────────────────┐  ┌─────────────────┐ │
│  │     Cloud SQL          │  │     Cloud SQL          │  │   Memorystore   │ │
│  │     (MySQL 5.7)        │  │   (PostgreSQL 10)      │  │     (Redis)     │ │
│  │     db-120011          │  │     pg-120011          │  │   redis-120011  │ │
│  │                        │  │                        │  │     1GB BASIC   │ │
│  │  Private: 10.117.160.3 │  │  Private: 10.117.160.5 │  │                 │ │
│  │  SSL: 未強制 ⚠️         │  │  SSL: 未強制 ⚠️         │  │                 │ │
│  └─────────────────────────┘  └─────────────────────────┘  └─────────────────┘ │
│                                                                                  │
│  【 Filestore 】: file-120011 (1TB BASIC_HDD, 10.117.161.2, /dorapita)          │
│    └─ img-dorapita-com (Cloud Run) が /mnt にマウント                           │
└─────────────────────────────────────────────────────────────────────────────────┘
```

**本番環境の特徴:**
- Compute Engine が主要なWebワークロードを処理
- Cloud Run は静的コンテンツ配信（help, img, maintenance）のみ
- 全DBアクセスはCompute Engine経由（Cloud RunはDB接続なし）

---

### 開発環境 (dorapita-core-dev)

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              【 Cloud DNS 】                                     │
│   gcp (private)          ── 内部名前解決用                                       │
│   rc.dorapita.com (private) ── RC環境用                                         │
│   ※ 公開DNSゾーンなし（本番で一元管理）                                           │
└────────────────────────────────────┬────────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        【 GCLB (Global HTTPS LB) 】                              │
│   + Cloud Armor (WAF)                                                           │
│   + SSLポリシー: 未設定 ⚠️                                                       │
│                                                                                  │
│   主要エンドポイント:                                                             │
│   ├── stg-dorapita-com: 34.36.242.183                                           │
│   ├── dev-dorapita-com: 34.107.130.29                                           │
│   ├── kanri-stg-dorapita-com: 35.244.223.57                                     │
│   └── api-rc-dorapita-com: 34.120.125.219                                       │
│                                                                                  │
│   バックエンドサービス → Compute Engine / Cloud Run (NEG)                        │
└────────────────────────────────────┬────────────────────────────────────────────┘
                                     │
         ┌───────────────────────────┴───────────────────────────┐
         ▼                                                       ▼
┌──────────────────────────────┐              ┌──────────────────────────────────────┐
│     【 Compute Engine 】      │              │           【 Cloud Run 】             │
│                              │              │                                      │
│  web-120011 (e2-medium)      │              │  dorapita-com      ← VPC:あり DB:可  │
│  web-120021 (n2d-standard-2) │              │  cadm-dorapita-com ← VPC:あり DB:可  │
│  web-120031 (n2d-standard-2) │              │  img-dorapita-com  ← VPC:あり DB:可  │
│  ftp-120011 (e2-small)       │              │                                      │
│  mailhog (e2-micro)          │              │  api-dorapita-com  ← VPC:なし ⚠️     │
│                              │              │    └─ Ingress: all (全公開) ⚠️       │
│  VPC: dorapita-vpc           │              │                                      │
│  Subnet: 10.121.101.0/24     │              │  edit-dorapita-com ← VPC:なし ⚠️     │
│                              │              │    └─ DBシークレットあるがDB接続不可  │
│  停止中:                      │              │                                      │
│  - gw-120011                 │              │  help-dorapita-com ← VPC:なし        │
│  - pg-120010                 │              │  dorapita-maintenance ← VPC:なし     │
│  - web-120012                │              │                                      │
└──────────────┬───────────────┘              └──────────────┬───────────────────────┘
               │                                             │
               │  プライベートIP接続                          │ VPC Direct Egress
               ▼                                             ▼ (VPC接続ありのみ)
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           【 データベース層 】                                   │
│                                                                                  │
│  ┌─────────────────────────┐  ┌─────────────────────────┐  ┌─────────────────┐ │
│  │     Cloud SQL          │  │     Cloud SQL          │  │   Memorystore   │ │
│  │     (MySQL 5.7)        │  │   (PostgreSQL 10)      │  │     (Redis)     │ │
│  │     db-120011          │  │     pg-120011          │  │   redis-120011  │ │
│  │                        │  │                        │  │     1GB BASIC   │ │
│  │  Private: 10.227.112.3 │  │  Private: 10.227.112.14│  │                 │ │
│  │  SSL: 未強制 ⚠️         │  │  SSL: 未強制 ⚠️         │  │                 │ │
│  └─────────────────────────┘  └─────────────────────────┘  └─────────────────┘ │
│                                                                                  │
│  【 Filestore 】: file-120011 (1TB BASIC_HDD, 10.227.113.2, /dorapita)          │
│    └─ dorapita-com, img-dorapita-com (Cloud Run) が /mnt にマウント             │
└─────────────────────────────────────────────────────────────────────────────────┘
```

**開発環境の特徴:**
- Cloud Run でAPI・管理画面もサーバーレス化を試行
- VPC接続の有無でDB接続可否が決まる
- api-dorapita-comはIngress `all`で直接アクセス可能（要修正）

---

### サービス接続マトリックス

| サービス | DNS解決 | GCLB経由 | Cloud SQL接続 | Redis接続 | Filestore |
|----------|---------|----------|---------------|-----------|-----------|
| **本番 Compute Engine** | ✅ | ✅ | ✅ (Private IP) | ✅ | - |
| **本番 Cloud Run (img)** | ✅ | ✅ | ❌ (不要) | ❌ | ✅ (NFS) |
| **本番 Cloud Run (他)** | ✅ | ✅ | ❌ (不要) | ❌ | - |
| **開発 Compute Engine** | ✅ | ✅ | ✅ (Private IP) | ✅ | - |
| **開発 Cloud Run (VPC有)** | ✅ | ✅ | ✅ (VPC Egress) | ✅ | ✅ (NFS) |
| **開発 Cloud Run (VPC無)** | ✅ | ✅/⚠️ | ❌ | ❌ | - |

※ Filestore: 両環境でCloud Run (img-dorapita-com等) がNFSマウントで使用

---

## プロジェクト一覧

| 環境 | Project ID | Project Number | 用途 |
|------|------------|----------------|------|
| 本番 | dorapita-core | 78522513121 | Production環境 |
| 開発 | dorapita-core-dev | 1010225065326 | Staging/Development環境 |

---

## 環境比較サマリー

### コンピューティングリソース

| リソース | dorapita-core (本番) | dorapita-core-dev (開発) |
|----------|---------------------|-------------------------|
| **Compute Engine** | 6台（全稼働） | 8台（5稼働/3停止） |
| **Cloud Run** | 3サービス | 7サービス |
| **Cloud Run Jobs** | なし | なし |
| **Cloud Functions** | 未使用 | 未使用 |
| **App Engine** | 未使用 | 未使用 |
| **GKE** | なし | なし |

### データベース・ストレージ

| リソース | dorapita-core (本番) | dorapita-core-dev (開発) |
|----------|---------------------|-------------------------|
| **Cloud SQL** | 2 (MySQL 5.7, PostgreSQL 10) | 2 (MySQL 5.7, PostgreSQL 10) |
| **Cloud Storage** | 4バケット | 6バケット |
| **Memorystore Redis** | 1 (1GB BASIC) | 1 (1GB BASIC) |
| **BigQuery** | なし | なし |
| **Filestore** | 1 (1TB BASIC_HDD) | 1 (1TB BASIC_HDD) |

### セキュリティ・IAM

| リソース | dorapita-core (本番) | dorapita-core-dev (開発) |
|----------|---------------------|-------------------------|
| **サービスアカウント** | 1 | 2 |
| **Secret Manager** | 2 (GA4関連) | 5 (環境変数, DB設定) |
| **SAキー (USER_MANAGED)** | なし ✅ | なし ✅ |

**IAMユーザー比較**:

| ユーザー | 本番 | 開発 |
|----------|------|------|
| nagai@zigexn.co.jp | Owner | Owner |
| kazuya.iwai@zigexn.co.jp | enableSudo | enableSudo, compute.instanceAdmin, run.developer 等 |
| h.sato@awesomegroup.co.jp | Viewer | enableSudo, compute.instanceAdmin, run.developer 等 |
| m.yuasa@awesomegroup.co.jp | Viewer | enableSudo, compute.instanceAdmin, run.developer 等 |
| k.hashimoto@awesomegroup.co.jp | - | enableSudo, cloudsql.client, iap.tunnelResourceAccessor |
| r.matsui@awesomegroup.co.jp | - | iap.httpsResourceAccessor |
| g-awesome-leader@zigexn.vn | - | disableSudo, viewer, cloudsql.client, osLogin 等 |
| g-awesome-team@zigexn.vn | - | disableSudo, viewer, cloudsql.client, osLogin, storage.admin 等 |

### ネットワーキング

| リソース | dorapita-core (本番) | dorapita-core-dev (開発) |
|----------|---------------------|-------------------------|
| **VPC** | 2 (default, dorapita-vpc) | 2 (default, dorapita-vpc) |
| **サブネット** | 4 (db, gw, web, web-osaka) | 6 (api, db, ftp, gw, mail, web) |
| **Cloud DNS** | 3ゾーン（公開DNS含む） | 2ゾーン（プライベートのみ） |
| **ロードバランサー** | 17転送ルール | 26転送ルール |
| **IPv6対応** | あり | あり |

---

## VMインスタンス詳細比較

### dorapita-core (本番)

| 名前 | マシンタイプ | ゾーン | 用途 |
|------|--------------|--------|------|
| gw-120011 | e2-small | asia-northeast1-a | Gateway |
| web-120011 | n2-standard-2 | asia-northeast1-a | Webサーバー |
| web-120021 | n2d-standard-2 | asia-northeast1-a | Webサーバー |
| web-120022 | n2d-standard-2 | asia-northeast1-a | Webサーバー |
| web-120023 | n2d-standard-2 | asia-northeast1-a | Webサーバー |
| web-120031 | t2d-standard-2 | asia-northeast1-a | Webサーバー |

### dorapita-core-dev (開発)

| 名前 | マシンタイプ | ゾーン | ステータス | 用途 |
|------|--------------|--------|------------|------|
| gw-120011 | e2-small | asia-northeast1-b | TERMINATED | Gateway |
| ftp-120011 | e2-small | asia-northeast1-b | RUNNING | FTPサーバー |
| mailhog | e2-micro | asia-northeast1-b | RUNNING | メールテスト |
| pg-120010 | n2-standard-2 | asia-northeast1-b | TERMINATED | DB（旧） |
| web-120011 | e2-medium | asia-northeast1-b | RUNNING | Webサーバー |
| web-120012 | e2-medium | asia-northeast1-b | TERMINATED | Webサーバー（予備） |
| web-120021 | n2d-standard-2 | asia-northeast1-b | RUNNING | Webサーバー |
| web-120031 | n2d-standard-2 | asia-northeast1-b | RUNNING | Webサーバー |

---

## Cloud Runサービス比較

### 両環境共通

| サービス名 | 説明 |
|------------|------|
| dorapita-maintenance | メンテナンスページ |
| help-dorapita-com | ヘルプサイト |
| img-dorapita-com | 画像配信 |

### dorapita-core-dev のみ

| サービス名 | 説明 |
|------------|------|
| api-dorapita-com | API |
| cadm-dorapita-com | 管理画面 |
| dorapita-com | メインサイト |
| edit-dorapita-com | 編集機能 |

### 開発環境 Cloud Run イメージ詳細

| サービス名 | ARリポジトリ | テクノロジー |
|------------|--------------|--------------|
| api-dorapita-com | api-dorapita-com (専用) | CakePHP 5.2, PHP 8.4, Debian 13 |
| cadm-dorapita-com | cloud-run-source-deploy | - |
| dorapita-com | cloud-run-source-deploy | - |
| edit-dorapita-com | cloud-run-source-deploy | - |
| help-dorapita-com | cloud-run-source-deploy | - |
| img-dorapita-com | homebrew-containers | - |
| dorapita-maintenance | cloud-run-source-deploy | - |

※ api-dorapita-comのみ専用のArtifact Registryリポジトリを使用

---

## ネットワーク構成

### IPアドレス範囲

| 環境 | VPC | CIDR範囲 |
|------|-----|----------|
| 本番 | dorapita-vpc | 10.120.0.0/16 |
| 開発 | dorapita-vpc | 10.121.0.0/16 |

### サブネット構成の違い

| サブネット | 本番 (10.120.x.x) | 開発 (10.121.x.x) |
|------------|-------------------|-------------------|
| gw | 10.120.1.0/24 | 10.121.1.0/24 |
| api | - | 10.121.3.0/24 |
| db | 10.120.7.0/24 | 10.121.7.0/24 |
| ftp | - | 10.121.13.0/24 |
| mail | - | 10.121.14.0/24 |
| web | 10.120.101.0/24 | 10.121.101.0/24 |
| web-osaka | 10.120.102.0/24 | - |

---

## 環境間の主な違い

### 1. インフラ規模
- **本番**: Webサーバー5台構成（高可用性）
- **開発**: Webサーバー3台（一部停止中）、開発用サービス（mailhog, FTP）あり

### 2. Cloud Run活用
- **本番**: 3サービス（静的コンテンツ中心）
- **開発**: 7サービス（API、管理画面もCloud Run化）

### 3. DNS管理
- **本番**: dorapita.comの公開DNSゾーンを管理
- **開発**: プライベートDNSのみ

### 4. Secret Manager
- **本番**: GA4関連のシークレットのみ（2個）
- **開発**: 各サービスの環境変数を管理（5個）

### 5. マシンタイプ
- **本番**: n2-standard-2, n2d-standard-2, t2d-standard-2（高性能）
- **開発**: e2-small, e2-medium, e2-micro（コスト最適化）

---

## アーキテクチャ概念図

```
                    ┌─────────────────────────────────────────┐
                    │              ユーザー                    │
                    └────────────────┬────────────────────────┘
                                     │
                    ┌────────────────┴────────────────────────┐
                    │         Cloud Load Balancer              │
                    │  (本番: 35.186.245.94 等)                │
                    │  (開発: 34.36.242.183 等)                │
                    └────────────────┬────────────────────────┘
                                     │
              ┌──────────────────────┼──────────────────────┐
              │                      │                      │
              ▼                      ▼                      ▼
    ┌─────────────────┐   ┌─────────────────┐   ┌─────────────────┐
    │    Cloud Run    │   │  Compute Engine │   │  Compute Engine │
    │  (img, help,    │   │   (web-12001x)  │   │   (web-12002x)  │
    │   maintenance)  │   │                 │   │                 │
    └─────────────────┘   └────────┬────────┘   └────────┬────────┘
                                   │                      │
                                   └──────────┬───────────┘
                                              │
                    ┌─────────────────────────┼─────────────────────────┐
                    │                         │                         │
                    ▼                         ▼                         ▼
          ┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
          │    Cloud SQL    │     │    Cloud SQL    │     │   Memorystore   │
          │    (MySQL)      │     │  (PostgreSQL)   │     │     (Redis)     │
          │   db-120011     │     │   pg-120011     │     │   redis-120011  │
          └─────────────────┘     └─────────────────┘     └─────────────────┘
                                              │
                                              ▼
                                  ┌─────────────────┐
                                  │  Cloud Storage  │
                                  │ (dorapita-public│
                                  │  等)            │
                                  └─────────────────┘
```

---

## 推奨事項・注意点

1. **データベースバージョン**: MySQL 5.7, PostgreSQL 10 は古いバージョン。アップグレード検討を推奨

2. **未有効化API**: Cloud Functions, Cloud Scheduler, Cloud Tasks が両環境で未有効化

3. **開発環境の停止VM**: gw-120011, pg-120010, web-120012 が停止中。不要であれば削除検討

4. **Secret管理の差異**: 本番と開発でSecret Managerの使い方が異なる

5. **IPv6対応**: 両環境でIPv6が有効化済み

---

## Cloud SQL接続ガイド

### 接続方式

Cloud SQLへの接続には**Cloud SQL Auth Proxy**を使用する。

```
[Client] ---(Auth Proxy/IAM認証)---> [Cloud SQL Instance] ---(DB認証)---> [Database]
```

### 前提条件

| 必要な権限/情報 | 用途 |
|----------------|------|
| `roles/cloudsql.client` | Auth Proxy接続（IAM認証） |
| DBユーザー名・パスワード | データベース認証 |

### Cloud SQLインスタンス情報

| 環境 | インスタンス | バージョン | 接続名 |
|------|-------------|-----------|--------|
| 開発 | db-120011 | MySQL 5.7 | `dorapita-core-dev:asia-northeast1:db-120011` |
| 開発 | pg-120011 | PostgreSQL 10 | `dorapita-core-dev:asia-northeast1:pg-120011` |
| 本番 | db-120011 | MySQL 5.7 | `dorapita-core:asia-northeast1:db-120011` |
| 本番 | pg-120011 | PostgreSQL 10 | `dorapita-core:asia-northeast1:pg-120011` |

### DB認証情報の取得

DB認証情報（ユーザー名・パスワード）はSecret Managerに格納されている。
取得には以下のいずれかが必要：

#### 方法1: Secret Manager権限を付与してもらう（推奨）

管理者に以下を依頼：

```bash
gcloud projects add-iam-policy-binding dorapita-core-dev \
  --member="user:<YOUR_EMAIL>" \
  --role="roles/secretmanager.secretAccessor"
```

権限付与後、以下でシークレットを取得：

```bash
gcloud secrets versions access latest \
  --secret=database-edit-dorapita-com \
  --project=dorapita-core-dev
```

#### 方法2: 管理者からDB認証情報を直接共有してもらう

必要な情報：
- DBユーザー名
- DBパスワード
- DB名

### 接続手順

1. **Auth Proxyをインストール**
   ```bash
   # macOS
   brew install cloud-sql-proxy
   ```

2. **Auth Proxyを起動**
   ```bash
   # MySQL (開発環境)
   cloud-sql-proxy dorapita-core-dev:asia-northeast1:db-120011 --port=3306

   # PostgreSQL (開発環境)
   cloud-sql-proxy dorapita-core-dev:asia-northeast1:pg-120011 --port=5432
   ```

3. **別ターミナルでDB接続**
   ```bash
   # MySQL
   mysql -h 127.0.0.1 -P 3306 -u <USER> -p <DATABASE>

   # PostgreSQL
   psql -h 127.0.0.1 -p 5432 -U <USER> -d <DATABASE>
   ```

### 関連シークレット一覧（開発環境）

| シークレット名 | 推定内容 |
|---------------|---------|
| database-edit-dorapita-com | DB接続情報 |
| shell-environment-variables-dorapita-com | 環境変数（DB情報含む可能性） |
| shell-environment-variables-api-dorapita-com | API用環境変数 |
| shell-environment-variables-cadm-dorapita-com | 管理画面用環境変数 |
| settings-edit-dorapita-com | アプリ設定 |

---

## 関連ドキュメント

- [dorapita-core-inventory.md](./dorapita-core-inventory.md) - 本番環境詳細
- [dorapita-core-dev-inventory.md](./dorapita-core-dev-inventory.md) - 開発環境詳細

---

*調査完了: 2025-12-21*
