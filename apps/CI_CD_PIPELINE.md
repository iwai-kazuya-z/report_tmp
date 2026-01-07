# CI/CDパイプライン

## 概要

ドラピタプロジェクトは**Google Cloud Build**を使用した自動ビルド・デプロイパイプラインを採用しています。

## アーキテクチャ

```
┌─────────────┐      ┌──────────────┐      ┌─────────────────┐      ┌─────────────┐
│   GitHub    │      │ Cloud Build  │      │ Artifact Reg    │      │ Cloud Run   │
│             │      │              │      │                 │      │             │
│  ZIGExN/    │──push│   Trigger    │─build│  docker.pkg.dev │─deploy│ dorapita-com│
│  dorapita   │──────▶   (auto)     ├─────▶│  /cloud-run-... ├─────▶│ (dev)       │
│             │      │              │      │                 │      │             │
│ dev.*.com   │      │  Dockerfile  │      │  :COMMIT_SHA    │      │             │
│  branch     │      │              │      │                 │      │             │
└─────────────┘      └──────────────┘      └─────────────────┘      └─────────────┘
```

## デプロイフロー

### 1. ブランチ戦略

| ブランチ名 | デプロイ先 | トリガー |
|-----------|-----------|---------|
| `dev.dorapita.com` | Cloud Run: dorapita-com (dev) | 自動 |
| `dev.cadm.dorapita.com` | Cloud Run: cadm-dorapita-com (dev) | 自動 |
| `dev.edit.dorapita.com` | Cloud Run: edit-dorapita-com (dev) | 自動 |
| `main` / `master` | 本番VM (手動) | 手動 |

### 2. Cloud Build Trigger設定

#### dorapita.com デプロイトリガー

```yaml
name: rmgpgab-dorapita-com-asia-northeast1-ZIGExN-dorapita--dev-dopjo
description: Build and deploy to Cloud Run service dorapita-com on push to "dev.dorapita.com"

github:
  owner: ZIGExN
  name: dorapita
  push:
    branch: ^dev.dorapita.com$

build:
  steps:
  - name: 'gcr.io/cloud-builders/docker'
    args:
    - build
    - --no-cache
    - -t
    - asia-northeast1-docker.pkg.dev/$PROJECT_ID/cloud-run-source-deploy/dorapita/dorapita-com:$COMMIT_SHA
    - .
    - -f
    - Dockerfile
    id: Build

  - name: 'gcr.io/cloud-builders/docker'
    args:
    - push
    - asia-northeast1-docker.pkg.dev/$PROJECT_ID/cloud-run-source-deploy/dorapita/dorapita-com:$COMMIT_SHA
    id: Push

  - name: 'gcr.io/google.com/cloudsdktool/cloud-sdk:slim'
    entrypoint: gcloud
    args:
    - run
    - services
    - update
    - dorapita-com
    - --platform=managed
    - --image=asia-northeast1-docker.pkg.dev/$PROJECT_ID/cloud-run-source-deploy/dorapita/dorapita-com:$COMMIT_SHA
    - --region=asia-northeast1
    - --quiet
    id: Deploy

  images:
  - asia-northeast1-docker.pkg.dev/$PROJECT_ID/cloud-run-source-deploy/dorapita/dorapita-com:$COMMIT_SHA

  substitutions:
    _SERVICE_NAME: dorapita-com
    _DEPLOY_REGION: asia-northeast1
```

## デプロイ手順

### 開発環境デプロイ（自動）

```bash
# 1. 開発ブランチをチェックアウト
git checkout dev.dorapita.com

# 2. 変更をコミット
git add .
git commit -m "feat: 新機能追加"

# 3. GitHubにプッシュ → 自動的にCloud Buildがトリガー
git push origin dev.dorapita.com
```

**自動実行される処理**:
1. Cloud Build TriggerがGitHub pushを検知
2. Dockerfileからイメージをビルド
3. Artifact Registryにpush（タグ: コミットSHA）
4. Cloud Runサービスを更新

### デプロイ確認

