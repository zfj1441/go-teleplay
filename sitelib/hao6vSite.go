package sitelib

import (
	"net/http"
	"io/ioutil"
	"strings"
	"github.com/antchfx/htmlquery"
	"log"
)

type Hao6v struct{}

func (p *Hao6v) GetHtml(url string) (html string, err error) {
	log.Println("Hao6v GetHtml")
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

func (p *Hao6v) ParseHtml(html string, tel Teleplay) (ret map[string]string, err error) {
	log.Println("Hao6v ParseHtml")
	ret = make(map[string]string)
	root, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		return ret, err
	}
	tr := htmlquery.Find(root, `//*[@id="endText"]/table/tbody/tr/td/a`)
	var name, url string
	for _, r := range tr {
		name = ConvertToString(htmlquery.InnerText(r), "gbk", "utf8")
		url = htmlquery.SelectAttr(r, `href`)
		ret[name] = url
	}

	// 筛选出未下载部分
	keys := make([]string, 0, len(ret))
	for k := range ret {
		keys = append(keys, k)
	}
	for _, key := range keys {
		if IsContain(tel.Download, key) {
			delete(ret, key)
		}
	}
	return ret, nil
}
