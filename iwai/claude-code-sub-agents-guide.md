# Claude Code Sub Agents 活用ガイド

作成日: 2025-01-07

---

## 1. Sub Agentsとは

**定義**: 特定のタスクに特化した専門的なAIアシスタント。メイン会話から独立したコンテキストウィンドウで動作し、カスタムシステムプロンプトとツール権限を持つ。

### 主なメリット

| メリット | 説明 |
|---------|------|
| **コンテキスト保護** | エージェント固有の操作がメイン会話を汚さない |
| **専門性向上** | ドメイン特定の詳細指示により成功率向上 |
| **再利用性** | 一度作成したエージェントを異なるプロジェクトで再利用可能 |
| **権限の柔軟性** | 各エージェントに異なるツール権限を設定可能 |

### auto-compact問題の解決

Sub Agentsを使うことで、メイン会話のコンテキストを消費せずに専門的なタスクを実行できる。これにより**コンテキスト枯渇**（auto-compact問題）を回避。

---

## 2. 設定方法

### ファイル配置場所

| タイプ | 配置場所 | スコープ | 優先度 |
|--------|----------|---------|--------|
| プロジェクト | `.claude/agents/` | 現在のプロジェクトのみ | 最高 |
| ユーザー | `~/.claude/agents/` | すべてのプロジェクト | 低 |

**注意**: 同名エージェントが複数存在する場合、プロジェクトレベルが優先される。

### ファイル形式

```markdown
---
name: agent-name
description: エージェントの説明（自動委譲トリガーに重要）
tools: Read, Edit, Bash, Grep, Glob
model: sonnet
---

ここにシステムプロンプトを記述
複数段落で詳細な指示を書く
```

### フロントマター完全リファレンス

| フィールド | 必須 | 型 | 説明 | 例 |
|-----------|------|-----|------|-----|
| `name` | Yes | string | 一意の識別子（小文字+ハイフン） | `code-reviewer` |
| `description` | Yes | string | 自然言語での目的説明（自動委譲用） | `"Expert code reviewer..."` |
| `tools` | No | string | 利用可能なツール（カンマ区切り） | `Read, Grep, Bash` |
| `model` | No | string | モデル指定 | `sonnet`, `opus`, `haiku`, `inherit` |
| `skills` | No | string | 起動時に自動ロードするスキル | `test-skill` |

---

## 3. 利用可能なツール一覧

```
Read           - ファイル読み取り
Edit           - ファイル編集（対象型）
Write          - ファイル作成/上書き
Bash           - シェルコマンド実行
Glob           - ファイルパターンマッチ
Grep           - ファイル内容検索（正規表現対応）
WebFetch       - URL内容取得
WebSearch      - Web検索
Task           - 別のサブエージェント呼び出し
TodoWrite      - タスクリスト管理
NotebookEdit   - Jupyterノートブック編集
AskUserQuestion- ユーザーへの質問
```

### ツール設定のポイント

```markdown
# 方法1: 全ツール継承（toolsフィールド省略）
---
name: general-agent
---

# 方法2: 必要なツールのみ指定（推奨）
---
name: analyzer
tools: Read, Grep, Glob
---
```

**セキュリティ推奨**: 必要最小限のツールを指定する。

---

## 4. モデル指定オプション

| エイリアス | 動作 | 用途 |
|-----------|------|------|
| `sonnet` | 最新Sonnet（現: 4.5） | 日常的なコーディングタスク |
| `opus` | Opusモデル（現: 4.5） | 複雑な推論と専門的タスク |
| `haiku` | 高速Haikuモデル | シンプルなタスク |
| `inherit` | メイン会話と同じモデル | 一貫性が必要な場合 |

---

## 5. 呼び出し方法

### 5.1 自動委譲（Automatic Delegation）

Claudeが`description`に基づいて自動的にエージェントを選択。

```markdown
# 自動委譲されやすいdescription
description: "Expert code reviewer. Use proactively after code changes."

# 自動委譲されにくいdescription
description: "Reviews code"
```

### 5.2 明示的な呼び出し

```
> Use the code-reviewer subagent to check my recent changes
> Have the debugger subagent investigate this error
```

### 5.3 /agentsコマンド（推奨）

```
/agents
```

