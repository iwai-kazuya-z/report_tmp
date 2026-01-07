# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## リポジトリ概要

Dorapita GCPインフラとアプリケーション構造の調査ドキュメントリポジトリ。インフラ構成、アプリケーションアーキテクチャ、開発プロセス、組織変革戦略を管理。

## ディレクトリ構成

```
.
├── gcp_inventory_reports/     # GCPインフラ調査レポート
│   ├── dorapita-inventory.md
│   ├── dorapita-core-inventory.md
│   ├── dorapita-core-dev-inventory.md
│   ├── dorapita-security-audit.md
│   └── dorapita-gcp-infrastructure-report.md  # (gitignore) Google Docs用
├── apps/                      # アプリケーション構造ドキュメント
│   ├── README.md
│   ├── LOCAL_DEVELOPMENT.md
│   ├── CI_CD_PIPELINE.md
│   ├── DEPLOYMENT_COMPARISON.md
│   ├── FILESTORE_IMAGE_MANAGEMENT.md
│   └── *.md                   # 各サービス詳細
├── strategy/                  # 組織変革戦略（RACI、Phase計画）
│   ├── dorapita-org-transformation-strategy.md
│   └── ventura/               # Ventura社（VNオフショア）開発者分析
├── business_context/          # ビジネスコンテキスト（企業情報、M&A、市場環境）
│   ├── awesome_gigexn_aquisition.md
│   └── gigexn_biz.md
├── agenda/                    # 1on1アジェンダ
│   └── MTsan_1on1_*.md        # 松本さんとの1on1アジェンダ
├── book_reading/              # 技術書籍の読書メモ
├── dorapita_code/             # アプリケーションコード（submodule: ZIGExN/dorapita）
├── gemini_tips.md             # Gemini CLI 実践ガイド
├── about_*.md                 # (gitignore) 社内情報メモ
├── CLAUDE.md
└── README.md                  # ドキュメントガイド
```

## ドキュメント構成

### GCPインフラ（gcp_inventory_reports/）

| ファイル | 内容 |
|----------|------|
| `dorapita-inventory.md` | 全プロジェクト概要、Cloud SQL接続ガイド |
| `dorapita-core-inventory.md` | 本番環境詳細インベントリ |
| `dorapita-core-dev-inventory.md` | 開発環境詳細インベントリ |
| `dorapita-security-audit.md` | セキュリティ監査レポート |

### アプリケーション（apps/）

| ファイル | 内容 |
|----------|------|
| `README.md` | アプリケーション全体構造、技術スタック |
| `as-is/` | 現状の開発環境・DB分析 |
| `to-be/` | 改善戦略（DevContainer、Fixture） |
| `important_issues/` | 重要課題レポート（コスト削減、DB統合） |
| `tools/` | 開発支援ツール・スクリプト |
| `CI_CD_PIPELINE.md` | Cloud Build自動デプロイ |
| `DEPLOYMENT_COMPARISON.md` | 開発vs本番デプロイ方法 |
| `FILESTORE_IMAGE_MANAGEMENT.md` | 画像管理の仕組み |

### 組織・戦略（strategy/）

| ファイル | 内容 |
|----------|------|
| `dorapita-org-transformation-strategy.md` | 組織変革戦略、RACI、Phase計画 |
| `ventura/annual-report-2025.md` | Ventura社（VNオフショア）年間活動分析（2025年1-12月） |
| `ventura/contract-discrepancy.md` | 契約と実態の乖離レポート |

### ビジネスコンテキスト（business_context/）

| ファイル | 内容 |
|----------|------|
| `awesome_gigexn_aquisition.md` | じげんによるオーサムエージェント買収の詳細分析、PMI、シナジー効果 |
| `gigexn_biz.md` | じげんの企業概要、事業セグメント、財務分析、M&A戦略 |

## GCPプロジェクト

| 環境 | Project ID | 用途 |
|------|------------|------|
| 本番 | dorapita-core | Production |
| 開発 | dorapita-core-dev | Staging/Development |

