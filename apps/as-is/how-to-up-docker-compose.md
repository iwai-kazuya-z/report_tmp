# Docker Compose起動ガイド

## 概要

dorapita_codeのDocker Compose環境を起動するための手順と、よくある問題の解決方法をまとめています。

---

## 前提条件

- Docker Desktop がインストール済み
- dorapita_code サブモジュールがクローン済み

---

## クイックスタート（全手順）

```bash
cd /path/to/dorapita_code

# 1. .envファイル作成
cp dorapita.com/config/.env.example dorapita.com/config/.env
cp cadm.dorapita.com/config/.env.example cadm.dorapita.com/config/.env
cp kanri.dorapita.com/config/.env.example kanri.dorapita.com/config/.env

# 2. .env修正（各ファイルで以下を変更）
# dorapita.com: DATABASE_URL, REDIS_HOST='redis', REDIS_PASSWORD=''
# cadm/kanri: SECURITY_COOKIE_KEY="任意の32文字", REDIS_HOST='redis', REDIS_PASSWORD=''

# 3. コンテナ起動
docker compose up -d web-dorapita web-cadm web-kanri web-dora-pt pgsql mysql redis mailhog

# 4. composer install
docker exec dorapita_code-web-dorapita-1 composer install -d /var/www/dorapita.com

# 5. DBスキーマ適用
docker exec -i dorapita_code-pgsql-1 psql -U postgres -d postgres < dorapita.com/config/schema/init.sql
docker exec dorapita_code-mysql-1 mysql -u root -e "CREATE DATABASE IF NOT EXISTS dorapita1804;"
docker exec -i dorapita_code-mysql-1 mysql -u root dorapita1804 < dorapita.com/config/schema/mysql-init.sql

# 6. マイグレーションをマーク済みに
docker exec dorapita_code-web-dorapita-1 bash -c "cd /var/www/dorapita.com && \
  php bin/cake.php migrations mark_migrated 20240508054115 && \
  php bin/cake.php migrations mark_migrated 20240605084935 && \
  php bin/cake.php migrations mark_migrated 20240620090310 && \
  php bin/cake.php migrations mark_migrated 20240710000000 && \
  php bin/cake.php migrations mark_migrated 20250320043616"

# 7. 動作確認
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/   # → 200
curl -s -o /dev/null -w "%{http_code}" http://localhost:8888/   # → 302
curl -s -o /dev/null -w "%{http_code}" http://localhost:8088/   # → 302
curl -s -o /dev/null -w "%{http_code}" http://localhost:9999/   # → 200
```

---

## 起動手順（詳細）

### 1. 基本的な起動コマンド

```bash
# dorapita_codeディレクトリに移動
cd /path/to/dorapita_code

# === 最小構成（dorapita.comのみ） ===
docker compose up -d web-dorapita pgsql mysql redis mailhog

# === 全アプリケーション起動（推奨） ===
docker compose up -d \
  web-dorapita web-cadm web-kanri web-dora-pt \
  pgsql mysql redis mailhog

# フォアグラウンドで起動する場合（ログ確認用）
docker compose up web-dorapita pgsql mysql redis mailhog
```

**注意**: dorapita.comはPostgreSQLとMySQL両方を使用するハイブリッド構成のため、両方のDBコンテナが必要。

### 2. .envファイルの作成（必須）

各アプリケーションで`.env`ファイルを作成：

```bash
# dorapita.com
cp dorapita.com/config/.env.example dorapita.com/config/.env

# cadm.dorapita.com
cp cadm.dorapita.com/config/.env.example cadm.dorapita.com/config/.env

# kanri.dorapita.com
cp kanri.dorapita.com/config/.env.example kanri.dorapita.com/config/.env

# dora-pt.jp（.envなしでも起動可能）
```

### 3. .envの修正（必須）

#### dorapita.com

```bash
# Database: PostgreSQLユーザー/パスワードをdocker-compose.ymlに合わせる
export DATABASE_URL='postgres://postgres:D3arMyFr1end5@pgsql/postgres?encoding=UTF8&timezone=UTC&cacheMetadata=true'

# Redis: ホスト名とパスワードを修正
export REDIS_HOST = 'redis'
export REDIS_PASSWORD = ''
```

#### cadm.dorapita.com / kanri.dorapita.com