このコマンドで:
- 利用可能なエージェント一覧表示
- 新規エージェント作成（Claudeがガイド）
- 既存エージェント編集・削除
- ツール権限管理

### 5.4 CLI `--agents`フラグ（動的定義）

```bash
claude --agents '{
  "code-reviewer": {
    "description": "Expert code reviewer...",
    "prompt": "You are a senior code reviewer...",
    "tools": ["Read", "Grep", "Glob", "Bash"],
    "model": "sonnet"
  }
}'
```

---

## 6. 組み込みサブエージェント

Claude Codeには以下の組み込みエージェントがある:

| エージェント | モデル | 用途 |
|-------------|--------|------|
| **Explore** | Haiku（高速） | コードベースの高速検索・分析（読み取り専用） |
| **Plan** | Sonnet | 計画作成、情報収集 |
| **General-purpose** | Sonnet | 複雑なマルチステップタスク |

---

## 7. 実践的な実装例

### 7.1 code-reviewer（コードレビュー）

```markdown
---
name: code-reviewer
description: Expert code reviewer. Use proactively after code changes to check quality, security, and best practices.
tools: Read, Grep, Glob, Bash
model: inherit
---

You are a senior code reviewer ensuring high standards.

When invoked:
1. Run git diff to see recent changes
2. Focus on modified files
3. Begin review immediately

Review checklist:
- Code is clear and readable
- Functions and variables are well-named
- No duplicated code
- Proper error handling
- No exposed secrets or API keys
- Input validation implemented
- Good test coverage

Provide feedback organized by priority:
- Critical issues (must fix)
- Warnings (should fix)
- Suggestions (consider improving)

Include specific examples of how to fix issues.
```

### 7.2 debugger（デバッグ専門家）

```markdown
---
name: debugger
description: Debugging specialist for errors, test failures, and unexpected behavior. Use proactively when encountering any issues.
tools: Read, Edit, Bash, Grep, Glob
model: sonnet
---

You are an expert debugger specializing in root cause analysis.

When invoked:
1. Capture error message and stack trace
2. Identify reproduction steps
3. Isolate the failure location
4. Implement minimal fix
5. Verify solution works

Debugging process:
- Analyze error messages and logs
- Check recent code changes
- Form and test hypotheses
- Add strategic debug logging
- Inspect variable states

For each issue, provide:
- Root cause explanation
- Evidence supporting the diagnosis
- Specific code fix
- Testing approach
- Prevention recommendations
```

### 7.3 task-decomposer（タスク分解）

```markdown
---
name: task-decomposer
description: Break down complex tasks into manageable subtasks with dependencies mapped.
tools: Read, Grep, Glob
model: sonnet
---

You are a task decomposition specialist.

When invoked:
1. Analyze the given task
2. Identify dependencies and shared components
3. Break down into manageable subtasks
4. Create execution order based on dependencies

Output format:
- Task overview
- Subtasks list with priorities
- Dependency graph
- Risk areas and mitigation strategies
- Estimated complexity per subtask
```

### 7.4 task-executor（タスク実行）

```markdown
---
name: task-executor
description: Execute decomposed subtasks sequentially. Resolve all errors before proceeding to next phase.
tools: Read, Edit, Write, Bash, Grep, Glob
model: sonnet
---

You are a task execution specialist.

Workflow:
1. Read the task from the progress file
2. Execute the task
3. Verify completion (run tests, type checks)
4. Update progress file immediately
5. Move to next task only when current is 100% complete

Critical rules:
- Never proceed with errors unresolved
- Update progress file after each action
- Document any blockers encountered
- Run verification after each change
```

### 7.5 issue-analyzer（GitHub Issue分析）

```markdown
---
name: issue-analyzer
description: Analyze GitHub issues from SQLite database. Extract patterns, trends, and insights.
tools: Bash, Read, Grep
model: sonnet
---

You are a GitHub issue analysis specialist.

Data source: issues.db (SQLite3)

Workflow:
1. Query issues.db using sqlite3 command
2. Analyze patterns and trends
3. Generate markdown reports

Analysis areas:
- Keyword frequency
- Monthly trends (open/closed)
- Team response times
- Issue clusters and dependencies
- Recurring problems

SQL query examples:
- SELECT number, title, state FROM issues WHERE title LIKE '%keyword%'
- SELECT strftime('%Y-%m', created_at) as month, COUNT(*) FROM issues GROUP BY month
```

