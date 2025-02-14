# geektime-downloader

geektime-downloader 目前支持下载指定课程为PDF文档、markdown、音频、视频

[![go report card](https://goreportcard.com/badge/github.com/namejlt/geektime-downloader "go report card")](https://goreportcard.com/report/github.com/namejlt/geektime-downloader)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

## Usage

### Prerequisites

Chrome installed
sqlite installed - sync geektime column id to local

### Install form source

```bash
# Go 1.17+
go install github.com/namejlt/geektime-downloader@latest

or

git clone github.com/namejlt/geektime-downloader
make build
make run xxx

```

### Download binary files

See release page

### Sample

```bash
## windows 为例

下载我的栏目课程

> geektime-downloader.exe columns -u "phone number"


下载指定栏目课程

> geektime-downloader.exe diy -u "phone number" -i "column id"

> geektime-downloader.exe diy -u 13800000000 -i 100056701
```

### Help

```bash
## windows 为例
> geektime-downloader.exe -h

Geektime-downloader is used to download geek time lessons

Usage:
geektime-downloader [flags]
geektime-downloader [command]

Available Commands:
columns     Geektime-downloader is used to download geek time lessons
completion  Generate the autocompletion script for the specified shell
diy         Geektime-downloader is used to download geek time lessons diy
help        Help about any command

Flags:
-h, --help   help for geektime-downloader


> geektime-downloader.exe columns -h

Geektime-downloader is used to download geek time lessons

Usage:
geektime-downloader columns [flags]

Flags:
-c, --concurrency int   下载文章的并发数 (default 4)
-f, --folder string     PDF 文件下载目标位置 (default "/Users/tynam/geektime-downloader")
-h, --help              help for columns
-u, --phone string      你的极客时间账号(手机号)(required)


> geektime-downloader.exe diy -h


Geektime-downloader is used to download geek time lessons diy

Usage:
geektime-downloader diy [flags]

Flags:
-i, --column_diy_id int   指定下载课程id
-c, --concurrency int     下载文章的并发数 0 代表不并发且有等待时间 (default 4)
-f, --folder string       PDF 文件下载目标位置 (default "/Users/tynam/geektime-downloader")
-h, --help                help for diy
-u, --phone string        你的极客时间账号(手机号)(required)
-r, --relogin             是否重新登录
-s, --sleep int           下载文章间隔时间 毫秒 (default 1000)



```

## Core logic

获取课程

通过接口全量获取课程id、课程类型、名称、

下载课程

pdf通过调用chrome生成整页
md拉取html转换，其中会下载image
audio直接下载音频地址
video涉及到aes加密和视频拆分，先获取key，然后拉取所有视频，解密拼接

## Note

1. PDF 文件下载目标位置可以通过 help 查看

2. 为了避免多次登录账户，在目录 [UserConfigDir](https://pkg.go.dev/os#UserConfigDir)/geektime-downloader
   下会存放用户的登录信息，如果不是在自己的电脑上执行，请在使用完毕程序后手动删除。

3. 登录后填入指定课程的开篇页面地址，自动解析下载，要保证有订阅权限

4. 下载完的pdf可能为空白，重试即可

5. 多次空白，请不要用并发下载，并发设置0，设置等待时间

## Inspired by

* [geektime-dl](https://github.com/mmzou/geektime-dl)
