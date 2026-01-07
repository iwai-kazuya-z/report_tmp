package cache

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hashimoto-kazuhiro-aa/dorapita-issue-analyzer/pkg/github"
	_ "github.com/mattn/go-sqlite3"
)

// Cache SQLite3キャッシュ
type Cache struct {
	db *sql.DB
}

// New キャッシュを作成
func New(dbPath string) (*Cache, error) {
	// ディレクトリ作成
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("ディレクトリ作成エラー: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("データベースオープンエラー: %w", err)
	}

	// スキーマ作成
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("スキーマ作成エラー: %w", err)
	}

	return &Cache{db: db}, nil
}

// Close データベースを閉じる
func (c *Cache) Close() error {
	return c.db.Close()
}

// SaveIssues issueを一括保存
func (c *Cache) SaveIssues(issues []github.Issue) error {
	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("トランザクション開始エラー: %w", err)
	}
	defer tx.Rollback()

	// 既存データをクリア
	tables := []string{"issue_links", "comments", "labels", "issues"}
	for _, table := range tables {
		if _, err := tx.Exec("DELETE FROM " + table); err != nil {
			return fmt.Errorf("%sクリアエラー: %w", table, err)
		}
	}

	// issue保存
	issueStmt, err := tx.Prepare(`
		INSERT INTO issues (number, title, body, state, created_at, updated_at, closed_at, author_login, author_type, assignee_login, url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("issueステートメント準備エラー: %w", err)
	}
	defer issueStmt.Close()

	// label保存
	labelStmt, err := tx.Prepare(`
		INSERT INTO labels (issue_number, name, color)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("labelステートメント準備エラー: %w", err)
	}
	defer labelStmt.Close()

	for _, issue := range issues {
		// closed_at の処理
		var closedAt *string
		if issue.ClosedAt != nil {
			t := issue.ClosedAt.Format(time.RFC3339)
			closedAt = &t
		}

		// assignee の処理（最初の担当者を使用）
		var assigneeLogin *string
		if len(issue.Assignees) > 0 {
			assigneeLogin = &issue.Assignees[0].Login
		}

		_, err := issueStmt.Exec(
			issue.Number,
			issue.Title,
			issue.Body,
			issue.State,
			issue.CreatedAt.Format(time.RFC3339),
			issue.UpdatedAt.Format(time.RFC3339),
			closedAt,
			issue.Author.Login,
			issue.Author.Type,
			assigneeLogin,
			issue.URL,
		)
		if err != nil {
			return fmt.Errorf("issue #%d 保存エラー: %w", issue.Number, err)
		}

		// labels保存
		for _, label := range issue.Labels {
			_, err := labelStmt.Exec(issue.Number, label.Name, label.Color)
			if err != nil {
				return fmt.Errorf("issue #%d label保存エラー: %w", issue.Number, err)
			}
		}
	}

	// メタデータ更新
	_, err = tx.Exec(`
		INSERT OR REPLACE INTO metadata (key, value)
		VALUES ('last_updated', ?)
	`, time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("メタデータ更新エラー: %w", err)
	}

	return tx.Commit()
}

// SaveComments コメントを保存
func (c *Cache) SaveComments(issueNumber int, comments []github.Comment) error {
	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("トランザクション開始エラー: %w", err)
	}
	defer tx.Rollback()

	// 既存コメントを削除
	_, err = tx.Exec("DELETE FROM comments WHERE issue_number = ?", issueNumber)
	if err != nil {
		return fmt.Errorf("既存コメント削除エラー: %w", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO comments (id, issue_number, author_login, body, created_at)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("コメントステートメント準備エラー: %w", err)
	}
	defer stmt.Close()

	for _, comment := range comments {
		_, err := stmt.Exec(
			comment.ID,
			issueNumber,
			comment.Author.Login,
			comment.Body,
			comment.CreatedAt.Format(time.RFC3339),
		)
		if err != nil {
			return fmt.Errorf("コメント保存エラー: %w", err)
		}
	}

	return tx.Commit()
}

// SaveLinks 参照リンクを保存
func (c *Cache) SaveLinks(links []github.Link) error {
	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("トランザクション開始エラー: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO issue_links (source_issue, target_issue, link_type)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("リンクステートメント準備エラー: %w", err)
	}
	defer stmt.Close()

	for _, link := range links {
		_, err := stmt.Exec(link.SourceIssue, link.TargetIssue, link.LinkType)
		if err != nil {
			return fmt.Errorf("リンク保存エラー: %w", err)
		}
	}

	return tx.Commit()
}

