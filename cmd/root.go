package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/namejlt/geektime-downloader/dao"
	audiodown "github.com/namejlt/geektime-downloader/internal/pkg/audio"
	"github.com/namejlt/geektime-downloader/internal/pkg/markdown"
	videodown "github.com/namejlt/geektime-downloader/internal/pkg/video"
	"github.com/namejlt/geektime-downloader/model"
	"github.com/namejlt/geektime-downloader/pconst"
	"os"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/namejlt/geektime-downloader/cmd/prompt"
	"github.com/namejlt/geektime-downloader/internal/geektime"
	"github.com/namejlt/geektime-downloader/internal/loader"
	"github.com/namejlt/geektime-downloader/internal/pkg/chromedp"
	"github.com/namejlt/geektime-downloader/internal/pkg/util"
	"github.com/spf13/cobra"
)

// rootCmd 执行命令
var rootCmd = &cobra.Command{
	Use:   "geektime-downloader",
	Short: "Geektime-downloader is used to download geek time lessons",
	Run: func(cmd *cobra.Command, args []string) {
		//不执行内容 具体执行在子命令
	},
}

// Execute func 执行函数
func Execute() {
	rootCmd.AddCommand(
		selectColumnCmd,
		selectDiyCmd,
		batchDownCmd,
		syncInitDb,
		syncAllColumnToLocal,
	)

	if err := rootCmd.Execute(); err != nil {
		printErrAndExit(err)
	}
}

// =============================================================================== download

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

	//检测课程db是否存在并初始化

	//检测下载类型
	switch columnType {
	case 1:
	case 2: //仅视频
		if isColumn {
			fmt.Println("图文课程类型跳过")
			return
		}
	case 3: //仅图文
		if isVideo {
			fmt.Println("视频课程类型跳过")
			return
		}
	default:

	}

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
			//检测是否已下载
			check, err := checkCourseDownload(columns[currentColumnIndex].CID, aid, pconst.DownloadTypePdf)
			checkError(err)
			if !check {
				loader.Run(l, prefix, func() {
					err := chromedp.PrintArticlePageToPDF(aid, filepath.Join(folder, util.FileName(title, "pdf")), client.Cookies)
					if err != nil {
						printErrAndExit(err)
					} else {
						atomic.AddUint64(&counter, 1)
					}
				})
				err = saveCourseDownload(columns[currentColumnIndex].CID, aid, cTitle, pconst.DownloadTypePdf)
				checkError(err)
			}
		}
		if md || audio {
			if isColumn {
				articleInfo, err := geektime.GetArticleInfo(a.AID, client)
				checkError(err)
				if md {
					check, err := checkCourseDownload(columns[currentColumnIndex].CID, aid, pconst.DownloadTypeMD)
					checkError(err)
					if !check {
						err = markdown.Download(ctx, articleInfo.ArticleContent, a.Title, folder, a.AID, 1)
						checkError(err)
						err = saveCourseDownload(columns[currentColumnIndex].CID, aid, cTitle, pconst.DownloadTypeMD)
						checkError(err)
					}
				}
				if audio {
					check, err := checkCourseDownload(columns[currentColumnIndex].CID, aid, pconst.DownloadTypeAudio)
					checkError(err)
					if !check {
						err = audiodown.DownloadAudio(ctx, articleInfo.AudioDownloadURL, folder, a.Title)
						checkError(err)
						err = saveCourseDownload(columns[currentColumnIndex].CID, aid, cTitle, pconst.DownloadTypeAudio)
						checkError(err)
					}
				}
			} else {
				println(cTitle, a.Title, "非课程")
			}
		}

		if video {
			if isVideo {
				check, err := checkCourseDownload(columns[currentColumnIndex].CID, aid, pconst.DownloadTypeVideoMap[quality])
				checkError(err)
				if !check {
					videoInfo, err := geektime.GetVideoInfo(a.AID, client, quality)
					checkError(err)
					err = videodown.DownloadVideo(ctx, videoInfo.M3U8URL, a.Title+quality, folder, int64(videoInfo.Size), 1)
					checkError(err)
					err = saveCourseDownload(columns[currentColumnIndex].CID, aid, cTitle, pconst.DownloadTypeVideoMap[quality]) //视频记录一次
					checkError(err)
				}
			} else {
				println(cTitle, a.Title, "无视频")
			}
		}
		//更新db

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

func printErrAndExit(err error) { // 打印报错
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func printMsgAndExit(msg string) { // 打印信息
	fmt.Fprintln(os.Stdout, msg)
	os.Exit(1)
}

func checkCourseDownload(columnId int, articleId int, downloadType int) (b bool, err error) {
	info, err := dao.NewDb(dbPath).GetCourseDownloadRecord(uint64(columnId), articleId, downloadType)
	if err != nil {
		return
	}
	if info.Id != 0 {
		b = true
	}
	return
}

func saveCourseDownload(columnId int, articleId int, articleName string, downloadType int) (err error) {
	addData := model.CourseDownloadRecord{}
	addData.CourseId = uint64(columnId)
	addData.ArticleId = uint64(articleId)
	addData.ArticleName = articleName
	addData.DownloadType = uint8(downloadType)
	now := time.Now().Unix()
	addData.CreatedAt = now
	addData.UpdatedAt = now
	err = dao.NewDb(dbPath).SaveCourseDownloadRecord(addData)
	return
}

func checkArgs() {
	checkQuality(quality)
}

func checkQuality(quality string) {
	if _, ok := pconst.DownloadTypeVideoMap[quality]; !ok {
		printMsgAndExit("quality is error")
	}
}
