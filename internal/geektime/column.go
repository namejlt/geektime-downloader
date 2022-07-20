package geektime

import (
	"errors"
	"fmt"
	pgt "github.com/namejlt/geektime-downloader/pconst"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
)

// ColumnSummary Mini column struct
type ColumnSummary struct {
	CID          int //课程id
	Title        string
	Desc         string
	AuthorName   string
	Type         string
	IsVideo      bool
	IsAudio      bool
	IsFinish     bool
	ArticleCount uint32
	Articles     []ArticleSummary
}

type ColumnLabel struct {
	CourseId   uint64
	CourseType uint8
}

func GetColumnLabels(client *resty.Client) (data []ColumnLabel, err error) {
	var result struct {
		Code int `json:"code"`
		Data struct {
			List []struct {
				ColumnSku  uint64 `json:"column_sku"`
				ColumnType uint8  `json:"column_type"`
			} `json:"list"`
		} `json:"data"`
		Error interface{} `json:"error"`
	}
	result.Code = -1

	client.SetHeader("Referer", pgt.GeekBang+"/resource?plus=0&m=0&d=0&c=0&sort=0&order=sort")
	_, err = client.R().
		SetBody(
			map[string]interface{}{
				"label_id": 0,
				"type":     0,
			}).
		SetResult(&result).
		Post("/serv/v1/column/label_skus") //获取专栏列表

	if err != nil {
		return nil, err
	}

	if result.Code == 0 {
		for _, v := range result.Data.List {
			data = append(data, ColumnLabel{
				CourseId:   v.ColumnSku,
				CourseType: v.ColumnType,
			})
		}
		return data, nil
	}
	return nil, errors.New("call geektime product api failed")
}

// GetColumnList call geektime api to get column list
func GetColumnList(client *resty.Client) ([]ColumnSummary, error) {
	var result struct {
		Code int `json:"code"`
		Data struct {
			Products []struct {
				ID     int    `json:"id"`
				Title  string `json:"title"`
				Author struct {
					Name string `json:"name"`
				} `json:"author"`
			} `json:"products"`
		} `json:"data"`
		Error struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		} `json:"error"`
	}

	client.SetHeader("Referer", pgt.GeekBang+"/dashboard/course") //我的课程页面
	_, err := client.R().
		SetBody(
			map[string]interface{}{
				"desc":             false,
				"expire":           1,
				"last_learn":       0,
				"learn_status":     0,
				"prev":             0,
				"size":             200,
				"sort":             1,
				"type":             "c1",
				"with_learn_count": 1,
			}).
		SetResult(&result).
		Post("/serv/v3/learn/product") //获取专栏列表

	if err != nil {
		return nil, err
	}

	if result.Code == 0 {
		var products []ColumnSummary
		for _, v := range result.Data.Products {
			products = append(products, ColumnSummary{
				CID:        v.ID,
				Title:      replaceSep(v.Title),
				AuthorName: v.Author.Name,
			})
		}
		return products, nil
	}
	return nil, errors.New("call geektime product api failed")
}

// GetColumnInfo call geektime api to get column info
func GetColumnInfo(client *resty.Client, columnId int) (data ColumnSummary, err error) {
	var result struct {
		Code int `json:"code"`
		Data struct {
			Author struct {
				Name string `json:"name"`
			} `json:"author"`

			Article struct {
				Count uint32 `json:"count"`
			} `json:"article"`

			Path struct {
				Desc string `json:"desc"`
			} `json:"path"`

			ID    int    `json:"id"`
			Title string `json:"title"`

			IsVideo  bool   `json:"is_video"`  // 是否视频
			IsAudio  bool   `json:"is_audio"`  // 是否音频
			IsFinish bool   `json:"is_finish"` // 是否完成
			Type     string `json:"type"`
		} `json:"data"`
		Error struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		} `json:"error"`
	}
	result.Code = -1

	client.SetHeader("Referer", fmt.Sprintf(pgt.GeekBang+"/column/intro/%d", columnId)) //课程页面
	_, err = client.R().
		SetBody(
			map[string]interface{}{
				"with_recommend_article": false,
				"product_id":             columnId,
			}).
		SetResult(&result).
		Post("/serv/v3/column/info") //获取专栏信息

	if err != nil {
		return
	}

	if result.Code == 0 {
		data = ColumnSummary{
			CID:          result.Data.ID,
			Title:        replaceSep(result.Data.Title),
			Desc:         result.Data.Path.Desc,
			AuthorName:   result.Data.Author.Name,
			IsVideo:      result.Data.IsVideo,
			IsAudio:      result.Data.IsAudio,
			IsFinish:     result.Data.IsFinish,
			ArticleCount: result.Data.Article.Count,
			Type:         result.Data.Type,
		}
		return
	}
	return
}

func replaceSep(str string) string {
	str = strings.ReplaceAll(str, string(filepath.Separator), "")
	return strings.TrimRight(str, " ")
}
