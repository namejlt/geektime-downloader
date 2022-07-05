package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	audiodown "github.com/namejlt/geektime-downloader/internal/pkg/audio"
	"github.com/namejlt/geektime-downloader/internal/pkg/markdown"
	videodown "github.com/namejlt/geektime-downloader/internal/pkg/video"
	"os"
	"path/filepath"
	"strconv"
	"sync/atomic"

	"github.com/briandowns/spinner"
	"github.com/namejlt/geektime-downloader/cmd/prompt"
	"github.com/namejlt/geektime-downloader/internal/geektime"
	"github.com/namejlt/geektime-downloader/internal/loader"
	"github.com/namejlt/geektime-downloader/internal/pkg/chromedp"
	"github.com/namejlt/geektime-downloader/internal/pkg/util"
	"github.com/spf13/cobra"
)

var (
	phone              string
	downloadFolder     string
	l                  *spinner.Spinner
	columnDiyId        int
	sleep              int
	sleepMax           int
	reLogin            bool
	columns            []geektime.ColumnSummary
	currentColumnIndex int
	quality            string

	//脚本批量下载
	columnIDsFile string //cid 每行一个

	//下载类型
	pdf   bool
	md    bool
	audio bool
	video bool
)

func init() {
	userHomeDir, _ := os.UserHomeDir()
	defaultDownloadFolder := filepath.Join(userHomeDir, util.GeektimeDownloaderFolder)

	//公共参数

	rootCmd.PersistentFlags().StringVarP(&phone, "phone", "u", "", "你的极客时间账号(手机号)(required)")
	_ = rootCmd.MarkFlagRequired("phone")
	rootCmd.PersistentFlags().StringVarP(&downloadFolder, "folder", "f", defaultDownloadFolder, "下载目标位置")
	rootCmd.PersistentFlags().BoolVarP(&pdf, "pdf", "p", true, "PDF是否下载")
	rootCmd.PersistentFlags().BoolVarP(&md, "md", "d", false, "markdown是否下载")
	rootCmd.PersistentFlags().BoolVarP(&audio, "audio", "a", false, "音频是否下载")
	rootCmd.PersistentFlags().BoolVarP(&video, "video", "v", false, "视频是否下载")
	rootCmd.PersistentFlags().StringVarP(&quality, "quality", "q", "sd", "下载视频清晰度(ld标清,sd高清,hd超清)")

	selectDiyCmd.Flags().IntVarP(&columnDiyId, "column_diy_id", "i", 0, "指定下载课程id")
	_ = selectDiyCmd.MarkFlagRequired("column_diy_id")
	selectDiyCmd.Flags().IntVarP(&sleep, "sleep", "s", 1000, "下载文章间隔时间 毫秒")
	selectDiyCmd.Flags().BoolVarP(&reLogin, "relogin", "r", false, "是否重新登录")

	batchDownCmd.Flags().StringVarP(&columnIDsFile, "column_ids_file", "i", "./doc/cid.txt", "指定下载课程id文件")
	_ = batchDownCmd.MarkFlagRequired("column_ids_file")
	batchDownCmd.Flags().IntVarP(&sleep, "sleep", "s", 1000, "下载文章间隔时间 毫秒")
	batchDownCmd.Flags().IntVarP(&sleepMax, "sleepmax", "m", 1200, "下载文章间隔时间 毫秒 max")
	batchDownCmd.Flags().BoolVarP(&reLogin, "relogin", "r", false, "是否重新登录")

	l = loader.NewSpinner()
}

// rootCmd 执行命令
var rootCmd = &cobra.Command{
	Use:   "geektime-downloader",
	Short: "Geektime-downloader is used to download geek time lessons",
	Run: func(cmd *cobra.Command, args []string) {
		//不执行内容 具体执行在子命令
	},
}

func selectColumn(client *resty.Client) {
	if len(columns) == 0 {
		if err := util.RemoveConfig(phone); err != nil {
			printErrAndExit(err)
		} else {
			fmt.Println("当前账户在其他设备登录, 请尝试重新登录")
			os.Exit(1)
		}
	}
	currentColumnIndex = prompt.SelectColumn(columns)
	handleSelectColumn(client)
}

func handleSelectColumn(client *resty.Client) {
	option := prompt.SelectDownLoadAllOrSelectArticles(columns[currentColumnIndex].Title)
	handleSelectDownloadAll(option, client)
}

/**

{"返回上一级", 0},
{"下载当前专栏所有文章", 1},
{"选择文章", 2},

*/

