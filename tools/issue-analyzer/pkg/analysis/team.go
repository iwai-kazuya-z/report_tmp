package analysis

import (
	"sort"
	"time"

	"github.com/hashimoto-kazuhiro-aa/dorapita-issue-analyzer/pkg/github"
)

// TeamAnalysis チーム対応分析結果
type TeamAnalysis struct {
	AssigneeStats   []AssigneeStat
	AuthorStats     []AuthorStat
	MonthlyActivity []MonthlyActivity
}

// AssigneeStat 担当者別統計
type AssigneeStat struct {
	Login             string
	TotalAssigned     int
	ClosedCount       int
	OpenCount         int
	CloseRate         float64
	AvgCloseTimeDays  float64
	AvgResponseHours  float64 // 最初のコメントまでの時間
}

// AuthorStat 作成者別統計
type AuthorStat struct {
	Login       string
	IssueCount  int
	ClosedCount int
	OpenCount   int
}

// MonthlyActivity 月別活動量
type MonthlyActivity struct {
	Month    string // YYYY-MM形式
	Assignee string
	Closed   int
	Assigned int
}

// AnalyzeTeam チーム対応分析を実行
func AnalyzeTeam(issues []github.Issue, comments map[int][]github.Comment) *TeamAnalysis {
	result := &TeamAnalysis{}

	// 担当者別統計
	result.AssigneeStats = analyzeAssigneeStats(issues, comments)

	// 作成者別統計
	result.AuthorStats = analyzeAuthorStats(issues)

	// 月別活動量
	result.MonthlyActivity = analyzeMonthlyActivity(issues)

	return result
}

// analyzeAssigneeStats 担当者別統計
func analyzeAssigneeStats(issues []github.Issue, comments map[int][]github.Comment) []AssigneeStat {
	assigneeData := make(map[string]*struct {
		total          int
		closed         int
		open           int
		totalCloseTime time.Duration
		totalResponse  time.Duration
		responseCount  int
	})

	for _, issue := range issues {
		// 担当者を取得
		var assignee string
		if len(issue.Assignees) > 0 {
			assignee = issue.Assignees[0].Login
		} else {
			assignee = "(未割り当て)"
		}

		if _, ok := assigneeData[assignee]; !ok {
			assigneeData[assignee] = &struct {
				total          int
				closed         int
				open           int
				totalCloseTime time.Duration
				totalResponse  time.Duration
				responseCount  int
			}{}
		}

		data := assigneeData[assignee]
		data.total++

		if issue.State == "closed" {
			data.closed++
			if issue.ClosedAt != nil {
				data.totalCloseTime += issue.ClosedAt.Sub(issue.CreatedAt)
			}
		} else {
			data.open++
		}

		// 最初のコメントまでの応答時間
		if issueComments, ok := comments[issue.Number]; ok && len(issueComments) > 0 {
			// 作成者以外の最初のコメントを探す
			for _, comment := range issueComments {
				if comment.Author.Login != issue.Author.Login {
					data.totalResponse += comment.CreatedAt.Sub(issue.CreatedAt)
					data.responseCount++
					break
				}
			}
		}
	}

	// 統計を計算
	var stats []AssigneeStat
	for login, data := range assigneeData {
		stat := AssigneeStat{
			Login:         login,
			TotalAssigned: data.total,
			ClosedCount:   data.closed,
			OpenCount:     data.open,
		}

		if data.total > 0 {
			stat.CloseRate = float64(data.closed) / float64(data.total) * 100
		}

		if data.closed > 0 {
			stat.AvgCloseTimeDays = data.totalCloseTime.Hours() / float64(data.closed) / 24
		}

		if data.responseCount > 0 {
			stat.AvgResponseHours = data.totalResponse.Hours() / float64(data.responseCount)
		}

		stats = append(stats, stat)
	}

	// アサイン数順でソート
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].TotalAssigned > stats[j].TotalAssigned
	})

	return stats
}

// analyzeAuthorStats 作成者別統計
func analyzeAuthorStats(issues []github.Issue) []AuthorStat {
	authorData := make(map[string]*AuthorStat)

	for _, issue := range issues {
		author := issue.Author.Login
		if author == "" {
			author = "(不明)"
		}

		if _, ok := authorData[author]; !ok {
			authorData[author] = &AuthorStat{Login: author}
		}

		authorData[author].IssueCount++
		if issue.State == "closed" {
			authorData[author].ClosedCount++
		} else {
			authorData[author].OpenCount++
		}
	}

	var stats []AuthorStat
	for _, stat := range authorData {
		stats = append(stats, *stat)
	}

	// issue数順でソート
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].IssueCount > stats[j].IssueCount
	})

	return stats
}

// analyzeMonthlyActivity 月別活動量
func analyzeMonthlyActivity(issues []github.Issue) []MonthlyActivity {
	type key struct {
		month    string
		assignee string
	}

	activityData := make(map[key]*MonthlyActivity)

	for _, issue := range issues {
		month := issue.CreatedAt.Format("2006-01")

		var assignee string
		if len(issue.Assignees) > 0 {
			assignee = issue.Assignees[0].Login
		} else {
			assignee = "(未割り当て)"
		}

		k := key{month: month, assignee: assignee}

		if _, ok := activityData[k]; !ok {
			activityData[k] = &MonthlyActivity{
				Month:    month,
				Assignee: assignee,
			}
		}

		activityData[k].Assigned++
		if issue.State == "closed" {
			activityData[k].Closed++
		}
	}

	var activities []MonthlyActivity
	for _, activity := range activityData {
		activities = append(activities, *activity)
	}

	// 月、担当者順でソート
	sort.Slice(activities, func(i, j int) bool {
		if activities[i].Month != activities[j].Month {
			return activities[i].Month < activities[j].Month
		}
		return activities[i].Assignee < activities[j].Assignee
	})

	return activities
}
