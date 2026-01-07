# Phase 1: DevContainer環境構築

## 目標

**VS Code DevContainerでdorapita_codeを起動し、コード修正ができる状態にする**

---

## 現状分析

### 存在するもの

| ファイル | 状態 | 備考 |
|----------|------|------|
| `docker-compose.yml` | ✅ 存在 | DBポート公開済み、全サービス定義あり |
| `Dockerfile*` | ✅ 存在 | 各サービス用（dorapita, cadm, kanri, dora-pt, legacy） |
| `.env.example` | ✅ 存在 | Docker用DB接続設定あり |
| `phpunit.xml.dist` | ✅ 存在 | CakePHP標準設定 |

### 不足するもの

| ファイル | 必要な作業 |
|----------|------------|
| `.devcontainer/` | 新規作成（各サービス用） |
| `.env` | `.env.example`からコピー＆ローカル用調整 |
| 開発用エントリーポイント | Dockerfileの調整 or devcontainer.jsonで対応 |

---

## アーキテクチャ決定

### Multiple DevContainer Definitions

```
dorapita_code/
├── .devcontainer/
│   ├── web-dorapita/
│   │   └── devcontainer.json    # dorapita.com用
│   ├── web-cadm/
│   │   └── devcontainer.json    # cadm.dorapita.com用
│   ├── web-kanri/
│   │   └── devcontainer.json    # kanri.dorapita.com用
│   ├── web-dora-pt/
│   │   └── devcontainer.json    # dora-pt.jp用
│   └── web-legacy/
│       └── devcontainer.json    # legacy.dorapita.com用（CakePHP 1.3）
```

**理由**: 各サービスが異なるディレクトリ・設定を持つため、個別のDevContainer定義が必要

### 開発時の動作モード

**現状の問題**: Dockerfileは`httpd -D FOREGROUND`で終端。開発時は対話的にコマンド実行したい。

**解決策（2つの選択肢）**:

| 選択肢 | 方法 | Pros | Cons |
|--------|------|------|------|
| A. `overrideCommand` | devcontainer.jsonで`"overrideCommand": true`＋`sleep infinity` | Dockerfile変更不要 | Webサーバーは手動起動 |
| B. 別Dockerfile | `Dockerfile.dev`を作成し開発用コマンド定義 | 明確な分離 | ファイル増加 |

**推奨: 選択肢A**（シンプルさ優先）

---

## 実装計画

### Step 1: .devcontainer/web-dorapita/devcontainer.json 作成

**最小構成**:

```json
{
  "name": "Dorapita Main Site (dorapita.com)",
  "dockerComposeFile": ["../../docker-compose.yml"],
  "service": "web-dorapita",
  "workspaceFolder": "/var/www/dorapita.com",
  "shutdownAction": "none",
  "overrideCommand": true,
  "customizations": {
    "vscode": {
      "settings": {
        "php.validate.executablePath": "/usr/bin/php"
      },
      "extensions": [
        "bmewburn.vscode-intelephense-client",
        "ms-azuretools.vscode-docker",
        "xdebug.php-debug",
        "esbenp.prettier-vscode"
      ]
    }
  },
  "postAttachCommand": "php -v && composer --version"
}
```

**パラメータ説明**:

| パラメータ | 値 | 説明 |
|-----------|-----|------|
| `dockerComposeFile` | `../../docker-compose.yml` | 既存のdocker-composeを利用 |
| `service` | `web-dorapita` | アタッチするサービス名 |
| `workspaceFolder` | `/var/www/dorapita.com` | VS Codeで開くディレクトリ |
| `shutdownAction` | `none` | VS Code閉じてもコンテナ維持 |
| `overrideCommand` | `true` | Dockerfileのcommandを無効化（sleep infinity相当） |
| `postAttachCommand` | PHP/Composer確認 | アタッチ後の動作確認 |

### Step 2: docker-compose.yml の調整（必要に応じて）

**現状で動作する可能性が高い**が、以下の点を確認:

1. **ボリュームマウント**: `./dorapita.com:/var/www/dorapita.com` ✅ 設定済み
2. **depends_on**: pgsql, redis ✅ 設定済み
3. **ポート公開**: 8080:8080 ✅ 設定済み

**追加検討**:

```yaml
# 開発時の追加環境変数（必要に応じて）
services:
  web-dorapita:
    environment:
      - PHP_OPCACHE_ENABLE=0  # 開発時OPCache無効化
```

### Step 3: .env ファイルの準備

```bash
# dorapita_code/dorapita.com/config/ にて
cp .env.example .env
```

**ローカル用の修正箇所**:

