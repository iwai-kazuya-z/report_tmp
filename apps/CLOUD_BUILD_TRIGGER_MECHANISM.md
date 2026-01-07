# Cloud Build Trigger の仕組み

## 概要

Cloud Build TriggerはGitHubと連携し、特定のブランチへのpushを検知して自動的にビルド・デプロイを実行します。

## アーキテクチャ図

```
┌─────────────────────────────────────────────────────────────────────┐
│                          GitHub (ZIGExN/dorapita)                    │
│                                                                       │
│  1. Developer pushes to "dev.dorapita.com" branch                   │
│     git push origin dev.dorapita.com                                 │
│                                                                       │
│                              │                                        │
│                              ▼                                        │
│  2. GitHub detects push event                                        │
│     - Branch: dev.dorapita.com                                       │
│     - Commit SHA: 4a6a231ca                                          │
│     - Author, Message, Files changed...                              │
└──────────────────────────────┬───────────────────────────────────────┘
                               │
                               │ Webhook (GitHub App経由)
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Google Cloud Build API                          │
│                                                                       │
│  3. Cloud Build receives webhook                                     │
│     POST https://cloudbuild.googleapis.com/v1/webhook               │
│     {                                                                 │
│       "ref": "refs/heads/dev.dorapita.com",                          │
│       "commits": [{...}],                                            │
│       "repository": {...}                                            │
│     }                                                                 │
│                                                                       │
│                              │                                        │
│                              ▼                                        │
│  4. Trigger Matching Engine                                          │
│     ┌─────────────────────────────────────┐                         │
│     │ Trigger: rmgpgab-dorapita-com-...   │                         │
│     │ Condition:                           │                         │
│     │   - Repo: ZIGExN/dorapita            │                         │
│     │   - Branch: ^dev.dorapita.com$       │ ← 正規表現マッチング    │
│     │   - Event: push                      │                         │
│     └─────────────────────────────────────┘                         │
│                                                                       │
│     条件マッチ！ → ビルド起動                                          │
│                                                                       │
│                              │                                        │
│                              ▼                                        │
│  5. Build Execution                                                  │
│     - Create Build ID: 10a4fb73-62fd-4211-8362-bc1492841791         │
│     - Allocate Build Worker                                          │
│     - Fetch Source from GitHub                                       │
│     - Run Build Steps                                                │
└──────────────────────────────┬───────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                       Build Worker (VM)                              │
│                                                                       │
│  6. Fetch Source                                                     │
│     gsutil cp gs://dorapita-core-dev_cloudbuild/source/...          │
│     ↓ 解凍してワークスペースに配置                                     │
│                                                                       │
│  7. Execute Build Steps                                              │
│     Step #0: docker build                                            │
│     Step #1: docker push                                             │
│     Step #2: gcloud run services update                              │
│                                                                       │
└───────────────────────────────────────────────────────────────────────┘
```

## 詳細フロー

### 1. GitHub App連携

Cloud BuildはGitHub Appとして登録されており、リポジトリへのアクセス権限を持っています。

#### GitHub App設定内容

| 項目 | 設定値 |
|------|--------|
| App名 | Google Cloud Build |
| 権限 | - Repository contents: Read<br>- Repository metadata: Read<br>- Pull requests: Read & Write<br>- Commit statuses: Read & Write |
| Webhook URL | `https://cloudbuild.googleapis.com/v1/projects/PROJECT_ID/triggers/TRIGGER_ID:webhook` |
| Webhook Events | - push<br>- pull_request<br>- release |

#### インストール先

- **Organization**: ZIGExN
- **Repository**: dorapita
- **インストール日**: トリガー作成時に自動

### 2. Webhookイベント送信

GitHubでpushが発生すると、以下のWebhookペイロードがCloud Buildに送信されます。

