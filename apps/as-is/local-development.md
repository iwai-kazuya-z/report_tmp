# ローカル開発環境セットアップガイド

## 概要

ドラピタプロジェクトはDocker Composeを使用したマルチコンテナ環境でローカル開発を行います。

## 前提条件

| ソフトウェア | 推奨バージョン |
|-------------|--------------|
| Docker | 最新版 |
| Docker Compose | v3.9以上 |
| Git | 最新版 |

## ディレクトリ構成

```
dorapita_code/
├── docker-compose.yml         # 開発環境定義
├── Dockerfile                 # dorapita.com用
├── Dockerfile_cadm.dorapita.com
├── Dockerfile_kanri.dorapita.com
├── Dockerfile_legacy.dorapita.com
├── Dockerfile_dora-pt.jp
├── docker-entrypoint.sh       # コンテナ起動スクリプト
├── docker-entrypoint-legacy.sh
├── dorapita.com/              # アプリコード
├── cadm.dorapita.com/
├── kanri.dorapita.com/
└── etc/httpd/conf.d/          # Apache設定
```

## セットアップ手順

### 1. リポジトリクローン

```bash
git clone https://github.com/ZIGExN/dorapita.git
cd dorapita
```

### 2. 環境変数設定

```bash
# dorapita.com
cp dorapita.com/config/.env.example dorapita.com/config/.env

# 必要に応じてRedis設定を編集
export REDIS_HOST=''
export REDIS_PORT=''
export REDIS_PASSWORD=''
```

### 3. コンテナビルド & 起動

```bash
docker compose build
docker compose up
```

### 4. データベース初期化

#### PostgreSQL

```bash
# PostgreSQLコンテナに接続
docker compose exec pgsql psql -U postgres

# データベース作成
CREATE DATABASE dorapita;
CREATE USER dorauser2022 WITH PASSWORD '<※.secret参照>';
\q
```

#### MySQL

```bash
# MySQLコンテナに接続
docker compose exec mysql mysql -u root

# データベース作成
CREATE DATABASE dorapita1804;
exit
```

### 5. ダンプデータリストア（オプション）

```bash
# PostgreSQL
gsutil cp gs://dorapita-infra-stg/masking-data-dorapita.zip .
unzip masking-data-dorapita.zip
docker compose exec -T pgsql psql -U postgres dorapita < masking_data_dorapita.backup

# MySQL
gsutil cp gs://dorapita-infra-stg/dump/dorapita1804-2023-11-07.dump.sql .
docker compose exec -T mysql mysql -u root < dorapita1804-2023-11-07.dump.sql

# アップロードファイル
gsutil -m cp -r gs://dorapita-infra-stg/assets/upfiles .
```

## サービス構成

### Webサービス

| サービス | ポート | URL | 備考 |
|---------|-------:|-----|------|
| dorapita.com | 8080 | http://localhost:8080 | メインサイト |
| cadm.dorapita.com | 8888 | http://localhost:8888 | 顧客管理画面 |
| kanri.dorapita.com | 8088 | http://localhost:8088 | 管理画面 |
| legacy.dorapita.com | 9090 | http://localhost:9090 | 旧システム (CakePHP 1.3) |
| dora-pt.jp | 9999 | http://localhost:9999 | 別ドメイン |
| メンテナンスページ | 11111 | http://localhost:11111 | メンテナンス表示 |
| MailHog | 8025 | http://localhost:8025 | メール確認UI |

### データベース・キャッシュ

| サービス | ポート | 認証情報 |
|---------|-------:|---------|
| PostgreSQL | 5432 | postgres / D3arMyFr1end5 |
| MySQL | 3306 | root / (パスワードなし) |
| Redis | 6379 | (環境変数で設定) |
| cloud-sql-proxy | 54321 (PG), 33306 (MySQL) | gcloud認証 |

## Docker Compose構成詳細

### サービス定義

```yaml
services:
  web-dorapita:       # dorapita.com
  web-cadm:           # cadm.dorapita.com
  web-kanri:          # kanri.dorapita.com
  web-legacy:         # legacy.dorapita.com (CakePHP 1.3, PHP 5.6)
  web-dora-pt:        # dora-pt.jp
  web-sorry:          # メンテナンスページ
  pgsql:              # PostgreSQL 10
  mysql:              # MySQL 5.7
  redis:              # Redis 7.2
  cloud-sql-proxy:    # GCP Cloud SQL接続用
  mailhog:            # ローカルメールサーバー
```

