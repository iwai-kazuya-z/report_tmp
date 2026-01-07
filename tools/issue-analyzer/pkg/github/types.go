package github

import "time"

// Issue GitHub issueの基本情報
type Issue struct {
	Number      int       `json:"number"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	State       string    `json:"state"` // open, closed
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ClosedAt    *time.Time `json:"closedAt"` // nilの場合あり
	Author      Author    `json:"author"`
	Assignees   []Assignee `json:"assignees"`
	Labels      []Label   `json:"labels"`
	URL         string    `json:"url"`
	Comments    []Comment `json:"comments,omitempty"` // issue詳細取得時のみ
}

// Author issue作成者
type Author struct {
	Login string `json:"login"`
	Type  string `json:"type"` // User, Bot
}

// Assignee issue担当者
type Assignee struct {
	Login string `json:"login"`
}

// Label issueラベル
type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// Comment issueコメント
type Comment struct {
	ID        int       `json:"id"`
	Author    Author    `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

// Link issue間の参照リンク
type Link struct {
	SourceIssue int
	TargetIssue int
	LinkType    string // mentions, closes, fixes, relates
}
