package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	"github.com/chyroc/goexec"
)

func main() {
	// 切换到 gh-pages 分支
	run1CheckoutBranch()

	// 生成 json 文件
	run2MakeJson()

	// 生成 index 文件
	run3MakeHTML()

	// 删除必要的文件
	run4CleanFile()
}

func run1CheckoutBranch() {
	_ = goexec.New("git", "branch", "-D", "gh-pages").RunInStream()
	assert(goexec.New("git", "checkout", "-b", "gh-pages").RunInStream())
}

func run2MakeJson() {
	fs, err := ioutil.ReadDir("./json")
	assert(err)
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".json") {
			assert(goexec.New("cp", "./json/"+f.Name(), f.Name()).RunInStream())
		}
	}
}

func run3MakeHTML() {
	postDateList := loadPosts()
	for _, postDate := range postDateList {
		assert(ioutil.WriteFile(postDate.Date.In(ZoneGMT).Format("2006-01-02")+".md", []byte(postDate.FormatMD()), 0o666))
	}
	assert(ioutil.WriteFile("README.md", []byte(postDateList.FormatMD()), 0o666))
}

func run4CleanFile() {
	// assert(goexec.New("rm", ".gitignore").RunInStream())
	assert(goexec.New("rm", "LICENSE").RunInStream())
	assert(goexec.New("rm", "-rf", "./.github/fetch").RunInStream())
	assert(goexec.New("rm", "-rf", "./.github/render").RunInStream())
	assert(goexec.New("rm", "-rf", "./.github/workflows").RunInStream())
	assert(goexec.New("rm", "-rf", "./json").RunInStream())
}

func loadPosts() PostDateList {
	fs, err := ioutil.ReadDir("./json")
	assert(err)
	res := PostDateList{}
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".json") {
			date, err := time.Parse("2006-01-02", f.Name()[:len(f.Name())-len(".json")])
			assert(err)
			bs, err := ioutil.ReadFile("./json/" + f.Name())
			assert(err)
			posts := []*Post{}
			assert(json.Unmarshal(bs, &posts))
			res = append(res, &PostDate{Date: date, Posts: posts})
		}
	}
	sort.Sort(sort.Reverse(res)) // 最新时间在前
	return res
}

type Post struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Count int    `json:"count"`
}

type PostList []*Post

type PostDate struct {
	Date  time.Time
	Posts PostList
}

type PostDateList []*PostDate

func (r PostDateList) Len() int {
	return len(r)
}

func (r PostDateList) Less(i, j int) bool {
	return r[i].Date.Before(r[j].Date)
}

func (r PostDateList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r PostDate) FormatMD() string {
	date := r.Date.In(ZoneGMT).Format("2006-01-02")
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("# 北邮人论坛十大热帖 %s\n\n", date))
	for _, v := range r.Posts {
		s.WriteString(fmt.Sprintf("- [%s](%s) %d\n", v.Title, v.URL, v.Count))
	}
	s.WriteString("\n\n")
	return s.String()
}

func (r PostDateList) FormatMD() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf(`# 北邮人论坛十大热帖

## 项目地址

[项目 GitHub 地址](https://github.com/chyroc/byr-topten)

## 历史数据

`))
	for _, v := range r {
		date := v.Date.In(ZoneGMT).Format("2006-01-02")
		s.WriteString(fmt.Sprintf("- [北邮人论坛十大热帖 - %s](./%s.md) ([JSON 格式](./%s.json))\n", date, date, date))
	}
	s.WriteString("\n\n")
	return s.String()
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

var ZoneGMT = time.FixedZone("GMT+8", 8*60*60)
