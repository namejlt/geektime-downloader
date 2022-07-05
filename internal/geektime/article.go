package geektime

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/namejlt/geektime-downloader/pconst"
	"strconv"
)

// ArticleSummary Mini article struct
type ArticleSummary struct {
	AID   int
	Title string
}

// GetArticles call geektime api to get article list
func GetArticles(cid string, client *resty.Client) ([]ArticleSummary, error) {
	var result struct {
		Code int `json:"code"`
		Data struct {
			List []struct {
				ID           int    `json:"id"`
				ArticleTitle string `json:"article_title"`
			} `json:"list"`
		} `json:"data"`
	}
	_, err := client.R().
		SetBody(
			map[string]interface{}{
				"cid":    cid,
				"order":  "earliest",
				"prev":   0,
				"sample": false,
				"size":   500,
			}).
		SetResult(&result).
		Post("/serv/v1/column/articles")

	if err != nil {
		return nil, err
	}

	if result.Code == 0 {
		var articles []ArticleSummary
		for _, v := range result.Data.List {
			articles = append(articles, ArticleSummary{
				AID:   v.ID,
				Title: v.ArticleTitle,
			})
		}
		return articles, nil
	}
	return nil, errors.New("call geektime articles api failed")
}

// ColumnResponse ...
type ColumnResponse struct {
	Code int `json:"code"`
	Data struct {
		ArticleTitle     string `json:"article_title"`
		ArticleContent   string `json:"article_content"`
		AudioDownloadURL string `json:"audio_download_url"`
	} `json:"data"`
}

// VideoResponse ...
type VideoResponse struct {
	Code int `json:"code"`
	Data struct {
		ArticleTitle string `json:"article_title"`
		HLSVideos    struct {
			SD struct {
				Size int    `json:"size"`
				URL  string `json:"url"`
			} `json:"sd"`
			HD struct {
				Size int    `json:"size"`
				URL  string `json:"url"`
			} `json:"hd"`
			LD struct {
				Size int    `json:"size"`
				URL  string `json:"url"`
			} `json:"ld"`
		} `json:"hls_videos"`
	} `json:"data"`
}

// GetArticleInfo call geektime api to get article info
func GetArticleInfo(aid int, client *resty.Client) (data ArticleInfo, err error) {
	var result = ColumnResponse{}
	resp, err := client.R().
		SetBody(
			map[string]interface{}{
				"id":                strconv.Itoa(aid),
				"include_neighbors": true,
				"is_freelyread":     true,
				"reverse":           false,
			}).
		SetResult(&result).
		Post(ArticleV1Path)

	if err != nil {
		return
	}
	if resp.RawResponse.StatusCode == 451 {
		return data, pconst.ErrGeekTimeRateLimit
	}

	if result.Code != 0 {
		err = ErrGeekTimeAPIBadCode{ArticleV1Path, result.Code, ""}
		return
	}
	return ArticleInfo{
		result.Data.ArticleContent,
		result.Data.AudioDownloadURL,
	}, err
}

func GetVideoInfo(aid int, client *resty.Client, quality string) (data VideoInfo, err error) {
	var result = VideoResponse{}
	resp, err := client.R().
		SetBody(
			map[string]interface{}{
				"id":                strconv.Itoa(aid),
				"include_neighbors": true,
				"is_freelyread":     true,
				"reverse":           false,
			}).
		SetResult(&result).
		Post(ArticleV1Path)

	if err != nil {
		return
	}
	if resp.RawResponse.StatusCode == 451 {
		return data, pconst.ErrGeekTimeRateLimit
	}

	if result.Code != 0 {
		err = ErrGeekTimeAPIBadCode{ArticleV1Path, result.Code, ""}
		return
	}
	if quality == "sd" {
		data = VideoInfo{
			M3U8URL: result.Data.HLSVideos.SD.URL,
			Size:    result.Data.HLSVideos.SD.Size,
		}
	} else if quality == "hd" {
		data = VideoInfo{
			M3U8URL: result.Data.HLSVideos.HD.URL,
			Size:    result.Data.HLSVideos.HD.Size,
		}
	} else if quality == "ld" {
		data = VideoInfo{
			M3U8URL: result.Data.HLSVideos.LD.URL,
			Size:    result.Data.HLSVideos.LD.Size,
		}
	}
	return
}
