package github

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Client GitHub API クライアント（gh-wrapper.sh経由）
type Client struct {
	ghWrapperPath string
	repo          string
	workDir       string // gh-wrapper.sh実行時の作業ディレクトリ
}

// NewClient クライアントを作成
func NewClient(repo string) (*Client, error) {
	// gh-wrapper.shのパスを解決
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("ホームディレクトリ取得エラー: %w", err)
	}
	ghWrapperPath := filepath.Join(homeDir, ".local", "bin", "gh-wrapper.sh")

	// 存在確認
	if _, err := os.Stat(ghWrapperPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("gh-wrapper.shが見つかりません: %s", ghWrapperPath)
	}

	return &Client{
		ghWrapperPath: ghWrapperPath,
		repo:          repo,
	}, nil
}

// SetWorkDir 作業ディレクトリを設定（PAT抽出用）
func (c *Client) SetWorkDir(dir string) {
	c.workDir = dir
}

// rawIssue GitHub APIからの生のissueデータ（stateが大文字）
type rawIssue struct {
	Number    int        `json:"number"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	State     string     `json:"state"` // OPEN, CLOSED (大文字)
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ClosedAt  *time.Time `json:"closedAt"`
	Author    Author     `json:"author"`
	Assignees []Assignee `json:"assignees"`
	Labels    []Label    `json:"labels"`
	URL       string     `json:"url"`
}

// ListIssues issue一覧を取得（直近N日間）
func (c *Client) ListIssues(days int) ([]Issue, error) {
	// gh issue list --state all --limit 1000 -R owner/repo --json ...
	args := []string{
		"--raw", // JSON変換をスキップしてそのまま出力
		"issue", "list",
		"--state", "all",
		"--limit", "1000",
		"-R", c.repo,
		"--json", "number,title,body,state,createdAt,updatedAt,closedAt,author,assignees,labels,url",
	}

	output, err := c.runGHCommand(args...)
	if err != nil {
		return nil, fmt.Errorf("issue一覧取得エラー: %w", err)
	}

	var rawIssues []rawIssue
	if err := json.Unmarshal(output, &rawIssues); err != nil {
		return nil, fmt.Errorf("JSONパースエラー: %w", err)
	}

	// 直近N日間でフィルタ + stateを小文字に変換
	cutoff := time.Now().AddDate(0, 0, -days)
	var filtered []Issue
	for _, ri := range rawIssues {
		if ri.CreatedAt.After(cutoff) {
			issue := Issue{
				Number:    ri.Number,
				Title:     ri.Title,
				Body:      ri.Body,
				State:     strings.ToLower(ri.State), // OPEN -> open, CLOSED -> closed
				CreatedAt: ri.CreatedAt,
				UpdatedAt: ri.UpdatedAt,
				ClosedAt:  ri.ClosedAt,
				Author:    ri.Author,
				Assignees: ri.Assignees,
				Labels:    ri.Labels,
				URL:       ri.URL,
			}
			filtered = append(filtered, issue)
		}
	}

	return filtered, nil
}

// GetIssue issue詳細を取得（コメント含む）
func (c *Client) GetIssue(number int) (*Issue, error) {
	// gh issue view <number> -R owner/repo
	args := []string{
		"issue", "view", strconv.Itoa(number),
		"-R", c.repo,
	}

	output, err := c.runGHCommand(args...)
	if err != nil {
		return nil, fmt.Errorf("issue詳細取得エラー: %w", err)
	}

	var issue Issue
	if err := json.Unmarshal(output, &issue); err != nil {
		return nil, fmt.Errorf("JSONパースエラー: %w", err)
	}

	// コメント取得
	comments, err := c.GetComments(number)
	if err != nil {
		// コメント取得失敗は警告のみ
		fmt.Fprintf(os.Stderr, "警告: issue #%d のコメント取得に失敗: %v\n", number, err)
	} else {
		issue.Comments = comments
	}

	return &issue, nil
}

// GetComments issueのコメント一覧を取得
func (c *Client) GetComments(issueNumber int) ([]Comment, error) {
	// gh api repos/{owner}/{repo}/issues/{number}/comments
	args := []string{
		"api",
		fmt.Sprintf("repos/%s/issues/%d/comments", c.repo, issueNumber),
	}

	output, err := c.runGHCommand(args...)
	if err != nil {
		return nil, fmt.Errorf("コメント取得エラー: %w", err)
	}

	// GitHub API形式のコメントをパース
	var rawComments []struct {
		ID        int    `json:"id"`
		Body      string `json:"body"`
		CreatedAt string `json:"created_at"`
		User      struct {
			Login string `json:"login"`
			Type  string `json:"type"`
		} `json:"user"`
	}

	if err := json.Unmarshal(output, &rawComments); err != nil {
		return nil, fmt.Errorf("コメントJSONパースエラー: %w", err)
	}

	var comments []Comment
	for _, rc := range rawComments {
		createdAt, _ := time.Parse(time.RFC3339, rc.CreatedAt)
		comments = append(comments, Comment{
			ID: rc.ID,
			Author: Author{
				Login: rc.User.Login,
				Type:  rc.User.Type,
			},
			Body:      rc.Body,
			CreatedAt: createdAt,
		})
	}

	return comments, nil
}

// runGHCommand gh-wrapper.shを実行
func (c *Client) runGHCommand(args ...string) ([]byte, error) {
	cmd := exec.Command(c.ghWrapperPath, args...)

	// 作業ディレクトリ設定（PAT抽出用）
	if c.workDir != "" {
		cmd.Dir = c.workDir
	}

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("gh-wrapper.sh実行エラー: %s", string(exitErr.Stderr))
		}
		return nil, err
	}

	return output, nil
}
