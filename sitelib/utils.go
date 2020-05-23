package sitelib

import (
	"github.com/anacrolix/torrent/metainfo"
	"log"
	"github.com/axgle/mahonia"
	"net/http"
	"net/url"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"strings"
	"net/smtp"
)

//种子转磁力
func Bt2magnet(path string) (magnet string) {
	type TFile struct {
		Name string
		Size int64
	}
	mi, err := metainfo.LoadFromFile(path)
	if err != nil {
		log.Println(err)
		return ""
	}
	info, err := mi.UnmarshalInfo()
	if err != nil {
		log.Println(err)
		return ""
	}
	sl := make([]TFile, len(info.Files))
	for index := 0; index < len(info.Files); index++ {
		filename := info.Files[index].DisplayPath(&info)
		sl[index] = TFile{
			Name: filename,
			Size: info.Files[index].Length,
		}
	}
	return "magnet:?xt=urn:btih:" + mi.HashInfoBytes().String()
}

//判断字符串包含
func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

//字符编码转换
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

// Server酱推送
func SendToServerChan(title string, message string, key string) (errno int, errmsg string){

	type retMsg struct {
		Errno    int  `json:"errno"`
		Errmsg   string  `json:"errmsg"`
	}
	ret := retMsg{}
	serverurl := fmt.Sprintf("http://sc.ftqq.com/%s.send", key)

	resp, err := http.PostForm(serverurl, url.Values{"text": {title}, "desp": {message}})
	if err != nil {
		log.Fatalln(err)
		return-1, `系统错误`
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return-1, `系统错误`
	}

	err = json.Unmarshal(body, &ret)
	if err != nil{
		return-1, `系统错误`
	}
	return ret.Errno, ret.Errmsg
}

//发送邮件
func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}