package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/namejlt/geektime-downloader/pconst"
	"log"
	"os"
)

type Geek struct {
	Path string
}

func NewDb(db string) *Geek {
	p := Geek{}
	p.Path = db
	return &p
}

func (p *Geek) getDb() (db *sql.DB, err error) {
	if p.Path != "" {
		p.Path = pconst.DbPath
	}
	return sql.Open("sqlite3", p.Path)
}

func CheckDb(dbPath string) (exist bool) {
	exist, err := pathExists(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// InitDb 初始化数据库和表
func InitDb(dbPath string, must bool) {
	if must {
		os.Remove(dbPath)
	} else {
		exist, err := pathExists(dbPath)
		if err != nil {
			log.Fatal(err)
		}
		if exist {
			fmt.Println("已存在db")
			return
		}
	}

	db, err := sql.Open("sqlite3", dbPath) //不存在 则自动创建
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//课程表
	sqlStmt := `
	create table course (
	    id integer not null primary key, 
	    course_id INTEGER not null,
	    course_name TEXT not null,
	    course_desc TEXT not null,
	    author TEXT not null,
	    course_type INTEGER not null,
	    article_count INTEGER not null,
	    is_finish INTEGER not null,
	    is_audio INTEGER not null,
	    is_video INTEGER not null,
	    created_at INTEGER not null,
	    updated_at INTEGER not null
	                    );
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	//课程下载明细表
	sqlStmt = `
	create table course_download_record (
	    id integer not null primary key, 
	    course_id INTEGER not null,
	    article_id INTEGER not null,
	    article_name TEXT not null,
	    download_type INTEGER not null,
	    created_at INTEGER not null,
	    updated_at INTEGER not null
	                    );
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	fmt.Println("初始化成功")
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
