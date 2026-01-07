# Dorapita Issue Analyzer

ZIGExN/dorapita リポジトリのissue履歴をSQLite3に保管し、パターン分析・関連付けを行うGo製CLIツール。

## 機能

### 1. Issue取得・保管
- GitHub issue履歴をgh-wrapper.sh経由で取得
- SQLite3データベースに保管
- 直近1年間のissueを対象（期間指定可能）

### 2. パターン分析
- 頻出キーワード抽出（タイトル・本文）
- ラベル別統計（issue数、平均クローズ時間）
- 月別issue作成トレンド
- open/closed分布

### 3. チーム対応分析
- 担当者別平均応答時間
- クローズ率
- 月別活動量（assignee別）

### 4. 参照リンク分析
- issue間の参照関係追跡（mentions、closes、fixesなど）
- 関連issueクラスタ検出
- Mermaidフォーマットでグラフ出力

## インストール

### 前提条件
- Go 1.23+
- Task (`go install github.com/go-task/task/v3/cmd/task@latest`)
- gh-wrapper.sh (`~/.local/bin/gh-wrapper.sh`)

### ビルド

```bash
# ビルド
task build

# ビルド成果物: bin/issue-analyzer
```

## 使い方

### Issue取得

```bash
# 直近1年間のissueを取得
./bin/issue-analyzer fetch \
    --repo ZIGExN/dorapita \
    --days 365 \
    --cache ./data/issues.db
```

### 分析実行

```bash
# 全分析を実行
./bin/issue-analyzer analyze \
    --cache ./data/issues.db \
    --output ./out/report.md

# パターン分析のみ
./bin/issue-analyzer analyze \
    --cache ./data/issues.db \
    --type pattern \
    --output ./out/pattern-report.md

# チーム対応分析のみ
./bin/issue-analyzer analyze \
    --cache ./data/issues.db \
    --type team \
    --output ./out/team-report.md

# 参照リンク分析のみ
./bin/issue-analyzer analyze \
    --cache ./data/issues.db \
    --type links \
    --output ./out/link-graph.md
```

## コマンドオプション

### fetchコマンド

| オプション | 説明 | デフォルト値 |
|-----------|------|-------------|
| `--repo` | リポジトリ名（owner/repo形式） | 必須 |
| `--days` | 取得期間（日数） | 365 |
| `--cache` | キャッシュファイルパス | `./issues.db` |
| `--gh-token` | GitHub Token | 自動取得 |

### analyzeコマンド

| オプション | 説明 | デフォルト値 |
|-----------|------|-------------|
| `--cache` | キャッシュファイルパス | `./issues.db` |
| `--type` | 分析タイプ（pattern, team, links, all） | `all` |
| `--output` | 出力ファイルパス | stdout |

## 開発

### テスト

```bash
# 全テスト
task test

# カバレッジ付きテスト
task test-coverage

# 品質チェック
task check
```

### コード修正後

```bash
# import整理 + フォーマット
task fix

# リント
task lint

# 自動修正付きリント
task lint-fix
```

## データベーススキーマ

### issues テーブル
| カラム | 型 | 説明 |
|-------|-----|------|
| number | INTEGER | issue番号（PK） |
| title | TEXT | タイトル |
| body | TEXT | 本文 |
| state | TEXT | 状態（open, closed） |
| created_at | TEXT | 作成日時 |
| updated_at | TEXT | 更新日時 |
| closed_at | TEXT | クローズ日時 |
| author_login | TEXT | 作成者 |
| author_type | TEXT | 作成者タイプ（User, Bot） |
| assignee_login | TEXT | 担当者 |
| url | TEXT | issue URL |

### labels テーブル
| カラム | 型 | 説明 |
|-------|-----|------|
| id | INTEGER | ID（PK） |
| issue_number | INTEGER | issue番号（FK） |
| name | TEXT | ラベル名 |
| color | TEXT | ラベル色 |

### comments テーブル
| カラム | 型 | 説明 |
|-------|-----|------|
| id | INTEGER | コメントID（PK） |
| issue_number | INTEGER | issue番号（FK） |
| author_login | TEXT | コメント作成者 |
| body | TEXT | コメント本文 |
| created_at | TEXT | 作成日時 |

### issue_links テーブル
| カラム | 型 | 説明 |
|-------|-----|------|
| id | INTEGER | ID（PK） |
| source_issue | INTEGER | リンク元issue（FK） |
| target_issue | INTEGER | リンク先issue（FK） |
| link_type | TEXT | リンクタイプ（mentions, closes, fixes, relates） |

## 出力例

### パターン分析レポート

```markdown
# Dorapita Issue分析レポート

生成日時: 2025-12-25 12:34:56
分析期間: 2024-12-25 〜 2025-12-25
総issue数: 342件

## 1. パターン分析

### 頻出キーワード Top 10
| キーワード | 出現回数 | 割合 |
|-----------|---------|------|
| エラー     | 45      | 13%  |
| 画像表示   | 32      | 9%   |

### ラベル別統計
| ラベル       | issue数 | 平均クローズ時間 |
|-------------|---------|----------------|
| bug         | 120     | 3.5日          |
| enhancement | 80      | 7.2日          |
```

## ライセンス

MIT License

## 参照

- [sysdig-vuls-utils](https://github.com/loglass/sysdig-vuls-utils) - 構造とパターンの参照元
- [gh-wrapper.sh](~/.local/bin/gh-wrapper.sh) - GitHub API呼び出し
