# Dorapita インフラ・アプリケーション調査ドキュメント

ドラピタプロジェクトのGCPインフラ構成、アプリケーション構造、開発プロセス、組織変革戦略をまとめた調査ドキュメント集です。

## 📚 ドキュメント構成

```
.
├── gcp_inventory_reports/     # GCPインフラ調査レポート
├── apps/                      # アプリケーション構造ドキュメント
├── strategy/                  # 組織変革戦略
├── business_context/          # ビジネスコンテキスト
├── agenda/                    # 1on1アジェンダ
├── book_reading/              # 技術書籍の読書メモ
├── dorapita_code/             # アプリケーションコード（submodule）
├── gemini_tips.md             # Gemini CLI 実践ガイド
└── README.md                  # このファイル
```

## 🎯 読むべきドキュメントガイド

### シナリオ別ガイド

#### 1. 「ドラピタのインフラ構成を知りたい」

**START**: [`gcp_inventory_reports/dorapita-inventory.md`](gcp_inventory_reports/dorapita-inventory.md)

→ 全体概要、開発・本番環境比較、アーキテクチャ図

**詳細を知りたい場合**:
- 本番環境: [`dorapita-core-inventory.md`](gcp_inventory_reports/dorapita-core-inventory.md)
- 開発環境: [`dorapita-core-dev-inventory.md`](gcp_inventory_reports/dorapita-core-dev-inventory.md)
- セキュリティ: [`dorapita-security-audit.md`](gcp_inventory_reports/dorapita-security-audit.md)

#### 2. 「ローカルで開発環境を立ち上げたい」

**START**: [`apps/as-is/local-development.md`](apps/as-is/local-development.md)

→ Docker Compose環境セットアップ手順

**関連**:
- アプリ構造: [`apps/README.md`](apps/README.md)
- Docker Compose詳細: [`apps/as-is/how-to-up-docker-compose.md`](apps/as-is/how-to-up-docker-compose.md)
- 画像管理: [`apps/FILESTORE_IMAGE_MANAGEMENT.md`](apps/FILESTORE_IMAGE_MANAGEMENT.md)
- DB接続: [`gcp_inventory_reports/dorapita-inventory.md#cloud-sql接続ガイド`](gcp_inventory_reports/dorapita-inventory.md)

#### 3. 「デプロイ方法を知りたい」

**START**: [`apps/DEPLOYMENT_COMPARISON.md`](apps/DEPLOYMENT_COMPARISON.md)

→ 開発・本番のデプロイ方法比較

**詳細**:
- CI/CD全体: [`apps/CI_CD_PIPELINE.md`](apps/CI_CD_PIPELINE.md)
- Cloud Build仕組み: [`apps/CLOUD_BUILD_TRIGGER_MECHANISM.md`](apps/CLOUD_BUILD_TRIGGER_MECHANISM.md)

#### 4. 「重要な課題・改善提案を知りたい」

**START**: [`apps/important_issues/`](apps/important_issues/)

→ コスト削減、DB統合、オートスケーリング戦略

**主要レポート**:
- [サーバーコスト削減](apps/important_issues/SERVER_COST_REDUCTION_REPORT.md) - 年間296〜392万円削減可能
- [DB二重管理リスク](apps/important_issues/DB_INTEGRATION_TECHNICAL_RISK_REPORT.md) - 瞬断・保守コスト問題
- [BOT対策・オートスケール](apps/important_issues/BOT_ACCESS_AUTOSCALING_REPORT.md)

**開発環境改善**:
- [現状分析（As-Is）](apps/as-is/) - ローカル環境・DB分析
- [改善戦略（To-Be）](apps/to-be/) - DevContainer・Fixture戦略（年間1,064万円削減可能）

#### 5. 「組織体制・課題を知りたい」

**START**: [`strategy/dorapita-org-transformation-strategy.md`](strategy/dorapita-org-transformation-strategy.md)

→ 現状課題、RACI、Phase計画、次のアクション

**関連**:
- VNオフショア分析: [`strategy/ventura/`](strategy/ventura/) - Ventura社開発者の年間活動・コスト分析

#### 6. 「ビジネスコンテキスト・市場環境を知りたい」

**START**: [`business_context/awesome_gigexn_aquisition.md`](business_context/awesome_gigexn_aquisition.md)

→ じげんによるオーサムエージェント買収の詳細、PMI、シナジー効果、2024年問題

**関連**:
- じげん企業分析: [`business_context/gigexn_biz.md`](business_context/gigexn_biz.md)

#### 7. 「1on1のアジェンダを確認したい」

**START**: [`agenda/`](agenda/)

→ 松本さんとの1on1アジェンダ、プロジェクト進捗共有

#### 8. 「AI駆動開発のTipsを知りたい」

**START**: [`gemini_tips.md`](gemini_tips.md)

→ Gemini CLI 実践ガイド（settings.json、ラッパースクリプト、SQLite3分析）

**関連**:
- Claude Code設定: [`CLAUDE.md`](CLAUDE.md)

