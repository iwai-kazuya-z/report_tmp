# dorapita-core インベントリ

**調査日時**: 2025-12-21
**環境**: 本番環境 (Production)

## プロジェクト概要

| 項目 | 値 |
|------|-----|
| Project ID | dorapita-core |
| Project Number | 78522513121 |
| Status | ACTIVE |

---

## 1. コンピューティング

### Compute Engine (VM)

| 名前 | ゾーン | マシンタイプ | ステータス | 内部IP | 外部IP |
|------|--------|--------------|------------|--------|--------|
| gw-120011 | asia-northeast1-a | e2-small | RUNNING | 10.120.1.11 | 34.84.187.36 |
| web-120011 | asia-northeast1-a | n2-standard-2 | RUNNING | 10.120.101.11 | 34.146.149.10 |
| web-120021 | asia-northeast1-a | n2d-standard-2 | RUNNING | 10.120.101.21 | 35.190.224.174 |
| web-120022 | asia-northeast1-a | n2d-standard-2 | RUNNING | 10.120.101.22 | 104.198.83.165 |
| web-120023 | asia-northeast1-a | n2d-standard-2 | RUNNING | 10.120.101.23 | 34.85.114.186 |
| web-120031 | asia-northeast1-a | t2d-standard-2 | RUNNING | 10.120.101.31 | 34.84.115.88 |

**稼働中**: 6台 / **停止中**: 0台

### Cloud Run サービス

| サービス名 | リージョン | URL |
|------------|------------|-----|
| dorapita-maintenance | asia-northeast1 | https://dorapita-maintenance-kvv7uew52a-an.a.run.app |
| help-dorapita-com | asia-northeast1 | https://help-dorapita-com-kvv7uew52a-an.a.run.app |
| img-dorapita-com | asia-northeast1 | https://img-dorapita-com-kvv7uew52a-an.a.run.app |

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
| dorapita-core_cloudbuild | US | STANDARD |
| dorapita-infra | ASIA-NORTHEAST1 | STANDARD |
| dorapita-public | ASIA-NORTHEAST1 | STANDARD |
| run-sources-dorapita-core-asia-northeast1 | ASIA-NORTHEAST1 | STANDARD |

### Memorystore (Redis)

| インスタンス名 | Tier | メモリサイズ |
|----------------|------|--------------|
| redis-120011 | BASIC | 1 GB |

### Filestore (NFS)

| インスタンス名 | Tier | 容量 | IPアドレス | 共有パス |
|----------------|------|------|------------|----------|
| file-120011 | BASIC_HDD | 1024 GB | 10.117.161.2 | /dorapita |

**使用サービス**:
- img-dorapita-com (Cloud Run) → `/mnt` にマウント

### BigQuery

データセットなし

---

## 3. メッセージング・スケジューリング

### Pub/Sub トピック

| トピック名 |
|------------|
| container-analysis-notes-v1beta1 |
| container-analysis-occurrences-v1beta1 |
| container-analysis-occurrences-v1 |
| container-analysis-notes-v1 |

**サブスクリプション**: 未確認

### Cloud Scheduler

API未有効化

### Cloud Tasks

未確認

---

## 4. セキュリティ・IAM

### IAMバインディング（ユーザー・グループ）

| メンバー | ロール | 説明 |
|----------|--------|------|
| user:nagai@zigexn.co.jp | roles/owner | プロジェクトオーナー |
| user:kazuya.iwai@zigexn.co.jp | organizations/.../engineer.enableSudo | sudo権限 |
| user:h.sato@awesomegroup.co.jp | roles/viewer | 閲覧者 |
| user:m.yuasa@awesomegroup.co.jp | roles/viewer | 閲覧者 |

### サービスアカウント

| メールアドレス | 表示名 | 無効化 |
|----------------|--------|--------|
| 78522513121-compute@developer.gserviceaccount.com | Compute Engine default service account | No |

**サービスアカウントに付与されたロール**:
| サービスアカウント | ロール |
|-------------------|--------|
| 78522513121-compute@developer.gserviceaccount.com | roles/editor |
| 78522513121@cloudservices.gserviceaccount.com | roles/editor |
| 78522513121@cloudbuild.gserviceaccount.com | roles/cloudbuild.builds.builder |

### サービスアカウントキー

| サービスアカウント | キータイプ | ユーザー発行キー |
|-------------------|-----------|-----------------|
| 78522513121-compute@developer.gserviceaccount.com | SYSTEM_MANAGED | **なし** ✅ |

※ ユーザー発行のJSONキーファイル（USER_MANAGED）は存在しません

### Secret Manager

