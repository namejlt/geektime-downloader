package geektime

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	pgt "github.com/namejlt/geektime-downloader/pconst"
	"net/http"
	"time"
)

type UserAuth struct {
	AppID int
	Uid   int
}

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

// CheckAuth 检测是否登录
func CheckAuth(client *resty.Client, columnId int) (data UserAuth, err error) {
	var result struct {
		Code int `json:"code"`
		Data struct {
			AppID int `json:"appid"`
			UID   int `json:"uid"`
		} `json:"data"`
	}
	if columnId == 0 {
		columnId = 100109401
	}
	client.SetHeader("Referer", fmt.Sprintf(pgt.GeekBang+"/column/intro/%d", columnId)) //课程页面
	resp, err := client.R().
		SetResult(&result).
		Get(fmt.Sprintf(pgt.GeekBangAccount+"/serv/v1/user/auth?t=%d", time.Now().UnixMicro()))
	if err != nil {
		fmt.Println("失去登录态", resp.Result())
		return
	}
	if result.Code == 0 {
		data = UserAuth{
			AppID: result.Data.AppID,
			Uid:   result.Data.UID,
		}
		return
	}
	return
}
