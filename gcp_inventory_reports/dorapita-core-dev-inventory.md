# dorapita-core-dev インベントリ

**調査日時**: 2025-12-21
**環境**: 開発環境 (Staging)

## プロジェクト概要

| 項目 | 値 |
|------|-----|
| Project ID | dorapita-core-dev |
| Project Number | 1010225065326 |
| Status | ACTIVE |

---

## 1. コンピューティング

### Compute Engine (VM)

| 名前 | ゾーン | マシンタイプ | ステータス | 内部IP | 外部IP |
|------|--------|--------------|------------|--------|--------|
| ftp-120011 | asia-northeast1-b | e2-small | RUNNING | 10.121.13.11 | 35.200.74.161 |
| gw-120011 | asia-northeast1-b | e2-small | TERMINATED | 10.121.1.11 | - |
| mailhog | asia-northeast1-b | e2-micro | RUNNING | 10.121.14.2 | - |
| pg-120010 | asia-northeast1-b | n2-standard-2 | TERMINATED | 10.121.7.10 | - |
| web-120011 | asia-northeast1-b | e2-medium | RUNNING | 10.121.101.11 | 35.187.213.112 |
| web-120012 | asia-northeast1-b | e2-medium | TERMINATED | 10.121.101.12 | - |
| web-120021 | asia-northeast1-b | n2d-standard-2 | RUNNING | 10.121.101.21 | 35.200.24.243 |
| web-120031 | asia-northeast1-b | n2d-standard-2 | RUNNING | 10.121.101.31 | 34.146.10.179 |

**稼働中**: 5台 / **停止中**: 3台

### Cloud Run サービス

| サービス名 | リージョン | URL |
|------------|------------|-----|
| api-dorapita-com | asia-northeast1 | https://api-dorapita-com-dr3sq5ocsq-an.a.run.app |
| cadm-dorapita-com | asia-northeast1 | https://cadm-dorapita-com-dr3sq5ocsq-an.a.run.app |
| dorapita-com | asia-northeast1 | https://dorapita-com-dr3sq5ocsq-an.a.run.app |
| dorapita-maintenance | asia-northeast1 | https://dorapita-maintenance-dr3sq5ocsq-an.a.run.app |
| edit-dorapita-com | asia-northeast1 | https://edit-dorapita-com-dr3sq5ocsq-an.a.run.app |
| help-dorapita-com | asia-northeast1 | https://help-dorapita-com-dr3sq5ocsq-an.a.run.app |
| img-dorapita-com | asia-northeast1 | https://img-dorapita-com-dr3sq5ocsq-an.a.run.app |

**Cloud Runイメージ詳細**:

| サービス名 | ARリポジトリ | イメージタグ |
|------------|--------------|--------------|
| api-dorapita-com | api-dorapita-com | 4a6a231ca (CakePHP 5.2, PHP 8.4) |
| cadm-dorapita-com | cloud-run-source-deploy | f1a0e41ebc60... |
| dorapita-com | cloud-run-source-deploy | 8ba1c526f20b... |
| dorapita-maintenance | cloud-run-source-deploy | sha256:4408ae7c... |
| edit-dorapita-com | cloud-run-source-deploy | bbe31cbe6fd2... |
| help-dorapita-com | cloud-run-source-deploy | sha256:89dcd379... |
| img-dorapita-com | homebrew-containers | latest |

※ api-dorapita-comのみ専用ARリポジトリを使用。他はCloud Run Source Deploy経由。

**api-dorapita-com イメージ詳細分析** (2025-12-21調査):

| 項目 | 値 |
|------|-----|
| ベースOS | Debian 13 (trixie) |
| PHP | 8.4.15 (PHP-FPM) |
| フレームワーク | CakePHP 5.2 |
| イメージサイズ | 1.15GB |
| Entrypoint | docker-entrypoint.sh (nginx起動後 php-fpm実行) |
| ポート | 80/tcp, 9000/tcp |

**コード構成**:
```
/var/myapp/
├── src/           # アプリケーションコード
│   ├── Controller/  # AppController, ErrorController, PagesController のみ
│   ├── Model/       # 空（.gitkeepのみ）
│   └── View/        # AppView, AjaxView
├── config/        # routes.php等（デフォルト設定）
├── templates/     # テンプレート
└── vendor/        # Composer依存
```

