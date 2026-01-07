package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashimoto-kazuhiro-aa/dorapita-issue-analyzer/pkg/analysis"
	"github.com/hashimoto-kazuhiro-aa/dorapita-issue-analyzer/pkg/github"
)

// GenerateReport レポートを生成
func GenerateReport(issues []github.Issue, pattern *analysis.PatternAnalysis, team *analysis.TeamAnalysis, link *analysis.LinkAnalysis) string {
	var sb strings.Builder

	// ヘッダー
	sb.WriteString("# Dorapita Issue分析レポート\n\n")
	sb.WriteString(fmt.Sprintf("生成日時: %s\n", time.Now().Format("2006-01-02 15:04:05")))

	// 分析期間を計算
	if len(issues) > 0 {
		var oldest, newest time.Time
		for i, issue := range issues {
			if i == 0 || issue.CreatedAt.Before(oldest) {
				oldest = issue.CreatedAt
			}
			if i == 0 || issue.CreatedAt.After(newest) {
				newest = issue.CreatedAt
			}
		}
		sb.WriteString(fmt.Sprintf("分析期間: %s 〜 %s\n", oldest.Format("2006-01-02"), newest.Format("2006-01-02")))
	}

	sb.WriteString(fmt.Sprintf("総issue数: %d件\n\n", len(issues)))

	// パターン分析
	if pattern != nil {
		sb.WriteString(generatePatternSection(pattern))
	}

	// チーム対応分析
	if team != nil {
		sb.WriteString(generateTeamSection(team))
	}

	// 参照リンク分析
	if link != nil {
		sb.WriteString(generateLinkSection(link))
	}

	return sb.String()
}

