package cmd

import (
	"errors"
	"fmt"
	"github.com/namejlt/geektime-downloader/cmd/prompt"
	"github.com/namejlt/geektime-downloader/internal/geektime"
	"github.com/namejlt/geektime-downloader/internal/loader"
	"github.com/namejlt/geektime-downloader/internal/pkg/util"
	"github.com/spf13/cobra"
	"os"
)

// selectDiyCmd 手工选择课程下载
var selectDiyCmd = &cobra.Command{
	Use:   "diy",
	Short: "Geektime-downloader is used to download geek time lessons diy",
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

		/**

		1、根据入参找到对应课程并检查权限
		2、批量下载课程

		*/
		loader.Run(l, "[ 正在检测课程是否有效... ]", func() {
			c, err := geektime.GetColumnInfo(client, columnDiyId)
			if err != nil {
				printErrAndExit(err)
			}
			if c.CID == 0 { //仅检测是否存在 要保证有权限获取全部文章 不然下载的是试读
				err = errors.New("栏目不存在")
				printErrAndExit(err)
			}
			columns = append(columns, c)
		})

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

		//检测课程是否在db中
		check, err := checkCourse(columnDiyId)
		if err != nil {
			printErrAndExit(err)
		}
		if !check {
			printMsgAndExit(fmt.Sprintln("id", columnDiyId, "课程在DB不存在"))
		}

		selectColumn(client)
	},
}
