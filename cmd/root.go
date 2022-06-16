package cmd

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync/atomic"

	"github.com/briandowns/spinner"
	"github.com/gammazero/workerpool"
	"github.com/namejlt/geektime-downloader/cmd/prompt"
	"github.com/namejlt/geektime-downloader/internal/geektime"
	"github.com/namejlt/geektime-downloader/internal/loader"
	"github.com/namejlt/geektime-downloader/internal/pkg/chromedp"
	"github.com/namejlt/geektime-downloader/internal/pkg/util"
	"github.com/spf13/cobra"
)

var (
	phone              string
	concurrency        int
	downloadFolder     string
	l                  *spinner.Spinner
	columnDiyId        int
	sleep              int
	reLogin            bool
	columns            []geektime.ColumnSummary
	currentColumnIndex int
)

func init() {
	userHomeDir, _ := os.UserHomeDir()
	defaultConcurency := int(math.Ceil(float64(runtime.NumCPU()) / 2.0))
	defaultDownloadFolder := filepath.Join(userHomeDir, util.GeektimeDownloaderFolder)
	selectColumnCmd.Flags().StringVarP(&phone, "phone", "u", "", "你的极客时间账号(手机号)(required)")
	_ = selectColumnCmd.MarkFlagRequired("phone")
	selectColumnCmd.Flags().StringVarP(&downloadFolder, "folder", "f", defaultDownloadFolder, "PDF 文件下载目标位置")
	selectColumnCmd.Flags().IntVarP(&concurrency, "concurrency", "c", defaultConcurency, "下载文章的并发数")

	selectDiyCmd.Flags().StringVarP(&phone, "phone", "u", "", "你的极客时间账号(手机号)(required)")
	_ = selectDiyCmd.MarkFlagRequired("phone")
	selectDiyCmd.Flags().StringVarP(&downloadFolder, "folder", "f", defaultDownloadFolder, "PDF 文件下载目标位置")
	selectDiyCmd.Flags().IntVarP(&concurrency, "concurrency", "c", defaultConcurency, "下载文章的并发数 0 代表不并发且有等待时间")
	selectDiyCmd.Flags().IntVarP(&columnDiyId, "column_diy_id", "i", 0, "指定下载课程id")
	selectDiyCmd.Flags().IntVarP(&sleep, "sleep", "s", 1000, "下载文章间隔时间 毫秒")
	selectDiyCmd.Flags().BoolVarP(&reLogin, "relogin", "r", false, "是否重新登录")
	_ = selectDiyCmd.MarkFlagRequired("column_diy_id")
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

// selectColumnCmd 选择栏目
var selectColumnCmd = &cobra.Command{
	Use:   "columns",
	Short: "Geektime-downloader is used to download geek time lessons",
	Run: func(cmd *cobra.Command, args []string) {
		readCookies, err := util.ReadCookieFromConfigFile(phone) //获取登录态
		if err != nil {
			printErrAndExit(err)
		}
		if readCookies == nil { // 不存在 则登录
			pwd := prompt.GetPwd()
			loader.Run(l, "[ 正在登录... ]", func() {
				errMsg, cookies := geektime.Login(phone, pwd)
				if errMsg != "" {
					fmt.Fprintln(os.Stderr, errMsg)
					os.Exit(1)
				}
				readCookies = cookies
				err := util.WriteCookieToConfigFile(phone, cookies) //保存登录态
				if err != nil {
					printErrAndExit(err)
				}
			})
			fmt.Println("登录成功")
		}
		client := geektime.NewTimeGeekRestyClient(readCookies) //封装包含登录态的http client

		/**

		根据已有栏目下载

		*/

		loader.Run(l, "[ 正在加载已购买专栏列表... ]", func() {
			c, err := geektime.GetColumnList(client)
			if err != nil {
				printErrAndExit(err)
			}
			columns = c
		})

		selectColumn(client)
	},
}

// selectDiyCmd 手工选择课程下载
var selectDiyCmd = &cobra.Command{
	Use:   "diy",
	Short: "Geektime-downloader is used to download geek time lessons diy",
	Run: func(cmd *cobra.Command, args []string) {
		readCookies, err := util.ReadCookieFromConfigFile(phone) //获取登录态
		if err != nil {
			printErrAndExit(err)
		}
		//是否重新登录
		if reLogin {
			if err := util.RemoveConfig(phone); err != nil {
				printErrAndExit(err)
			} else {
				fmt.Println("清空登录态, 尝试重新登录")
			}
			readCookies = nil
		}
		if readCookies == nil { // 不存在 则登录
			pwd := prompt.GetPwd()
			loader.Run(l, "[ 正在登录... ]", func() {
				errMsg, cookies := geektime.Login(phone, pwd)
				if errMsg != "" {
					fmt.Fprintln(os.Stderr, errMsg)
					os.Exit(1)
				}
				readCookies = cookies
				err := util.WriteCookieToConfigFile(phone, cookies) //保存登录态
				if err != nil {
					printErrAndExit(err)
				}
			})
			fmt.Println("登录成功")
		}
		client := geektime.NewTimeGeekRestyClient(readCookies) //封装包含登录态的http client

		/**

		1、根据入参找到对应课程并检查权限
		2、批量下载课程

		*/
		loader.Run(l, "[ 正在检测课程是否有权限... ]", func() {
			c, err := geektime.GetColumnInfo(client, columnDiyId)
			if err != nil {
				printErrAndExit(err)
			}
			if c.CID == 0 { //仅检测是否存在 要保证有权限获取全部文章 不然下载的是试读
				err = errors.New("栏目不存在 或 请重新登录")
				printErrAndExit(err)
			}
			columns = append(columns, c)
		})

		selectColumn(client)
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
		handleDownloadAll(client)
	case 2:
		selectArticle(client)
	}
}

func selectArticle(client *resty.Client) {
	articles := loadArticles(client)
	index := prompt.SelectArticles(articles)
	handleSelectArticle(articles, index, client)
}

func handleSelectArticle(articles []geektime.ArticleSummary, index int, client *resty.Client) {
	if index == 0 {
		handleSelectColumn(client)
	}
	a := articles[index-1]
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
	})
	selectArticle(client)
}

