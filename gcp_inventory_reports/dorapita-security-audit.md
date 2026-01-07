# Dorapita GCPインフラ セキュリティ監査レポート

**調査日時**: 2025-12-21
**対象プロジェクト**: dorapita-core (本番), dorapita-core-dev (開発)

---

## エグゼクティブサマリー

| 重大度 | 件数 | 主な課題 |
|--------|------|----------|
| 高 | 2 | Cloud SQL SSL未強制、DBバージョンEOL |
| 中 | 7 | SSLポリシー未設定、FTP公開、IAM過剰権限、**Cloud Runセキュリティ問題**等 |
| 低 | 3 | 外部IP、ログ転送差異等 |

---

## 1. VPCとファイアウォールのリスク評価

### default VPCの使用状況

| 環境 | default VPC接続リソース | リスク判定 |
|------|------------------------|------------|
| **本番** | **なし** - 全VM・DB・RedisはdorapitaVPCに接続 | **低** |
| **開発** | **なし** - 全VM・DB・RedisはdorapitaVPCに接続 | **低** |

#### 詳細

**本番環境 (dorapita-core)**:
- VM 6台: すべて `dorapita-vpc` に接続
- Cloud SQL 2台: すべて `dorapita-vpc` に接続
- Redis 1台: `dorapita-vpc` に接続

**開発環境 (dorapita-core-dev)**:
- VM 8台: すべて `dorapita-vpc` に接続
- Cloud SQL 2台: すべて `dorapita-vpc` に接続
- Redis 1台: `dorapita-vpc` に接続

#### 結論

default VPCのファイアウォールルール（SSH/RDP全世界公開等）は**現時点でリスクなし**。
ただし、将来の誤用防止のため以下を推奨：
- default VPCの削除、または
- 危険なファイアウォールルールの削除

---

### ファイアウォールルール評価

#### dorapita-vpc（使用中）

| ルール | ソース | 許可 | リスク | 備考 |
|--------|--------|------|--------|------|
| allow-health-check | GCP Health Check IP | tcp | **低** | 正常 |
| allow-http | 0.0.0.0/0 | tcp:80 | **低** | LB経由想定 |
| allow-ingress-from-iap | 35.235.240.0/20 | tcp:22 | **低** | IAP経由SSH（推奨構成） |
| allow-internal | 10.120.0.0/16 or 10.121.0.0/16 | tcp/all | **低** | VPC内通信 |
| **allow-ftp (dev)** | **0.0.0.0/0** | **tcp全ポート** | **中** | 全世界からFTP可能 |

#### default VPC（未使用）

| ルール | ソース | 許可 | 現状リスク | 潜在リスク |
|--------|--------|------|------------|------------|
| allow-ssh | 0.0.0.0/0 | tcp:22 | **低** | 高（誤接続時） |
| allow-rdp | 0.0.0.0/0 | tcp:3389 | **低** | 高（誤接続時） |
| allow-icmp | 0.0.0.0/0 | icmp | **低** | 低 |
| allow-internal | 10.128.0.0/9 | all | **低** | 中 |

---

## 2. 重大度：高

### 2.1 Cloud SQL - SSL/TLS未強制

| 環境 | インスタンス | requireSsl | sslMode |
|------|-------------|------------|---------|
| 本番 | db-120011 (MySQL) | false | ALLOW_UNENCRYPTED_AND_ENCRYPTED |
| 本番 | pg-120011 (PostgreSQL) | false | ALLOW_UNENCRYPTED_AND_ENCRYPTED |
| 開発 | db-120011 (MySQL) | false | ALLOW_UNENCRYPTED_AND_ENCRYPTED |
| 開発 | pg-120011 (PostgreSQL) | false | ALLOW_UNENCRYPTED_AND_ENCRYPTED |

**リスク**:
- データベース接続が暗号化されていない可能性
- VPC内でも中間者攻撃のリスクあり
- 機密データ（個人情報等）が平文で通信される可能性

**推奨対応**:
```bash
# SSL強制化
gcloud sql instances patch INSTANCE_NAME \
  --require-ssl \
  --ssl-mode=ENCRYPTED_ONLY
```

---

### 2.2 データベースバージョン EOL

| インスタンス | 現行バージョン | EOL日 | 状態 |
|-------------|---------------|-------|------|
| db-120011 | MySQL 5.7 | 2023年10月 | **EOL済み（1年以上経過）** |
| pg-120011 | PostgreSQL 10 | 2022年11月 | **EOL済み（2年以上経過）** |

