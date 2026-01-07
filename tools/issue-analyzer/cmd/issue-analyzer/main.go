package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hashimoto-kazuhiro-aa/dorapita-issue-analyzer/pkg/analysis"
	"github.com/hashimoto-kazuhiro-aa/dorapita-issue-analyzer/pkg/cache"
	"github.com/hashimoto-kazuhiro-aa/dorapita-issue-analyzer/pkg/github"
	"github.com/hashimoto-kazuhiro-aa/dorapita-issue-analyzer/pkg/output"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "fetch":
		if err := runFetch(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
	case "analyze":
		if err := runAnalyze(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "不明なコマンド: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Dorapita Issue Analyzer

Usage:
  issue-analyzer <command> [options]

Commands:
  fetch    GitHub issueを取得してキャッシュに保存
  analyze  キャッシュを分析してレポートを生成
  help     このヘルプを表示

Examples:
  # issue取得（直近1年間）
  issue-analyzer fetch --repo ZIGExN/dorapita --days 365 --cache ./data/issues.db

  # 全分析を実行
  issue-analyzer analyze --cache ./data/issues.db --output ./out/report.md

  # パターン分析のみ
  issue-analyzer analyze --cache ./data/issues.db --type pattern

Use "issue-analyzer <command> --help" for more information about a command.`)
}

func runFetch(args []string) error {
	fs := flag.NewFlagSet("fetch", flag.ExitOnError)
	repo := fs.String("repo", "", "リポジトリ名（owner/repo形式）")
	days := fs.Int("days", 365, "取得期間（日数）")
	cachePath := fs.String("cache", "./issues.db", "キャッシュファイルパス")
	workDir := fs.String("workdir", "", "gh-wrapper.sh実行用作業ディレクトリ（PAT抽出用）")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *repo == "" {
		return fmt.Errorf("--repo オプションは必須です")
	}

	fmt.Printf("issue取得開始: %s (直近%d日間)\n", *repo, *days)

	// GitHubクライアント作成
	client, err := github.NewClient(*repo)
	if err != nil {
		return fmt.Errorf("GitHubクライアント作成エラー: %w", err)
	}

	if *workDir != "" {
		client.SetWorkDir(*workDir)
	}

	// issue一覧取得
	fmt.Println("issue一覧を取得中...")
	issues, err := client.ListIssues(*days)
	if err != nil {
		return fmt.Errorf("issue一覧取得エラー: %w", err)
	}
	fmt.Printf("取得したissue数: %d\n", len(issues))

	// キャッシュに保存
	fmt.Println("キャッシュに保存中...")
	c, err := cache.New(*cachePath)
	if err != nil {
		return fmt.Errorf("キャッシュ作成エラー: %w", err)
	}
	defer c.Close()

	if err := c.SaveIssues(issues); err != nil {
		return fmt.Errorf("issue保存エラー: %w", err)
	}

	// コメントと参照リンクを取得・保存
	fmt.Println("コメントと参照リンクを取得中...")
	var allLinks []github.Link
	for i, issue := range issues {
		if (i+1)%10 == 0 {
			fmt.Printf("進捗: %d/%d\n", i+1, len(issues))
		}

		// コメント取得
		comments, err := client.GetComments(issue.Number)
		if err != nil {
			fmt.Fprintf(os.Stderr, "警告: issue #%d のコメント取得に失敗: %v\n", issue.Number, err)
			continue
		}

		if len(comments) > 0 {
			if err := c.SaveComments(issue.Number, comments); err != nil {
				fmt.Fprintf(os.Stderr, "警告: issue #%d のコメント保存に失敗: %v\n", issue.Number, err)
			}
		}

		// 参照リンク抽出
		issue.Comments = comments
		links := github.ParseAllLinks(&issue)
		allLinks = append(allLinks, links...)
	}

	// 参照リンク保存
	if len(allLinks) > 0 {
		fmt.Printf("参照リンク数: %d\n", len(allLinks))
		if err := c.SaveLinks(allLinks); err != nil {
			return fmt.Errorf("リンク保存エラー: %w", err)
		}
	}

	fmt.Println("完了!")
	fmt.Printf("キャッシュファイル: %s\n", *cachePath)

	return nil
}

func runAnalyze(args []string) error {
	fs := flag.NewFlagSet("analyze", flag.ExitOnError)
	cachePath := fs.String("cache", "./issues.db", "キャッシュファイルパス")
	analysisType := fs.String("type", "all", "分析タイプ（pattern, team, links, all）")
	outputPath := fs.String("output", "", "出力ファイルパス（空の場合は標準出力）")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// キャッシュを開く
	c, err := cache.New(*cachePath)
	if err != nil {
		return fmt.Errorf("キャッシュオープンエラー: %w", err)
	}
	defer c.Close()

	// issueを読み込み
	issues, err := c.LoadIssues()
	if err != nil {
		return fmt.Errorf("issue読み込みエラー: %w", err)
	}

	if len(issues) == 0 {
		return fmt.Errorf("キャッシュにissueがありません。先に fetch コマンドを実行してください")
	}

	// コメントを読み込み（チーム分析用）
	comments := make(map[int][]github.Comment)
	for _, issue := range issues {
		issueComments, err := c.LoadComments(issue.Number)
		if err == nil && len(issueComments) > 0 {
			comments[issue.Number] = issueComments
		}
	}

	// リンクを読み込み
	links, err := c.LoadLinks()
	if err != nil {
		return fmt.Errorf("リンク読み込みエラー: %w", err)
	}

	// 分析実行
	var patternResult *analysis.PatternAnalysis
	var teamResult *analysis.TeamAnalysis
	var linkResult *analysis.LinkAnalysis

	switch *analysisType {
	case "pattern":
		patternResult = analysis.AnalyzePattern(issues)
	case "team":
		teamResult = analysis.AnalyzeTeam(issues, comments)
	case "links":
		linkResult = analysis.AnalyzeLinks(links, issues)
	case "all":
		patternResult = analysis.AnalyzePattern(issues)
		teamResult = analysis.AnalyzeTeam(issues, comments)
		linkResult = analysis.AnalyzeLinks(links, issues)
	default:
		return fmt.Errorf("不明な分析タイプ: %s", *analysisType)
	}

	// レポート生成
	report := output.GenerateReport(issues, patternResult, teamResult, linkResult)

	// 出力
	if *outputPath == "" {
		fmt.Println(report)
	} else {
		if err := os.WriteFile(*outputPath, []byte(report), 0644); err != nil {
			return fmt.Errorf("ファイル出力エラー: %w", err)
		}
		fmt.Printf("レポートを出力しました: %s\n", *outputPath)
	}

	return nil
}
