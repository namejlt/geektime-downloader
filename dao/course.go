package dao

import "fmt"

func AddCourse(columnId int) (err error) {
	db, err := getDb()
	if err != nil {
		return
	}
	fmt.Println(db)
	return
}