**リスク**:
- セキュリティパッチが提供されない
- 既知の脆弱性が修正されない
- コンプライアンス違反の可能性

**推奨対応**:
- MySQL 8.0 へのアップグレード計画策定
- PostgreSQL 15以降へのアップグレード計画策定
- アップグレード前のテスト環境での検証

---

## 3. 重大度：中

### 3.1 ロードバランサー SSLポリシー未設定

| 環境 | HTTPSプロキシ数 | SSLポリシー設定 |
|------|-----------------|-----------------|
| 本番 | 9個 | **なし** |
| 開発 | 13個 | **なし** |

**リスク**:
- TLS 1.0/1.1が許可されている可能性
- 弱い暗号スイート（RC4、3DES等）が使用される可能性
- PCI DSS等のコンプライアンス違反の可能性

**推奨対応**:
```bash
# SSLポリシー作成
gcloud compute ssl-policies create dorapita-ssl-policy \
  --profile=MODERN \
  --min-tls-version=1.2

# HTTPSプロキシに適用
gcloud compute target-https-proxies update PROXY_NAME \
  --ssl-policy=dorapita-ssl-policy
```

---

### 3.2 FTPファイアウォールルール（開発環境）

| ルール | ネットワーク | ソース | 許可 |
|--------|--------------|--------|------|
| dorapita-vpc-allow-ftp | dorapita-vpc | 0.0.0.0/0 | tcp（全ポート） |

**リスク**:
- 全世界からFTPサーバーにアクセス可能
- ブルートフォース攻撃のリスク
- FTPは平文プロトコル（認証情報漏洩リスク）

**推奨対応**:
- ソースIPを許可されたIPアドレスのみに制限
- SFTP/SCP への移行検討
- FTP不要であればルール削除

---

### 3.3 IAM - Owner権限の個人付与

| 環境 | ロール | 付与先 |
|------|--------|--------|
| 本番 | roles/owner | nagai@zigexn.co.jp |
| 開発 | roles/owner | nagai@zigexn.co.jp |

**リスク**:
- 個人アカウントへの過剰な権限付与
- アカウント侵害時に全リソースが危険に
- 監査・追跡の困難さ

**推奨対応**:
- グループ経由での権限付与に変更
- 最小権限の原則に基づくロール設計
- 定期的な権限棚卸し

---

### 3.4 外部IPを持つVM

| 環境 | VM名 | 外部IP | 用途 |
|------|------|--------|------|
| 本番 | gw-120011 | 34.84.187.36 | Gateway |
| 本番 | web-120011 | 34.146.149.10 | Web |
| 本番 | web-120021 | 35.190.224.174 | Web |
| 本番 | web-120022 | 104.198.83.165 | Web |
| 本番 | web-120023 | 34.85.114.186 | Web |
| 本番 | web-120031 | 34.84.115.88 | Web |
| 開発 | ftp-120011 | 35.200.74.161 | FTP |
| 開発 | web-120011 | 35.187.213.112 | Web |
| 開発 | web-120021 | 35.200.24.243 | Web |
| 開発 | web-120031 | 34.146.10.179 | Web |

**リスク**:
- 攻撃対象面（Attack Surface）の増大
- 直接攻撃を受けるリスク

**推奨対応**:
- ロードバランサー経由のみでアクセスするVMは外部IP削除
- Cloud NAT経由でのアウトバウンド接続に変更
- 必要最小限の外部IP保持

---

### 3.5 サービスアカウント管理

| 環境 | SA数 | 内容 |
|------|------|------|
| 本番 | 1 | Compute Engineデフォルトのみ |
| 開発 | 2 | デフォルト + dorapita-gcs |

**リスク**:
- デフォルトサービスアカウントは過剰な権限（roles/editor）を持つ
- ワークロード間の権限分離ができていない

**推奨対応**:
- 最小権限のカスタムサービスアカウント作成
- ワークロード別にSAを分離
- Workload Identity の検討

---

### 3.6 Cloud Run セキュリティ問題（統合）

#### 3.6.1 サービス一覧と設定状況

**開発環境 (dorapita-core-dev)** - 7サービス:

