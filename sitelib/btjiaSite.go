package sitelib

import (
	"net/http"
	"io/ioutil"
	"os"
	"strings"
	"github.com/antchfx/htmlquery"
	"path/filepath"
	"log"
	"errors"
	"runtime"
)

//Btjia
type Btjia struct{}


func getTmp() string {
	var tmpPath string
	osType := runtime.GOOS
	if osType == `windows` {
		tmpPath = os.Getenv("TMP")
	}else {
		tmpPath = `/tmp`
	}
	return tmpPath
}

func (p *Btjia) GetHtml(url string) (html string, err error) {
	log.Println("btjia GetHtml")
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()                // 函数结束时关闭Body
	body, err := ioutil.ReadAll(resp.Body) // 读取Body
	if err != nil {
		return "", err
	}
	html = string(body)
	return html, nil
}

func (p *Btjia) ParseHtml(html string, tel Teleplay) (ret map[string]string, err error) {
	log.Println("btjia ParseHtml")
	tmp := make(map[string]string)
	ret = make(map[string]string)
	root, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		return ret, err
	}
	tr := htmlquery.Find(root, `//*[@id="body"]/div/table/tbody/tr/td[@class="post_td"]/div/div/table[@class="noborder"]/tbody/tr/td/a`)
	var name, url string
	for _, r := range tr {
		name = htmlquery.InnerText(r)
		url = htmlquery.SelectAttr(r, `href`)
		tmp[name] = strings.Replace(strings.Replace(url, "dialog", "download", 1), "-ajax-1", "", 1)
	}
	// 筛选出未下载部分
	keys := make([]string, 0, len(tmp))
	for k := range tmp {
		keys = append(keys, k)
	}
	for _, key := range keys {
		if IsContain(tel.Download, key) {
			delete(tmp, key)
		}
	}
	//下载种子并转成磁力链接
	for k, v := range tmp {
		fname := getBt(k, v)
		if fname != "" {
			tmp_name := filepath.Join(getTmp(), fname)
			ret[fname] = Bt2magnet(tmp_name)
		} else {
			return nil, errors.New("种子下载失败")
		}
	}
	return ret, nil
}

//bt种子下载
func getBt(name string, url string) (fname string) {
	tmp_name := filepath.Join(getTmp(), name)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("ERROR:", err)
		return ""
	}
	resCode := resp.StatusCode
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resCode == 200 {
		f, err := os.Create(tmp_name)
		if err != nil {
			log.Println("ERROR:", err)
			return ""
		}
		f.Write(body)
		f.Close()
	}else {
		log.Println("下载种子失败")
		return ""
	}
	return name
}
