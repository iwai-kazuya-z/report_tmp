package cache

// SQLite3スキーマ定義
const schema = `
-- issues テーブル
CREATE TABLE IF NOT EXISTS issues (
    number INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT,
    state TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT,
    closed_at TEXT,
    author_login TEXT,
    author_type TEXT,
    assignee_login TEXT,
    url TEXT NOT NULL
);

-- labels テーブル
CREATE TABLE IF NOT EXISTS labels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    issue_number INTEGER NOT NULL,
    name TEXT NOT NULL,
    color TEXT,
    FOREIGN KEY (issue_number) REFERENCES issues(number) ON DELETE CASCADE
);

-- comments テーブル
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY,
    issue_number INTEGER NOT NULL,
    author_login TEXT,
    body TEXT,
    created_at TEXT NOT NULL,
    FOREIGN KEY (issue_number) REFERENCES issues(number) ON DELETE CASCADE
);

-- issue_links テーブル（参照関係）
CREATE TABLE IF NOT EXISTS issue_links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_issue INTEGER NOT NULL,
    target_issue INTEGER NOT NULL,
    link_type TEXT NOT NULL,
    FOREIGN KEY (source_issue) REFERENCES issues(number) ON DELETE CASCADE,
    UNIQUE(source_issue, target_issue, link_type)
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_issues_state ON issues(state);
CREATE INDEX IF NOT EXISTS idx_issues_created_at ON issues(created_at);
CREATE INDEX IF NOT EXISTS idx_issues_author ON issues(author_login);
CREATE INDEX IF NOT EXISTS idx_issues_assignee ON issues(assignee_login);
CREATE INDEX IF NOT EXISTS idx_labels_issue ON labels(issue_number);
CREATE INDEX IF NOT EXISTS idx_labels_name ON labels(name);
CREATE INDEX IF NOT EXISTS idx_comments_issue ON comments(issue_number);
CREATE INDEX IF NOT EXISTS idx_comments_author ON comments(author_login);
CREATE INDEX IF NOT EXISTS idx_links_source ON issue_links(source_issue);
CREATE INDEX IF NOT EXISTS idx_links_target ON issue_links(target_issue);
CREATE INDEX IF NOT EXISTS idx_links_type ON issue_links(link_type);

-- メタデータ テーブル（取得日時などの情報）
CREATE TABLE IF NOT EXISTS metadata (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
`
