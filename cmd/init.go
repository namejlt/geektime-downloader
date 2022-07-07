package cmd

import (
	"github.com/briandowns/spinner"
	"github.com/namejlt/geektime-downloader/internal/geektime"
	"github.com/namejlt/geektime-downloader/internal/loader"
	"github.com/namejlt/geektime-downloader/internal/pkg/util"
	"os"
	"path/filepath"
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
	columnType int
	pdf        bool
	md         bool
	audio      bool
	video      bool

	//db
	dbPath string
	dbMust bool
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
	rootCmd.PersistentFlags().IntVarP(&columnType, "columntype", "t", 1, "课程类型 1 全部 2 视频 3 图文")

	selectDiyCmd.Flags().IntVarP(&columnDiyId, "column_diy_id", "i", 0, "指定下载课程id")
	_ = selectDiyCmd.MarkFlagRequired("column_diy_id")
	selectDiyCmd.Flags().IntVarP(&sleep, "sleep", "s", 1000, "下载文章间隔时间 毫秒")
	selectDiyCmd.Flags().BoolVarP(&reLogin, "relogin", "r", false, "是否重新登录")

	batchDownCmd.Flags().StringVarP(&columnIDsFile, "column_ids_file", "i", "./doc/cid.txt", "指定下载课程id文件")
	_ = batchDownCmd.MarkFlagRequired("column_ids_file")
	batchDownCmd.Flags().IntVarP(&sleep, "sleep", "s", 1000, "下载文章间隔时间 毫秒")
	batchDownCmd.Flags().IntVarP(&sleepMax, "sleepmax", "m", 1200, "下载文章间隔时间 毫秒 max")
	batchDownCmd.Flags().BoolVarP(&reLogin, "relogin", "r", false, "是否重新登录")

	syncInitDb.Flags().StringVarP(&dbPath, "dbpath", "o", "./db/geek.db", "db路径")
	syncInitDb.Flags().BoolVarP(&dbMust, "dbmust", "m", false, "db强制初始化")
	_ = syncInitDb.MarkFlagRequired("dbpath")

	syncAllColumnToLocal.Flags().StringVarP(&dbPath, "dbpath", "o", "./db/geek.db", "db路径，空则取默认值")
	syncAllColumnToLocal.Flags().IntVarP(&sleep, "sleep", "s", 1000, "下载文章间隔时间 毫秒")

	l = loader.NewSpinner()
}
