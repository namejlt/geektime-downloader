package test

import (
	"github.com/namejlt/geektime-downloader/dao"
	"testing"
)

func TestInitDb(t *testing.T) {
	dao.InitDb("./test.db")
}
