package geektime

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	pgt "github.com/namejlt/geektime-downloader/pconst"
	"net/http"
	"time"
)

const (
	// ProductPath ...
	ProductPath = "/serv/v3/learn/product"
	// ArticlesPath ...
	ArticlesPath = "/serv/v1/column/articles"
	// ArticleV1Path ...
	ArticleV1Path = "/serv/v1/article"
	// ColumnInfoV3Path ...
	ColumnInfoV3Path = "/serv/v3/column/info"
)

// UserAgent is Web browser User Agent
const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.92 Safari/537.36"

// NewTimeGeekRestyClient new Http Client with user auth cookies
func NewTimeGeekRestyClient(cookies []*http.Cookie) *resty.Client {
	return resty.New().
		SetTimeout(30*time.Second).
		SetHeader("User-Agent", UserAgent).
		SetHeader("Origin", pgt.GeekBang).
		SetBaseURL(pgt.GeekBang).
		SetCookies(cookies)
}

// ErrGeekTimeAPIBadCode ...
type ErrGeekTimeAPIBadCode struct {
	Path string
	Code int
	Msg  string
}

// Error implements error interface
func (e ErrGeekTimeAPIBadCode) Error() string {
	return fmt.Sprintf("请求极客时间接口 %s 失败, code %d, msg %s", e.Path, e.Code, e.Msg)
}

// Product ...
type Product struct {
	Access   bool
	ID       int
	Title    string
	Type     string
	Articles []Article
}

// Article ...
type Article struct {
	AID   int
	Title string
}

// VideoInfo ...
type VideoInfo struct {
	M3U8URL string
	Size    int
}

// ArticleInfo ...
type ArticleInfo struct {
	ArticleContent   string
	AudioDownloadURL string
}
