package analysis

import (
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/hashimoto-kazuhiro-aa/dorapita-issue-analyzer/pkg/github"
)

// PatternAnalysis パターン分析結果
type PatternAnalysis struct {
	TotalIssues       int
	OpenCount         int
	ClosedCount       int
	TopKeywords       []KeywordCount
	LabelStats        []LabelStat
	MonthlyTrend      []MonthlyCount
	StateDistribution map[string]int
	AverageCloseTime  time.Duration
}

// KeywordCount キーワード出現回数
type KeywordCount struct {
	Keyword string
	Count   int
	Percent float64
}

// LabelStat ラベル統計
type LabelStat struct {
	Name             string
	Count            int
	Percent          float64
	AvgCloseTimeDays float64
}

// MonthlyCount 月別カウント
type MonthlyCount struct {
	Month       string // YYYY-MM形式
	OpenCount   int
	ClosedCount int
	TotalCount  int
}

// AnalyzePattern パターン分析を実行
func AnalyzePattern(issues []github.Issue) *PatternAnalysis {
	result := &PatternAnalysis{
		TotalIssues:       len(issues),
		StateDistribution: make(map[string]int),
	}

	// 状態別カウント
	for _, issue := range issues {
		result.StateDistribution[issue.State]++
		if issue.State == "open" {
			result.OpenCount++
		} else if issue.State == "closed" {
			result.ClosedCount++
		}
	}

	// キーワード分析
	result.TopKeywords = analyzeKeywords(issues, 20)

	// ラベル統計
	result.LabelStats = analyzeLabelStats(issues)

	// 月別トレンド
	result.MonthlyTrend = analyzeMonthlyTrend(issues)

	// 平均クローズ時間
	result.AverageCloseTime = calcAverageCloseTime(issues)

	return result
}

// analyzeKeywords キーワード分析
func analyzeKeywords(issues []github.Issue, topN int) []KeywordCount {
	wordCount := make(map[string]int)
	totalWords := 0

	// ストップワード（除外する一般的な単語）
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "is": true, "are": true, "was": true, "were": true,
		"be": true, "been": true, "being": true, "have": true, "has": true, "had": true,
		"do": true, "does": true, "did": true, "will": true, "would": true, "could": true, "should": true,
		"may": true, "might": true, "must": true, "shall": true,
		"this": true, "that": true, "these": true, "those": true,
		"i": true, "you": true, "he": true, "she": true, "it": true, "we": true, "they": true,
		"what": true, "which": true, "who": true, "whom": true, "whose": true,
		"where": true, "when": true, "why": true, "how": true,
		"all": true, "each": true, "every": true, "both": true, "few": true, "more": true, "most": true,
		"other": true, "some": true, "such": true, "no": true, "nor": true, "not": true, "only": true,
		"same": true, "so": true, "than": true, "too": true, "very": true,
		"can": true, "just": true, "don": true, "now": true,
		"and": true, "or": true, "but": true, "if": true, "because": true, "as": true, "until": true,
		"while": true, "of": true, "at": true, "by": true, "for": true, "with": true, "about": true,
		"against": true, "between": true, "into": true, "through": true, "during": true, "before": true,
		"after": true, "above": true, "below": true, "to": true, "from": true, "up": true, "down": true,
		"in": true, "out": true, "on": true, "off": true, "over": true, "under": true,
		"again": true, "further": true, "then": true, "once": true,
		// 日本語のストップワード
		"の": true, "に": true, "は": true, "を": true, "た": true, "が": true, "で": true, "て": true,
		"と": true, "し": true, "れ": true, "さ": true, "ある": true, "いる": true, "も": true,
		"する": true, "から": true, "な": true, "こと": true, "として": true, "い": true, "や": true,
		"れる": true, "など": true, "なっ": true, "ない": true, "この": true, "ため": true, "その": true,
		"あっ": true, "よう": true, "また": true, "もの": true, "という": true, "あり": true,
		"まで": true, "られ": true, "なる": true, "へ": true, "か": true, "だ": true, "これ": true,
		"によって": true, "により": true, "おり": true, "より": true, "による": true,
		"ず": true, "なり": true, "られる": true, "において": true, "ば": true, "なかっ": true,
		"なく": true, "しかし": true, "について": true, "せ": true, "だっ": true, "その他": true,
		"できる": true, "それ": true, "ここ": true, "ところ": true, "ので": true, "です": true, "ます": true,
	}

	for _, issue := range issues {
		// タイトルと本文を分析
		text := issue.Title + " " + issue.Body
		words := tokenize(text)

		for _, word := range words {
			lower := strings.ToLower(word)
			if len(lower) > 1 && !stopWords[lower] {
				wordCount[lower]++
				totalWords++
			}
		}
	}

	// ソート
	var keywords []KeywordCount
	for word, count := range wordCount {
		keywords = append(keywords, KeywordCount{
			Keyword: word,
			Count:   count,
			Percent: float64(count) / float64(totalWords) * 100,
		})
	}

	sort.Slice(keywords, func(i, j int) bool {
		return keywords[i].Count > keywords[j].Count
	})

	// 上位N件を返す
	if len(keywords) > topN {
		keywords = keywords[:topN]
	}

	return keywords
}