## 📁 ディレクトリ詳細

### 1. gcp_inventory_reports/ - GCPインフラ

| ファイル | 用途 | 読むタイミング |
|---------|------|---------------|
| **dorapita-inventory.md** | 全体概要・環境比較 | 最初に読む |
| dorapita-core-inventory.md | 本番環境詳細 | 本番運用時 |
| dorapita-core-dev-inventory.md | 開発環境詳細 | 開発環境構築時 |
| dorapita-security-audit.md | セキュリティ監査 | セキュリティレビュー時 |

**主な内容**:
- Compute Engine, Cloud Run, Cloud SQL, Filestore構成
- ネットワーク（VPC, LB, Cloud Armor）
- IAM権限マトリクス
- Cloud SQL接続ガイド

### 2. apps/ - アプリケーション構造

| ディレクトリ/ファイル | 用途 | 読むタイミング |
|-------------------|------|---------------|
| **README.md** | アプリ全体構造 | 最初に読む |
| **as-is/** | 現状の開発環境・DB分析 | 環境理解時 |
| **to-be/** | 改善戦略（DevContainer・Fixture） | 改善計画時 |
| **important_issues/** | 重要課題レポート | コスト削減・最適化時 |
| **tools/** | 開発支援ツール・スクリプト | DB分析・Fixture生成時 |
| CI_CD_PIPELINE.md | CI/CDパイプライン | デプロイ理解時 |
| DEPLOYMENT_COMPARISON.md | デプロイ方法比較 | 環境差異理解時 |
| FILESTORE_IMAGE_MANAGEMENT.md | 画像管理 | ファイル処理理解時 |

**主な内容**:
- CakePHP 4.4, 4.5, 5.0構成
- Docker Compose環境
- Cloud Build自動デプロイ
- 開発・本番のインフラ差異
- **As-Is**: ローカル開発環境、DB分析（113テーブル、24未使用）
- **To-Be**: DevContainer移行、Fixture戦略、年間1,064万円削減
- **重要課題**: コスト削減（年間296〜392万円）、DB統合、BOT対策

### 3. strategy/ - 組織変革戦略

| ファイル | 用途 | 読むタイミング |
|---------|------|---------------|
| **dorapita-org-transformation-strategy.md** | 組織変革戦略 | 体制・プロセス理解時 |
| **ventura/annual-report-2025.md** | VNオフショア年間分析 | オフショア評価・コスト分析時 |
| **ventura/contract-discrepancy.md** | 契約と実態の乖離 | 契約見直し時 |

**主な内容**:
- 現状課題8項目（PdM不在、デプロイ主体分裂等）
- 目標設定
- RACIマトリクス（開発・本番環境別）
- Phase計画と進捗（Phase 1-6）
- 次のアクション（ローカル開発セットアップ等）
- **Ventura社（VNオフショア）分析**:
  - 年間活動分析（2025年1-12月、5,735コミット）
  - 契約と実態の乖離（名簿外開発者、稼働率43%減少）
  - コスト効率分析（実質単価77万円/人月、相場の約2倍）
  - 推奨アクション（年間最大1,920万円の削減可能性）

### 4. business_context/ - ビジネスコンテキスト

| ファイル | 用途 | 読むタイミング |
|---------|------|---------------|
| **awesome_gigexn_aquisition.md** | じげん×オーサムエージェントM&A分析 | ビジネス背景・戦略理解時 |
| **gigexn_biz.md** | じげん企業・財務分析 | 親会社の事業戦略理解時 |

**主な内容**:
- 2024年問題（物流業界の構造的課題）
- じげんによるオーサムエージェント買収の戦略的合理性
- PMI（買収後統合）プロセスと価値創造
- じげんの事業セグメント（Vertical HR、Living Tech、Life Service）
- M&A戦略とロールアップ手法
- 財務パフォーマンス（FY2021-FY2026）

### 5. dorapita_code/ - アプリケーションコード（submodule）

**submodule**: `https://github.com/ZIGExN/dorapita.git`

**含まれるもの**:
- dorapita.com（CakePHP 4.4）
- cadm.dorapita.com（CakePHP 4.5）
- img.dorapita.com（CakePHP 5.0）
- kanri.dorapita.com（CakePHP 4.5）
- legacy.dorapita.com（CakePHP 1.3.21）
- その他

**更新方法**:
```bash
git submodule update --remote dorapita_code
```

## 🚀 クイックスタート

### 新しくプロジェクトに参加した場合

1. **全体像を把握**
   - [`gcp_inventory_reports/dorapita-inventory.md`](gcp_inventory_reports/dorapita-inventory.md) を読む
   - [`apps/README.md`](apps/README.md) を読む

2. **開発環境をセットアップ**
   - [`apps/LOCAL_DEVELOPMENT.md`](apps/LOCAL_DEVELOPMENT.md) の手順に従う

3. **組織体制を理解**
   - [`strategy/dorapita-org-transformation-strategy.md`](strategy/dorapita-org-transformation-strategy.md) を読む

### デプロイしたい場合

1. **開発環境へのデプロイ**
   - [`apps/CI_CD_PIPELINE.md`](apps/CI_CD_PIPELINE.md#開発環境デプロイ自動) を参照
   - `dev.*` ブランチにpush → 自動デプロイ

2. **本番環境へのデプロイ**
   - [`apps/DEPLOYMENT_COMPARISON.md`](apps/DEPLOYMENT_COMPARISON.md#本番環境dorapita-core) を参照
   - じげん担当の手動デプロイ（VM SSH）

### Cloud SQLに接続したい場合

[`gcp_inventory_reports/dorapita-inventory.md#cloud-sql接続ガイド`](gcp_inventory_reports/dorapita-inventory.md) を参照

## 🔧 技術スタック概要

| 技術 | バージョン | 用途 |
|------|----------|------|
| PHP | 7.4, 8.1 (5.6 legacy) | バックエンド |
| CakePHP | 4.4, 4.5, 5.0 (1.3.21 legacy) | フレームワーク |
| PostgreSQL | 10 | 主DB（dorapita, cadm, edit） |
| MySQL | 5.7 | 副DB（kanri, legacy） |
| Redis | 7.2 | キャッシュ・セッション |
| Docker | - | ローカル開発 |
| Cloud Run | - | 開発環境デプロイ |
| Compute Engine | - | 本番環境 |
| Cloud Build | - | CI/CD |
| Filestore | - | 画像ストレージ（NFS） |

## 📊 環境情報

| 環境 | GCPプロジェクト | インフラ | デプロイ方式 |
|------|----------------|---------|-------------|
| **本番** | dorapita-core | Compute Engine (VM) | 手動（じげん担当） |
| **開発** | dorapita-core-dev | Cloud Run | 自動（Cloud Build） |
| **ローカル** | - | Docker Compose | - |

## ⚠️ 重要な注意事項

1. **DB認証情報**
   - Secret Managerに格納
   - `roles/secretmanager.secretAccessor` 権限が必要
   - 詳細: [`gcp_inventory_reports/dorapita-inventory.md#cloud-sql接続ガイド`](gcp_inventory_reports/dorapita-inventory.md)

2. **画像ファイル**
   - リポジトリに含まれない（`.gitignore`除外）
   - Filestore（NFS）に保存
   - ローカル開発: `gs://dorapita-infra-stg/assets/upfiles/` から取得
   - 詳細: [`apps/FILESTORE_IMAGE_MANAGEMENT.md`](apps/FILESTORE_IMAGE_MANAGEMENT.md)

3. **submodule**
   - `dorapita_code/` は別リポジトリ（ZIGExN/dorapita）
   - clone後は `git submodule update --init --recursive`

## 🔍 調査状況

| カテゴリ | 進捗 | ステータス |
|---------|-----:|----------|
| GCPインフラ | 100% | ✅ 完了 |
| アプリケーション構造 | 90% | ✅ ほぼ完了 |
| CI/CD | 100% | ✅ 完了 |
| データベーススキーマ | 40% | 🟡 ブロッカーあり |
| ユーザージャーニー | 0% | 🔴 次のステップ |
| ビジネスロジック | 0% | 🔴 未着手 |

**残調査項目**:
- Filestore実ファイル調査（VM経由）
- Cloud SQL スキーマ調査（Secret Manager権限必要）
- ローカル開発環境セットアップ・動作確認
- ユーザージャーニー調査

詳細: [`strategy/dorapita-org-transformation-strategy.md#phase詳細タスク`](strategy/dorapita-org-transformation-strategy.md)

## 🤝 コントリビューション

このリポジトリは調査ドキュメントの集積です。新しい発見や更新があれば、該当するmdファイルを更新してください。

### ドキュメント更新ガイドライン

- **GCPリソース変更**: `gcp_inventory_reports/` 配下を更新
- **アプリ構造変更**: `apps/` 配下を更新
- **組織・プロセス変更**: `strategy/` 配下を更新
- **新しいサービス追加**: `apps/` に新規mdファイル作成

## 📞 問い合わせ

- **GCP権限**: k.hashimoto@awesomegroup.co.jp
- **組織・体制**: strategy/dorapita-org-transformation-strategy.md 参照

## 📝 更新履歴

| 日付 | 内容 |
|------|------|
| 2025-12-21 | 初版作成（GCPインフラ調査） |
| 2025-12-21 | Filestore修正、IAM追加 |
| 2025-12-21 | apps/追加（アプリ構造調査）、strategy/追加（組織変革戦略） |
| 2025-12-23 | README.md作成、CLAUDE.md更新 |
| 2025-12-23 | VN開発者アクティビティ分析追加 |
| 2025-12-28 | business_context/追加（じげん×オーサムエージェントM&A、企業分析） |
| 2026-01-07 | agenda/追加、gemini_tips.md追加、settings.local.json更新 |
