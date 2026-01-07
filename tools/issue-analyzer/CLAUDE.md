# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリのコードを扱う際のガイダンスを提供します。

## プロジェクト概要

Dorapita Issue分析ツール - ZIGExN/dorapita リポジトリのissue履歴をSQLite3に保管し、パターン分析・関連付けを行うGo製CLIツール。

## ビルドシステム

- **Task** (Taskfile.yml) を使用
- **重要**: SQLite3使用のため、`CGO_ENABLED=1` でビルド必須

```bash
# ビルド
task build

# テスト（必ずtaskコマンド経由で実行）
task test

# 品質チェック
task check
```

## 依存関係

- `github.com/mattn/go-sqlite3` v1.14.32 のみ
- その他は標準ライブラリのみ使用

## コーディング規約

### 言語・ドキュメント
- すべてのドキュメント、コメント、gitコミットは日本語で記述
- セキュリティ: OWASP Top 10脆弱性の導入禁止

### 設計原則（過剰エンジニアリング回避）
- 要求された変更のみ実装
- 3行の類似コードは抽象化しない（早すぎる抽象化を回避）
- 発生し得ないシナリオのためのエラーハンドリング・バリデーションは削除
- 後方互換性ハックなし（未使用変数のリネーム、再エクスポート、削除コードのコメントなど）

### import整理
```bash
task fix  # goimports + go fmt 実行（コード修正後必須）
```

## テスト

### テスト実行

```bash
# 全テスト
task test

# カバレッジ付きテスト
task test-coverage

# 個別パッケージテスト
CGO_ENABLED=1 go test -v ./pkg/github/...
```

**重要**: `go test` を直接実行すると `CGO_ENABLED=0` でビルドされ、SQLiteテストが失敗します。必ず `task test` を使用してください。

### テスト構造

- ユニットテスト: `*_test.go`
- テストデータ: `testdata/` ディレクトリ
- モックサーバ: 外部API依存を排除

## GitHub API統合

- `~/.local/bin/gh-wrapper.sh` を使用してissue取得
- PAT自動抽出
- JSON形式の出力をパース

## SQLite3スキーマ

### テーブル構成

- `issues` - issue基本情報
- `labels` - issueラベル
- `comments` - issueコメント
- `issue_links` - issue間の参照関係

詳細は `.claude/rules/sqlite-schema.md` を参照。

## 分析機能

1. **パターン分析** (`pkg/analysis/pattern.go`)
   - 頻出キーワード抽出
   - ラベル別統計
   - 月別トレンド

2. **チーム対応分析** (`pkg/analysis/team.go`)
   - 担当者別応答時間
   - クローズ率
   - 月別活動量

3. **参照リンク分析** (`pkg/analysis/links.go`)
   - issue間の参照関係
   - クラスタ検出
   - Mermaidグラフ出力

## コマンド例

```bash
# issue取得（直近1年間）
./bin/issue-analyzer fetch \
    --repo ZIGExN/dorapita \
    --days 365 \
    --cache ./data/issues.db

# 分析実行
./bin/issue-analyzer analyze \
    --cache ./data/issues.db \
    --output ./out/report.md

# パターン分析のみ
./bin/issue-analyzer analyze \
    --cache ./data/issues.db \
    --type pattern \
    --output ./out/pattern-report.md
```

## ディレクトリ構造

```
.
├── cmd/
│   └── issue-analyzer/
│       └── main.go                    # CLIエントリポイント
├── pkg/
│   ├── github/                        # GitHub API統合
│   ├── cache/                         # SQLite3キャッシュ
│   ├── analysis/                      # 分析ロジック
│   └── output/                        # レポート出力
├── scripts/                           # ヘルパースクリプト
├── testdata/                          # テストデータ
├── .claude/rules/                     # Claude Codeルール
├── bin/                               # ビルド成果物
├── data/                              # データベースファイル
├── go.mod
├── Taskfile.yml
└── README.md
```

## 参照プロジェクト

このプロジェクトの構造とパターンは `/Users/kaz/git/loglass/sysdig-vuls-utils` を参考にしています。
