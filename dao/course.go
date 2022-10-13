package dao

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/namejlt/geektime-downloader/model"
	"log"
)

// SaveCourse 保存课程到本地 增加或更新
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
			"is_finish,is_audio,is_video,created_at,updated_at) "+
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

func (p *Geek) SaveCourseDownloadRecord(addData model.CourseDownloadRecord) (err error) {
	db, err := p.getDb()
	if err != nil {
		return
	}
	defer db.Close()

	var id uint32
	err = db.QueryRow("SELECT id FROM course_download_record WHERE course_id=? and article_id=? and download_type=?",
		addData.CourseId, addData.ArticleId, addData.DownloadType).Scan(&id)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	if id == 0 {
		//新增
		_, err = db.Exec("INSERT INTO course_download_record(course_id, article_id, article_name, download_type,created_at,updated_at) "+
			"VALUES( ?, ?, ?, ?, ?, ? )",
			addData.CourseId,
			addData.ArticleId,
			addData.ArticleName,
			addData.DownloadType,
			addData.CreatedAt,
			addData.UpdatedAt,
		)
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}

func (p *Geek) GetCourseDownloadRecord(courseId uint64, articleId int, downloadType int) (data model.CourseDownloadRecord, err error) {
	db, err := p.getDb()
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow("SELECT id,course_id,article_id,article_name,download_type  FROM course_download_record WHERE course_id=? and article_id=? and download_type=?",
		courseId, articleId, downloadType).Scan(&data.Id, &data.CourseId, &data.ArticleId, &data.ArticleName, &data.DownloadType)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}

func (p *Geek) GetCourse(courseId uint64) (data model.Course, err error) {
	db, err := p.getDb()
	if err != nil {
		return
	}
	defer db.Close()
	err = db.QueryRow("SELECT id, course_id,course_name,course_type,article_count,is_finish,is_audio,is_video FROM course WHERE course_id=?", courseId).
		Scan(&data.Id, &data.CourseId, &data.CourseName, &data.CourseType, &data.ArticleCount, &data.IsFinish, &data.IsAudio, &data.IsVideo)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}