// tokenize テキストをトークン化
func tokenize(text string) []string {
	var tokens []string
	var current strings.Builder

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current.WriteRune(r)
		} else {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

// analyzeLabelStats ラベル統計
func analyzeLabelStats(issues []github.Issue) []LabelStat {
	labelIssues := make(map[string][]github.Issue)

	for _, issue := range issues {
		for _, label := range issue.Labels {
			labelIssues[label.Name] = append(labelIssues[label.Name], issue)
		}
	}

	var stats []LabelStat
	totalIssues := len(issues)

	for name, labeledIssues := range labelIssues {
		var closedCount int
		var totalCloseTime time.Duration

		for _, issue := range labeledIssues {
			if issue.State == "closed" && issue.ClosedAt != nil {
				closedCount++
				totalCloseTime += issue.ClosedAt.Sub(issue.CreatedAt)
			}
		}

		var avgCloseTimeDays float64
		if closedCount > 0 {
			avgCloseTimeDays = totalCloseTime.Hours() / float64(closedCount) / 24
		}

		stats = append(stats, LabelStat{
			Name:             name,
			Count:            len(labeledIssues),
			Percent:          float64(len(labeledIssues)) / float64(totalIssues) * 100,
			AvgCloseTimeDays: avgCloseTimeDays,
		})
	}

	// カウント順でソート
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Count > stats[j].Count
	})

	return stats
}

// analyzeMonthlyTrend 月別トレンド
func analyzeMonthlyTrend(issues []github.Issue) []MonthlyCount {
	monthlyData := make(map[string]*MonthlyCount)

	for _, issue := range issues {
		month := issue.CreatedAt.Format("2006-01")

		if _, ok := monthlyData[month]; !ok {
			monthlyData[month] = &MonthlyCount{Month: month}
		}

		monthlyData[month].TotalCount++
		if issue.State == "open" {
			monthlyData[month].OpenCount++
		} else if issue.State == "closed" {
			monthlyData[month].ClosedCount++
		}
	}

	// スライスに変換してソート
	var trend []MonthlyCount
	for _, mc := range monthlyData {
		trend = append(trend, *mc)
	}

	sort.Slice(trend, func(i, j int) bool {
		return trend[i].Month < trend[j].Month
	})

	return trend
}

// calcAverageCloseTime 平均クローズ時間を計算
func calcAverageCloseTime(issues []github.Issue) time.Duration {
	var total time.Duration
	var count int

	for _, issue := range issues {
		if issue.State == "closed" && issue.ClosedAt != nil {
			total += issue.ClosedAt.Sub(issue.CreatedAt)
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return total / time.Duration(count)
}