func handleDownloadAll(client *resty.Client) {
	cTitle := columns[currentColumnIndex].Title
	articles := loadArticles(client)
	var counter uint64
	folder, err := mkColumnDownloadFolder(phone, cTitle)
	if err != nil {
		printErrAndExit(err)
	}
	if concurrency > 0 {
		wp := workerpool.New(concurrency)
		for _, a := range articles {
			aid := a.AID
			title := a.Title
			wp.Submit(func() {
				prefix := fmt.Sprintf("[ 正在下载专栏 《%s》 中的所有文章, 已完成下载%d/%d ... ]", cTitle, counter, len(articles))
				loader.Run(l, prefix, func() {
					err := chromedp.PrintArticlePageToPDF(aid, filepath.Join(folder, util.FileName(title, "pdf")), client.Cookies)
					if err != nil {
						printErrAndExit(err)
					} else {
						atomic.AddUint64(&counter, 1)
					}
				})
			})
		}
		wp.StopWait()
	} else {
		//单个处理 并sleep
		for _, a := range articles {
			aid := a.AID
			title := a.Title
			prefix := fmt.Sprintf("[ 正在下载专栏 《%s》 中的所有文章, 已完成下载%d/%d ... ]", cTitle, counter, len(articles))
			loader.Run(l, prefix, func() {
				err := chromedp.PrintArticlePageToPDF(aid, filepath.Join(folder, util.FileName(title, "pdf")), client.Cookies)
				if err != nil {
					printErrAndExit(err)
				} else {
					atomic.AddUint64(&counter, 1)
				}
			})
			util.SleepMS(sleep)
		}
	}
	selectColumn(client)
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
	rootCmd.AddCommand(selectColumnCmd, selectDiyCmd)

	if err := rootCmd.Execute(); err != nil {
		printErrAndExit(err)
	}
}

func printErrAndExit(err error) { // 打印报错
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
