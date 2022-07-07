package cmd

import (
	"fmt"
	"github.com/namejlt/geektime-downloader/cmd/prompt"
	"github.com/namejlt/geektime-downloader/dao"
	"github.com/namejlt/geektime-downloader/internal/geektime"
	"github.com/namejlt/geektime-downloader/internal/loader"
	"github.com/namejlt/geektime-downloader/internal/pkg/util"
	"github.com/namejlt/geektime-downloader/model"
	"github.com/namejlt/geektime-downloader/pconst"
	"github.com/spf13/cobra"
	"os"
	"time"
)

/**

从网站同步课程到本地数据库

从本地目录课程状态同步到本地数据库

手工修改本地数据库数据


*/

//syncInitDb 初始化db
var syncInitDb = &cobra.Command{
	Use:   "syncinitdb",
	Short: "Geektime-downloader syncinitdb",
	Run: func(cmd *cobra.Command, args []string) {
		dao.InitDb(dbPath, dbMust)
	},
}

//syncAllColumnToLocal 同步线上课程到本地
var syncAllColumnToLocal = &cobra.Command{
	Use:   "syncgeek2local",
	Short: "Geektime-downloader is used to sync geek lessons to local",
	Run: func(cmd *cobra.Command, args []string) {
		/**

		检查本地数据库

		拉取所有课程id，保存本地

		更新课程信息，不影响课程下载标记

		*/

		exist := dao.CheckDb(dbPath)
		if !exist {
			printMsgAndExit("db不存在，请初始化")
		}

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
		client := geektime.NewTimeGeekRestyClient(readCookies) //封装包含登录态的http client 未登录也可

		//获取所有课程
		columnList, err := geektime.GetColumnLabels(client)
		if err != nil {
			printErrAndExit(err)
		}
		//循环获取课程信息并写入本地
		for _, label := range columnList {
			info, err := geektime.GetColumnInfo(client, int(label.CourseId))
			if err != nil {
				printErrAndExit(err)
			}
			fmt.Println("写入", info.CID, info.Title)
			addData := model.Course{}
			addData.CourseId = label.CourseId
			addData.CourseName = info.Title
			addData.CourseDesc = info.Desc
			addData.Author = info.AuthorName
			addData.CourseType = label.CourseType
			addData.ArticleCount = info.ArticleCount
			addData.IsFinish = getISFlag(info.IsFinish)
			addData.IsAudio = getISFlag(info.IsAudio)
			addData.IsVideo = getISFlag(info.IsVideo)
			now := time.Now().Unix()
			addData.CreatedAt = now
			addData.UpdatedAt = now
			err = dao.NewDb(dbPath).SaveCourse(addData)
			if err != nil {
				printErrAndExit(err)
			}
			util.SleepMS(sleep, sleep)
		}

		fmt.Println("写入完成")
	},
}

//syncLocalColumnToDb 同步本地存储到数据库
var syncLocalColumnToDb = &cobra.Command{
	Use:   "synclocal2db",
	Short: "Geektime-downloader synclocal2db",
	Run: func(cmd *cobra.Command, args []string) {
		/**
		todo
				遍历本地存储目录
				根据课程名称匹配找到课程id
				判断本地存储课程文章文件更新下载明细记录
		*/

	},
}

func getISFlag(b bool) uint8 {
	if b {
		return pconst.CommonTrue
	}
	return pconst.CommonFalse
}