| サービス | Ingress | VPC | Secrets | Env直接 | DB接続 |
|----------|---------|-----|---------|---------|--------|
| dorapita-com | internal+LB | あり | SM(vol) | なし | 可能 |
| cadm-dorapita-com | internal+LB | あり | SM(vol) | なし | 可能 |
| img-dorapita-com | internal+LB | あり | なし | DEBUG, CORS | 可能 |
| api-dorapita-com | **all** ⚠️ | **なし** | SM(env) | APP_ENV, LOG | **不可** |
| edit-dorapita-com | internal+LB | **なし** ⚠️ | SM(vol) **※DB** | なし | **不可** |
| help-dorapita-com | internal+LB | なし | なし | なし | 不要 |
| dorapita-maintenance | internal+LB | なし | なし | なし | 不要 |

**本番環境 (dorapita-core)** - 3サービス:

| サービス | Ingress | VPC | Secrets | Env直接 | DB接続 |
|----------|---------|-----|---------|---------|--------|
| img-dorapita-com | internal+LB | あり | なし | DEBUG, CORS | 不要 |
| help-dorapita-com | internal+LB | なし | なし | なし | 不要 |
| dorapita-maintenance | internal+LB | なし | なし | なし | 不要 |

凡例: SM=Secret Manager, vol=ボリュームマウント, env=環境変数参照

---

#### 3.6.2 問題1: Ingress設定 `all`（api-dorapita-com）

```yaml
run.googleapis.com/ingress: all  # 全公開 ⚠️
```

**影響**: 開発環境の api-dorapita-com のみ

**リスク**:
- Cloud Run URL (`https://api-dorapita-com-dr3sq5ocsq-an.a.run.app`) に直接アクセス可能
- Load Balancer の Cloud Armor をバイパス可能
- 意図しないアクセス経路が存在

**推奨対応**:
```bash
gcloud run services update api-dorapita-com \
  --region=asia-northeast1 \
  --ingress=internal-and-cloud-load-balancing
```

---

#### 3.6.3 問題2: VPC接続なしでDBシークレット保持（edit-dorapita-com）

```yaml
volumes:
- name: secret-2
  secret:
    secretName: database-edit-dorapita-com  # DB接続情報
# VPC設定なし → DBにアクセス不可
```

**影響**: 開発環境の edit-dorapita-com

**リスク**:
- DBシークレットをマウントしているが、VPC接続がないためDBにアクセス不可
- 設定ミスの可能性（意図した動作か不明）
- シークレットが無駄に露出している

**推奨対応**:
1. DB接続が必要な場合: VPC Direct Egress を追加
2. DB接続が不要な場合: シークレットマウントを削除

---

#### 3.6.4 問題3: 環境変数のハードコード

**該当サービス**:

| サービス | 環境 | ハードコードされた値 |
|----------|------|---------------------|
| api-dorapita-com | 開発 | `APP_ENV=production`, `LOG_CHANNEL=stack` |
| img-dorapita-com | 開発 | `DEBUG=false`, `ACCESS_CONTROL_ALLOW_ORIGIN=*` |
| img-dorapita-com | 本番 | `DEBUG=false`, `ACCESS_CONTROL_ALLOW_ORIGIN=*` |

**リスク**:
- 設定変更時にサービス再デプロイが必要
- マニフェスト/コンソールに平文で表示
- 監査ログで追跡困難
- 環境間での値の不整合リスク

**推奨対応**: Secret Manager への移行

```yaml
# Before: ハードコード
env:
- name: APP_ENV
  value: production

# After: Secret Manager参照
env:
- name: APP_ENV
  valueFrom:
    secretKeyRef:
      name: app-config-api-dorapita-com
      key: latest
```

---

#### 3.6.5 問題4: CORS設定 `*`（全オリジン許可）

```yaml
env:
- name: ACCESS_CONTROL_ALLOW_ORIGIN
  value: '*'  # 全オリジン許可
```

**影響**: 両環境の img-dorapita-com

**リスク**:
- 任意のオリジンからのリクエストを許可
- 意図しないサイトからの画像参照が可能
- CSRF攻撃のリスク増加

**推奨対応**:
- 許可するオリジンを明示的に指定
- 例: `https://dorapita.com,https://stg.dorapita.com`

---

#### 3.6.6 問題5: Secret Manager 管理の不統一

**Secret Manager 使用状況**:

| 環境 | シークレット名 | 用途 | 使用サービス |
|------|---------------|------|-------------|
| 本番 | GA4_API_SECRET | GA4 API | 不明（Cloud Runでは未使用） |
| 本番 | GA4_MEASUREMENT_ID | GA4 | 不明（Cloud Runでは未使用） |
| 開発 | shell-environment-variables-dorapita-com | 環境変数 | dorapita-com |
| 開発 | shell-environment-variables-cadm-dorapita-com | 環境変数 | cadm-dorapita-com |
| 開発 | shell-environment-variables-api-dorapita-com | 環境変数 | api-dorapita-com |
| 開発 | settings-edit-dorapita-com | 設定 | edit-dorapita-com |
| 開発 | database-edit-dorapita-com | DB接続 | edit-dorapita-com |