## アーキテクチャ

- **リージョン**: asia-northeast1 (東京)
- **VPC**: dorapita-vpc（本番: 10.120.0.0/16、開発: 10.121.0.0/16）
- **コンピューティング**: Compute Engine + Cloud Run ハイブリッド構成
- **データベース**: Cloud SQL (MySQL 5.7, PostgreSQL 10) + Memorystore Redis
- **ロードバランサー**: Global HTTPS LB + Cloud Armor (WAF)
- **認証**: IAP経由SSH、OS Login有効

## 主要なセキュリティ考慮事項

- Cloud SQL: SSL未強制、DBバージョンEOL（MySQL 5.7, PostgreSQL 10）
- Cloud Run: 一部サービスでIngress `all`、VPC接続の不整合あり
- ロードバランサー: SSLポリシー未設定（TLS 1.0/1.1許可の可能性）

## ビジネスコンテキスト

### 運営会社

| 会社名 | 役割 | 備考 |
|--------|------|------|
| 株式会社オーサムエージェント (AA) | ドラピタ運営会社 | 2016年設立、名古屋拠点 |
| 株式会社じげん (ZIGExN) | 親会社 | 東証プライム上場 (3679) |

- **買収**: 2022年11月、じげんがAAを完全子会社化（約2.4億円）
- **創業者**: 竹村優氏（マイナビ出身）→ 現代表: 代田晴久氏

### ドラピタとは

**運送・物流業界特化の求人メディア**

- ドライバー、倉庫作業員、配車係などに特化
- 東海エリア（愛知・岐阜・三重）で創業、全国展開中
- 2,000社以上の物流企業が顧客基盤

**検索軸の特徴**:
- 保有免許（中型、大型、牽引など）
- 車種（ウイング車、平ボディ、冷凍車など）
- 積み荷の種類

### 事業ポートフォリオ（AAグループ）

| サービス | 種別 | 概要 |
|----------|------|------|
| ドラピタ | 求人メディア | 掲載課金型（フロー収益） |
| ドラピタエージェント | 人材紹介 | 成功報酬型（2025年本格稼働） |
| ビズコミ | 求人メディア | エッセンシャルワーカー領域（製造、建設、警備、医療介護） |
| ドラウェブ | HP制作 | 物流企業向け採用サイト制作 |
| +α Design | クリエイティブ | パンフレット、ロゴ等の販促物制作 |

### 親会社じげんについて

**企業理念**: 「生活機会の最大化」

**事業セグメント**:

| セグメント | 領域 | 主要ブランド |
|------------|------|--------------|
| Vertical HR | 業界特化型人材 | リジョブ（美容）、ミラクス（介護）、AA（物流） |
| Living Tech | 不動産・住まい | 賃貸スモッカ、リショップナビ |
| Life Service | 生活サービス | アップルワールド（旅行）、TCV（自動車輸出） |

**じげんの強み**:
- M&A戦略（創業以来30社以上買収）
- PMI（買収後統合）ノウハウ
- デジタルマーケティング・SEO最適化

### 市場環境: 2024年問題

2024年4月より働き方改革関連法適用:
- ドライバーの時間外労働上限: **年960時間**
- 輸送能力の減少 → ドライバー需要増大
- 若年層志望者減少により慢性的な人手不足

→ 物流HR市場は構造的な成長市場

### 開発体制

| 担当 | 役割 |
|------|------|
| Ventura社 (VN) | オフショア開発（ベトナム） |
| AA社 | 受け入れ検査・要件定義（※責任分界点が曖昧） |

**技術スタック**: PHP（詳細はGitHubアクセス後に確認予定）

## gcloudコマンド例

```bash
# プロジェクト切り替え
gcloud config set project dorapita-core        # 本番
gcloud config set project dorapita-core-dev    # 開発

# VM一覧
gcloud compute instances list

# Cloud Run一覧
gcloud run services list --region=asia-northeast1

# Cloud SQL一覧
gcloud sql instances list
```