func handleSelectDownloadAll(option int, client *resty.Client) {
	switch option {
	case 0:
		selectColumn(client)
	case 1:
		handleDownloadAll(client, true, []geektime.ArticleSummary{}, 0)
	case 2:
		selectArticle(client)
	}
}

func selectArticle(client *resty.Client) {
	articles := loadArticles(client)
	index := prompt.SelectArticles(articles)
	handleSelectArticle(articles, index, client)
}

//下载单个文章
func handleSelectArticle(articles []geektime.ArticleSummary, index int, client *resty.Client) {
	if index == 0 {
		handleSelectColumn(client)
	}
	/*a := articles[index-1]
	folder, err := mkColumnDownloadFolder(phone, columns[currentColumnIndex].Title)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	loader.Run(l, fmt.Sprintf("[ 正在下载文章 《%s》... ]", a.Title), func() {
		err := chromedp.PrintArticlePageToPDF(a.AID, filepath.Join(folder, util.FileName(a.Title, "pdf")), client.Cookies)
		if err != nil {
			printErrAndExit(err)
		}
	})*/

	handleDownloadAll(client, false, articles, index)

	selectArticle(client)
}

//下载所有
func handleDownloadAll(client *resty.Client, pause bool, articles []geektime.ArticleSummary, index int) {
	cTitle := columns[currentColumnIndex].Title
	isColumn := isColumn(columns[currentColumnIndex].Type)
	isVideo := isVideo(columns[currentColumnIndex].Type)
	if len(articles) == 0 {
		articles = loadArticles(client)
	}

	if index != 0 { //只获取指定文章
		articles = articles[index-1 : index]
	}

	var counter uint64
	folder, err := mkColumnDownloadFolder(phone, cTitle)
	if err != nil {
		printErrAndExit(err)
	}
	ctx := context.Background()
	//单个处理 并sleep
	for _, a := range articles {
		aid := a.AID
		title := a.Title
		prefix := fmt.Sprintf("[ 正在下载专栏 《%s》 中的所有文章或视频, 已完成下载%d/%d ... ]", cTitle, counter, len(articles))

		if pdf { //无论课程还是视频都可以生成pdf
			loader.Run(l, prefix, func() {
				err := chromedp.PrintArticlePageToPDF(aid, filepath.Join(folder, util.FileName(title, "pdf")), client.Cookies)
				if err != nil {
					printErrAndExit(err)
				} else {
					atomic.AddUint64(&counter, 1)
				}
			})
		}
		if md || audio {
			if isColumn {
				articleInfo, err := geektime.GetArticleInfo(a.AID, client)
				checkError(err)
				if md {
					err = markdown.Download(ctx, articleInfo.ArticleContent, a.Title, folder, a.AID, 1)
				}
				if audio {
					err = audiodown.DownloadAudio(ctx, articleInfo.AudioDownloadURL, folder, a.Title)
				}
			} else {
				println(cTitle, a.Title, "非课程")
			}
		}

		if video {
			if isVideo {
				videoInfo, err := geektime.GetVideoInfo(a.AID, client, quality)
				checkError(err)
				err = videodown.DownloadVideo(ctx, videoInfo.M3U8URL, a.Title, folder, int64(videoInfo.Size), 1)
				checkError(err)
			}
			println(cTitle, a.Title, "无视频")
		}

		util.SleepMS(sleep, sleepMax)
	}

	if pause {
		selectColumn(client)
	}
}

func loadArticles(client *resty.Client) []geektime.ArticleSummary {
	c := columns[currentColumnIndex]
	if len(c.Articles) <= 0 {
		loader.Run(l, "[ 正在加载文章列表...]", func() {
			articles, err := geektime.GetArticles(strconv.Itoa(c.CID), client)
			if err != nil {
				printErrAndExit(err)
			}
			columns[currentColumnIndex].Articles = articles
		})
	}
	return columns[currentColumnIndex].Articles
}

func mkColumnDownloadFolder(phone, columnName string) (string, error) {
	path := filepath.Join(downloadFolder, phone, columnName)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return path, nil
}

// Execute func 执行函数
func Execute() {
	rootCmd.AddCommand(selectColumnCmd, selectDiyCmd, batchDownCmd)

	if err := rootCmd.Execute(); err != nil {
		printErrAndExit(err)
	}
}

func printErrAndExit(err error) { // 打印报错
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func printMsgAndExit(msg string) { // 打印信息
	fmt.Fprintln(os.Stdout, msg)
	os.Exit(1)
}