```json
POST https://cloudbuild.googleapis.com/v1/webhook
Content-Type: application/json
X-GitHub-Event: push
X-Hub-Signature-256: sha256=...

{
  "ref": "refs/heads/dev.dorapita.com",
  "before": "abc123...",
  "after": "4a6a231ca...",
  "repository": {
    "id": 123456789,
    "name": "dorapita",
    "full_name": "ZIGExN/dorapita",
    "owner": {
      "name": "ZIGExN",
      "login": "ZIGExN"
    }
  },
  "commits": [
    {
      "id": "4a6a231ca...",
      "message": "feat: 新機能追加",
      "author": {...},
      "timestamp": "2025-12-16T08:21:15Z",
      "added": [],
      "modified": ["dorapita.com/src/Controller/..."],
      "removed": []
    }
  ],
  "pusher": {...},
  "head_commit": {...}
}
```

### 3. Trigger Matching（トリガー判定）

Cloud BuildはWebhookを受信すると、登録されている全Triggerの条件をチェックします。

#### マッチング条件

```yaml
github:
  owner: ZIGExN           # リポジトリオーナー
  name: dorapita          # リポジトリ名
  push:
    branch: ^dev.dorapita.com$   # 正規表現マッチング
```

#### マッチングロジック

```python
# 擬似コード
def should_trigger_build(webhook_event, trigger_config):
    # 1. リポジトリチェック
    if webhook_event.repository.full_name != f"{trigger_config.github.owner}/{trigger_config.github.name}":
        return False

    # 2. イベントタイプチェック
    if webhook_event.event_type != "push":
        return False

    # 3. ブランチチェック（正規表現）
    branch = webhook_event.ref.replace("refs/heads/", "")
    if not re.match(trigger_config.github.push.branch, branch):
        return False

    # 4. includedFiles / ignoredFiles チェック（設定されている場合）
    if trigger_config.includedFiles:
        if not any(fnmatch(file, pattern) for file in webhook_event.commits.modified for pattern in trigger_config.includedFiles):
            return False

    if trigger_config.ignoredFiles:
        if all(fnmatch(file, pattern) for file in webhook_event.commits.modified for pattern in trigger_config.ignoredFiles):
            return False

    return True  # 全条件マッチ → ビルド実行
```

### 4. ビルド起動

マッチング成功後、Cloud Buildはビルドジョブを作成します。

#### Build ID生成

```
10a4fb73-62fd-4211-8362-bc1492841791
```

UUIDv4形式で一意のビルドIDを生成

#### ソースコード取得

```bash
# GitHubからソースコードをfetch
git clone --depth=1 --branch=dev.dorapita.com https://github.com/ZIGExN/dorapita.git

# Cloud Storageにアップロード（キャッシュ）
tar czf source.tgz dorapita/
gsutil cp source.tgz gs://dorapita-core-dev_cloudbuild/source/1765873273.236597-*.tgz
```

#### Build Substitutions（変数置換）

```yaml
substitutions:
  PROJECT_ID: dorapita-core-dev
  COMMIT_SHA: 4a6a231ca
  BRANCH_NAME: dev.dorapita.com
  REPO_NAME: dorapita
  _SERVICE_NAME: dorapita-com
  _DEPLOY_REGION: asia-northeast1
  _AR_HOSTNAME: asia-northeast1-docker.pkg.dev
```

これらの変数はビルドステップで使用されます：

```yaml
- name: 'gcr.io/cloud-builders/docker'
  args:
  - build
  - -t
  - ${_AR_HOSTNAME}/${PROJECT_ID}/cloud-run-source-deploy/${REPO_NAME}/${_SERVICE_NAME}:${COMMIT_SHA}
  - .
```

### 5. ビルド実行

#### Build Worker

- **VM**: Google管理の一時VM（n1-standard-1相当）
- **OS**: Container-Optimized OS
- **Docker**: プリインストール
- **ネットワーク**: GCP内部ネットワーク

#### ステップ実行