**所見**: CakePHP 5.2のスケルトン（雛形）アプリケーション。カスタムコントローラー・モデルは未実装。開発初期段階またはテスト用ベースイメージと推測される。

### Cloud Run Jobs

なし

### Cloud Functions

API未有効化

### App Engine

未使用

### GKE (Kubernetes Engine)

クラスターなし

---

## 2. データベース・ストレージ

### Cloud SQL

| インスタンス名 | データベースバージョン | リージョン |
|----------------|------------------------|------------|
| db-120011 | MySQL 5.7 | asia-northeast1 |
| pg-120011 | PostgreSQL 10 | asia-northeast1 |

### Cloud Storage バケット

| バケット名 | ロケーション | ストレージクラス |
|------------|--------------|------------------|
| dorapita-core-dev-archives | ASIA-NORTHEAST1 | STANDARD |
| dorapita-core-dev_cloudbuild | US | STANDARD |
| dorapita-document | ASIA-NORTHEAST1 | STANDARD |
| dorapita-infra-stg | ASIA-NORTHEAST1 | STANDARD |
| dorapita-stg-public | ASIA-NORTHEAST1 | STANDARD |
| run-sources-dorapita-core-dev-asia-northeast1 | ASIA-NORTHEAST1 | STANDARD |

### Memorystore (Redis)

| インスタンス名 | Tier | メモリサイズ |
|----------------|------|--------------|
| redis-120011 | BASIC | 1 GB |

### Filestore (NFS)

| インスタンス名 | Tier | 容量 | IPアドレス | 共有パス |
|----------------|------|------|------------|----------|
| file-120011 | BASIC_HDD | 1024 GB | 10.227.113.2 | /dorapita |

**使用サービス**:
- dorapita-com (Cloud Run) → `/mnt` にマウント
- img-dorapita-com (Cloud Run) → `/mnt` にマウント

### BigQuery

データセットなし

---

## 3. メッセージング・スケジューリング

### Pub/Sub トピック

| トピック名 |
|------------|
| container-analysis-occurrences-v1beta1 |
| container-analysis-notes-v1beta1 |
| container-analysis-occurrences-v1 |
| container-analysis-notes-v1 |

**サブスクリプション**: なし

### Cloud Scheduler

API未有効化

### Cloud Tasks

API未有効化

---

## 4. セキュリティ・IAM

### IAMバインディング（ユーザー・グループ）

| メンバー | ロール | 説明 |
|----------|--------|------|
| user:nagai@zigexn.co.jp | roles/owner | プロジェクトオーナー |
| user:kazuya.iwai@zigexn.co.jp | organizations/.../engineer.enableSudo | sudo権限 |
| user:h.sato@awesomegroup.co.jp | organizations/.../engineer.enableSudo | sudo権限 |
| user:k.hashimoto@awesomegroup.co.jp | organizations/.../engineer.enableSudo | sudo権限 |
| user:m.yuasa@awesomegroup.co.jp | organizations/.../engineer.enableSudo | sudo権限 |
| group:g-awesome-leader@zigexn.vn | organizations/.../engineer.disableSudo | グループ |
| group:g-awesome-team@zigexn.vn | organizations/.../engineer.disableSudo | グループ |

**詳細ロール一覧**:

