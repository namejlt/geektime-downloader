package geektime

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// ColumnSummary Mini column struct
type ColumnSummary struct {
	CID        int //课程id
	Title      string
	AuthorName string
	Articles   []ArticleSummary
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

	client.SetHeader("Referer", "https://time.geekbang.org/dashboard/course") //我的课程页面
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
				Title:      v.Title,
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

			ID    int    `json:"id"`
			Title string `json:"title"`
		} `json:"data"`
		Error struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		} `json:"error"`
	}

	client.SetHeader("Referer", fmt.Sprintf("https://time.geekbang.org/column/intro/%d", columnId)) //课程页面
	_, err = client.R().
		SetBody(
			map[string]interface{}{
				"with_recommend_article": true,
				"product_id":             columnId,
			}).
		SetResult(&result).
		Post("/serv/v3/column/info") //获取专栏信息

	if err != nil {
		return
	}

	if result.Code == 0 {
		data = ColumnSummary{
			CID:        result.Data.ID,
			Title:      result.Data.Title,
			AuthorName: result.Data.Author.Name,
		}
		return
	}
	return
}
