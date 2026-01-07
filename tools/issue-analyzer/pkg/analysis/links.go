package analysis

import (
	"fmt"
	"sort"

	"github.com/hashimoto-kazuhiro-aa/dorapita-issue-analyzer/pkg/github"
)

// LinkAnalysis 参照リンク分析結果
type LinkAnalysis struct {
	TotalLinks      int
	LinkTypeCounts  map[string]int
	TopReferenced   []ReferencedIssue
	TopReferencing  []ReferencingIssue
	Clusters        []IssueCluster
	MermaidGraph    string
}

// ReferencedIssue 被参照issueの統計
type ReferencedIssue struct {
	Number         int
	Title          string
	ReferenceCount int
}

// ReferencingIssue 参照issueの統計
type ReferencingIssue struct {
	Number    int
	Title     string
	LinkCount int
}

// IssueCluster 関連issueクラスタ
type IssueCluster struct {
	Name    string
	Issues  []int
	Size    int
}

// AnalyzeLinks 参照リンク分析を実行
func AnalyzeLinks(links []github.Link, issues []github.Issue) *LinkAnalysis {
	result := &LinkAnalysis{
		TotalLinks:     len(links),
		LinkTypeCounts: make(map[string]int),
	}

	// issueタイトルのマップを作成
	issueTitles := make(map[int]string)
	for _, issue := range issues {
		issueTitles[issue.Number] = issue.Title
	}

	// リンクタイプ別カウント
	for _, link := range links {
		result.LinkTypeCounts[link.LinkType]++
	}

	// 被参照issue統計
	result.TopReferenced = analyzeReferencedIssues(links, issueTitles)

	// 参照issue統計
	result.TopReferencing = analyzeReferencingIssues(links, issueTitles)

	// クラスタ検出
	result.Clusters = detectClusters(links, issueTitles)

	// Mermaidグラフ生成
	result.MermaidGraph = generateMermaidGraph(links, issueTitles)

	return result
}

// analyzeReferencedIssues 被参照issue統計
func analyzeReferencedIssues(links []github.Link, titles map[int]string) []ReferencedIssue {
	refCount := make(map[int]int)

	for _, link := range links {
		refCount[link.TargetIssue]++
	}

	var referenced []ReferencedIssue
	for number, count := range refCount {
		referenced = append(referenced, ReferencedIssue{
			Number:         number,
			Title:          titles[number],
			ReferenceCount: count,
		})
	}

	// 参照回数順でソート
	sort.Slice(referenced, func(i, j int) bool {
		return referenced[i].ReferenceCount > referenced[j].ReferenceCount
	})

	// 上位10件を返す
	if len(referenced) > 10 {
		referenced = referenced[:10]
	}

	return referenced
}

// analyzeReferencingIssues 参照issue統計
func analyzeReferencingIssues(links []github.Link, titles map[int]string) []ReferencingIssue {
	linkCount := make(map[int]int)

	for _, link := range links {
		linkCount[link.SourceIssue]++
	}

	var referencing []ReferencingIssue
	for number, count := range linkCount {
		referencing = append(referencing, ReferencingIssue{
			Number:    number,
			Title:     titles[number],
			LinkCount: count,
		})
	}

	// リンク数順でソート
	sort.Slice(referencing, func(i, j int) bool {
		return referencing[i].LinkCount > referencing[j].LinkCount
	})

	// 上位10件を返す
	if len(referencing) > 10 {
		referencing = referencing[:10]
	}

	return referencing
}

// detectClusters Union-Find法でクラスタを検出
func detectClusters(links []github.Link, titles map[int]string) []IssueCluster {
	if len(links) == 0 {
		return nil
	}

	// Union-Find構造
	parent := make(map[int]int)

	// 再帰呼び出しのためvar宣言が必要
	var find func(x int) int
	find = func(x int) int {
		if _, ok := parent[x]; !ok {
			parent[x] = x
		}
		if parent[x] != x {
			parent[x] = find(parent[x]) // 経路圧縮
		}
		return parent[x]
	}

	union := func(x, y int) {
		px, py := find(x), find(y)
		if px != py {
			parent[px] = py
		}
	}

	// リンクを元にUnion
	for _, link := range links {
		union(link.SourceIssue, link.TargetIssue)
	}

	// クラスタをグループ化
	clusterMap := make(map[int][]int)
	for issue := range parent {
		root := find(issue)
		clusterMap[root] = append(clusterMap[root], issue)
	}

	// 2つ以上のissueを含むクラスタのみ
	var clusters []IssueCluster
	clusterNum := 1
	for _, issues := range clusterMap {
		if len(issues) >= 2 {
			sort.Ints(issues)
			clusters = append(clusters, IssueCluster{
				Name:   fmt.Sprintf("クラスタ%d", clusterNum),
				Issues: issues,
				Size:   len(issues),
			})
			clusterNum++
		}
	}

	// サイズ順でソート
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Size > clusters[j].Size
	})

	// 上位5件を返す
	if len(clusters) > 5 {
		clusters = clusters[:5]
	}

	return clusters
}

// generateMermaidGraph Mermaidフォーマットのグラフを生成
func generateMermaidGraph(links []github.Link, titles map[int]string) string {
	if len(links) == 0 {
		return ""
	}

	// リンク数の多い上位20件のリンクを使用
	selectedLinks := links
	if len(selectedLinks) > 20 {
		// リンクタイプの優先度でソート（closes > fixes > resolves > mentions）
		priority := map[string]int{
			"closes":   1,
			"fixes":    2,
			"resolves": 3,
			"mentions": 4,
		}
		sort.Slice(selectedLinks, func(i, j int) bool {
			return priority[selectedLinks[i].LinkType] < priority[selectedLinks[j].LinkType]
		})
		selectedLinks = selectedLinks[:20]
	}

	var graph string
	graph = "```mermaid\ngraph TD\n"

	// ノード定義
	usedNodes := make(map[int]bool)
	for _, link := range selectedLinks {
		usedNodes[link.SourceIssue] = true
		usedNodes[link.TargetIssue] = true
	}

	for node := range usedNodes {
		title := titles[node]
		if len(title) > 20 {
			title = title[:20] + "..."
		}
		// 特殊文字をエスケープ
		title = escapeForMermaid(title)
		graph += fmt.Sprintf("    N%d[\"#%d: %s\"]\n", node, node, title)
	}

	// エッジ定義
	for _, link := range selectedLinks {
		var label string
		switch link.LinkType {
		case "closes":
			label = "closes"
		case "fixes":
			label = "fixes"
		case "resolves":
			label = "resolves"
		default:
			label = "mentions"
		}
		graph += fmt.Sprintf("    N%d -->|%s| N%d\n", link.SourceIssue, label, link.TargetIssue)
	}

	graph += "```"
	return graph
}

// escapeForMermaid Mermaid用に特殊文字をエスケープ
func escapeForMermaid(s string) string {
	replacer := map[string]string{
		"\"": "'",
		"[":  "(",
		"]":  ")",
		"{":  "(",
		"}":  ")",
		"<":  "(",
		">":  ")",
	}

	for old, new := range replacer {
		for {
			newS := replaceOnce(s, old, new)
			if newS == s {
				break
			}
			s = newS
		}
	}

	return s
}

func replaceOnce(s, old, new string) string {
	for i := 0; i <= len(s)-len(old); i++ {
		if s[i:i+len(old)] == old {
			return s[:i] + new + s[i+len(old):]
		}
	}
	return s
}
