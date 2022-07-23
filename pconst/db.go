package pconst

const (
	DbPath = "./db/geek.db"
)

const (
	CommonTrue  = 1
	CommonFalse = 0

	//下载类型 下载类型 1 pdf 2 md 3 audio 4 video 标清 5 video 高清 6 video 超清
	// ld标清,sd高清,hd超清
	DownloadTypePdf     = 1
	DownloadTypeMD      = 2
	DownloadTypeAudio   = 3
	DownloadTypeVideoLD = 4
	DownloadTypeVideoSD = 5
	DownloadTypeVideoHD = 6

	DownloadTypeVideoLDTag = "ld"
	DownloadTypeVideoSDTag = "sd"
	DownloadTypeVideoHDTag = "hd"
)

var (
	DownloadTypeVideoMap = map[string]int{
		DownloadTypeVideoLDTag: DownloadTypeVideoLD,
		DownloadTypeVideoSDTag: DownloadTypeVideoSD,
		DownloadTypeVideoHDTag: DownloadTypeVideoHD,
	}
)