// generatePatternSection パターン分析セクション
func generatePatternSection(p *analysis.PatternAnalysis) string {
	var sb strings.Builder

	sb.WriteString("## 1. パターン分析\n\n")

	// サマリー
	sb.WriteString("### サマリー\n\n")
	sb.WriteString(fmt.Sprintf("- 総issue数: %d件\n", p.TotalIssues))
	sb.WriteString(fmt.Sprintf("- Open: %d件 (%.1f%%)\n", p.OpenCount, float64(p.OpenCount)/float64(p.TotalIssues)*100))
	sb.WriteString(fmt.Sprintf("- Closed: %d件 (%.1f%%)\n", p.ClosedCount, float64(p.ClosedCount)/float64(p.TotalIssues)*100))
	sb.WriteString(fmt.Sprintf("- 平均クローズ時間: %.1f日\n\n", p.AverageCloseTime.Hours()/24))

	// 頻出キーワード
	sb.WriteString("### 頻出キーワード Top 20\n\n")
	sb.WriteString("| キーワード | 出現回数 | 割合 |\n")
	sb.WriteString("|-----------|---------|------|\n")
	for _, kw := range p.TopKeywords {
		sb.WriteString(fmt.Sprintf("| %s | %d | %.1f%% |\n", kw.Keyword, kw.Count, kw.Percent))
	}
	sb.WriteString("\n")

	// ラベル統計
	if len(p.LabelStats) > 0 {
		sb.WriteString("### ラベル別統計\n\n")
		sb.WriteString("| ラベル | issue数 | 割合 | 平均クローズ時間 |\n")
		sb.WriteString("|-------|---------|------|----------------|\n")
		for _, label := range p.LabelStats {
			avgClose := "-"
			if label.AvgCloseTimeDays > 0 {
				avgClose = fmt.Sprintf("%.1f日", label.AvgCloseTimeDays)
			}
			sb.WriteString(fmt.Sprintf("| %s | %d | %.1f%% | %s |\n", label.Name, label.Count, label.Percent, avgClose))
		}
		sb.WriteString("\n")
	}

	// 月別トレンド
	if len(p.MonthlyTrend) > 0 {
		sb.WriteString("### 月別issue作成数\n\n")
		sb.WriteString("| 月 | 総数 | Open | Closed |\n")
		sb.WriteString("|-----|------|------|--------|\n")
		for _, mt := range p.MonthlyTrend {
			sb.WriteString(fmt.Sprintf("| %s | %d | %d | %d |\n", mt.Month, mt.TotalCount, mt.OpenCount, mt.ClosedCount))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// generateTeamSection チーム対応分析セクション
func generateTeamSection(t *analysis.TeamAnalysis) string {
	var sb strings.Builder

	sb.WriteString("## 2. チーム対応分析\n\n")

	// 担当者別統計
	sb.WriteString("### 担当者別統計\n\n")
	sb.WriteString("| 担当者 | アサイン数 | クローズ数 | クローズ率 | 平均クローズ時間 | 平均応答時間 |\n")
	sb.WriteString("|-------|----------|----------|----------|----------------|------------|\n")
	for _, stat := range t.AssigneeStats {
		avgClose := "-"
		if stat.AvgCloseTimeDays > 0 {
			avgClose = fmt.Sprintf("%.1f日", stat.AvgCloseTimeDays)
		}
		avgResponse := "-"
		if stat.AvgResponseHours > 0 {
			if stat.AvgResponseHours < 24 {
				avgResponse = fmt.Sprintf("%.1f時間", stat.AvgResponseHours)
			} else {
				avgResponse = fmt.Sprintf("%.1f日", stat.AvgResponseHours/24)
			}
		}
		sb.WriteString(fmt.Sprintf("| %s | %d | %d | %.1f%% | %s | %s |\n",
			stat.Login, stat.TotalAssigned, stat.ClosedCount, stat.CloseRate, avgClose, avgResponse))
	}
	sb.WriteString("\n")

	// 作成者別統計
	if len(t.AuthorStats) > 0 {
		sb.WriteString("### issue作成者別統計\n\n")
		sb.WriteString("| 作成者 | issue数 | Closed | Open |\n")
		sb.WriteString("|-------|---------|--------|------|\n")
		// 上位10件のみ
		limit := 10
		if len(t.AuthorStats) < limit {
			limit = len(t.AuthorStats)
		}
		for i := 0; i < limit; i++ {
			stat := t.AuthorStats[i]
			sb.WriteString(fmt.Sprintf("| %s | %d | %d | %d |\n",
				stat.Login, stat.IssueCount, stat.ClosedCount, stat.OpenCount))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// generateLinkSection 参照リンク分析セクション
func generateLinkSection(l *analysis.LinkAnalysis) string {
	var sb strings.Builder

	sb.WriteString("## 3. 参照リンク分析\n\n")

	// サマリー
	sb.WriteString("### サマリー\n\n")
	sb.WriteString(fmt.Sprintf("- 総リンク数: %d\n", l.TotalLinks))

	// リンクタイプ別
	if len(l.LinkTypeCounts) > 0 {
		sb.WriteString("\n**リンクタイプ別:**\n")
		for linkType, count := range l.LinkTypeCounts {
			sb.WriteString(fmt.Sprintf("- %s: %d件\n", linkType, count))
		}
	}
	sb.WriteString("\n")

	// 被参照issue Top 10
	if len(l.TopReferenced) > 0 {
		sb.WriteString("### 被参照issue Top 10\n\n")
		sb.WriteString("| issue | タイトル | 被参照回数 |\n")
		sb.WriteString("|-------|---------|----------|\n")
		for _, ref := range l.TopReferenced {
			title := ref.Title
			if len(title) > 50 {
				title = title[:50] + "..."
			}
			sb.WriteString(fmt.Sprintf("| #%d | %s | %d |\n", ref.Number, title, ref.ReferenceCount))
		}
		sb.WriteString("\n")
	}

	// 参照issue Top 10
	if len(l.TopReferencing) > 0 {
		sb.WriteString("### 参照issue Top 10\n\n")
		sb.WriteString("| issue | タイトル | 参照数 |\n")
		sb.WriteString("|-------|---------|-------|\n")
		for _, ref := range l.TopReferencing {
			title := ref.Title
			if len(title) > 50 {
				title = title[:50] + "..."
			}
			sb.WriteString(fmt.Sprintf("| #%d | %s | %d |\n", ref.Number, title, ref.LinkCount))
		}
		sb.WriteString("\n")
	}

	// クラスタ
	if len(l.Clusters) > 0 {
		sb.WriteString("### 関連issueクラスタ\n\n")
		for _, cluster := range l.Clusters {
			sb.WriteString(fmt.Sprintf("**%s** (サイズ: %d)\n", cluster.Name, cluster.Size))
			sb.WriteString("- ")
			for i, num := range cluster.Issues {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(fmt.Sprintf("#%d", num))
			}
			sb.WriteString("\n\n")
		}
	}

	// Mermaidグラフ
	if l.MermaidGraph != "" {
		sb.WriteString("### 参照グラフ\n\n")
		sb.WriteString(l.MermaidGraph)
		sb.WriteString("\n\n")
	}

	return sb.String()
}