| メンバー | ロール |
|----------|--------|
| user:nagai@zigexn.co.jp | roles/owner |
| user:kazuya.iwai@zigexn.co.jp | roles/compute.instanceAdmin.v1, roles/iam.serviceAccountUser, roles/run.developer, roles/run.sourceDeveloper, roles/serviceusage.serviceUsageConsumer, roles/storage.admin |
| user:h.sato@awesomegroup.co.jp | roles/viewer, roles/cloudsql.client, roles/compute.instanceAdmin.v1, roles/iam.serviceAccountUser, roles/iap.tunnelResourceAccessor, roles/redis.viewer, roles/run.developer, roles/run.sourceDeveloper, roles/serviceusage.serviceUsageConsumer |
| user:k.hashimoto@awesomegroup.co.jp | roles/cloudsql.client, roles/iap.tunnelResourceAccessor, roles/serviceusage.serviceUsageViewer |
| user:m.yuasa@awesomegroup.co.jp | roles/viewer, roles/compute.instanceAdmin.v1, roles/iam.serviceAccountUser, roles/redis.viewer, roles/run.developer, roles/run.sourceDeveloper, roles/serviceusage.serviceUsageConsumer |
| user:r.matsui@awesomegroup.co.jp | roles/iap.httpsResourceAccessor |
| group:g-awesome-leader@zigexn.vn | roles/viewer, roles/cloudsql.client, roles/compute.osLogin, roles/iam.serviceAccountUser, roles/iap.tunnelResourceAccessor |
| group:g-awesome-team@zigexn.vn | roles/viewer, roles/cloudsql.client, roles/compute.osLogin, roles/iam.serviceAccountUser, roles/iap.tunnelResourceAccessor, roles/recaptchaenterprise.admin, roles/storage.admin |

### サービスアカウント

| メールアドレス | 表示名 | 無効化 |
|----------------|--------|--------|
| 1010225065326-compute@developer.gserviceaccount.com | Compute Engine default service account | No |
| dorapita-gcs@dorapita-core-dev.iam.gserviceaccount.com | dorapita-gcs | No |

**サービスアカウントに付与されたロール**:
| サービスアカウント | ロール |
|-------------------|--------|
| 1010225065326-compute@developer.gserviceaccount.com | roles/editor, roles/artifactregistry.writer, roles/iam.serviceAccountUser, roles/logging.logWriter, roles/secretmanager.secretAccessor |
| 1010225065326@cloudservices.gserviceaccount.com | roles/editor |
| 1010225065326@cloudbuild.gserviceaccount.com | roles/cloudbuild.builds.builder, roles/run.admin |

### サービスアカウントキー

| サービスアカウント | キータイプ | ユーザー発行キー |
|-------------------|-----------|-----------------|
| 1010225065326-compute@developer.gserviceaccount.com | SYSTEM_MANAGED | **なし** ✅ |
| dorapita-gcs@dorapita-core-dev.iam.gserviceaccount.com | SYSTEM_MANAGED | **なし** ✅ |

※ ユーザー発行のJSONキーファイル（USER_MANAGED）は存在しません

### Secret Manager

| シークレット名 | 作成日 |
|----------------|--------|
| database-edit-dorapita-com | 2025-02-15 |
| settings-edit-dorapita-com | 2025-02-15 |
| shell-environment-variables-api-dorapita-com | 2025-12-16 |
| shell-environment-variables-cadm-dorapita-com | 2025-02-14 |
| shell-environment-variables-dorapita-com | 2025-02-06 |

---

## 5. ネットワーキング

### VPC Networks

| ネットワーク名 | サブネットモード | BGPルーティングモード |
|----------------|------------------|----------------------|
| default | AUTO | REGIONAL |
| dorapita-vpc | CUSTOM | REGIONAL |

### サブネット (dorapita-vpc)

| サブネット名 | リージョン | CIDR |
|--------------|------------|------|
| api | asia-northeast1 | 10.121.3.0/24 |
| db | asia-northeast1 | 10.121.7.0/24 |
| ftp | asia-northeast1 | 10.121.13.0/24 |
| gw | asia-northeast1 | 10.121.1.0/24 |
| mail | asia-northeast1 | 10.121.14.0/24 |
| web | asia-northeast1 | 10.121.101.0/24 |

### ロードバランサー（転送ルール）

複数のグローバルHTTPS転送ルールが設定されています（計26個）。

主要なエンドポイント:
- stg-dorapita-com: 34.36.242.183
- kanri-stg-dorapita-com: 35.244.223.57
- cadm-stg-dorapita-com: 34.160.129.119
- dev-dorapita-com: 34.107.130.29
- api-rc-dorapita-com: 34.120.125.219

### バックエンドサービス

