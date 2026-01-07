# Claude Code実践ガイド: 大規模プロジェクトの設計パターン

**最終更新**: 2025-12-25

---

## 目次

1. [はじめに](#1-はじめに)
2. [初期セットアップ: プロジェクト構成の基本設計](#2-初期セットアップ-プロジェクト構成の基本設計)
3. [ラッパースクリプトによる標準化](#3-ラッパースクリプトによる標準化)
4. [データ分析基盤の設計パターン](#4-データ分析基盤の設計パターン)
5. [ドキュメント駆動調査の実践](#5-ドキュメント駆動調査の実践)
6. [データ分析ワークフロー](#6-データ分析ワークフロー)
7. [応用例とユースケース](#7-応用例とユースケース)
8. [参考資料](#8-参考資料)

---

## 1. はじめに

### 対象読者

- Claude Code初心者〜中級者
- 大規模インフラ調査・ドキュメント作成プロジェクトを担当する方
- ラッパースクリプトやカスタムコマンドの活用に興味がある方

### このガイドで学べること

このガイドは、実践プロジェクト（約5日間、24ファイル・総行数10,000行超）での知見を基に、以下の実践パターンを提供します:

1. **初期セットアップ**: グローバル/プロジェクトCLAUDE.mdの設計、.claude/rules/による知識ベース構築
2. **ラッパースクリプト**: git/gh操作の標準化、安全性の向上、Claude Code統合
3. **データ分析基盤**: SQLite3によるローカルキャッシュ、柔軟な分析環境の構築
4. **ドキュメント駆動**: 大規模調査プロジェクトの設計、シナリオ別読み方ガイドの作成

### 実践プロジェクトの成果概要

- **ドキュメント**: 24ファイル、総行数10,000行超（インフラ調査、アプリ構造、組織戦略）
- **データベース**: SQLite3による分析基盤（387件のデータ、1,154件の詳細情報）
- **カスタムツール**: 独自開発のCLIツール、ラッパースクリプト2種
- **知識ベース**: .claude/rules/による組織・技術知識の外部化

---

## 2. 初期セットアップ: プロジェクト構成の基本設計

### 2.1 グローバルCLAUDE.mdの設計

**配置**: `~/.claude/CLAUDE.md`

グローバル設定は全プロジェクト共通の動作を定義します。

#### Tool Usage Priority

```markdown
## Tool Usage Priority

**重要: 専用ラッパーツールを優先的に使用する**

1. **Git操作**: `~/.local/bin/git-wrapper.sh`を使用（`git`コマンドを直接使わない）
2. **GitHub操作**: `~/.local/bin/gh-wrapper.sh`を使用（`gh`コマンドを直接使わない）
```

**ポイント**:
- ラッパースクリプトの**強制使用**により操作の安全性とログ記録を実現
- git/ghコマンドを直接使わず、必ずラッパー経由で実行

#### 日本語会話ルール

```markdown
## Conversation Guidelines

- 常に日本語で会話する
- gitコメントも日本語で記述する
```

#### Claude Code Permissions

```markdown
## Claude Code Permissions

以下のツールは自動承認対象として設定する：

- `Bash(~/.local/bin/git-wrapper.sh:*)` - Git操作ラッパー（安全な操作のみ自動実行）
- `Bash(~/.local/bin/gh-wrapper.sh:*)` - GitHub CLI操作ラッパー
```

---

## 3. ラッパースクリプトによる標準化

### 3.1 git-wrapper.sh: 安全なGit操作

#### 3.1.1 目的と設計思想

**なぜラッパーが必要か**:

1. **安全性**: 危険な操作（push, reset等）を実行前に警告
2. **ログ記録**: すべてのGit操作を記録（デバッグ・監査用）
3. **統一インターフェース**: チーム全体で一貫した操作方法
4. **Claude Code統合**: 安全な操作は自動承認、危険な操作のみインタラクティブ

**設計原則**:

- **読み取り専用操作**: 自動承認（status, diff, log, show, branch）
- **ローカル変更操作**: 自動承認（add, commit, checkout -b）
- **リモート操作**: インタラクティブ（push, pull, fetch）
- **破壊的操作**: インタラクティブ（reset, clean, rebase, merge）

#### 3.1.2 実装の完全ガイド

**スクリプト全体構造**（疑似コード）:

```bash
#!/bin/bash
#
# git-wrapper.sh - Safe Git operation wrapper
#

# 色定義
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m'

# ログ記録
log_operation() {
    local operation="$1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] git $operation" >> ~/.git-wrapper.log
}

# 安全な操作かチェック
is_safe_operation() {
    local cmd="$1"

    # 読み取り専用操作
    if [[ "$cmd" =~ ^(status|diff|log|show|branch|rev-parse)$ ]]; then
        return 0
    fi

    # ローカル変更操作
    if [[ "$cmd" == "add" || "$cmd" == "commit" ]]; then
        return 0
    fi

    # checkout -b（新規ブランチ作成）
    if [[ "$cmd" == "checkout" && "$2" == "-b" ]]; then
        return 0
    fi

    # config --get（設定取得）
    if [[ "$cmd" == "config" && "$2" == "--get" ]]; then
        return 0
    fi

    # remote -v（リモート一覧）
    if [[ "$cmd" == "remote" && "$2" == "-v" ]]; then
        return 0
    fi

    # それ以外は危険な操作
    return 1
}

# インタラクティブ警告
show_warning() {
    local operation="$*"
    echo -e "${YELLOW}[WARN]${NC} 危険な操作: $operation"
    echo -e "${YELLOW}[警告]${NC} 危険な操作を実行しようとしています: git $operation"
    echo -e "${YELLOW}[警告]${NC} この操作はインタラクティブに実行されます"
}

# 特別コマンド: new-pr-branch
new_pr_branch() {
    local timestamp=$(date +%Y%m%d%H%M%S)
    local branch_name="auto-pr-${timestamp}"
    echo "ブランチ '$branch_name' を作成し、切り替えました"
    git checkout -b "$branch_name"
}

# メイン処理
main() {
    # ヘルプ表示
    if [[ "$1" == "--help" || "$1" == "-h" ]]; then
        echo "Usage: git-wrapper.sh <git-command> [args...]"
        echo ""
        echo "Safe operations (auto-approved):"
        echo "  status, diff, log, show, branch, rev-parse"
        echo "  add, commit, checkout -b"
        echo ""
        echo "Interactive operations (requires confirmation):"
        echo "  push, pull, fetch, reset, clean, rebase, merge"
        echo ""
        echo "Special commands:"
        echo "  new-pr-branch  - Create auto-pr-TIMESTAMP branch"
        exit 0
    fi

    # 特別コマンド処理
    if [[ "$1" == "new-pr-branch" ]]; then
        new_pr_branch
        exit $?
    fi

    # 操作をログ記録
    log_operation "$*"

    # 安全な操作かチェック
    if is_safe_operation "$@"; then
        # 自動実行
        exec git "$@"
    else
        # 警告を表示してから実行
        show_warning "$@"
        exec git "$@"
    fi
}

main "$@"
```

#### 3.1.3 is_safe_operation関数の詳細

```bash
is_safe_operation() {
    local cmd="$1"
    shift
    local args="$*"

    # 読み取り専用操作
    case "$cmd" in
        status|diff|log|show|branch|rev-parse)
            return 0
            ;;
    esac

    # ローカル変更操作
    case "$cmd" in
        add|commit)
            return 0
            ;;
        checkout)
            # checkout -bは安全（新規ブランチ作成）
            if [[ "$1" == "-b" ]]; then
                return 0
            fi
            # 既存ブランチへの切り替えは警告
            return 1
            ;;
        config)
            # config --getは安全（取得のみ）
            if [[ "$1" == "--get" ]]; then
                return 0
            fi
            # config設定変更は警告
            return 1
            ;;
        remote)
            # remote -vは安全（一覧表示）
            if [[ "$1" == "-v" ]]; then
                return 0
            fi
            return 1
            ;;
    esac

    # それ以外は危険な操作として扱う
    return 1
}
```

#### 3.1.4 特別コマンド: new-pr-branch

**目的**: auto-pr-TIMESTAMP形式のブランチを自動作成

```bash
new_pr_branch() {
    local timestamp=$(date +%Y%m%d%H%M%S)
    local branch_name="auto-pr-${timestamp}"

    # 現在のブランチを確認
    local current_branch=$(git branch --show-current)
    echo "現在のブランチ: $current_branch"

    # 新しいブランチを作成
    echo "ブランチ '$branch_name' を作成し、切り替えました"
    git checkout -b "$branch_name"

    return $?
}
```

**使用例**:

```bash
# mainブランチから新しいPRブランチを作成
~/.local/bin/git-wrapper.sh checkout main
~/.local/bin/git-wrapper.sh new-pr-branch
# → auto-pr-20251225193332 が作成される
```

#### 3.1.5 Claude Code統合

**~/.claude/CLAUDE.mdへの記載**:

```markdown
## Claude Code Permissions

以下のツールは自動承認対象として設定する：

- `Bash(~/.local/bin/git-wrapper.sh:*)` - Git操作ラッパー
```

**動作確認**:

```bash
# 1. スクリプトを配置
cp git-wrapper.sh ~/.local/bin/
chmod +x ~/.local/bin/git-wrapper.sh

# 2. 安全な操作（自動実行）
~/.local/bin/git-wrapper.sh status
~/.local/bin/git-wrapper.sh diff HEAD
~/.local/bin/git-wrapper.sh add .
~/.local/bin/git-wrapper.sh commit -m "test"

# 3. 危険な操作（警告表示）
~/.local/bin/git-wrapper.sh push origin main
# → [WARN] 危険な操作: push origin main
# → [警告] 危険な操作を実行しようとしています...
```

### 3.2 gh-wrapper.sh: GitHub API統合

#### 3.2.1 目的と設計思想

**なぜgh-wrapper.shが必要か**:

1. **PAT自動抽出**: `.git/config`からPersonal Access Tokenを自動抽出、手動設定不要
2. **JSON形式出力**: 構造化データで`jq`等と連携しやすい
3. **統一インターフェース**: PR、Issue、リポジトリ操作を統一された方法で実行
4. **Claude Code統合**: GitHub操作を標準化

#### 3.2.2 PAT自動抽出ロジックの完全実装

**スクリプト構造**（疑似コード）:

```bash
#!/bin/bash
#
# gh-wrapper.sh - GitHub CLI wrapper with PAT auto-extraction
#

# .git/configからPAT抽出
extract_pat_from_git_config() {
    local git_config_url
    git_config_url=$(git config --get remote.origin.url 2>/dev/null)

    if [[ -z "$git_config_url" ]]; then
        return 1
    fi

    # 対応フォーマット:
    # https://TOKEN@github.com/owner/repo.git
    # https://username:TOKEN@github.com/owner/repo.git

    local pat=""

    # パターン1: https://TOKEN@github.com
    if [[ "$git_config_url" =~ https://([^@:]+)@github\.com ]]; then
        pat="${BASH_REMATCH[1]}"
    # パターン2: https://username:TOKEN@github.com
    elif [[ "$git_config_url" =~ https://[^:]+:([^@]+)@github\.com ]]; then
        pat="${BASH_REMATCH[1]}"
    else
        return 1
    fi

    echo "$pat"
    return 0
}

# PAT設定
setup_gh_token() {
    local verbose="$1"

    # 1. .git/configから抽出（優先）
    local pat
    pat=$(extract_pat_from_git_config)

    if [[ -n "$pat" ]]; then
        export GH_TOKEN="$pat"
        if [[ "$verbose" == "true" ]]; then
            echo "[DEBUG] PAT extracted from .git/config" >&2
        fi
        return 0
    fi

    # 2. フォールバック: gh auth token
    if command -v gh &>/dev/null; then
        pat=$(gh auth token 2>/dev/null)
        if [[ -n "$pat" ]]; then
            export GH_TOKEN="$pat"
            if [[ "$verbose" == "true" ]]; then
                echo "[DEBUG] PAT extracted from gh auth token" >&2
            fi
            return 0
        fi
    fi

    # 3. 失敗
    if [[ "$verbose" == "true" ]]; then
        echo "[ERROR] Failed to extract PAT" >&2
        echo "[ERROR] Please check .git/config or run 'gh auth login'" >&2
    fi
    return 1
}

# メイン処理
main() {
    local verbose=false
    local raw_mode=false
    local args=()

    # オプション解析
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --verbose)
                verbose=true
                shift
                ;;
            --raw)
                raw_mode=true
                shift
                ;;
            --help|-h)
                echo "Usage: gh-wrapper.sh [--verbose] [--raw] <gh-command> [args...]"
                echo ""
                echo "Options:"
                echo "  --verbose  Show debug information"
                echo "  --raw      Use gh CLI default output format"
                echo ""
                echo "Examples:"
                echo "  gh-wrapper.sh pr list --state all"
                echo "  gh-wrapper.sh --raw pr diff 123"
                exit 0
                ;;
            *)
                args+=("$1")
                shift
                ;;
        esac
    done

    # PAT設定
    setup_gh_token "$verbose" || exit 1

    # gh CLI実行
    if [[ "$raw_mode" == "true" ]]; then
        # 生出力
        gh "${args[@]}"
    else
        # JSON形式（デフォルト）
        # 注: すべてのコマンドがJSONをサポートしているわけではない
        gh "${args[@]}"
    fi
}

main "$@"
```

#### 3.2.3 JSON形式出力の実装

**--rawオプションなし（デフォルト）**:

```bash
# PR一覧（JSON形式）
~/.local/bin/gh-wrapper.sh pr list --state all --json number,title,state,author

# 出力例（JSON）
[
  {
    "number": 12,
    "title": "Add feature X",
    "state": "OPEN",
    "author": {"login": "user1"}
  },
  {
    "number": 11,
    "title": "Fix bug Y",
    "state": "MERGED",
    "author": {"login": "user2"}
  }
]

# jqと連携
~/.local/bin/gh-wrapper.sh pr list --state all --json number,title,state \
  | jq '.[] | select(.state == "OPEN") | .number'
```

**--rawオプションあり（表形式）**:

```bash
# PR一覧（表形式）
~/.local/bin/gh-wrapper.sh --raw pr list --state all

# 出力例（表形式）
#12  Add feature X      OPEN   2025-12-25
#11  Fix bug Y          MERGED 2025-12-24
```

#### 3.2.4 使用例集

**PR操作**:

```bash
# PR一覧（オープン中）
~/.local/bin/gh-wrapper.sh pr list --state open --json number,title,author

# PR一覧（クローズ済み）
~/.local/bin/gh-wrapper.sh pr list --state closed --limit 10

# PR一覧（全状態）
~/.local/bin/gh-wrapper.sh pr list --state all --json number,title,state,createdAt

# PR詳細
~/.local/bin/gh-wrapper.sh pr view 12 --json number,title,body,state

# PR差分
~/.local/bin/gh-wrapper.sh --raw pr diff 12

# PR作成
~/.local/bin/gh-wrapper.sh pr create --title "Add feature X" --body "Description..."

# PR編集
~/.local/bin/gh-wrapper.sh pr edit 12 --title "Updated title"

# PRコメント
~/.local/bin/gh-wrapper.sh pr comment 12 --body "LGTM"
```

**Issue操作**:

```bash
# Issue一覧
~/.local/bin/gh-wrapper.sh issue list --json number,title,state

# Issue詳細
~/.local/bin/gh-wrapper.sh issue view 100 --json number,title,body,comments

# Issue作成
~/.local/bin/gh-wrapper.sh issue create --title "Bug report" --body "..."
```

**リポジトリ情報**:

```bash
# リポジトリ情報
~/.local/bin/gh-wrapper.sh repo view --json name,description,url

# GitHub API直接呼び出し
~/.local/bin/gh-wrapper.sh api repos/owner/repo/pulls/123
```

#### 3.2.5 Claude Code統合

**~/.claude/CLAUDE.mdへの記載**:

```markdown
## Claude Code Permissions

- `Bash(~/.local/bin/gh-wrapper.sh:*)` - GitHub CLI操作ラッパー
```

**Claude Code使用時の推奨パターン**:

| ユーザー指示 | 実行コマンド |
|--------------|--------------|
| 「PRを確認して」「PRを見て」 | `~/.local/bin/gh-wrapper.sh pr list --state all --json number,title,state` |
| 「クローズ済みPRを確認」 | `~/.local/bin/gh-wrapper.sh pr list --state closed` |
| 「PR #12の詳細を見せて」 | `~/.local/bin/gh-wrapper.sh pr view 12 --json number,title,body` |
| 「PRを作成して」 | 変更内容確認後、`~/.local/bin/gh-wrapper.sh pr create --title "..." --body "..."` |
| 「PRを編集して」 | `~/.local/bin/gh-wrapper.sh pr edit <番号> --title "..." または --body "..."` |
| 「リポジトリ情報を見せて」 | `~/.local/bin/gh-wrapper.sh repo view --json name,description` |
| 「Issue一覧を見せて」 | `~/.local/bin/gh-wrapper.sh issue list --json number,title,state` |

#### 3.2.6 トラブルシューティング

**問題1: PAT抽出失敗**

```bash
# デバッグモードで実行
~/.local/bin/gh-wrapper.sh --verbose pr list
# → [DEBUG] PAT extracted from .git/config
# または
# → [ERROR] Failed to extract PAT
```

**解決方法**:

```bash
# .git/configを確認
git config --get remote.origin.url
# → https://TOKEN@github.com/owner/repo.git
#    または
#    https://username:TOKEN@github.com/owner/repo.git

# URLが正しくない場合は修正
git config remote.origin.url https://YOUR_TOKEN@github.com/owner/repo.git

# または gh auth loginで再認証
gh auth login
```

**問題2: GitHub API rate limit**

```bash
# 認証状態を確認
gh auth status

# レート制限を確認
~/.local/bin/gh-wrapper.sh api rate_limit
```

### 3.3 カスタムスラッシュコマンドとの統合

**配置**: `~/.claude/commands/git:pr.md`

```markdown
---
allowed-tools:
  - Bash(~/.local/bin/git-wrapper.sh:*)
  - Bash(~/.local/bin/gh-wrapper.sh:*)
description: gitの差分を元にコミット、新しいブランチを作成、プッシュ、PR作成までを自動化します。
---

# Git差分を元にPRを作成する

## Variables
- PR_TITLE: $ARGUMENTS
- TARGET_BRANCH: ${GH_TARGET_BRANCH:-main}

## Instruction
1. 現在のGitの状況を確認してください。
   !~/.local/bin/git-wrapper.sh status --short

2. mainブランチにいる場合は、新しいPRブランチを作成してください。
   !if [ "$(git branch --show-current)" = "main" ]; then ~/.local/bin/git-wrapper.sh new-pr-branch; fi

3. 変更内容を理解するために、差分を分析してください。
   !~/.local/bin/git-wrapper.sh diff HEAD

4. すべての変更をステージングしてください。
   !~/.local/bin/git-wrapper.sh add .

5. ステージングされた変更に基づいて、変更内容を簡潔に要約したコミットメッセージを提案してください。

6. AIが提案したコミットメッセージでコミットしてください。
   !~/.local/bin/git-wrapper.sh commit -m "<コミットメッセージ>"

7. 新しいブランチをプッシュしてください。
   !git push -u origin $(git branch --show-current)

8. gh-wrapper.shを使って、PRを作成してください。
   !~/.local/bin/gh-wrapper.sh pr create --title "$PR_TITLE" --body "$ARGUMENTS" --base $TARGET_BRANCH
```

**使用例**:

```
/git:pr 新機能Xを追加
```

**Variables（環境変数）の活用**:

- `$ARGUMENTS`: ユーザーが入力した引数（PR_TITLEに代入）
- `${GH_TARGET_BRANCH:-main}`: 環境変数GH_TARGET_BRANCHが設定されていればそれを使用、なければmainをデフォルト値として使用

---

## 4. データ分析基盤の設計パターン

### 4.1 SQLite3によるローカルキャッシュ

#### 4.1.1 なぜSQLite3か

**メリット**:

1. **ポータビリティ**: 単一ファイルで完結、どの環境でも動作
2. **標準ツール**: `sqlite3`コマンドが標準搭載、追加インストール不要
3. **高速**: ローカルディスク上で動作、ネットワーク遅延なし
4. **軽量**: 数MB〜数十MBのデータ分析に最適
5. **SQL標準**: 複雑な集計クエリが簡単に書ける

**適用シーン**:

- GitHub Issue/PRの分析
- ログ分析
- コミット履歴の分析
- API応答のキャッシュ

#### 4.1.2 スキーマ設計の原則

**例: GitHub Issue分析DB**

```sql
-- issues.db スキーマ設計例

-- Issue本体テーブル
CREATE TABLE issues (
    number INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT,
    state TEXT NOT NULL,          -- open, closed
    author TEXT,
    created_at TEXT,
    closed_at TEXT,
    labels TEXT,                  -- JSON形式で保存
    fetched_at TEXT DEFAULT CURRENT_TIMESTAMP
);

-- コメントテーブル
CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    issue_number INTEGER NOT NULL,
    author TEXT,
    body TEXT,
    created_at TEXT,
    FOREIGN KEY (issue_number) REFERENCES issues(number)
);

-- 参照リンクテーブル（Issue間の関係）
CREATE TABLE links (
    source_number INTEGER NOT NULL,
    target_number INTEGER NOT NULL,
    link_type TEXT,               -- references, closes, duplicate
    PRIMARY KEY (source_number, target_number)
);

-- 検索用インデックス
CREATE INDEX idx_issues_state ON issues(state);
CREATE INDEX idx_issues_author ON issues(author);
CREATE INDEX idx_issues_created ON issues(created_at);
CREATE INDEX idx_comments_issue ON comments(issue_number);
```

**設計ポイント**:

- **正規化**: 重複データを避ける（Issue本体とコメントを分離）
- **インデックス**: よく検索するカラムにインデックスを作成
- **TEXT型**: SQLite3ではTEXT型を積極的に使用（日付もTEXT可）
- **JSON格納**: 配列データはJSON文字列で保存

#### 4.1.3 キャッシュ戦略

**初回取得**:

```bash
# GitHub APIからデータ取得してSQLite3に保存
~/.local/bin/gh-wrapper.sh issue list --state all --limit 1000 --json number,title,state,author,createdAt \
  | jq -r '.[] | [.number, .title, .state, .author.login, .createdAt] | @tsv' \
  | while IFS=$'\t' read -r num title state author created; do
      sqlite3 issues.db "INSERT OR REPLACE INTO issues (number, title, state, author, created_at) VALUES ($num, '$title', '$state', '$author', '$created')"
    done
```

**差分更新**:

```bash
# 最終取得日時以降のデータのみ取得
last_fetched=$(sqlite3 issues.db "SELECT MAX(fetched_at) FROM issues")
# APIで since パラメータを使用
~/.local/bin/gh-wrapper.sh issue list --state all --since "$last_fetched" ...
```

### 4.2 Claude Codeとの統合

#### 4.2.1 「データを分析して」指示パターン

**ユーザーの指示**:
```
「Issueを分析して、月別のオープン数を教えて」
```

**Claude Codeの動作**:

```bash
sqlite3 issues.db "
SELECT strftime('%Y-%m', created_at) as month,
       COUNT(*) as total
FROM issues
WHERE state = 'open'
GROUP BY month
ORDER BY month DESC
LIMIT 12"
```

#### 4.2.2 .claude/rules/にクエリ例を記載

**配置**: `.claude/rules/tool-usage.md`

```markdown
## データ分析ツール

### ユーザーが「Issueを分析して」と指示した場合

**重要: `issues.db`を使ってSQLite3クエリで直接分析する**

1. **データベース**: `./issues.db`
2. **クエリ方法**: `sqlite3 issues.db "SELECT ..."`を使用
3. **分析視点**: ユーザーの質問に応じて適切なSQLクエリを実行

#### クエリ例

**月別オープン数**:
\```bash
sqlite3 issues.db "
SELECT strftime('%Y-%m', created_at) as month,
       COUNT(*) as total
FROM issues
WHERE state = 'open'
GROUP BY month
ORDER BY month DESC"
\```

**開発者別アクティビティ**:
\```bash
sqlite3 issues.db "
SELECT author,
       COUNT(*) as issue_count,
       SUM(CASE WHEN state = 'closed' THEN 1 ELSE 0 END) as closed_count
FROM issues
GROUP BY author
ORDER BY issue_count DESC"
\```

**キーワード検索**:
\```bash
sqlite3 issues.db "
SELECT number, title, state
FROM issues
WHERE title LIKE '%keyword%' OR body LIKE '%keyword%'
ORDER BY created_at DESC
LIMIT 20"
\```
```

#### 4.2.3 柔軟な分析が可能

**複雑な集計クエリ**:

```bash
# Issue のライフサイクル分析（作成からクローズまでの日数）
sqlite3 issues.db "
SELECT author,
       AVG(JULIANDAY(closed_at) - JULIANDAY(created_at)) as avg_days_to_close,
       COUNT(*) as total_closed
FROM issues
WHERE state = 'closed' AND closed_at IS NOT NULL
GROUP BY author
HAVING total_closed >= 5
ORDER BY avg_days_to_close ASC"
```

**JOIN による関連分析**:

```bash
# コメント数の多いIssue分析
sqlite3 issues.db "
SELECT i.number,
       i.title,
       COUNT(c.id) as comment_count
FROM issues i
LEFT JOIN comments c ON i.number = c.issue_number
GROUP BY i.number
ORDER BY comment_count DESC
LIMIT 10"
```

### 4.3 カスタムCLIツール開発のヒント

#### 4.3.1 Go/Python等での実装パターン

**Go言語での実装例**（概念のみ）:

```go
// tools/data-analyzer/main.go
package main

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, _ := sql.Open("sqlite3", "./issues.db")
    defer db.Close()

    // データ取得
    rows, _ := db.Query("SELECT number, title FROM issues WHERE state = 'open'")
    defer rows.Close()

    // 分析処理
    for rows.Next() {
        var number int
        var title string
        rows.Scan(&number, &title)
        // 分析ロジック
    }
}
```

**ビルド時の注意点**（Go + SQLite3）:

```bash
# CGO_ENABLED=1 が必須
CGO_ENABLED=1 go build -o bin/analyzer main.go
```

#### 4.3.2 gh-wrapper.sh連携によるGitHub API呼び出し

**Goからgh-wrapper.shを呼び出す**:

```go
cmd := exec.Command("~/.local/bin/gh-wrapper.sh", "issue", "list", "--json", "number,title,state")
output, _ := cmd.Output()

var issues []Issue
json.Unmarshal(output, &issues)
```

#### 4.3.3 Taskfile.yml/Makefileによるビルド自動化

**Taskfile.yml の例**:

```yaml
version: '3'

tasks:
  build:
    desc: Build the analyzer tool
    cmds:
      - CGO_ENABLED=1 go build -o bin/analyzer cmd/analyzer/main.go

  fetch:
    desc: Fetch latest data from GitHub
    cmds:
      - ./bin/analyzer fetch

  analyze:
    desc: Analyze data
    cmds:
      - ./bin/analyzer analyze --report daily

  all:
    desc: Build, fetch, and analyze
    cmds:
      - task: build
      - task: fetch
      - task: analyze
```

**使用例**:

```bash
# ビルド
task build

# データ取得
task fetch

# 分析
task analyze

# 一括実行
task all
```

---

## 5. ドキュメント駆動調査の実践

### 5.1 複数レイヤーのドキュメント構成

大規模プロジェクトでは、ドキュメントをレイヤーに分けて管理することで、読みやすさと保守性を向上させます。

#### 5.1.1 Layer 1: インフラレイヤー（500-1600行/ファイル）

**対象**: GCPリソース、AWS環境、ネットワーク構成、セキュリティ設定

```
infrastructure/
├── inventory.md              # 全体概要（800行）
├── network-architecture.md   # ネットワーク構成（600行）
├── security-audit.md         # セキュリティ監査（1200行）
└── resource-list.md          # リソース一覧（1600行）
```

**記載内容**:

- プロジェクト一覧
- Compute Engine / Cloud Run / Cloud Functions
- Cloud SQL / Memorystore
- VPC / Firewall / Load Balancer
- IAM / Secret Manager
- Cloud Storage / Filestore

**ドキュメント構造例**:

```markdown
# インフラ調査レポート

## 1. プロジェクト概要

| 環境 | Project ID | リージョン | 用途 |
|------|-----------|----------|------|
| 本番 | project-prod | asia-northeast1 | Production |
| 開発 | project-dev | asia-northeast1 | Staging |

## 2. コンピューティング

### 2.1 Compute Engine

| VM名 | マシンタイプ | 役割 | 備考 |
|------|-------------|------|------|
| web-01 | n1-standard-2 | Webサーバー | Nginx, PHP |
| db-01 | n1-highmem-4 | DBサーバー | PostgreSQL 10 |

### 2.2 Cloud Run

| サービス名 | リージョン | CPU | メモリ | 備考 |
|-----------|----------|-----|--------|------|
| api-service | asia-northeast1 | 1 | 512Mi | REST API |

## 3. データベース

### 3.1 Cloud SQL

| インスタンス名 | バージョン | マシンタイプ | HA | 備考 |
|---------------|----------|-------------|-----|------|
| main-db | PostgreSQL 10 | db-n1-standard-2 | 有効 | EOL注意 |

## 4. セキュリティ

### 4.1 IAM設定

（詳細な権限マトリクス）

### 4.2 ファイアウォール

（ルール一覧と評価）
```

#### 5.1.2 Layer 2: アプリケーションレイヤー（300-800行/ファイル）

**対象**: アプリケーション構造、API仕様、デプロイフロー、開発環境

```
apps/
├── README.md                      # 全体構造（400行）
├── local-development.md           # ローカル環境（500行）
├── ci-cd-pipeline.md              # CI/CD詳細（600行）
├── deployment-comparison.md       # 開発vs本番デプロイ（400行）
├── main-site.md                   # メインサイト詳細（800行）
├── admin-site.md                  # 管理画面詳細（700行）
└── api-service.md                 # APIサービス詳細（600行）
```

**記載内容**:

- 技術スタック（フレームワーク、言語、ライブラリ）
- ディレクトリ構造
- ローカル開発環境セットアップ
- Docker Compose構成
- CI/CDパイプライン（Cloud Build, GitHub Actions）
- デプロイ方法（自動vs手動、開発vs本番）
- API仕様（エンドポイント、認証、レート制限）

**ドキュメント構造例**:

```markdown
# アプリケーション全体構造

## 1. 技術スタック

| コンポーネント | 技術 | バージョン |
|---------------|------|-----------|
| バックエンド | PHP, CakePHP | 8.1, 4.4 |
| フロントエンド | JavaScript, jQuery | - |
| データベース | PostgreSQL, MySQL | 10, 5.7 |
| キャッシュ | Redis | 6.x |

## 2. アプリケーション構成

```mermaid
graph LR
    A[User] --> B[Load Balancer]
    B --> C[Web Server]
    C --> D[App Container]
    D --> E[Cloud SQL]
    D --> F[Redis]
```
```

#### 5.1.3 Layer 3: 組織・戦略レイヤー（500-900行/ファイル）

**対象**: 組織構造、チーム分析、変革戦略、RACI

```
strategy/
├── transformation-strategy.md     # 組織変革戦略（900行）
├── team-analysis/
│   ├── developer-activity.md      # 開発者アクティビティ（800行）
│   ├── technical-evaluation.md    # 技術力評価（700行）
│   └── contract-roster.md         # 契約名簿分析（500行）
└── phase-plan.md                  # Phase計画と進捗（600行）
```

**記載内容**:

- 会社組織図（親会社・子会社・ベンダー関係）
- ステークホルダーマップ
- RACIマトリクス（現状 vs 目標）
- チーム実態評価（技術力、稼働率、問題点）
- 変革戦略（Phase計画、優先順位、ブロッカー）
- コスト分析（実質単価、削減効果）

**ドキュメント構造例**:

```markdown
# 組織変革戦略

## 1. 現状の組織構造

```mermaid
graph TD
    A[Parent Company] --> B[Company A]
    B --> C[Vendor X]
```

## 2. RACIマトリクス（現状）

| フェーズ | Company A | Vendor X |
|----------|-----------|----------|
| 企画 | R/A | I |
| 設計 | C | R/A |
| 開発 | I | R/A |
| デプロイ | R/A | C |

## 3. 課題一覧

| ID | 課題 | 症状 | 影響度 |
|----|------|------|--------|
| C1 | 専任PdM不在 | ビジネス全体像を把握している人がいない | 高 |
| C2 | adhoc調査で忙殺 | 本来のビジネス価値に向き合えない | 高 |

## 4. Phase計画

| Phase | 領域 | 状態 | 進捗 |
|-------|------|------|------|
| 1 | インフラ層把握 | 完了 | 100% |
| 2 | データ層把握 | 部分完了 | 40% |
| 3 | アプリケーション層把握 | 完了 | 90% |

## 5. チーム実態評価

### 5.1 開発者評価

| 名前 | 役割 | 技術力 | 実態評価 | 推奨アクション |
|------|------|--------|---------|---------------|
| Developer A | Senior Dev | ★★★★★ | ✅ フル稼働 | 継続・昇給検討 |
| Developer B | Tech Lead | ☆☆☆☆☆ | ❌ 不活性 | 契約解除検討 |
```

#### 5.1.4 Layer 4: データ分析レイヤー（500-1000行/ファイル）

**対象**: Issue分析、PR統計、コミット履歴、コスト分析

```
analysis/
├── issue-analysis.md              # Issue分析（800行）
├── pr-statistics.md               # PR統計（600行）
├── commit-history.md              # コミット履歴分析（700行）
└── cost-analysis.md               # コスト分析（500行）
```

**記載内容**:

- Issue/PR統計（月別、開発者別、タグ別）
- コミット履歴分析（LOC変更、ファイル数、頻度）
- コスト分析（名目稼働、実質稼働、単価）
- 技術的負債の定量化
- データ可視化（グラフ、チャート）

### 5.2 調査結果のMarkdown化戦略

#### 5.2.1 テーブル駆動

**複雑な情報を表で整理**:

```markdown
## Compute Engine一覧

| VM名 | マシンタイプ | CPU | メモリ | ディスク | OS | 用途 | 備考 |
|------|-------------|-----|--------|---------|-----|------|------|
| web-01 | n1-standard-2 | 2 | 7.5GB | 100GB | Ubuntu 20.04 | Webサーバー | Nginx, PHP-FPM |
| web-02 | n1-standard-2 | 2 | 7.5GB | 100GB | Ubuntu 20.04 | Webサーバー | Nginx, PHP-FPM |
| db-01 | n1-highmem-4 | 4 | 26GB | 500GB | Ubuntu 20.04 | DBサーバー | PostgreSQL 10 |
```

**評価・判定をシンボルで表現**:

```markdown
| 項目 | 状態 | 評価 | 備考 |
|------|------|------|------|
| SSL証明書 | 有効 | ✅ | Let's Encrypt, 2025-12-31まで |
| TLS 1.2+ | 有効 | ✅ | TLS 1.0/1.1無効化済み |
| Cloud Armor | 未設定 | ⚠️ | DDoS対策が不十分 |
| Secret Manager | 未使用 | ❌ | 環境変数に平文保存 |
```

#### 5.2.2 コードブロック

**gcloudコマンド、SQLクエリ、設定ファイル**:

````markdown
### Cloud SQL接続方法

```bash
# Cloud SQL Proxyを使用
cloud_sql_proxy -instances=project-prod:asia-northeast1:main-db=tcp:5432

# 別ターミナルでpsqlで接続
psql -h 127.0.0.1 -U dbuser -d maindb
```

### スキーマ確認

```sql
-- テーブル一覧
\dt

-- users テーブルの構造
\d users

-- データ件数
SELECT COUNT(*) FROM users;
```
````

#### 5.2.3 Mermaid図の活用

**アーキテクチャ図**:

```markdown
```mermaid
graph LR
    User[ユーザー] --> LB[Load Balancer]
    LB --> VM1[web-01]
    LB --> VM2[web-02]
    VM1 --> DB[(Cloud SQL)]
    VM2 --> DB
    VM1 --> Redis[(Memorystore)]
    VM2 --> Redis
```
```

**フローチャート**:

```markdown
```mermaid
flowchart TD
    A[PR作成] --> B{CIパス?}
    B -->|Yes| C[レビュー依頼]
    B -->|No| A
    C --> D{Approve?}
    D -->|Yes| E[マージ]
    D -->|No| F[修正]
    F --> A
    E --> G[自動デプロイ]
```
```

**シーケンス図**:

```markdown
```mermaid
sequenceDiagram
    participant User
    participant Web
    participant API
    participant DB

    User->>Web: リクエスト
    Web->>API: API呼び出し
    API->>DB: クエリ実行
    DB-->>API: 結果返却
    API-->>Web: JSON返却
    Web-->>User: HTML表示
```
```

**ER図（Entity-Relationship）**:

```markdown
```mermaid
erDiagram
    USERS ||--o{ ORDERS : places
    USERS {
        int id PK
        string name
        string email UK
    }
    ORDERS {
        int id PK
        int user_id FK
        datetime created_at
    }
```
```

---

## 6. データ分析ワークフロー

### 6.1 Claude Codeでの指示パターン

#### 6.1.1 基本的な分析指示

**ユーザー指示**: 「Issueを分析して」

**Claude Codeの動作**:

1. `.claude/rules/tool-usage.md`を参照
2. `issues.db`が存在することを確認
3. SQLite3クエリで直接分析
4. 結果を整形して提示

#### 6.1.2 段階的な分析

```
ユーザー: 「Issueの傾向を分析して」

Claude:
1. 総数・状態別カウント
2. 月別の推移
3. 開発者別の分布
4. キーワード頻度分析
```

### 6.2 SQLite3クエリ例集

#### 6.2.1 キーワード検索

```bash
# タイトルまたは本文にキーワードを含むIssue
sqlite3 issues.db "
SELECT number, title, state, author
FROM issues
WHERE title LIKE '%bug%' OR body LIKE '%bug%'
ORDER BY created_at DESC
LIMIT 20"
```

#### 6.2.2 月別集計

```bash
# 月別のIssue作成数
sqlite3 issues.db "
SELECT strftime('%Y-%m', created_at) as month,
       COUNT(*) as total,
       SUM(CASE WHEN state = 'open' THEN 1 ELSE 0 END) as open_count,
       SUM(CASE WHEN state = 'closed' THEN 1 ELSE 0 END) as closed_count
FROM issues
GROUP BY month
ORDER BY month DESC
LIMIT 12"
```

#### 6.2.3 開発者別統計

```bash
# 開発者別のIssue作成・クローズ数
sqlite3 issues.db "
SELECT author,
       COUNT(*) as total_issues,
       SUM(CASE WHEN state = 'open' THEN 1 ELSE 0 END) as open_issues,
       SUM(CASE WHEN state = 'closed' THEN 1 ELSE 0 END) as closed_issues,
       ROUND(100.0 * SUM(CASE WHEN state = 'closed' THEN 1 ELSE 0 END) / COUNT(*), 1) as close_rate
FROM issues
GROUP BY author
ORDER BY total_issues DESC
LIMIT 10"
```

#### 6.2.4 参照関係の抽出

```bash
# Issue間の参照関係を分析
sqlite3 issues.db "
SELECT source_number,
       target_number,
       link_type,
       (SELECT title FROM issues WHERE number = source_number) as source_title,
       (SELECT title FROM issues WHERE number = target_number) as target_title
FROM links
WHERE link_type = 'closes'
LIMIT 20"
```

#### 6.2.5 コメント数の多いIssue

```bash
# ディスカッションが活発なIssue
sqlite3 issues.db "
SELECT i.number,
       i.title,
       COUNT(c.id) as comment_count,
       i.state
FROM issues i
LEFT JOIN comments c ON i.number = c.issue_number
GROUP BY i.number
ORDER BY comment_count DESC
LIMIT 10"
```

### 6.3 分析レポート作成フロー

#### 6.3.1 フェーズ1: データ収集

```bash
# GitHub APIからデータ取得
~/.local/bin/gh-wrapper.sh issue list --state all --limit 500 --json number,title,state,author,createdAt,closedAt,labels \
  | jq -c '.[]' \
  | while read issue; do
      # SQLite3にINSERT
      # （スクリプトで実装）
    done
```

#### 6.3.2 フェーズ2: 分析

```
ユーザー: 「過去3ヶ月のIssue傾向をレポートして」

Claude Code:
1. 月別集計クエリ実行
2. 開発者別クエリ実行
3. ラベル別クエリ実行
4. 結果をMarkdownで整形
```

#### 6.3.3 フェーズ3: レポート作成

**自動生成されるレポート例**:

```markdown
# Issue分析レポート

**分析期間**: 2025-10-01 〜 2025-12-25
**総Issue数**: 387件（Open: 52件、Closed: 335件）

## 月別推移

| 月 | 新規作成 | クローズ | 純増 |
|----|---------|---------|------|
| 2025-12 | 28 | 31 | -3 |
| 2025-11 | 35 | 40 | -5 |
| 2025-10 | 42 | 38 | +4 |

## 開発者別統計

| 開発者 | 作成 | クローズ | クローズ率 |
|--------|------|---------|----------|
| Developer A | 120 | 115 | 95.8% |
| Developer B | 85 | 72 | 84.7% |

## 主要な Issue

1. #3977 - Performance improvement for main page (15 comments)
2. #3896 - Fix login bug (12 comments)
```

### 6.4 グラフ・可視化

#### 6.4.1 Mermaid図でのトレンド表示

```markdown
```mermaid
xychart-beta
    title "月別Issue推移"
    x-axis [Oct, Nov, Dec]
    y-axis "Issue数" 0 --> 50
    bar [42, 35, 28]
    line [38, 40, 31]
```
```

#### 6.4.2 外部ツールとの連携

```bash
# CSVエクスポートしてExcel/Google Sheetsで可視化
sqlite3 -csv issues.db "
SELECT strftime('%Y-%m', created_at) as month, COUNT(*) as total
FROM issues
GROUP BY month
ORDER BY month" > monthly_issues.csv
```

---

## 7. 応用例とユースケース

### 7.1 インフラ調査プロジェクトの進め方（4フェーズ）

#### Phase 1: リソース把握 → インフラ調査

**目的**: クラウドリソース（GCP/AWS等）の全体像を把握

**成果物**:
- インフラ構成ドキュメント（リソース一覧、ネットワーク図）
- セキュリティ評価レポート

**期間**: 1-2日

#### Phase 2: 構造理解 → アプリケーション調査

**目的**: アプリケーション構造・技術スタックの理解

**成果物**:
- アプリケーション構造ドキュメント
- ローカル開発環境セットアップガイド
- CI/CDパイプライン詳細ドキュメント

**期間**: 2-3日

#### Phase 3: プロセス分析 → 開発プロセス・組織調査

**目的**: 開発体制・責任分界点の明確化

**成果物**:
- 組織変革戦略ドキュメント
- チーム実態評価レポート
- RACIマトリクス（現状 vs 目標）

**期間**: 1-2日

#### Phase 4: データドリブン → データ分析・レポート作成

**目的**: 定量的なデータに基づく意思決定

**成果物**:
- 分析用データベース（SQLite3等）
- Issue/PR統計レポート
- 開発者アクティビティ分析レポート

**期間**: 1-2日

### 7.2 開発者稼働分析の実施方法

```bash
# git logによるアクティビティ抽出
git log --all --since="2025-01-01" --until="2025-12-31" \
  --author="Developer A\|Developer B" \
  --pretty=format:"%H|%an|%ae|%ad|%s" --date=short --numstat

# PR統計（gh-wrapper.sh活用）
~/.local/bin/gh-wrapper.sh pr list --state all --limit 200 --json number,title,author,createdAt,additions,deletions
```

### 7.3 組織変革戦略ドキュメントの作成

**手順**:

1. 現状課題の整理（ヒアリング・調査）
2. RACIマトリクス作成（現状 vs 目標）
3. Phase計画策定（優先順位・マイルストーン）
4. 進捗管理テーブル更新

**RACIマトリクステンプレート**:

```markdown
| フェーズ | Company A | Company B | Vendor |
|----------|-----------|-----------|--------|
| 企画 | R/A | C/I | I |
| 設計 | C | I | R/A |
| 開発 | I | I | R/A |
| デプロイ | R/A | C | C |
```

### 7.4 ドキュメントベースのナレッジマネジメント

**原則**:

- `.claude/rules/`に知識を外部化
- 定期的にドキュメントを更新（週次/月次）
- シナリオ別読み方ガイド（README.md）を維持

**更新フロー**:

```
1. 新しい知見を得る
   ↓
2. .claude/rules/{category}.md に追記
   ↓
3. CLAUDE.md に反映（必要に応じて）
   ↓
4. README.md のシナリオを更新
```

---

## 8. 参考資料

### 8.1 ファイル配置例

#### グローバル設定

```
~/.claude/
├── CLAUDE.md              # グローバル設定
└── commands/              # カスタムスラッシュコマンド
    ├── git:pr.md
    ├── git:rebase.md
    └── git-sync.md

~/.local/bin/
├── git-wrapper.sh         # Git操作ラッパー
└── gh-wrapper.sh          # GitHub CLI操作ラッパー
```

### 8.2 コマンドチートシート

#### ラッパースクリプト操作

```bash
# Git操作
~/.local/bin/git-wrapper.sh status
~/.local/bin/git-wrapper.sh diff HEAD
~/.local/bin/git-wrapper.sh add .
~/.local/bin/git-wrapper.sh commit -m "メッセージ"
~/.local/bin/git-wrapper.sh new-pr-branch

# GitHub操作
~/.local/bin/gh-wrapper.sh pr list --state all
~/.local/bin/gh-wrapper.sh pr view 12 --json number,title
~/.local/bin/gh-wrapper.sh issue list
~/.local/bin/gh-wrapper.sh repo view
```

#### SQLite3クエリ

```bash
# テーブル一覧
sqlite3 issues.db ".tables"

# スキーマ確認
sqlite3 issues.db ".schema issues"

# CSV出力
sqlite3 -csv issues.db "SELECT * FROM issues LIMIT 10" > output.csv

# ヘッダー付きCSV
sqlite3 -header -csv issues.db "SELECT * FROM issues" > output.csv
```

#### ドキュメント統計コマンド

```bash
# 総行数
find . -name "*.md" -exec wc -l {} + | sort -n

# ディレクトリ別行数
for dir in infrastructure apps strategy; do
  echo "$dir: $(find $dir -name "*.md" -exec wc -l {} + | tail -1)"
done
```

#### git log分析コマンド

```bash
# コミット数（開発者別）
git log --all --pretty=format:"%an" | sort | uniq -c | sort -rn

# LOC変更（開発者別）
git log --all --author="Developer A" --numstat --pretty="" | awk '{add+=$1; del+=$2} END {print "Add:", add, "Del:", del}'
```

### 8.3 プロジェクトテンプレート

**グローバルCLAUDE.mdのテンプレート**: 本ドキュメント§2.1を参照

**.gitignoreのテンプレート**:

```gitignore
# カスタムツール成果物
tools/*/bin/
tools/*/out/

# ローカル設定
.claude/*.local.json
.claude/CLAUDE.local.md

# 分析DB（チーム共有の場合はコメントアウト）
# issues.db

# macOS
.DS_Store
```

### 8.4 外部リソース

- **Claude Code公式ドキュメント**: https://docs.anthropic.com/en/docs/claude-code
- **GitHub CLI公式ドキュメント**: https://cli.github.com/manual/
- **SQLite3公式ドキュメント**: https://www.sqlite.org/docs.html
- **Mermaid公式ドキュメント**: https://mermaid.js.org/

---

## おわりに

このガイドで紹介した実践パターンを活用することで、Claude Codeを使った大規模プロジェクトの調査・ドキュメント作成・データ分析を効率的に進めることができます。

**重要なポイント**:

1. **ラッパースクリプト**: git-wrapper.sh、gh-wrapper.shによる標準化
2. **知識の外部化**: .claude/rules/による組織知識の蓄積
3. **データ分析基盤**: SQLite3によるローカルキャッシュとClaude Code統合
4. **ドキュメント駆動**: 複数レイヤーのドキュメント構成とシナリオ別ガイド

これらのパターンは、同様のインフラ調査・ドキュメント作成・データ分析プロジェクトに応用可能です。

**最終更新**: 2025-12-25
