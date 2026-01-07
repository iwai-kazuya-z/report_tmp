package github

import (
	"regexp"
	"strconv"
)

// リンクタイプの定義
const (
	LinkTypeCloses   = "closes"
	LinkTypeFixes    = "fixes"
	LinkTypeResolves = "resolves"
	LinkTypeMentions = "mentions"
)

// 正規表現パターン
var (
	// closes #123, close #123, closed #123 などにマッチ
	closesPattern = regexp.MustCompile(`(?i)\b(close[sd]?)\s+#(\d+)`)
	// fixes #123, fix #123, fixed #123 などにマッチ
	fixesPattern = regexp.MustCompile(`(?i)\b(fix(?:e[sd])?)\s+#(\d+)`)
	// resolves #123, resolve #123, resolved #123 などにマッチ
	resolvesPattern = regexp.MustCompile(`(?i)\b(resolve[sd]?)\s+#(\d+)`)
	// 単純な #123 参照にマッチ（closes/fixes/resolvesを含まない）
	mentionsPattern = regexp.MustCompile(`(?:^|[^\w])#(\d+)`)
)

// ParseLinks issue本文から参照リンクを抽出
func ParseLinks(sourceIssue int, body string) []Link {
	if body == "" {
		return nil
	}

	var links []Link
	seen := make(map[string]bool) // 重複排除用

	// closes パターン
	for _, match := range closesPattern.FindAllStringSubmatch(body, -1) {
		if len(match) >= 3 {
			targetIssue, _ := strconv.Atoi(match[2])
			key := linkKey(sourceIssue, targetIssue, LinkTypeCloses)
			if !seen[key] && targetIssue != sourceIssue {
				seen[key] = true
				links = append(links, Link{
					SourceIssue: sourceIssue,
					TargetIssue: targetIssue,
					LinkType:    LinkTypeCloses,
				})
			}
		}
	}

	// fixes パターン
	for _, match := range fixesPattern.FindAllStringSubmatch(body, -1) {
		if len(match) >= 3 {
			targetIssue, _ := strconv.Atoi(match[2])
			key := linkKey(sourceIssue, targetIssue, LinkTypeFixes)
			if !seen[key] && targetIssue != sourceIssue {
				seen[key] = true
				links = append(links, Link{
					SourceIssue: sourceIssue,
					TargetIssue: targetIssue,
					LinkType:    LinkTypeFixes,
				})
			}
		}
	}

	// resolves パターン
	for _, match := range resolvesPattern.FindAllStringSubmatch(body, -1) {
		if len(match) >= 3 {
			targetIssue, _ := strconv.Atoi(match[2])
			key := linkKey(sourceIssue, targetIssue, LinkTypeResolves)
			if !seen[key] && targetIssue != sourceIssue {
				seen[key] = true
				links = append(links, Link{
					SourceIssue: sourceIssue,
					TargetIssue: targetIssue,
					LinkType:    LinkTypeResolves,
				})
			}
		}
	}

	// mentions パターン（closes/fixes/resolvesを除外）
	for _, match := range mentionsPattern.FindAllStringSubmatch(body, -1) {
		if len(match) >= 2 {
			targetIssue, _ := strconv.Atoi(match[1])
			// 既に他のタイプで登録済みの場合はスキップ
			closesKey := linkKey(sourceIssue, targetIssue, LinkTypeCloses)
			fixesKey := linkKey(sourceIssue, targetIssue, LinkTypeFixes)
			resolvesKey := linkKey(sourceIssue, targetIssue, LinkTypeResolves)
			mentionsKey := linkKey(sourceIssue, targetIssue, LinkTypeMentions)

			if !seen[closesKey] && !seen[fixesKey] && !seen[resolvesKey] && !seen[mentionsKey] && targetIssue != sourceIssue {
				seen[mentionsKey] = true
				links = append(links, Link{
					SourceIssue: sourceIssue,
					TargetIssue: targetIssue,
					LinkType:    LinkTypeMentions,
				})
			}
		}
	}

	return links
}

// ParseAllLinks issueとそのコメントから全ての参照リンクを抽出
func ParseAllLinks(issue *Issue) []Link {
	var allLinks []Link

	// issue本文から抽出
	bodyLinks := ParseLinks(issue.Number, issue.Body)
	allLinks = append(allLinks, bodyLinks...)

	// コメントから抽出
	for _, comment := range issue.Comments {
		commentLinks := ParseLinks(issue.Number, comment.Body)
		allLinks = append(allLinks, commentLinks...)
	}

	return allLinks
}

// linkKey 重複チェック用のキー生成
func linkKey(source, target int, linkType string) string {
	return strconv.Itoa(source) + "-" + strconv.Itoa(target) + "-" + linkType
}
