# To-Be: 開発環境改善戦略

現状（As-Is）からの脱却と、理想的な開発環境への移行戦略をまとめています。

---

## 📋 戦略ドキュメント構成

### 1. DevContainer移行計画（Phase 1）
**ファイル**: [`devcontainer-setup.md`](devcontainer-setup.md)

**目的**: ローカル開発環境の標準化と効率化

**主要な改善**:
- ✅ **開発環境の完全な再現性**: Docker-in-Docker不要、VS Code統合
- ✅ **オンボーディング時間短縮**: 1日 → 30分
- ✅ **環境差異の解消**: 「私の環境では動く」問題の撲滅

**技術スタック**:
- Dev Container (VS Code)
- Docker Compose統合
- PostgreSQL 15 (最新版)
- MySQL 8.0 (延長サポート解消)

**移行ステップ**:

| Phase | タスク | 期間 | 担当 |
|-------|--------|------|------|
| Phase 1.1 | `.devcontainer/` 作成、PostgreSQL 15環境構築 | 1週間 | Infra + Dev |
| Phase 1.2 | Fixture作成（Phase 2と並行） | 2週間 | Dev |
| Phase 1.3 | MySQL 8.0統合、二重管理解消 | 2週間 | Dev + Infra |
| Phase 1.4 | ドキュメント整備・トレーニング | 1週間 | Dev |

**期待効果**:
- 新規開発者のセットアップ時間: **95%削減**（1日 → 30分）
- 環境起因のトラブルシューティング時間: **80%削減**
- チーム全体の生産性: **20〜30%向上**

---

### 2. Fixture戦略
**ファイル**: [`fixture-strategy.md`](fixture-strategy.md)

**目的**: テスト用データの標準化とローカル開発の効率化

**Fixture作成方針**:

#### 優先度: 高（56テーブル）

**求人関連（最優先）**:
- `recruitments`: 100件（多様な条件）
- `recruitment_details`: 100件
- `companies`: 50社
- `recruitment_images`: 200件

**応募関連**:
- `applications`: 500件
- `application_histories`: 1,000件
- `application_messages`: 200件

**認証・ユーザー**:
- `admin_users`: 10件（各ロール）
- `company_users`: 30件
- `users`: 100件

**マスタデータ**:
- `prefectures`: 47件（全都道府県）
- `cities`: 1,896件（全市区町村）
- `stations`: 1,000件（主要駅）

#### 優先度: 中（17テーブル）

- 画像・ファイル関連
- 履歴・ログ（直近1ヶ月分のみ）

#### 優先度: 低（16テーブル）

- 設定・マスタ（デフォルト値のみ）
- 統計・集計用テーブル

#### 削除候補（24テーブル）

- `backup_*`, `tmp_*`, `legacy_*`
- 未使用テーブル

**Fixture生成方法**:

```bash
# 1. STG環境から匿名化データを抽出
./tools/analyze-table-usage.sh dump-fixtures

# 2. Fixtureファイル生成
php bin/cake.php bake fixture Recruitments --records 100

# 3. テストで使用
php bin/cake.php migrations seed --seed FixtureSeeder
```

**期待効果**:
- ローカル開発のリアリティ向上
- テスト実行時間の短縮
- 新規開発者の学習効率向上

---

## 🎯 統合改善ロードマップ

### Phase 1: DevContainer + Fixture（1-2ヶ月）

```
Week 1-2: DevContainer環境構築
├─ .devcontainer/ 作成
├─ PostgreSQL 15 統合
├─ MySQL 8.0 統合（延長サポート解消）
└─ docker-compose.yml 最適化

Week 3-4: Fixture作成（高優先度）
├─ 求人関連 Fixture
├─ 応募関連 Fixture
├─ 認証・ユーザー Fixture
└─ マスタデータ Fixture

Week 5-6: MySQL 8.0移行 + 二重管理解消
├─ PostgreSQL → MySQL 移行
├─ 同期バッチ廃止
└─ N+1クエリ解消

Week 7-8: ドキュメント・トレーニング
├─ セットアップガイド更新
├─ トラブルシューティング整備
└─ チームトレーニング
```

### Phase 2: テスト自動化（2-3ヶ月）

- Unit Test整備
- Integration Test整備
- E2E Test整備（Playwright）
- CI/CDパイプライン統合

### Phase 3: 本番環境改善（3-6ヶ月）

- Cloud Run移行
- オートスケーリング導入
- BOT対策強化

---

## 💰 コスト削減効果

| 施策 | 年間削減額 | 備考 |
|------|-----------|------|
| **MySQL 8.0アップグレード** | **¥164万円** | 延長サポート費用解消 |
| DB二重管理解消 | ¥700万円相当 | 保守工数削減 |
| DevContainer導入 | ¥200万円相当 | オンボーディング効率化 |
| **合計** | **約1,064万円/年** | |

---

## 📈 KPI（成功指標）

### 開発効率

| 指標 | 現状 | 目標 | 測定方法 |
|------|------|------|---------|
| 新規開発者セットアップ時間 | 1日 | 30分 | オンボーディング記録 |
| 環境起因トラブル件数 | 週5-10件 | 週1件以下 | Slackログ |
| ローカル環境起動成功率 | 60% | 95%以上 | アンケート |

### システム安定性

| 指標 | 現状 | 目標 | 測定方法 |
|------|------|------|---------|
| 瞬断発生頻度 | 週10-20回 | 月1回以下 | Cloud Monitoring |
| DB同期エラー率 | 5-10% | 0%（廃止） | ログ分析 |
| エラー率 | 最大80% | 1%以下 | APM |

### コスト

| 指標 | 現状 | 目標 | 測定方法 |
|------|------|------|---------|
| MySQL延長サポート費用 | ¥13.7万円/月 | ¥0 | 請求書 |
| 保守工数 | +40% | 標準 | 工数管理 |

---

## 🚀 次のアクション

### 最優先（今週中）
1. DevContainer環境のプロトタイプ作成
2. Fixture優先度の最終確認
3. MySQL 8.0検証環境セットアップ

### 高優先度（1ヶ月以内）
4. DevContainer正式版リリース
5. 高優先度Fixture作成完了
6. チームトレーニング実施

### 中優先度（3ヶ月以内）
7. MySQL 8.0本番適用
8. PostgreSQL → MySQL移行完了
9. 同期バッチ廃止

---

## 📚 関連ドキュメント

### 現状分析（As-Is）
- [`../as-is/`](../as-is/) - 現状の開発環境・構成

### 重要課題
- [`../important_issues/DB_INTEGRATION_TECHNICAL_RISK_REPORT.md`](../important_issues/DB_INTEGRATION_TECHNICAL_RISK_REPORT.md) - DB二重管理リスク
- [`../important_issues/SERVER_COST_REDUCTION_REPORT.md`](../important_issues/SERVER_COST_REDUCTION_REPORT.md) - サーバーコスト削減

### 戦略全体
- [`../../strategy/dorapita-org-transformation-strategy.md`](../../strategy/dorapita-org-transformation-strategy.md) - 組織変革戦略

---

**最終更新**: 2025-12-28  
**作成者**: Claude Code  
**レビュー**: 要確認
