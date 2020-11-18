package wkycore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	//URL
	URL_Login       = "http://account.onethingpcs.com/user/login"
	URL_ListPeer    = "http://control.onethingpcs.com/listPeer"
	URL_UsbInfo     = "http://control.onethingpcs.com/getUSBInfo"
	URL_RemoteLogin = "http://control-remotedl.onethingpcs.com/login"
	URL_RemoteInfo  = "http://control-remotedl.onethingpcs.com/list"
	URL_UrlResolve  = "http://control-remotedl.onethingpcs.com/urlResolve"
	URL_CreateTasks = "http://control-remotedl.onethingpcs.com/createTask"

	WkAppVersion = "2.6.0"
	AccountType  = "4"
)

var (
	DefaultHeaders = map[string]string{
		"User-Agent":      "Chrome/51.0.2704.103",
		"ContentType":     "application/json", //"text/html; charset=utf-8",
		"Connection":      "keep-alive",
		"Accept-Encoding": "gzip, deflate",
		"Accept-Language": "zh-CN,zh;q=0.8",
	}
)

// 玩客云配置文件
type WkyConfig struct {
	Phone        string
	Pass         string
	CookieFile   string
	UserinfoFile string
}

//UserInfo defines user information
type UserInfo struct {
	Userid          string `json:"userid"`
	Phone           string `json:"phone"`
	PhoneArea       string `json:"phone_area"`
	AccountType     string `json:"account_type"`
	BindPwd         string `json:"bind_pwd"`
	NickName        string `json:"nickname"`
	SessionID       string `json:"sessionid"`
	EnableHomeShare uint   `json:"enable_homeshare"`
}

// WkyCore wrap jing dong operation
type WkyCore struct {
	WkyConfig WkyConfig
	client    *http.Client
	jar       *SimpleJar
	Userinfo  *UserInfo
	Peers     *Peers
}

// NewWkyCore create an object to wrap WkyCore related operation
//
func NewWkyCore(cf WkyConfig) *WkyCore {
	wky := &WkyCore{
		WkyConfig: cf,
	}

	wky.jar = NewSimpleJar(JarOption{
		JarType:  JarJson,
		Filename: wky.WkyConfig.CookieFile,
	})

	// 装载cookies
	if err := wky.jar.Load(); err != nil {
		log.Println("加载Cookies失败: %s", err)
		wky.jar.Clean()
	}
	// 状态用户信息
	var uinfo UserInfo
	fd, err := os.Open(wky.WkyConfig.UserinfoFile)
	if err == nil {
		err = json.NewDecoder(fd).Decode(&uinfo)
	} else if os.IsNotExist(err) {
		err = nil
	}
	wky.Userinfo = &uinfo

	wky.client = &http.Client{
		Timeout: time.Minute,
		Jar:     wky.jar,
	}
	return wky
}

// Release the resource opened
//
func (wky *WkyCore) Release() {
	if wky.jar != nil {
		if err := wky.jar.Persist(); err != nil {
			log.Printf("Failed to persist cookiejar. error %+v.", err)
		}
	}

	if wky.Userinfo != nil {
		fd, err := os.Create(wky.WkyConfig.UserinfoFile)
		if err == nil {
			err = json.NewEncoder(fd).Encode(wky.Userinfo)
		}
	}
}

func applyCustomHeader(req *http.Request, header map[string]string) {
	if req == nil || len(header) == 0 {
		return
	}

	for key, val := range header {
		req.Header.Set(key, val)
	}
}

//校验是否需要登录
func (wky *WkyCore) validateLogin() bool {
	err := wky.GetPeerList()
	if err != nil {
		log.Print(err.Error())
		return false
	} else {
		log.Print("success")
		return true
	}
	return false
}

// 使用用户密码登录
func (wky *WkyCore) Login() error {

	type respLogin struct {
		Rtn  int      `json:"rtn"`
		Msg  string   `json:"msg"`
		Data UserInfo `json:"data"`
	}

	var (
		err  error
		resp *http.Response
	)

	sign := GetSign(false, map[string]string{
		"deviceid":     GetDevID(wky.WkyConfig.Phone),
		"imeiid":       GetIMEI(wky.WkyConfig.Phone),
		"phone":        wky.WkyConfig.Phone,
		"pwd":          GetPWD(wky.WkyConfig.Pass),
		"account_type": AccountType,
	}, "")

	body := map[string]string{
		"deviceid":     GetDevID(wky.WkyConfig.Phone),
		"imeiid":       GetIMEI(wky.WkyConfig.Phone),
		"phone":        wky.WkyConfig.Phone,
		"pwd":          GetPWD(wky.WkyConfig.Pass),
		"account_type": "4",
		"sign":         sign,
	}

	query := map[string]string{
		"appversion": WkAppVersion,
	}

	u, _ := url.Parse(URL_Login)
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	var values []string
	for k, v := range body {
		values = append(values, fmt.Sprintf("%s=%s", k, v))
	}

	resp, err = wky.client.Post(u.String(), "", strings.NewReader(strings.Join(values, "&")))
	if err != nil {
		log.Printf("请求（%+v）失败: %+v", u.String(), err)
		return err
	}
	if resp.StatusCode == http.StatusOK {
		log.Printf("请求成功")
	} else {
		log.Printf("登陆失败:%+v", resp.Status)
		return fmt.Errorf("登陆失败:%+v", resp.Status)
	}
	contentBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应错误 %+v", err.Error())
		return err
	}
	var v respLogin
	if err := json.Unmarshal([]byte(contentBytes), &v); err != nil {
		log.Printf("响应解析错误 %+v", err.Error())
		return err
	}
	if v.Rtn != 0 {
		return fmt.Errorf("Invalid response result[%+v][%+v]", v.Rtn, v.Msg)
	}

	wky.Userinfo = &v.Data
	return nil
}

func (wky *WkyCore) LoginEx(args ...interface{}) bool {
	if wky.validateLogin() {
		log.Print("无需重新登录")
		return true
	} else {
		log.Print("新登录中")
		wky.jar.Clean()

		if err := wky.Login(); err != nil {
			log.Printf(err.Error())
			return false
		} else {
			if wky.validateLogin() {
				return true
			} else {
				log.Print("玩客云登录失败")
				return false
			}
		}
	}
}
