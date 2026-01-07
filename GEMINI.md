# GEMINI.md

このファイルは、Gemini CLI がこのリポジトリで作業する際の共通ルールとプロジェクト概要を定義します。

---

## 共通ルール

### 1. 言語設定
- 常に**日本語**で会話する。
- コミットメッセージ、プルリクエストのタイトルおよび説明も**日本語**で記述する。

### 2. ツール使用の優先順位
**重要: 専用ラッパーツールを優先的に使用する**
直接 `git` や `gh` コマンドを叩くのではなく、以下のラッパースクリプトが利用可能な場合はそれを使用してください。

- **Git操作**: `~/.local/bin/git-wrapper.sh` を使用。
- **GitHub操作**: `~/.local/bin/gh-wrapper.sh` を使用。

#### 使用例:
```bash
~/.local/bin/git-wrapper.sh status
~/.local/bin/gh-wrapper.sh pr view 4079
```

---

## プロジェクト概要

ドラピタ（dorapita）は複数のCakePHPアプリケーションで構成された求人サービスシステム。

## 技術スタック

- **メインアプリ**: PHP 8.1 / CakePHP 4.x
- **レガシー**: PHP 5.6 / CakePHP 1.3.21
- **DB**: PostgreSQL 10 / MySQL 5.7
- **環境**: Docker (DevContainer対応)

## アプリケーション構成

| アプリ | ポート | 用途 |
|--------|--------|------|
| dorapita.com | 8080 | メインサイト |
| cadm.dorapita.com | 8888 | 企業管理画面 |
| kanri.dorapita.com | 8088 | 運営管理画面 |
| dora-pt.jp | 9999 | PT連携サイト |
| legacy.dorapita.com | 9090 | 旧システム |

## 主要コマンド

### CakePHP 4.x (dorapita.com, cadm, kanri, dora-pt.jp)
```bash
cd {アプリディレクトリ}
vendor/bin/phpunit                    # テスト実行
vendor/bin/phpcs --colors -p src/ tests/    # CSチェック
vendor/bin/phpcbf --colors -p src/ tests/   # CS自動修正
bin/cake migrations migrate            # マイグレーション
```

### CakePHP 1.3 (legacy.dorapita.com)
```bash
cd legacy.dorapita.com
cake/console/cake -app app {command}
```

---

## 開発ガイドライン

- **ブランチ命名**: `feature/内容` または `fix/内容`。
- **環境変数**: ローカル開発では `.env.local` を使用する（PR #4079 で導入）。
- **テスト**: 新機能追加やバグ修正の際は、可能な限り PHPUnit テストを追加または更新する。

---

## 推奨ワークフロー (Custom Commands相当)

### 1. リモートと同期する (Sync)
作業開始前やPR作成前に実行してください。
1. `~/.local/bin/git-wrapper.sh status --short` で現状確認。
2. 変更がある場合は `stash` に退避。
3. `git fetch origin` を実行。
4. `main` ブランチなら `pull`、それ以外なら `rebase origin/main`。
5. 退避した変更があれば `stash pop`。

### 2. PR作成フロー (Auto PR)
機能実装完了後のPR作成手順です。
1. `~/.local/bin/git-wrapper.sh status` で確認。
2. `main` ブランチにいる場合は `~/.local/bin/git-wrapper.sh new-pr-branch` でブランチ作成。
3. `~/.local/bin/git-wrapper.sh add .` ですべてステージング。
4. `~/.local/bin/git-wrapper.sh commit -m "日本語のコミットメッセージ"`。
5. `git push -u origin $(git branch --show-current)`。
6. `~/.local/bin/gh-wrapper.sh pr create --title "タイトル" --body "説明" --base main`。