```bash
# Step #0: Build
docker build \
  --no-cache \
  -t asia-northeast1-docker.pkg.dev/dorapita-core-dev/cloud-run-source-deploy/dorapita/dorapita-com:4a6a231ca \
  -f Dockerfile \
  .

# Step #1: Push
docker push asia-northeast1-docker.pkg.dev/dorapita-core-dev/cloud-run-source-deploy/dorapita/dorapita-com:4a6a231ca

# Step #2: Deploy
gcloud run services update dorapita-com \
  --image=asia-northeast1-docker.pkg.dev/dorapita-core-dev/cloud-run-source-deploy/dorapita/dorapita-com:4a6a231ca \
  --platform=managed \
  --region=asia-northeast1
```

### 6. ビルド結果通知

#### Cloud Build UI

- **URL**: `https://console.cloud.google.com/cloud-build/builds/BUILD_ID`
- **ステータス**: SUCCESS / FAILURE / TIMEOUT
- **ログ**: リアルタイムストリーミング

#### GitHub Status Check

Cloud BuildはGitHubにビルド結果を報告します：

```
✅ Google Cloud Build — dorapita-com
   Build succeeded in 3m 42s
   View logs →
```

コミットページに表示されるチェックマーク

## 設定の確認方法

### Trigger一覧

```bash
gcloud builds triggers list --project=dorapita-core-dev
```

### Trigger詳細

```bash
gcloud builds triggers describe TRIGGER_NAME --project=dorapita-core-dev
```

### GitHub App確認

1. GitHubリポジトリ → Settings → Integrations → GitHub Apps
2. "Google Cloud Build" を確認
3. Repository access をチェック

## セキュリティ

### 認証フロー

```
GitHub → Cloud Build
  ↓ Webhook署名検証
  X-Hub-Signature-256: sha256=HMAC(payload, secret)

Cloud Build → GitHub
  ↓ GitHub App JWT認証
  Authorization: Bearer <JWT_TOKEN>
```

### サービスアカウント

| SA | 役割 | 権限 |
|----|------|------|
| `1010225065326-compute@developer.gserviceaccount.com` | ビルド実行 | Cloud Run Admin, Storage Admin |
| `service-1010225065326@gcp-sa-cloudbuild.iam.gserviceaccount.com` | Cloud Build Agent | Secret Manager Accessor |

### GitHub App Token

- **有効期限**: 1時間
- **スコープ**: Repository contents, metadata
- **更新**: Cloud Buildが自動更新

## トラブルシューティング

### Webhookが届かない

```bash
# GitHubでWebhook配信履歴を確認
# Repository → Settings → Webhooks → Recent Deliveries

# エラー例:
# - 401 Unauthorized: GitHub App認証失敗
# - 404 Not Found: Trigger削除済み
# - 500 Internal Error: Cloud Build API障害
```

### Triggerが起動しない

```bash
# ブランチ名の正規表現チェック
# ^dev.dorapita.com$ は "dev.dorapita.com" のみマッチ
# ^dev.*$ なら "dev.xxx" 全てマッチ

# Trigger条件を再確認
gcloud builds triggers describe TRIGGER_NAME --project=dorapita-core-dev
```

### ビルドログが見れない

```bash
# Cloud Loggingで確認
gcloud logging read "resource.type=build AND resource.labels.build_id=BUILD_ID" \
  --project=dorapita-core-dev \
  --limit=100
```

## 料金

### Cloud Build

- **無料枠**: 1日あたり120ビルド分（n1-standard-1）
- **超過**: $0.003/分

### ネットワーク

- **GCP内部**: 無料
- **外部（GitHub fetch）**: $0.12/GB

## 関連ドキュメント

- [CI/CDパイプライン全体](CI_CD_PIPELINE.md)
- [ローカル開発環境](LOCAL_DEVELOPMENT.md)
- [Cloud Build公式ドキュメント](https://cloud.google.com/build/docs)
