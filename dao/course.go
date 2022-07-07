package dao

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/namejlt/geektime-downloader/model"
	"github.com/namejlt/geektime-downloader/pconst"
	"log"
)

//SaveCourse 保存课程到本地 增加或更新
func (p *Geek) SaveCourse(addData model.Course) (err error) {
	db, err := p.getDb()
	if err != nil {
		return
	}
	defer db.Close()

	var id uint32
	err = db.QueryRow("SELECT id FROM course WHERE course_id=?", addData.CourseId).Scan(&id)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	if id == 0 {
		//新增
		_, err = db.Exec("INSERT INTO course(course_id, course_name, course_desc, author,course_type,article_count,"+
			"is_finish,is_audio,is_video,is_pdf_download,is_markdown_download,is_video_download,is_audio_download,created_at,updated_at) "+
			"VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )",
			addData.CourseId,
			addData.CourseName,
			addData.CourseDesc,
			addData.Author,
			addData.CourseType,
			addData.ArticleCount,
			addData.IsFinish,
			addData.IsAudio,
			addData.IsVideo,
			pconst.CommonFalse,
			pconst.CommonFalse,
			pconst.CommonFalse,
			pconst.CommonFalse,
			addData.CreatedAt,
			addData.UpdatedAt,
		)
	} else {
		//更新
		_, err = db.Exec(`update course set is_finish=?,updated_at=? where course_id=?`,
			addData.IsFinish,
			addData.UpdatedAt,
			addData.CourseId,
		)
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}
