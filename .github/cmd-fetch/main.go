package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	"github.com/chyroc/gorequests"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

func main() {
	date := time.Now().In(ZoneGMT).Format("2006-01-02")
	jsonFilename := fmt.Sprintf("./json/%s.json", date)
	posts := getPostList()

	spew.Dump(posts)

	assert(os.MkdirAll("./json", 0o777))

	assert(ioutil.WriteFile(jsonFilename, []byte(posts.FormatJSON()), 0o666))
	assert(ioutil.WriteFile("README.md", []byte(posts.README()), 0o666))
}

type Post struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Count int    `json:"count"`
}

type PostList []*Post

func (r PostList) Len() int {
	return len(r)
}

func (r PostList) Less(i, j int) bool {
	return r[i].Count < r[j].Count
}

func (r PostList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r PostList) FormatMD() string {
	date := time.Now().In(ZoneGMT).Format("2006-01-02")
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("# V2EX 每日N贴 %s\n\n", date))
	for _, v := range r {
		s.WriteString(fmt.Sprintf("- [%s](%s) %d\n", v.Title, v.URL, v.Count))
	}
	s.WriteString("\n\n")
	return s.String()
}

func (r PostList) FormatJSON() string {
	bs, _ := json.MarshalIndent(r, "", "  ")
	return string(bs)
}

func (r PostList) README() string {
	return `# 北邮人论坛十大热帖

> 存储北邮人论坛十大热帖历史数据

## 项目地址

[项目 GitHub 地址](https://github.com/chyroc/byr-topten)

## 网页使用

访问 [主页](https://chyroc.github.io/byr-topten/) 即可访问所有数据

## 接口使用

以日期为例:

- HTML 数据： https://chyroc.github.io/byr-topten/2021-01-01.html
- JSON 数据： https://chyroc.github.io/byr-topten/2021-01-01.json

`
}

func (r PostList) TopTen() PostList {
	count := 0
	res := []*Post{}
	sort.Sort(sort.Reverse(r))
	for idx, v := range r {
		if idx < 10 {
			res = append(res, v)
			count = v.Count // count 存储最小的 count
		} else {
			if v.Count >= count {
				res = append(res, v)
			}
		}
	}
	return res
}

func getPostList() PostList {
	logrus.SetOutput(io.Discard)

	// 获取密码
	username := os.Getenv("BYR_USERNAME")
	password := os.Getenv("BYR_PASSWORD")

	// 登录
	cookies := ""
	{
		res := gorequests.New(http.MethodPost, "https://bbs.byr.cn/user/ajax_login.json").WithHeaders(map[string]string{"x-requested-with": "XMLHttpRequest", "content-type": "application/x-www-form-urlencoded"}).WithBody(fmt.Sprintf("id=%s&passwd=%s", username, password))
		resp, err := res.Response()
		assert(err)
		text, err := res.Text()
		assert(err)
		text = mahonia.NewDecoder("gbk").ConvertString(string(text))
		if !strings.Contains(text, username) {
			panic("登录失败")
		}
		for k, v := range resp.Header {
			if k == "Set-Cookie" {
				for _, vv := range v {
					cookies += ";" + strings.Split(vv, ";")[0]
				}
			}
		}
	}

	text, err := gorequests.New(http.MethodGet, "https://bbs.byr.cn/default").WithHeaders(map[string]string{"x-requested-with": "XMLHttpRequest", "cookie": cookies}).Text()
	assert(err)
	bodystr := mahonia.NewDecoder("gbk").ConvertString(string(text))
	if strings.Contains(bodystr, "您未登录") {
		panic("您未登录")
	}
	bodystr = strings.Split(bodystr, "近期热门话题")[1]
	bodystr = strings.Split(bodystr, "热门话题")[0]
	bodystr = strings.Join(strings.Split(bodystr, "<li"), "\n")

	posts := []*Post{}
	for _, v := range regPost.FindAllStringSubmatch(bodystr, -1) {
		if len(v) != 5 {
			continue
		}
		path := strings.TrimSpace(v[1])
		title := strings.TrimSpace(v[2])
		count, _ := strconv.ParseInt(strings.TrimSpace(v[4]), 10, 64)
		url := "https://bbs.byr.cn" + path

		if !strings.HasPrefix(path, "/article/") {
			continue
		}

		posts = append(posts, &Post{
			Title: title,
			URL:   url,
			Count: int(count),
		})
	}

	sort.Sort(sort.Reverse(PostList(posts)))
	return posts
}

var (
	ZoneGMT = time.FixedZone("GMT+8", 8*60*60)
	regPost = regexp.MustCompile(`(?m)<a href="(.*?)">(.*?)(\(<span.*?>(\d*)<\/span>\))?<\/a>`)
)

func assert(err error) {
	if err != nil {
		panic(err)
	}
}