**問題点**:
- 本番環境のCloud RunはSecret Manager未使用
- 環境変数が直接ハードコードされている
- 本番と開発で管理方法が大きく異なる

**推奨対応**:
- 全環境でSecret Manager使用を統一
- 命名規則の統一（例: `{env}-{service}-{type}`）

---

#### 3.6.7 Cloud SQL 接続経路

**Cloud SQL IP設定状況**:

| 環境 | インスタンス | パブリックIP | プライベートIP | 認可ネットワーク |
|------|-------------|-------------|----------------|-----------------|
| 本番 | db-120011 (MySQL) | 34.146.103.51 | 10.117.160.3 | **なし** |
| 本番 | pg-120011 (PostgreSQL) | 35.221.105.252 | 10.117.160.5 | **なし** |
| 開発 | db-120011 (MySQL) | 34.146.52.246 | 10.227.112.3 | **なし** |
| 開発 | pg-120011 (PostgreSQL) | 35.194.97.196 | 10.227.112.14 | **なし** |

**現状**: パブリックIPは有効だが、認可ネットワーク未設定のため外部からの直接接続は不可。
実質的に**プライベートアクセスのみ**で運用。

**接続経路図**:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              dorapita-vpc                               │
│                                                                         │
│  【DBアクセス可能】                                                      │
│                                                                         │
│  ┌──────────────────┐      プライベートIP        ┌──────────────────┐  │
│  │  Compute Engine  │ ─────────────────────────> │    Cloud SQL     │  │
│  │  (web-120011等)   │      10.x.x.x              │  db-120011       │  │
│  │  10.120.101.x     │                            │  pg-120011       │  │
│  └──────────────────┘                             │                  │  │
│                                                   │  プライベートIP:  │  │
│  ┌──────────────────┐      VPC Direct Egress      │  本番:10.117.160.x│  │
│  │  Cloud Run       │ ─────────────────────────> │  開発:10.227.112.x│  │
│  │  (VPC接続あり)    │      dorapita-vpc/web      │                  │  │
│  │  - dorapita-com  │                             └──────────────────┘  │
│  │  - cadm-dorapita │                                                   │
│  │  - img-dorapita  │                                                   │
│  └──────────────────┘                                                   │
│                                                                         │
│  【DBアクセス不可】                                                      │
│                                                                         │
│  ┌──────────────────┐                                                   │
│  │  Cloud Run       │ ─────── X ─────── VPCなし、Cloud SQL Proxyなし   │
│  │  (VPC接続なし)    │        ↑                                         │
│  │  - api-dorapita  │   プライベートIPに到達不可                         │
│  │  - edit-dorapita │   パブリックIPは認可ネットワークなし               │
│  └──────────────────┘                                                   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**補足**:
- Cloud SQL Auth Proxy (`run.googleapis.com/cloudsql-instances`) は両環境とも未使用
- 本番Cloud RunからのDBアクセスは設計上不要（静的コンテンツのみ）

---

#### 3.6.8 Cloud Run セキュリティ問題まとめ

| # | 問題 | 環境 | 重大度 | 対応 |
|---|------|------|--------|------|
| 1 | Ingress `all` | 開発 api | 中 | `internal-and-cloud-load-balancing` に変更 |
| 2 | VPCなし+DBシークレット | 開発 edit | 中 | VPC追加 or シークレット削除 |
| 3 | 環境変数ハードコード | 両環境 | 低 | Secret Manager 移行 |
| 4 | CORS `*` | 両環境 img | 低 | オリジン明示 |
| 5 | Secret Manager不統一 | 両環境 | 低 | 管理方法統一 |
| 6 | Cloud SQL パブリックIP | 両環境 | 低 | 無効化推奨 |

---

## 4. 重大度：低

### 4.1 環境間の設定差異

| 項目 | 本番 | 開発 | 影響 |
|------|------|------|------|
| アクセスログ転送 | BigQuery転送あり | なし | 開発でのログ分析困難 |
| Secret Manager | 2シークレット（GA4関連） | 5シークレット（環境変数等） | 管理方法の不統一 |
| Shielded VM | 有効 | 未確認 | セキュリティ機能の差 |

