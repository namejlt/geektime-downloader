package db

type Course struct {
	Id         uint32 `json:"id"`          // id
	CourseId   uint32 `json:"course_id"`   // 课程id
	Name       string `json:"name"`        // 名称
	Author     string `json:"author"`      // 作者
	IsSync     uint32 `json:"is_sync"`     // 是否同步 0 否 1 是
	IsFinish   uint32 `json:"is_finish"`   // 是否完成 0 未完成 1 已完成
	IsVideo    uint32 `json:"is_video"`    // 是否视频 0 否 1 是
	IsDownload uint32 `json:"is_download"` // 是否下载 0 否 1 是
	CreatedAt  int64  `json:"created_at"`  // 创建时间
	UpdatedAt  int64  `json:"updated_at"`  // 更新时间
}