// LoadIssues 全issueを取得
func (c *Cache) LoadIssues() ([]github.Issue, error) {
	rows, err := c.db.Query(`
		SELECT number, title, body, state, created_at, updated_at, closed_at, author_login, author_type, assignee_login, url
		FROM issues
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("issue取得エラー: %w", err)
	}
	defer rows.Close()

	var issues []github.Issue
	for rows.Next() {
		var issue github.Issue
		var createdAt, updatedAt string
		var closedAt, authorType, assigneeLogin sql.NullString

		err := rows.Scan(
			&issue.Number,
			&issue.Title,
			&issue.Body,
			&issue.State,
			&createdAt,
			&updatedAt,
			&closedAt,
			&issue.Author.Login,
			&authorType,
			&assigneeLogin,
			&issue.URL,
		)
		if err != nil {
			return nil, fmt.Errorf("issue読み込みエラー: %w", err)
		}

		issue.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		issue.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		if closedAt.Valid {
			t, _ := time.Parse(time.RFC3339, closedAt.String)
			issue.ClosedAt = &t
		}
		if authorType.Valid {
			issue.Author.Type = authorType.String
		}
		if assigneeLogin.Valid {
			issue.Assignees = []github.Assignee{{Login: assigneeLogin.String}}
		}

		// labels取得
		labels, err := c.loadLabelsForIssue(issue.Number)
		if err != nil {
			return nil, err
		}
		issue.Labels = labels

		issues = append(issues, issue)
	}

	return issues, nil
}

// loadLabelsForIssue 指定issueのラベルを取得
func (c *Cache) loadLabelsForIssue(issueNumber int) ([]github.Label, error) {
	rows, err := c.db.Query(`
		SELECT name, color FROM labels WHERE issue_number = ?
	`, issueNumber)
	if err != nil {
		return nil, fmt.Errorf("label取得エラー: %w", err)
	}
	defer rows.Close()

	var labels []github.Label
	for rows.Next() {
		var label github.Label
		var color sql.NullString
		if err := rows.Scan(&label.Name, &color); err != nil {
			return nil, fmt.Errorf("label読み込みエラー: %w", err)
		}
		if color.Valid {
			label.Color = color.String
		}
		labels = append(labels, label)
	}

	return labels, nil
}

// LoadComments 指定issueのコメントを取得
func (c *Cache) LoadComments(issueNumber int) ([]github.Comment, error) {
	rows, err := c.db.Query(`
		SELECT id, author_login, body, created_at
		FROM comments
		WHERE issue_number = ?
		ORDER BY created_at ASC
	`, issueNumber)
	if err != nil {
		return nil, fmt.Errorf("コメント取得エラー: %w", err)
	}
	defer rows.Close()

	var comments []github.Comment
	for rows.Next() {
		var comment github.Comment
		var createdAt string
		var authorLogin sql.NullString

		err := rows.Scan(&comment.ID, &authorLogin, &comment.Body, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("コメント読み込みエラー: %w", err)
		}

		comment.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		if authorLogin.Valid {
			comment.Author.Login = authorLogin.String
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// LoadLinks 全参照リンクを取得
func (c *Cache) LoadLinks() ([]github.Link, error) {
	rows, err := c.db.Query(`
		SELECT source_issue, target_issue, link_type
		FROM issue_links
	`)
	if err != nil {
		return nil, fmt.Errorf("リンク取得エラー: %w", err)
	}
	defer rows.Close()

	var links []github.Link
	for rows.Next() {
		var link github.Link
		err := rows.Scan(&link.SourceIssue, &link.TargetIssue, &link.LinkType)
		if err != nil {
			return nil, fmt.Errorf("リンク読み込みエラー: %w", err)
		}
		links = append(links, link)
	}

	return links, nil
}

// GetMetadata メタデータを取得
func (c *Cache) GetMetadata(key string) (string, error) {
	var value string
	err := c.db.QueryRow("SELECT value FROM metadata WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("メタデータ取得エラー: %w", err)
	}
	return value, nil
}

// GetDB データベース接続を取得（分析用クエリ実行用）
func (c *Cache) GetDB() *sql.DB {
	return c.db
}