```bash
# Cookie暗号化キー（必須！空だとエラーになる）
export SECURITY_COOKIE_KEY="local_dev_cookie_key_32chars!!"

# Redis設定
export REDIS_HOST = 'redis'
export REDIS_PASSWORD = ''
```

**注意**: `SECURITY_COOKIE_KEY`が空のままだと以下のエラーが発生：
```
TypeError: Argument 2 passed to EncryptedCookieMiddleware::__construct() must be of the type string, null given
```

### 4. composer install（必須）

ボリュームマウントでホスト側のvendor/が使われるため、コンテナ内で実行：

```bash
docker exec dorapita_code-web-dorapita-1 composer install -d /var/www/dorapita.com
```

### 5. 動作確認

```bash
# HTTPステータス確認
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/

# CakePHP CLI確認
docker exec dorapita_code-web-dorapita-1 php /var/www/dorapita.com/bin/cake.php
```

---

## サービス一覧とポート

| サービス | ポート | 用途 | 使用DB | 動作確認 |
|----------|--------|------|--------|----------|
| web-dorapita | 8080 | メインサイト (dorapita.com) | PostgreSQL + MySQL | ✅ HTTP 200 |
| web-cadm | 8888 | 顧客管理画面 | MySQL | ✅ HTTP 302→200 |
| web-kanri | 8088 | 管理画面 | MySQL | ✅ HTTP 302→200 |
| web-dora-pt | 9999 | dora-pt.jp | MySQL | ✅ HTTP 200 |
| web-legacy | 9090 | 旧システム | MySQL | 未確認 |
| pgsql | 5432 | PostgreSQL 10 | - | - |
| mysql | 3306 | MySQL 5.7 | - | - |
| redis | 6379 | Redis 7.2 | - | - |
| mailhog | 1025/8025 | メールテスト | - | - |

**備考**:
- cadm/kanriは認証必須のため、ルートアクセスでログインページ（302）へリダイレクト
- 302リダイレクト後のログインページで200が返れば正常

---

## よくある問題と解決方法

### 問題1: HTTP 500エラー - vendor/が不完全

**症状**:
```
PHP Fatal error: Interface "Authentication\AuthenticationServiceProviderInterface" not found
```

**原因**:
Dockerfileで`composer install`が実行されるが、ボリュームマウントでホスト側のvendor/が上書きされる。

**解決策**:
```bash
docker exec dorapita_code-web-dorapita-1 composer install -d /var/www/dorapita.com
```

---

### 問題2: HTTP 500エラー - DB認証失敗

**症状**:
```
FATAL: password authentication failed for user "dorauser2022"
```

**原因**:
`.env`のDB接続情報がdocker-compose.ymlのPostgreSQL設定と一致していない。

**解決策**:
`.env`を以下のように修正：
```bash
export DATABASE_URL='postgres://postgres:D3arMyFr1end5@pgsql/postgres?encoding=UTF8&timezone=UTC&cacheMetadata=true'
```

**docker-compose.ymlの設定**:
```yaml
pgsql:
  image: postgres:10-alpine
  environment:
    POSTGRES_PASSWORD: D3arMyFr1end5  # ← このパスワードを使う
```

---

### 問題3: HTTP 500エラー - エンコーディングエラー

**症状**:
```
invalid value for parameter "client_encoding": "utf8mb4"
```

**原因**:
`utf8mb4`はMySQL用のエンコーディング。PostgreSQLでは`UTF8`を使う。

**解決策**:
```bash
# Before
encoding=utf8mb4

# After
encoding=UTF8
```

---

### 問題4: HTTP 500エラー - テーブルが存在しない

**症状**:
```
Cannot describe recruits. It has 0 columns.
```

**原因**:
PostgreSQLコンテナは空のDBで起動する。テーブルが存在しない。

**解決策**:
初期スキーマを適用する必要がある。詳細は「初期スキーマのセットアップ」セクション参照。

---

### 問題5: mailhog接続エラー

**症状**:
```
php_network_getaddresses: getaddrinfo for mailhog failed
```

**原因**:
mailhogサービスが起動していない。エラー発生時にメール通知しようとして失敗。

**解決策**:
```bash
# mailhogも一緒に起動
docker compose up web-dorapita pgsql redis mailhog
```

---