### ボリューム

```yaml
volumes:
  mysql:              # MySQL永続化
  pgsql:              # PostgreSQL永続化
  redis_data:         # Redis永続化
```

### ホストマウント

```yaml
./dorapita.com -> /var/www/dorapita.com          # ホットリロード有効
./cadm.dorapita.com -> /var/www/cadm.dorapita.com
./upfiles -> /var/www/*/upfiles                   # アップロードファイル共有
```

## Dockerfileの特徴

### dorapita.com (Dockerfile)

```dockerfile
FROM rockylinux:9
# PHP 8.1 (Remi Repository)
# Apache httpd
# PostgreSQL client
# Redis extension

WORKDIR /var/www/dorapita.com
RUN composer install

# Secret Manager対応
RUN mkdir /mnt/secret-manager/
RUN ln -sf /mnt/secret-manager/shell-environment-variables-dorapita-com \
           /var/www/dorapita.com/config/.env

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["httpd", "-D", "FOREGROUND"]
```

**特徴**:
- RockyLinux 9ベース
- Remi Repoから PHP 8.1をインストール
- Secret Managerマウント対応（Cloud Run用）
- `.env`ファイルをシンボリックリンクで参照

### legacy.dorapita.com (Dockerfile_legacy.dorapita.com)

```dockerfile
FROM rockylinux:9
# PHP 5.6 (Remi Repository)
# CakePHP 1.3.21対応
```

**注意**: PHP 5.6は超レガシー。セキュリティリスクあり。

## Cloud SQL Proxyの使用

ローカルから本番/開発DBに接続する場合:

```bash
# 環境変数設定
export INSTANCE_CONNECTION_NAME=dorapita-core-dev:asia-northeast1:db-120011

# Docker Compose起動
docker compose up cloud-sql-proxy

# 別ターミナルから接続
mysql -h 127.0.0.1 -P 33306 -u <USER> -p
psql -h 127.0.0.1 -p 54321 -U <USER> -d dorapita
```

## トラブルシューティング

### ポート競合

既に使用中のポートがある場合、`docker-compose.yml`のポートマッピングを変更:

```yaml
ports:
  - "8081:8080"  # 8080→8081に変更
```

### データベース接続エラー

```bash
# コンテナが起動しているか確認
docker compose ps

# ログ確認
docker compose logs pgsql
docker compose logs mysql

# コンテナ再起動
docker compose restart pgsql
```

### Composer依存解決エラー

```bash
# コンテナ内でComposer更新
docker compose exec web-dorapita composer update
```

### ファイルパーミッションエラー

```bash
# upfilesディレクトリの権限修正
chmod -R 777 upfiles/
```

## 開発ワークフロー

### 1. コード変更

ホストの`dorapita.com/`等を編集 → 自動的にコンテナに反映（ボリュームマウント）

### 2. 依存パッケージ追加

```bash
# dorapita.com
docker compose exec web-dorapita composer require <package>

# cadm
docker compose exec web-cadm composer require <package>
```

### 3. マイグレーション

```bash
# dorapita.com
docker compose exec web-dorapita bin/cake migrations migrate

# cadm
docker compose exec web-cadm bin/cake migrations migrate
```

### 4. キャッシュクリア

```bash
docker compose exec web-dorapita bin/cake cache clear_all
```

### 5. ログ確認

```bash
# リアルタイムログ
docker compose logs -f web-dorapita

# Apacheエラーログ
docker compose exec web-dorapita tail -f /var/log/httpd/error_log
```

## 次のステップ

- [CI/CDパイプライン](CI_CD_PIPELINE.md) - デプロイメント方法
- [アプリケーション構造](README.md) - サービス概要
- [Cloud SQL接続ガイド](../gcp_inventory_reports/dorapita-inventory.md#cloud-sql接続ガイド)

## よくある質問

### Q: 本番データを使いたい
A: Cloud SQL Proxyを使用してGCP本番DBに直接接続できますが、推奨しません。ダンプを取得して使用してください。

### Q: M1/M2 Macで動作しない
A: `docker-compose.yml`で`platform: linux/amd64`を有効化してください（コメントアウト解除）

### Q: メール送信をテストしたい
A: MailHog (http://localhost:8025) で送信メールを確認できます
