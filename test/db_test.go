package test

import (
	"github.com/namejlt/geektime-downloader/dao"
	"testing"
)

func TestInitDb(t *testing.T) {
	dao.InitDb("./test.db")
}

func TestArr(t *testing.T) {
	a := []int{1, 2, 3, 4}
	a = a[2:3]
	t.Log(a)
}