---

## 5. 良好な設定（対策済み項目）

| 項目 | 本番 | 開発 | 評価 |
|------|------|------|------|
| OS Login | 有効 | 有効 | SSH鍵のIAM統合管理 |
| block-project-ssh-keys | - | true | プロジェクトSSH鍵の無効化 |
| Shielded VM | 有効 | - | セキュアブート、vTPM、整合性監視 |
| Cloud Armor | 8ポリシー | 8ポリシー | WAF保護 |
| IAP経由SSH | 設定済み | 設定済み | ゼロトラストアクセス |
| 監査ログ | 保持中 | 保持中 | 監査証跡の確保 |
| バックアップ | 7日保持 | - | データ保護 |
| VPC分離 | dorapita-vpc使用 | dorapita-vpc使用 | ネットワーク分離 |

---

## 6. 対応優先順位

| 優先度 | 項目 | 影響範囲 | 難易度 | 推定工数 |
|--------|------|----------|--------|----------|
| **1** | Cloud SQL SSL強制化 | 両環境 | 低 | 1日 |
| **2** | SSLポリシー設定 | 両環境 | 低 | 1日 |
| **3** | default VPC整理 | 両環境 | 低 | 半日 |
| **4** | FTPファイアウォール制限 | 開発 | 低 | 半日 |
| **5** | Cloud Run VPC設定確認・修正 | 開発 | 低 | 半日 |
| **6** | Cloud SQL パブリックIP無効化 | 両環境 | 低 | 半日 |
| **7** | DBアップグレード計画策定 | 両環境 | 高 | 計画:1週間、実施:要検証 |
| **8** | IAM権限見直し | 両環境 | 中 | 1週間 |
| **9** | 外部IP削減 | 両環境 | 中 | 1週間 |
| **10** | サービスアカウント整理 | 両環境 | 中 | 1週間 |

---

## 7. 推奨アクションプラン

### 即時対応（1週間以内）

1. **Cloud SQL SSL強制化**
   - 両環境の全Cloud SQLインスタンスでSSLを強制
   - アプリケーション側のSSL接続設定確認

2. **SSLポリシー設定**
   - TLS 1.2以上を強制するSSLポリシー作成
   - 全HTTPSプロキシに適用

3. **default VPC整理**
   - 使用していないdefault VPCのファイアウォールルール削除
   - または default VPC自体を削除

4. **Cloud Run VPC設定確認・修正**
   - edit-dorapita-com: DBシークレットがあるがVPC接続なし → 設定確認・修正
   - api-dorapita-com: DB接続必要性を確認し、必要なら VPC Direct Egress 追加

5. **Cloud SQL パブリックIP無効化**
   - 認可ネットワーク未設定で実質未使用のため無効化を検討
   - `gcloud sql instances patch INSTANCE --no-assign-ip`

### 短期対応（1ヶ月以内）

6. **FTPファイアウォール制限**
   - 許可IPアドレスのリスト作成
   - ソースIP制限の適用

7. **IAM権限棚卸し**
   - 全ユーザー・サービスアカウントの権限確認
   - Owner権限の見直し

### 中期対応（3ヶ月以内）

8. **データベースアップグレード計画**
   - 互換性テスト計画策定
   - ステージング環境でのテスト
   - 本番アップグレードスケジュール策定

9. **外部IP削減**
   - Cloud NAT導入検討
   - 不要な外部IPの削除

---

## 8. 参考：Cloud Armor設定状況

### 本番環境

| ポリシー名 | 適用先 |
|------------|--------|
| allow-dorapita-system | 特定IPのみ許可、デフォルト拒否 |
| default-security-policy-for-backend-service-* | 各バックエンドサービス |

**allow-dorapita-systemのルール**:
- Priority 10000: ドラピタサーバー群 (118.69.34.46) - 許可
- Priority 20000: リクオプ連携 (3.115.169.209等) - 許可
- Priority 21000: ビズコミ連携 (18.181.109.194等) - 許可
- Priority 2147483647: デフォルト - 拒否(403)

---

## 9. 関連ドキュメント

- [dorapita-inventory.md](./dorapita-inventory.md) - インフラ概要
- [dorapita-core-inventory.md](./dorapita-core-inventory.md) - 本番環境詳細
- [dorapita-core-dev-inventory.md](./dorapita-core-dev-inventory.md) - 開発環境詳細

---

*監査完了: 2025-12-21*