| シークレット名 | 作成日 |
|----------------|--------|
| GA4_API_SECRET | 2025-10-02 |
| GA4_MEASUREMENT_ID | 2025-10-02 |

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
| db | asia-northeast1 | 10.120.7.0/24 |
| gw | asia-northeast1 | 10.120.1.0/24 |
| web | asia-northeast1 | 10.120.101.0/24 |
| web-osaka | asia-northeast2 | 10.120.102.0/24 |

### ロードバランサー（転送ルール）

主要なエンドポイント:
| 名前 | IPアドレス |
|------|------------|
| kanri-dorapita-com-global | 35.186.245.94 |
| dorapita-com-global | 34.111.187.230 |
| img-dorapita-com-global | 35.190.92.252 |
| cadm-dorapita-com-global | 34.49.139.134 |
| edit-dorapita-com-global | 34.49.225.237 |
| help-dorapita-com-global | 34.36.124.133 |
| dora-pt-dorapita-com-global | 34.49.134.144 |

IPv6対応:
- dorapita-com-ipv6: 2600:1901:0:b576::
- img-dorapita-com-ipv6: 2600:1901:0:4eab::

### バックエンドサービス

| バックエンドサービス名 | グループ | プロトコル |
|------------------------|----------|------------|
| web-dora-pt-backend | web-php7 | HTTP |
| web-dorapita-admin-backend | web-admin-php5 | HTTP |
| web-dorapita-backend | web-dorapita | HTTP |
| web-dorapita-cached-backend | web-dorapita | HTTP |
| web-dorapita-cadm-backend | web-cadm | HTTP |
| web-dorapita-kanri-backend | web-kanri | HTTP |
| web-dorapita-kanri-restricted-backend | web-php7 | HTTP |
| web-dorapita-legacy-backend | web-php5 | HTTP |
| web-help-dorapita-backend | help-dorapita-com-neg | HTTPS |
| web-img-dorapita-backend | img-dorapita-com-neg | HTTPS |
| web-maintenance-dorapita-com | dorapita-maintenance-neg | HTTPS |

### Cloud DNS

| ゾーン名 | DNSネーム | 可視性 |
|----------|-----------|--------|
| dorapita-com | dorapita.com. | public |
| dorapita-com-private | dorapita.com. | private |
| gcp | gcp. | private |

### ファイアウォールルール (dorapita-vpc)

| ルール名 | 方向 | 許可 |
|----------|------|------|
| dorapita-vpc-allow-health-check | INGRESS | tcp |
| dorapita-vpc-allow-http | INGRESS | tcp:80 |
| dorapita-vpc-allow-ingress-from-iap | INGRESS | tcp:22 |
| dorapita-vpc-allow-internal | INGRESS | tcp |

---

## 6. 開発・デプロイツール

### Artifact Registry

| リポジトリ名 | 形式 |
|--------------|------|
| cloud-run-source-deploy | DOCKER |
| homebrew-containers | DOCKER |

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
| maps-backend.googleapis.com | Google Maps |
| geocoding-backend.googleapis.com | Geocoding |
| directions-backend.googleapis.com | Directions |
| recaptchaenterprise.googleapis.com | reCAPTCHA Enterprise |

**未有効化のAPI**:
- cloudfunctions.googleapis.com (Cloud Functions)
- cloudscheduler.googleapis.com (Cloud Scheduler)

---

## 8. アーキテクチャ概要

```
[ユーザー]
    |
    v
[Cloud Load Balancer]
  - kanri.dorapita.com (35.186.245.94)
  - dorapita.com (34.111.187.230)
  - img.dorapita.com (35.190.92.252)
  - etc.
    |
    +---> [Cloud Run]
    |       - help-dorapita-com
    |       - img-dorapita-com
    |       - dorapita-maintenance
    |
    +---> [Compute Engine VMs]
            - web-120011 (n2-standard-2) - メインWeb
            - web-120021/22/23 (n2d-standard-2) - Webクラスタ
            - web-120031 (t2d-standard-2)
            - gw-120011 (e2-small) - Gateway
                |
                v
            [Cloud SQL]
              - db-120011 (MySQL 5.7)
              - pg-120011 (PostgreSQL 10)
                |
                v
            [Memorystore Redis]
              - redis-120011 (1GB)
                |
                v
            [Cloud Storage]
              - dorapita-public
              - dorapita-infra
```

---

## 9. 本番環境の特徴

- **高可用性**: Webサーバー5台稼働
- **公開DNS**: dorapita.comの公開DNSゾーンを管理
- **IPv6対応**: 主要サービスがIPv6対応
- **Google Maps API**: Maps、Geocoding、Directions等のAPIが有効

---

*調査完了: 2025-12-21*