```bash
# Before (.env.example)
export DATABASE_URL='postgres://dorauser2022:<※.secret参照>@pgsql/dorapita?...'
export REDIS_HOST = 'host.docker.internal'

# After (.env - Docker内からアクセスする場合はそのまま)
# ホスト名はdocker-composeのサービス名を使用
export DATABASE_URL='postgres://dorauser2022:<※.secret参照>@pgsql/dorapita?encoding=utf8mb4&timezone=UTC&cacheMetadata=true'
export REDIS_HOST = 'redis'
export REDIS_PASSWORD = ''  # docker-compose.ymlのREDIS_PASSWORDと合わせる
```

### Step 4: 動作確認

DevContainerアタッチ後に以下を実行:

```bash
# 1. PHP/Composerバージョン確認
php -v
composer --version

# 2. 依存パッケージインストール（必要な場合）
composer install

# 3. DB接続確認
bin/cake migrations status

# 4. Webサーバー起動（手動）
httpd -D FOREGROUND &
# または
php -S 0.0.0.0:8080 -t webroot/

# 5. ブラウザで確認
# http://localhost:8080
```

---

## 作成するファイル一覧

### Phase 1 成果物

| ファイル | 説明 | 優先度 |
|----------|------|--------|
| `.devcontainer/web-dorapita/devcontainer.json` | メインサイト用DevContainer設定 | **最高** |
| `dorapita.com/config/.env` | ローカル環境変数 | **最高** |
| `.devcontainer/web-cadm/devcontainer.json` | 顧客管理画面用 | 高 |
| `.devcontainer/web-kanri/devcontainer.json` | 管理画面用 | 高 |
| `.devcontainer/web-dora-pt/devcontainer.json` | dora-pt.jp用 | 中 |
| `.devcontainer/web-legacy/devcontainer.json` | 旧システム用（CakePHP 1.3） | 低 |

### 各devcontainer.jsonの差分

| サービス | workspaceFolder | 備考 |
|----------|-----------------|------|
| web-dorapita | `/var/www/dorapita.com` | CakePHP 4.4, PostgreSQL |
| web-cadm | `/var/www/cadm.dorapita.com` | CakePHP 4.5, MySQL |
| web-kanri | `/var/www/kanri.dorapita.com` | CakePHP 4.5, MySQL |
| web-dora-pt | `/var/www/dora-pt.jp` | CakePHP 4.5, MySQL |
| web-legacy | `/var/www/cake_1_3` | CakePHP 1.3, MySQL（触らない方針だが参照用） |

---

## 検証チェックリスト

### 必須項目

- [ ] VS Code で `.devcontainer/web-dorapita/` を選択してコンテナにアタッチできる
- [ ] アタッチ後、`/var/www/dorapita.com` がワークスペースとして開く
- [ ] `php -v` でPHP 8.1が表示される
- [ ] `composer --version` でComposerが利用可能
- [ ] `bin/cake` コマンドが実行できる
- [ ] `bin/cake migrations status` でDB接続が成功する
- [ ] ソースコードを編集して保存できる（volume mount動作）

### 追加確認

- [ ] ホストのブラウザから `http://localhost:8080` でアクセスできる
- [ ] 他サービス（redis, mysql, pgsql）がDocker内部ネットワーク経由でアクセスできる
- [ ] PHPエラーがログに出力される

---

## トラブルシューティング

### 問題1: コンテナがすぐに終了する

**原因**: `overrideCommand: true` が効いていない、または Dockerfile の CMD が実行されている

**解決策**:

```json
// devcontainer.json に追加
"overrideCommand": true,
"runArgs": ["--init"]
```

---

### 問題2: DB接続エラー

**原因**: `.env` の DATABASE_URL ホスト名が間違っている

**解決策**:

```bash
# Docker内からはサービス名でアクセス
export DATABASE_URL='postgres://dorauser2022:<※.secret参照>@pgsql/dorapita...'
#                                                    ^^^^^ サービス名
```

---

### 問題3: Permission denied

**原因**: ボリュームマウントのファイル権限

**解決策**:

```json
// devcontainer.json に追加
"remoteUser": "root"  // または適切なユーザー
```

または

```bash
# コンテナ内で
chown -R apache:apache /var/www/dorapita.com
```

---

### 問題4: composer install で memory エラー

**原因**: PHPのメモリ制限

**解決策**:

```bash
php -d memory_limit=-1 /usr/local/bin/composer install
```

---

## 次のフェーズへの接続

Phase 1 完了後、以下が可能になる:

1. **Phase 2-3 (Fixture整備)**: DevContainer内で `bin/cake bake fixture` や SQLクエリ実行が可能
2. **Phase 4 (Unit Test)**: DevContainer内で `vendor/bin/phpunit` 実行が可能
3. **コード修正**: リアルタイムでコード編集・動作確認

---

## 参照ドキュメント

- [VS Code Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers)
- [devcontainer.json reference](https://containers.dev/implementors/json_reference/)
- [local-development-architecture.md](./local-development-architecture.md) - 全体アーキテクチャ設計

---

*作成日: 2025-12-26*
*最終更新: 2025-12-26*