### 問題6: Redis認証エラー

**症状**:
```
NOAUTH Authentication required
```

**原因**:
`.env`のREDIS_PASSWORDがdocker-compose.ymlの設定と一致していない。

**解決策**:
```bash
# .envを修正
export REDIS_HOST = 'redis'
export REDIS_PASSWORD = ''  # 空にする
```

---

### 問題7: Cookie暗号化キーエラー（cadm/kanri）

**症状**:
```
TypeError: Argument 2 passed to Cake\Http\Middleware\EncryptedCookieMiddleware::__construct() must be of the type string, null given
```

**原因**:
`.env`の`SECURITY_COOKIE_KEY`が空または未設定。

**解決策**:
```bash
# .envに32文字以上のキーを設定
export SECURITY_COOKIE_KEY="local_dev_cookie_key_32chars!!"
```

---

### 問題8: MySQLテーブル不足（dorapita.com）

**症状**:
```
Table 'dorapita1804.information' doesn't exist
```

**原因**:
MySQLコンテナは空のDBで起動する。dorapita.comはMySQLのテーブルも参照する。

**解決策**:
```bash
# MySQLスキーマを適用
docker exec dorapita_code-mysql-1 mysql -u root -e "CREATE DATABASE IF NOT EXISTS dorapita1804;"
docker exec -i dorapita_code-mysql-1 mysql -u root dorapita1804 < dorapita.com/config/schema/mysql-init.sql
```

---

## 初期スキーマのセットアップ

### 現状の問題

- マイグレーションファイルは2024年5月以降の差分変更のみ
- 初期スキーマ（CREATE TABLE）はコードに存在しない
- **PostgreSQL**: スキーマ取得済み（`config/schema/init.sql`）
- **MySQL**: スキーマ未取得（`information`テーブル等が必要）

### データベース構成

dorapita.comはハイブリッドDB構成を採用:

| DB | 用途 | テーブル例 |
|----|------|-----------|
| PostgreSQL 10 | 求人データ、エントリー | recruits, entries, companies |
| MySQL 5.7 | システム情報、連携データ | information 等 |

### PostgreSQLスキーマのセットアップ（実施済み）

```bash
# 1. init.sqlをローカルDBに適用
docker exec -i dorapita_code-pgsql-1 psql -U postgres -d postgres < dorapita.com/config/schema/init.sql

# 2. マイグレーションを適用済みにマーク（個別に実行）
docker exec dorapita_code-web-dorapita-1 bash -c "
cd /var/www/dorapita.com && \
php bin/cake.php migrations mark_migrated 20240508054115 && \
php bin/cake.php migrations mark_migrated 20240605084935 && \
php bin/cake.php migrations mark_migrated 20240620090310 && \
php bin/cake.php migrations mark_migrated 20240710000000 && \
php bin/cake.php migrations mark_migrated 20250320043616
"

# 3. ステータス確認
docker exec dorapita_code-web-dorapita-1 php /var/www/dorapita.com/bin/cake.php migrations status
```

### MySQLスキーマのセットアップ（実施済み）

```bash
# 1. MySQLコンテナを起動
docker compose -f /path/to/dorapita_code/docker-compose.yml up -d mysql

# 2. データベースを作成
docker exec dorapita_code-mysql-1 mysql -u root -e "CREATE DATABASE IF NOT EXISTS dorapita1804;"

# 3. スキーマを適用
docker exec -i dorapita_code-mysql-1 mysql -u root dorapita1804 < dorapita.com/config/schema/mysql-init.sql

# 4. 確認
docker exec dorapita_code-mysql-1 mysql -u root -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='dorapita1804';"
# → 159 テーブル/ビュー
```

**注意**: STGのデータベース名は `dorapita1804_db` だが、ローカルでは `dorapita1804` を使用

### STG DBからスキーマを取得する手順（参考）

