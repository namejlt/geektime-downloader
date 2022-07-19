package model

type Course struct {
	Id           uint32 `json:"id"`            // id
	CourseId     uint64 `json:"course_id"`     // 课程id
	CourseName   string `json:"course_name"`   // 名称
	CourseDesc   string `json:"course_desc"`   // 描述
	Author       string `json:"author"`        // 作者
	CourseType   uint8  `json:"course_type"`   // 类型 1 图文 3 视频
	ArticleCount uint32 `json:"article_count"` // 文章数目
	IsFinish     uint8  `json:"is_finish"`     // 是否完成 0 未完成 1 已完成
	IsAudio      uint8  `json:"is_audio"`      // 是否音频 0 否 1 是
	IsVideo      uint8  `json:"is_video"`      // 是否视频 0 否 1 是
	CreatedAt    int64  `json:"created_at"`    // 创建时间
	UpdatedAt    int64  `json:"updated_at"`    // 更新时间
}
