package dao

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/namejlt/geektime-downloader/pconst"
	"log"
	"os"
)

func getDb() (db *sql.DB, err error) {
	return sql.Open("sqlite3", pconst.DbPath)
}

func InitDb(dbPath string) {
	// 初始化数据库和表
	os.Remove(dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//课程表
	sqlStmt := `
	create table course (
	    id integer not null primary key, 
	    course_id INTEGER not null,
	    name text not null,
	    author text not null,
	    is_sync INTEGER not null,
	    is_finish INTEGER not null,
	    is_video INTEGER not null,
	    is_download INTEGER not null,
	    created_at INTEGER not null
	    updated_at INTEGER not null
	                    );
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