### 7.6 gcp-auditor（GCPインフラ監査）

```markdown
---
name: gcp-auditor
description: Audit GCP infrastructure for security issues and best practices.
tools: Bash, Read, Grep, Glob
model: sonnet
---

You are a GCP infrastructure auditor.

Focus areas:
- Cloud SQL security (SSL enforcement, DB versions)
- Cloud Run ingress settings (should not be 'all')
- IAM permissions (principle of least privilege)
- Network configuration (VPC, firewall rules)
- Load balancer SSL policies

Commands:
- gcloud compute instances list
- gcloud run services list --region=asia-northeast1
- gcloud sql instances list

Output:
- Security findings by severity (Critical/High/Medium/Low)
- Remediation recommendations
- Compliance status
```

---

## 8. ベストプラクティス

### 8.1 descriptionは具体的に

```markdown
# 良い例（自動委譲がトリガーされやすい）
description: "Expert code reviewer. Use proactively after code changes."

# 悪い例
description: "Reviews code"
```

### 8.2 責任の分離（Single Responsibility）

```markdown
# 良い例：1つの責任
name: code-reviewer
description: Review code for quality and security

# 悪い例：複数の責任（機能が複雑すぎる）
name: super-agent
description: Review, test, deploy, monitor
```

### 8.3 段階的実行

各フェーズでエラーを完全解消してから次へ進む:
- 問題の早期発見
- デバッグの容易化
- 安定した進捗

### 8.4 進捗管理（ファイルベース）

```markdown
# progress.md
- [x] Setup environment
- [x] Create database schema
- [ ] Implement API endpoints  ← 現在実行中
- [ ] Add tests
```

### 8.5 ツールは最小限に

```markdown
# 推奨：必要なツールのみ
tools: Read, Grep, Glob

# 避ける：すべて継承（セキュリティリスク）
tools: # 省略
```

### 8.6 プロジェクトエージェントはGit管理

```bash
# .claude/agents/ をGitで追跡
git add .claude/agents/
git commit -m "Add code-reviewer subagent"
```

チームメンバーが自動的にアクセス可能になる。

---

## 9. 導入手順

```bash
# 1. ディレクトリ作成
mkdir -p .claude/agents

# 2. エージェントファイル作成
cat > .claude/agents/code-reviewer.md << 'EOF'
---
name: code-reviewer
description: Expert code reviewer. Use proactively after code changes.
tools: Read, Grep, Glob, Bash
---

You are a senior code reviewer...
EOF

# 3. 確認（Claude Code内で実行）
/agents
```

### /agentsコマンドでの作成（推奨）

```
/agents
→ "Create New Agent" を選択
→ "Generate with Claude first" を選択
→ エージェントの説明を入力
→ Claudeが初期案を生成
→ 必要に応じてカスタマイズ
```

---

## 10. トラブルシューティング

### エージェントが自動委譲されない

**原因**: `description`が曖昧

**解決**:
```markdown
# 修正前
description: "Reviewer"

# 修正後
description: "Expert code reviewer. Use proactively after code changes."
```

### ツールが利用できない

**解決**: `/agents`コマンドでエージェント編集時にツールを確認・追加

### エージェントが起動しない

**確認項目**:
1. ファイル名が正しいMarkdown（`.md`）
2. YAMLフロントマターが正しく閉じられている（`---`）
3. `name`フィールドに小文字とハイフンのみ使用
4. ファイルが正しいディレクトリにある

---

## 11. 高度な機能

### サブエージェントのチェーン

```
> First use the code-analyzer subagent to find performance issues,
> then use the optimizer subagent to fix them
```

### 再開可能なサブエージェント

```
# 初期呼び出し
> Use the code-analyzer agent to start reviewing the authentication module
[Agent returns: agentId: "abc123"]

# 再開（コンテキスト保持）
> Resume agent abc123 and now analyze the authorization logic as well
```

---

## 12. 参考資料

- **Zenn記事**: https://zenn.dev/tacoms/articles/552140c84aaefa
- **Claude Code公式ドキュメント**: https://code.claude.com/docs/en/sub-agents.md
- **CLI Reference**: https://code.claude.com/docs/en/cli-reference.md

---

**作成者**: Claude Code
**最終更新**: 2025-01-07