| バックエンドサービス名 | グループ | プロトコル |
|------------------------|----------|------------|
| api-dorapita-backend | api-dorapita-com-neg | HTTPS |
| web-dorapita-api-backend | web-dorapita | HTTP |
| web-dorapita-backend | web-dorapita | HTTP |
| web-dorapita-dev-backend | web-dorapita-neg | HTTPS |
| web-dorapita-kanri-backend | web-php7 | HTTP |
| web-dorapita-kanri-restricted-backend | web-php7 | HTTP |
| web-dorapita-legacy-backend | web-php5 | HTTP |
| web-dorapita-restricted-backend | web-dorapita | HTTP |
| web-dorapita-serverless-backend | web-dorapita-neg | HTTPS |
| web-dorapita-users-backend | web-dorapita-users | HTTP |
| web-help-dorapita-backend | help-dorapita-com-neg | HTTPS |
| web-img-dorapita-backend | img-dorapita-com-neg | HTTPS |
| web-mailhog | mailhog | HTTP |

### Cloud DNS

| ゾーン名 | DNSネーム | 可視性 |
|----------|-----------|--------|
| gcp | gcp. | private |
| rc-dorapita-com | rc.dorapita.com. | private |

### ファイアウォールルール

| ルール名 | ネットワーク | 方向 | 許可 |
|----------|--------------|------|------|
| default-allow-http | default | INGRESS | tcp:80 |
| default-allow-icmp | default | INGRESS | icmp |
| default-allow-internal | default | INGRESS | tcp:0-65535,udp:0-65535,icmp |
| default-allow-rdp | default | INGRESS | tcp:3389 |
| default-allow-ssh | default | INGRESS | tcp:22 |
| default-allow-ssh-ipv6 | default | INGRESS | tcp:22 |
| dorapita-vpc-allow-ftp | dorapita-vpc | INGRESS | tcp |
| dorapita-vpc-allow-health-check | dorapita-vpc | INGRESS | tcp |
| dorapita-vpc-allow-http | dorapita-vpc | INGRESS | tcp:80 |
| dorapita-vpc-allow-ingress-from-iap | dorapita-vpc | INGRESS | tcp:22 |
| dorapita-vpc-allow-internal | dorapita-vpc | INGRESS | all |

---

## 6. 開発・デプロイツール

### Artifact Registry

| リポジトリ名 | 形式 |
|--------------|------|
| api-dorapita-com | DOCKER |
| cloud-run-source-deploy | DOCKER |
| homebrew-containers | DOCKER |
| gcr.io | DOCKER |

### Cloud Build

有効（cloudbuild.googleapis.com）

---

## 7. 有効化されているAPI（主要なもの）

| API | 説明 |
|-----|------|
| compute.googleapis.com | Compute Engine |
| run.googleapis.com | Cloud Run |
| sqladmin.googleapis.com | Cloud SQL Admin |
| redis.googleapis.com | Memorystore Redis |
| storage.googleapis.com | Cloud Storage |
| secretmanager.googleapis.com | Secret Manager |
| pubsub.googleapis.com | Pub/Sub |
| dns.googleapis.com | Cloud DNS |
| artifactregistry.googleapis.com | Artifact Registry |
| cloudbuild.googleapis.com | Cloud Build |
| container.googleapis.com | Kubernetes Engine |
| iam.googleapis.com | IAM |
| logging.googleapis.com | Cloud Logging |
| monitoring.googleapis.com | Cloud Monitoring |
| vpcaccess.googleapis.com | Serverless VPC Access |

**未有効化のAPI**:
- cloudfunctions.googleapis.com (Cloud Functions)
- file.googleapis.com (Filestore)
- cloudscheduler.googleapis.com (Cloud Scheduler)
- cloudtasks.googleapis.com (Cloud Tasks)

---

## 8. アーキテクチャ概要

```
[ユーザー]
    |
    v
[Cloud Load Balancer] (34.36.242.183 等)
    |
    +---> [Cloud Run] (api, cadm, edit, help, img, dorapita-com)
    |
    +---> [Compute Engine VMs]
          - web-120011/21/31 (Webサーバー)
          - ftp-120011 (FTPサーバー)
          - mailhog (メールテスト)
              |
              v
          [Cloud SQL]
          - db-120011 (MySQL 5.7)
          - pg-120011 (PostgreSQL 10)
              |
              v
          [Memorystore Redis]
          - redis-120011
              |
              v
          [Cloud Storage]
          - dorapita-stg-public
          - dorapita-document
```

---

*調査完了: 2025-12-21*