```bash
# 1. Cloud SQL Proxyを起動
pkill -f cloud-sql-proxy 2>/dev/null
cloud-sql-proxy "dorapita-core-dev:asia-northeast1:pg-120011" --port=35432 --gcloud-auth &

# 2. PostgreSQLスキーマをダンプ（権限問題でpg_catalog経由が必要）
export PATH="/opt/homebrew/opt/postgresql@16/bin:$PATH"
PGPASSWORD=<※.secret参照> psql -h 127.0.0.1 -p 35432 -U dorauser2022 -d dorapita -c "
SELECT 'CREATE TABLE IF NOT EXISTS ' || tablename || ' (' ||
  string_agg(column_name || ' ' || data_type ||
    CASE WHEN character_maximum_length IS NOT NULL
         THEN '(' || character_maximum_length || ')' ELSE '' END ||
    CASE WHEN is_nullable = 'NO' THEN ' NOT NULL' ELSE '' END ||
    CASE WHEN column_default IS NOT NULL
         THEN ' DEFAULT ' || column_default ELSE '' END, ', ') ||
  ');'
FROM information_schema.columns
WHERE table_schema = 'public'
GROUP BY tablename;
"

# 3. Proxyを停止
pkill -f cloud-sql-proxy
```

---

## キャッシュクリア

設定変更後はキャッシュクリアが必要な場合がある：

```bash
docker exec dorapita_code-web-dorapita-1 php /var/www/dorapita.com/bin/cake.php cache clear_all
```

---

## コンテナ操作

```bash
# コンテナに入る
docker exec -it dorapita_code-web-dorapita-1 bash

# ログ確認
docker logs dorapita_code-web-dorapita-1

# コンテナ再起動
docker restart dorapita_code-web-dorapita-1

# 全停止
docker compose down

# ボリューム含めて全削除（DBデータも消える）
docker compose down -v
```

---

## DB接続確認

```bash
# PostgreSQL
docker exec -it dorapita_code-pgsql-1 psql -U postgres

# MySQL
docker exec -it dorapita_code-mysql-1 mysql -u root
```

---

## トラブルシューティングの流れ

```
HTTP 500エラー発生
       ↓
1. エラーログ確認
   docker exec dorapita_code-web-dorapita-1 tail -50 /var/www/dorapita.com/logs/error.log
       ↓
2. CLIで動作確認
   docker exec dorapita_code-web-dorapita-1 php /var/www/dorapita.com/bin/cake.php
       ↓
3. 原因に応じて対処
   - vendor/不完全 → composer install
   - DB接続エラー → .env修正
   - テーブルなし → スキーマ適用
```

---

## 参照ドキュメント

- [phase1-devcontainer-setup.md](./phase1-devcontainer-setup.md) - DevContainer設定
- [local-development-architecture.md](./local-development-architecture.md) - 全体アーキテクチャ

---

*作成日: 2025-12-26*
*最終更新: 2025-12-26*

---

## 完了した作業

### スキーマファイル

| 項目 | 状態 | ファイル |
|------|------|----------|
| PostgreSQLスキーマ | ✅ 完了 | `dorapita.com/config/schema/init.sql` (394行, 20テーブル) |
| MySQLスキーマ | ✅ 完了 | `dorapita.com/config/schema/mysql-init.sql` (3,590行, 159テーブル/ビュー) |

### アプリケーション動作確認

| アプリケーション | ポート | 状態 | 備考 |
|-----------------|--------|------|------|
| dorapita.com | 8080 | ✅ HTTP 200 | メインサイト |
| cadm.dorapita.com | 8888 | ✅ HTTP 302→200 | ログインページへリダイレクト |
| kanri.dorapita.com | 8088 | ✅ HTTP 302→200 | ログインページへリダイレクト |
| dora-pt.jp | 9999 | ✅ HTTP 200 | - |

### .env設定

| アプリケーション | 主な設定項目 |
|-----------------|-------------|
| dorapita.com | DATABASE_URL, REDIS_HOST, REDIS_PASSWORD |
| cadm.dorapita.com | SECURITY_COOKIE_KEY, REDIS_HOST, REDIS_PASSWORD |
| kanri.dorapita.com | SECURITY_COOKIE_KEY, REDIS_HOST, REDIS_PASSWORD |
| dora-pt.jp | （デフォルトで動作） |

## 次のステップ（TODO）

1. **テストデータ（Fixtures）の整備**
   - STGからサンプルデータを抽出
   - テスト用のシードデータを作成

2. **DevContainer設定の作成**
   - `.devcontainer/devcontainer.json`の作成
   - VS Code拡張機能の設定

3. **web-legacy（旧システム）の動作確認**
   - CakePHP 1.3環境の確認
