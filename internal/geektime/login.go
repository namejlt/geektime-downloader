package geektime

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"

	pgt "github.com/namejlt/geektime-downloader/internal/pkg/geektime"
)

// Login call geektime login api and return auth cookies
func Login(phone, password string) (string, []*http.Cookie) {
	client := resty.New().
		SetTimeout(5*time.Second).
		SetHeader("User-Agent", UserAgent).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("Connection", "keep-alive").
		SetBaseURL(pgt.GeekBangAccount)

	var result struct {
		Code int `json:"code"`
		Data struct {
			UID  int    `json:"uid"`
			Name string `json:"nickname"`
		} `json:"data"`
		Error struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		} `json:"error"`
	}

	loginResponse, err := client.R().
		SetHeader("Referer", pgt.GeekBangAccount+"/signin?redirect=https%3A%2F%2Ftime.geekbang.org%2F").
		SetBody(
			map[string]interface{}{
				"country":   86,
				"appid":     1,
				"platform":  3,
				"cellphone": phone,
				"password":  password,
			}).
		SetResult(&result).
		Post("/account/ticket/login")

	if err != nil {
		return err.Error(), nil
	}

	if result.Code == 0 {
		var cookies []*http.Cookie
		for _, c := range loginResponse.Cookies() {
			if c.Name == "GCID" || c.Name == "GCESS" || c.Name == "SERVERID" {
				cookies = append(cookies, c)
			}
		}
		return "", cookies
	}
	return result.Error.Msg, nil
}
