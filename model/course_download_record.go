package model

type CourseDownloadRecord struct {
	Id           uint32 `json:"id"`            // id
	CourseId     uint64 `json:"course_id"`     // 课程id
	ArticleName  string `json:"article_name"`  // 文章名称
	DownloadType uint8  `json:"download_type"` // 下载类型 1 pdf 2 md 3 audio 4 video
	CreatedAt    int64  `json:"created_at"`    // 创建时间
	UpdatedAt    int64  `json:"updated_at"`    // 更新时间
}