```bash
# Cloud Buildステータス確認
gcloud builds list --project=dorapita-core-dev --limit=5

# ビルドログ確認
gcloud builds log <BUILD_ID> --project=dorapita-core-dev

# Cloud Runサービス確認
gcloud run services describe dorapita-com \
  --project=dorapita-core-dev \
  --region=asia-northeast1
```

### 本番環境デプロイ（手動）

⚠️ **本番環境はCloud Run未使用。Compute Engine VM上で稼働。**

本番デプロイは現状、手動でVMにSSH接続して実施している模様。

```bash
# 1. 本番VMにSSH
gcloud compute ssh web-120011 --project=dorapita-core --zone=asia-northeast1-b

# 2. コードをpull
cd /var/www/dorapita.com
git pull origin main

# 3. 依存パッケージ更新
composer install --no-dev --optimize-autoloader

# 4. マイグレーション実行
bin/cake migrations migrate

# 5. キャッシュクリア
bin/cake cache clear_all

# 6. Apache再起動
sudo systemctl restart httpd
```

## ビルド詳細

### Dockerイメージビルド

```bash
# ビルドコマンド（Cloud Buildで自動実行）
docker build \
  --no-cache \
  -t asia-northeast1-docker.pkg.dev/dorapita-core-dev/cloud-run-source-deploy/dorapita/dorapita-com:$COMMIT_SHA \
  -f Dockerfile \
  .
```

### ビルドステップ詳細

1. **ベースイメージ取得**
   - `composer/composer:latest-bin` (マルチステージビルド用)
   - `rockylinux:9`

2. **PHP 8.1インストール**
   - Remi Repositoryから
   - 拡張: gd, intl, mbstring, pgsql, redis等

3. **Apache httpd設定**
   - `/etc/httpd/conf.d/dorapita.conf`
   - ログを stdout/stderr にリダイレクト

4. **アプリコードコピー**
   - `COPY dorapita.com /var/www/dorapita.com`

5. **Composer依存インストール**
   - `composer install -n`

6. **Secret Manager対応**
   - `/mnt/secret-manager/`ディレクトリ作成
   - `.env`をシンボリックリンクで参照

7. **Entrypoint設定**
   - `docker-entrypoint.sh`
   - httpd起動

### ビルド時間

| サービス | ビルド時間（目安） |
|---------|----------------:|
| dorapita-com | 3-5分 |
| cadm-dorapita-com | 3-5分 |
| api-dorapita-com | 2-3分 |

## Artifact Registry構成

### リポジトリ

| リポジトリ | 用途 | サイズ |
|-----------|------|-------:|
| cloud-run-source-deploy | 開発環境全サービス | 53.6 GB |
| api-dorapita-com | api専用（テスト） | 2.8 GB |
| homebrew-containers | img-dorapita-com | 7.8 GB |

### イメージタグ戦略

```
asia-northeast1-docker.pkg.dev/dorapita-core-dev/cloud-run-source-deploy/dorapita/dorapita-com:4a6a231ca
                                                                                              ↑
                                                                                    GitコミットSHA（短縮）
```

- **タグ**: GitコミットSHA（短縮形、8文字）
- **latest**: なし（明示的なバージョン管理）
- **保持**: 無期限（削除ポリシーなし）

## CI/CD設定ファイル

### リポジトリ内

| ファイル | 用途 |
|---------|------|
| `.github/workflows/pr-review.yml` | PRレビュー自動化（Devin AI） |
| `cloudbuild.yaml` | **❌ なし** |

### GCP Cloud Build側

- **Trigger**: GitHubリポジトリと連携
- **設定**: GCPコンソールで管理（YAMLファイルなし）
- **権限**: デフォルトCompute Engine SA使用

## GitHub Actions

### pr-review.yml

```yaml
name: Automated PR Review
on:
  issue_comment:
    types: [created]

jobs:
  review-pr:
    if: github.event.issue.pull_request && github.event.comment.body == 'Devin review'
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
      - name: Get PR files
      - name: Create Devin Review Session  # Devin AIによる自動レビュー
```

**用途**: PRにコメント「Devin review」でDevin AIがコードレビューを実行

## ロールバック

### Cloud Run（開発環境）

