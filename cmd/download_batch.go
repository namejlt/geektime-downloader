package cmd

import (
	"errors"
	"fmt"
	"github.com/namejlt/geektime-downloader/cmd/prompt"
	"github.com/namejlt/geektime-downloader/dao"
	"github.com/namejlt/geektime-downloader/internal/geektime"
	"github.com/namejlt/geektime-downloader/internal/loader"
	"github.com/namejlt/geektime-downloader/internal/pkg/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

/**

通过外置文件批量下载

通过本地数据库批量下载，同时更新数据库内容


*/

// batchDownCmd 批量下载
var batchDownCmd = &cobra.Command{
	Use:   "batch",
	Short: "Geektime-downloader is used to download geek time lessons batch download by column ids",
	Run: func(cmd *cobra.Command, args []string) {
		checkArgs()

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

		loader.Run(l, "[ 正在检测是否登录... ]", func() {
			c, err := geektime.CheckAuth(client, 0)
			if err != nil {
				printErrAndExit(err)
			}
			if c.Uid == 0 {
				err = errors.New("登录态失效，请重新登录")
				printErrAndExit(err)
			}
		})

		/**

		1、根据入参找到对应课程id的文件
		2、通过一个cid来检查权限
		3、自动批量下载课程

		*/
		cidData, err := ioutil.ReadFile(columnIDsFile)
		if err != nil {
			printErrAndExit(err)
		}
		columnIdS := strings.Split(string(cidData), util.GetOsLineSep())
		if len(columnIdS) == 0 {
			printMsgAndExit("课程id为空")
		}

		columnDiyId, _ = strconv.Atoi(columnIdS[0])
		loader.Run(l, "[ 正在检测课程是否有效... ]", func() {
			c, err := geektime.GetColumnInfo(client, columnDiyId)
			if err != nil {
				printErrAndExit(err)
			}
			if c.CID == 0 { //仅检测是否存在 要保证有权限获取全部文章 不然下载的是试读
				err = errors.New("栏目不存在")
				printErrAndExit(err)
			}
		})

		//开始循环批量获取column 循环下载文章
		currentColumnIndex = 0
		for _, v := range columnIdS {
			columnDiyId, _ = strconv.Atoi(v)
			loader.Run(l, "[ 正在检测是否登录... ]", func() {
				c, err := geektime.CheckAuth(client, columnDiyId)
				if err != nil {
					printErrAndExit(err)
				}
				if c.Uid == 0 {
					err = errors.New("登录态失效，请重新登录")
					printErrAndExit(err)
				}
			})
			//获取课程
			column, err := geektime.GetColumnInfo(client, columnDiyId)
			if err != nil {
				printErrAndExit(err)
			}
			if column.CID == 0 {
				fmt.Println("id", v, "课程不存在")
				continue
			}
			//检测课程是否在db中
			check, err := checkCourse(column.CID)
			if err != nil {
				printErrAndExit(err)
			}
			if !check {
				fmt.Println("id", v, "课程在DB不存在")
				continue
			}
			if len(columns) == 0 {
				columns = append(columns, column)
			} else {
				columns[0] = column
			}

			//download
			handleDownloadAll(client, false, []geektime.ArticleSummary{}, 0)
		}
	},
}

func isColumn(columnType string) bool {
	return columnType == "c1"
}

func isVideo(columnType string) bool {
	return columnType == "c3"
}

func checkCourse(columnId int) (b bool, err error) {
	info, err := dao.NewDb(dbPath).GetCourse(uint64(columnId))
	if err != nil {
		return
	}
	if info.Id != 0 {
		b = true
	}
	return
}
