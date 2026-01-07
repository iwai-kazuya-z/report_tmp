# Tools: 開発支援ツール・スクリプト

開発環境の分析・管理に使用するツールとスクリプトを格納しています。

---

## 📋 ツール一覧

### 1. テーブル使用状況分析スクリプト
**ファイル**: [`analyze-table-usage.sh`](analyze-table-usage.sh)

**用途**: データベーステーブルの使用状況を分析

**主要機能**:

#### 1.1 テーブル一覧取得
```bash
./analyze-table-usage.sh list-tables
```
- PostgreSQL上の全テーブルを取得
- レコード数、サイズ情報を表示

#### 1.2 コードベース参照チェック
```bash
./analyze-table-usage.sh check-usage
```
- 各テーブルがコード内で参照されているかチェック
- CakePHP Model/Table クラスの存在確認
- SQL文での直接参照を検索

#### 1.3 未使用テーブル検出
```bash
./analyze-table-usage.sh find-unused
```
- コードから参照されていないテーブルをリスト化
- 削除候補の特定

#### 1.4 Fixture優先度判定
```bash
./analyze-table-usage.sh prioritize-fixtures
```
- テーブルのFixture作成優先度を判定
- 出力: 高/中/低/不要

#### 1.5 Fixtureデータダンプ（実装予定）
```bash
./analyze-table-usage.sh dump-fixtures --anonymize
```
- STG環境からFixture用データを抽出
- 個人情報の匿名化処理
- CakePHP Fixture形式で出力

**使用例**:
```bash
# 環境変数設定
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=dorapita
export DB_USER=postgres
export DB_PASS=postgres

# テーブル分析実行
./analyze-table-usage.sh check-usage > table-usage-report.txt

# 未使用テーブル確認
./analyze-table-usage.sh find-unused
```

**出力フォーマット**:
```
Table: recruitments
  Records: 12,345
  Size: 15.2 MB
  Status: ACTIVE
  References:
    - src/Model/Table/RecruitmentsTable.php
    - src/Controller/RecruitmentsController.php
    - templates/Recruitments/*.php (15 files)
  Fixture Priority: HIGH
```

---

### 2. 認証情報管理
**ファイル**: [`.secret`](.secret)

**⚠️ 重要**: このファイルは `.gitignore` に追加済み

**用途**: ローカル開発・分析用の認証情報

**形式**:
```bash
# Database Credentials
DB_HOST=localhost
DB_PORT=5432
DB_NAME=dorapita
DB_USER=dorauser2022
DB_PASS=xxxxx

# GCP Credentials
GCP_PROJECT_ID=dorapita-core-dev
GCP_SA_KEY=/path/to/service-account-key.json

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

**使用方法**:
```bash
# 環境変数として読み込み
source tools/.secret

# または、スクリプト内で読み込み
if [ -f tools/.secret ]; then
  source tools/.secret
fi
```

**セキュリティ**:
- ✅ `.gitignore` に追加済み
- ✅ `.secret.example` をテンプレートとして用意（推奨）
- ❌ **絶対にコミットしない**

---

## 🛠️ 開発中のツール

### Fixture生成ツール（計画中）
**予定ファイル**: `generate-fixtures.sh`

**機能**:
- STG環境からデータ抽出
- 個人情報の自動匿名化
- CakePHP Fixture形式で出力
- Seed用SQLも生成

### スキーマ比較ツール（計画中）
**予定ファイル**: `compare-schemas.sh`

**機能**:
- 本番 vs 開発環境のスキーマ差分検出
- Migration未適用の検出
- スキーマドリフトの警告

### パフォーマンス分析ツール（計画中）
**予定ファイル**: `analyze-performance.sh`

**機能**:
- スロークエリの検出
- N+1クエリの検出
- インデックス最適化の提案

---

## 📊 分析結果の活用

### 1. テーブル使用状況レポート
→ [`../as-is/database-analysis/table-usage-report.md`](../as-is/database-analysis/table-usage-report.md)

- `analyze-table-usage.sh` の実行結果
- 未使用テーブルのリスト
- 削除候補の特定

### 2. Fixture戦略
→ [`../to-be/fixture-strategy.md`](../to-be/fixture-strategy.md)

- 優先度判定結果の活用
- Fixture作成計画

### 3. データベース分析レポート
→ [`../as-is/database-analysis/stg-database-analysis-report.md`](../as-is/database-analysis/stg-database-analysis-report.md)

- 総合的なDB分析結果
- 改善提案

---

## 🚀 使い方（クイックスタート）

### Step 1: 認証情報設定
```bash
cp tools/.secret.example tools/.secret
vi tools/.secret  # 認証情報を編集
```

### Step 2: テーブル分析実行
```bash
cd /path/to/dorapita_code
source tools/.secret
./tools/analyze-table-usage.sh check-usage
```

### Step 3: 結果確認
```bash
# 未使用テーブル確認
./tools/analyze-table-usage.sh find-unused

# Fixture優先度確認
./tools/analyze-table-usage.sh prioritize-fixtures
```

---

## 📝 ツール追加ガイドライン

新しいツールを追加する際は、以下の要件を満たしてください：

1. **README.mdを更新**: 用途、使用方法を明記
2. **実行権限**: `chmod +x` で実行可能にする
3. **エラーハンドリング**: 適切なエラーメッセージを表示
4. **ヘルプ表示**: `--help` オプションを実装
5. **認証情報**: `.secret` を使用し、ハードコーディングしない
6. **出力形式**: 可能な限り構造化（JSON, CSV等）

---

## ⚠️ 注意事項

1. **本番環境での実行禁止**
   - これらのツールは開発・STG環境専用です
   - 本番環境での実行は厳禁

2. **認証情報の管理**
   - `.secret` は絶対にコミットしない
   - 必要に応じて暗号化ツール（sops, Vault等）の使用を検討

3. **リソース消費**
   - 大量データの分析はDB負荷に注意
   - 本番環境のレプリカで実行推奨

---

**最終更新**: 2025-12-28  
**メンテナー**: Infra Team