```bash
# 前のリビジョンにロールバック
gcloud run services update-traffic dorapita-com \
  --to-revisions=dorapita-com-00042-abc=100 \
  --project=dorapita-core-dev \
  --region=asia-northeast1
```

### 本番VM

```bash
# VMにSSH
gcloud compute ssh web-120011 --project=dorapita-core --zone=asia-northeast1-b

# Gitで前のコミットに戻す
cd /var/www/dorapita.com
git reset --hard <COMMIT_SHA>
composer install --no-dev
sudo systemctl restart httpd
```

## モニタリング

### ビルドステータス

```bash
# 最近のビルド一覧
gcloud builds list --project=dorapita-core-dev --limit=10

# ビルドログ
gcloud builds log <BUILD_ID> --project=dorapita-core-dev

# ビルド失敗の確認
gcloud builds list --project=dorapita-core-dev --filter="status=FAILURE"
```

### デプロイ履歴

```bash
# Cloud Runリビジョン一覧
gcloud run revisions list \
  --service=dorapita-com \
  --project=dorapita-core-dev \
  --region=asia-northeast1

# 現在のトラフィック割り当て
gcloud run services describe dorapita-com \
  --project=dorapita-core-dev \
  --region=asia-northeast1 \
  --format="value(status.traffic)"
```

## トラブルシューティング

### ビルド失敗

#### 1. Dockerビルドエラー

```bash
# ログから原因を特定
gcloud builds log <BUILD_ID> --project=dorapita-core-dev

# よくある原因:
# - Dockerfileの構文エラー
# - composer.jsonの依存解決失敗
# - ベースイメージのpullエラー
```

#### 2. デプロイエラー

```bash
# Cloud Runサービスのログ確認
gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=dorapita-com" \
  --project=dorapita-core-dev \
  --limit=50
```

#### 3. Secret Managerアクセスエラー

Cloud RunサービスアカウントにSecret Manager権限が必要:

```bash
gcloud projects add-iam-policy-binding dorapita-core-dev \
  --member="serviceAccount:<SA_EMAIL>" \
  --role="roles/secretmanager.secretAccessor"
```

### ビルド時間が長い

- `--no-cache`オプションを外す（キャッシュ利用）
- マルチステージビルド最適化
- 不要な依存パッケージ削減

## セキュリティ

### シークレット管理

- **DB認証情報**: Secret Manager
- **API Key**: Secret Manager
- **環境変数**: Cloud Run環境変数 or Secret Mount

### イメージスキャン

```bash
# Artifact Registryで脆弱性スキャン有効化
gcloud artifacts vulnerabilities enable \
  --project=dorapita-core-dev
```

## 改善提案

### 1. cloudbuild.yamlをリポジトリに追加

現状GCPコンソールで管理 → GitOpsへ移行

```yaml
# cloudbuild.yaml (提案)
steps:
- name: 'gcr.io/cloud-builders/docker'
  args: ['build', '-t', '$_IMAGE_NAME:$COMMIT_SHA', '-f', 'Dockerfile', '.']
- name: 'gcr.io/cloud-builders/docker'
  args: ['push', '$_IMAGE_NAME:$COMMIT_SHA']
- name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
  entrypoint: gcloud
  args: ['run', 'services', 'update', '$_SERVICE_NAME', '--image=$_IMAGE_NAME:$COMMIT_SHA']
```

### 2. 本番環境のCloud Run化

VM → Cloud Runへ移行でCI/CD統一

### 3. ステージング環境追加

dev → staging → production

### 4. 自動テスト追加

```yaml
# ビルド前にテスト実行
steps:
- name: 'composer'
  args: ['install']
- name: 'php'
  args: ['vendor/bin/phpunit']
- name: 'gcr.io/cloud-builders/docker'
  args: ['build', ...]
```

### 5. Blue-Greenデプロイ

```bash
# トラフィックを段階的に移行
gcloud run services update-traffic dorapita-com \
  --to-revisions=NEW_REVISION=10,OLD_REVISION=90
```

## 関連ドキュメント

- [ローカル開発環境](LOCAL_DEVELOPMENT.md)
- [アプリケーション構造](README.md)
- [GCPインフラ詳細](../gcp_inventory_reports/dorapita-core-dev-inventory.md)
