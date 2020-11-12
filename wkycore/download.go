package wkycore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

//Peers defines device formation
type Peers struct {
	Devices []PeerInfo `json:"devices"`
}

type PeerInfo struct {
	BindTime int `json:"bind_time"`
	Features struct {
		OnecloudCoin int `json:"onecloud_coin"`
		Miner        int `json:"miner"`
	} `json:"features"`
	IPInfo struct {
		Province string `json:"province"`
		Country  string `json:"country"`
		Isp      string `json:"isp"`
		City     string `json:"city"`
	} `json:"ip_info"`
	IP                string        `json:"ip"`
	CoturnOnline      int           `json:"coturn_online"`
	Paused            bool          `json:"paused"`
	DcdnUpnpMessage   string        `json:"dcdn_upnp_message"`
	DeviceSn          string        `json:"device_sn"`
	ExceptionMessage  string        `json:"exception_message"`
	AccountName       string        `json:"account_name"`
	DcdnClients       []interface{} `json:"dcdn_clients"`
	Hibernated        bool          `json:"hibernated"`
	Imported          int64         `json:"imported"`
	ExceptionName     string        `json:"exception_name"`
	HardwareModel     string        `json:"hardware_model"`
	MacAddress        string        `json:"mac_address"`
	Status            string        `json:"status"`
	LanIP             string        `json:"lan_ip"`
	AccountType       string        `json:"account_type"`
	AccountID         string        `json:"account_id"`
	Upgradeable       bool          `json:"upgradeable"`
	DcdnDownloadSpeed int           `json:"dcdn_download_speed"`
	DcdnUploadSpeed   int           `json:"dcdn_upload_speed"`
	DiskQuota         int64         `json:"disk_quota"`
	Peerid            string        `json:"peerid"`
	Licence           string        `json:"licence"`
	DcdnID            string        `json:"dcdn_id"`
	SystemVersion     string        `json:"system_version"`
	DeviceID          string        `json:"device_id"`
	SystemName        string        `json:"system_name"`
	DcdnUpnpStatus    string        `json:"dcdn_upnp_status"`
	ProductID         int           `json:"product_id"`
	DeviceName        string        `json:"device_name"`
	ScheduleHours     []struct {
		To     int `json:"to"`
		From   int `json:"from"`
		Params struct {
		} `json:"params"`
		Type string `json:"type"`
	} `json:"schedule_hours"`
}

//Peers defines device formation
type Usbs struct {
	Partitions []UsbInfo `json:"partitions"`
}

type UsbInfo struct {
	Label      string `json:"label"`
	Id         int    `json:"id"`
	PartSymbol string `json:"part_symbol"`
	DiskSn     int32  `json:"disk_sn"`
	Uuid       string `json:"uuid"`
	Used       string `json:"used"`
	Unique     int    `json:"unique"`
	Path       string `json:"path"`
	PartLabel  string `json:"part_label"`
	Capacity   string `json:"capacity"`
	DiskId     int    `json:"disk_id"`
	FsType     string `json:"fs_type"`
	Type       string `json:"type"`
}

type TaskInfo struct {
	CreateTime   int32         `json:"createTime"`
	State        int           `json:"state"`
	Id           string        `json:"id"`
	Refer        string        `json:"refer"`
	Progress     int           `json:"progress"`
	Exist        int           `json:"exist"`
	Url          string        `json:"url"`
	RemainTime   int           `json:"remainTime"`
	Size         string        `json:"size"`
	Name         string        `json:"name"`
	FailCode     int           `json:"failCode"`
	From         int           `json:"from"`
	SubList      []interface{} `json:"subList"`
	Speed        int           `json:"speed"`
	DownTime     int           `json:"downTime"`
	CompleteTime int           `json:"completeTime"`
	Type         int           `json:"type"`
	DcdnChannel  struct {
		DlBytes   int    `json:"dlBytes"`
		Available int    `json:"available"`
		DlSize    string `json:"dlSize"`
		State     int    `json:"state"`
		FailCode  int    `json:"failCode"`
		Speed     int    `json:"speed"`
	} `json:"dcdnChannel"`
	LixianChannel struct {
		DlBytes        int    `json:"dlBytes"`
		State          int    `json:"state"`
		ServerSpeed    int    `json:"serverSpeed"`
		DlSize         string `json:"dlSize"`
		ServerProgress int    `json:"serverProgress"`
		FailCode       int    `json:"failCode"`
		Speed          int    `json:"speed"`
	} `json:"lixianChannel"`
	VipChannel struct {
		Opened    int    `json:"opened"`
		DlBytes   int    `json:"dlBytes"`
		Available int    `json:"available"`
		DlSize    string `json:"dlSize"`
		Speed     int    `json:"speed"`
		FailCode  int    `json:"failCode"`
		Type      int    `json:"type"`
	} `json:"vipChannel"`
}

type Job struct {
	Filesize string `json:"filesize"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Infohash string `json:"infohash"`
}

func (wky *WkyCore) GetPeerList() error {
	type respListPeer struct {
		Rtn    int           `json:"rtn"`
		Msg    string        `json:"msg"`
		Result []interface{} `json:"result"`
	}

	var (
		err  error
		req  *http.Request
		resp *http.Response
	)

	sign := GetSign(true, map[string]string{
		"X-LICENCE-PUB": "1",
		"appversion":    WkAppVersion,
		"v":             "2",
		"ct":            "9",
	}, wky.Userinfo.SessionID)

	query := map[string]string{
		"X-LICENCE-PUB": "1",
		"appversion":    WkAppVersion,
		"v":             "2",
		"ct":            "9",
		"sign":          sign,
	}

	u, _ := url.Parse(URL_ListPeer)
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	if req, err = http.NewRequest("GET", u.String(), nil); err != nil {
		log.Printf("请求（%+v）失败: %+v", URL_ListPeer, err)
		return err
	}
	if resp, err = wky.client.Do(req); err != nil {
		log.Printf("二维码登陆校验失败: %+v", err)
		return err
	}
	content, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(content))
	var v respListPeer
	if err := json.Unmarshal(content, &v); err != nil {
		return err
	}
	if v.Rtn != 0 {
		return fmt.Errorf("Invalid response result[%+v]", v.Rtn)
	}
	if len(v.Result) != 2 {
		return fmt.Errorf("Invalid response result[Result:%+v]", len(v.Result))
	}

	//Convert interface to json string and then convert to struct Peers
	var b []byte
	if b, err = json.Marshal(v.Result[1]); err != nil {
		return err
	}
	var list Peers
	if err := json.Unmarshal(b, &list); err != nil {
		return err
	}
	if len(list.Devices) != 1 {
		log.Printf("Found %d peers!", len(list.Devices))
	}

	wky.Peers = &list
	log.Printf("PeerList %s", string(content))
	return nil
}

func (wky *WkyCore) GetUSBInfos() ([]UsbInfo, error) {
	type respListUsb struct {
		Rtn    int           `json:"rtn"`
		Msg    string        `json:"msg"`
		Result []interface{} `json:"result"`
	}

	var (
		err  error
		req  *http.Request
		resp *http.Response
	)
	sign := GetSign(true, map[string]string{
		"X-LICENCE-PUB": "1",
		"appversion":    WkAppVersion,
		"v":             "2",
		"ct":            "9",
		"deviceid":      wky.Peers.Devices[0].DeviceID,
	}, wky.Userinfo.SessionID)

	query := map[string]string{
		"X-LICENCE-PUB": "1",
		"appversion":    WkAppVersion,
		"v":             "2",
		"ct":            "9",
		"deviceid":      wky.Peers.Devices[0].DeviceID,
		"sign":          sign,
	}

	u, _ := url.Parse(URL_UsbInfo)
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	if req, err = http.NewRequest("GET", u.String(), nil); err != nil {
		log.Printf("请求（%+v）失败: %+v", URL_UsbInfo, err)
		return nil, err
	}
	if resp, err = wky.client.Do(req); err != nil {
		log.Printf("获取USB信息失败: %+v", err)
		return nil, err
	}
	content, _ := ioutil.ReadAll(resp.Body)
	log.Printf("USB信息:%+v", string(content))

	var v respListUsb
	if err := json.Unmarshal(content, &v); err != nil {
		return nil, err
	}
	if v.Rtn != 0 {
		return nil, errors.New("Invalid response result")
	}
	if len(v.Result) != 2 {
		return nil, errors.New("Invalid response result")
	}
	var b []byte
	if b, err = json.Marshal(v.Result[1]); err != nil {
		return nil, err
	}
	var list Usbs
	if err := json.Unmarshal(b, &list); err != nil {
		return nil, err
	}
	if len(list.Partitions) != 1 {
		return nil, fmt.Errorf("Found %d peers!", len(list.Partitions))
	}
	return list.Partitions, nil
}

//远端登录
func (wky *WkyCore) RemoteDlLogin() error {
	var (
		err  error
		req  *http.Request
		resp *http.Response
	)
	sign := GetSign(true, map[string]string{
		"pid":        wky.Peers.Devices[0].Peerid,
		"appversion": WkAppVersion,
		"v":          "1",
		"ct":         "32",
	}, wky.Userinfo.SessionID)

	query := map[string]string{
		"pid":        wky.Peers.Devices[0].Peerid,
		"appversion": WkAppVersion,
		"v":          "1",
		"ct":         "32",
		"sign":       sign,
	}

	u, _ := url.Parse(URL_RemoteLogin)
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	if req, err = http.NewRequest("GET", u.String(), nil); err != nil {
		log.Printf("请求（%+v）失败: %+v", URL_ListPeer, err)
		return err
	}
	if resp, err = wky.client.Do(req); err != nil {
		log.Printf("登录远程下载失败: %+v", err)
		return err
	}
	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(content))
	return nil
}

//获取远端下载任务信息
func (wky *WkyCore) GetRemoteDlInfo() error {
	var (
		err  error
		req  *http.Request
		resp *http.Response
	)
	sign := GetSign(true, map[string]string{
		"pid":     wky.Peers.Devices[0].Peerid,
		"ct":      "31",
		"v":       "2",
		"pos":     "0",
		"number":  "100",
		"type":    "4",
		"needUrl": "0",
	}, wky.Userinfo.SessionID)

	query := map[string]string{
		"pid":     wky.Peers.Devices[0].Peerid,
		"ct":      "31",
		"v":       "2",
		"pos":     "0",
		"number":  "100",
		"type":    "4",
		"needUrl": "0",
		"sign":    sign,
	}

	u, _ := url.Parse(URL_RemoteInfo)
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	if req, err = http.NewRequest("GET", u.String(), nil); err != nil {
		log.Printf("请求（%+v）失败: %+v", URL_ListPeer, err)
		return err
	}
	if resp, err = wky.client.Do(req); err != nil {
		log.Printf("获取远端下载信息失败: %+v", err)
		return err
	}
	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(content))
	return nil
}

//解析下载链接
func (wky *WkyCore) UrlResolve(urlstr string) (taskinfo TaskInfo, infoHash string, e error) {
	type respUrlInfo struct {
		Rtn      int    `json:"rtn"`
		Msg      string `json:"msg"`
		Infohash string `json:"infohash"`
		TaskInfo `json:"taskInfo"`
	}
	var (
		err  error
		resp *http.Response
	)
	body := map[string]string{
		"url": urlstr,
	}

	query := map[string]string{
		"pid": wky.Peers.Devices[0].Peerid,
		"ct":  "31",
		"v":   "2",
	}
	u, _ := url.Parse(URL_UrlResolve)
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	jsonbytes, _ := json.Marshal(body)
	buffer := bytes.NewBuffer(jsonbytes)
	if resp, err = wky.client.Post(u.String(), "", buffer); err != nil {
		log.Printf("解析链接（%+v）失败: %+v", URL_UrlResolve, err)
		return taskinfo, "", err
	}
	content, _ := ioutil.ReadAll(resp.Body)
	log.Printf("下载连接信息:%+v", string(content))
	var v respUrlInfo
	if err := json.Unmarshal(content, &v); err != nil {
		log.Printf("解析返回值失败：%+v", err.Error())
		return taskinfo, "", err
	}
	if v.Rtn != 0 {
		return taskinfo, "", errors.New("Invalid response result")
	}
	return v.TaskInfo, v.Infohash, nil
}

//创建下载任务
func (wky *WkyCore) CreateTasks(urlstr string) error {
	var (
		jb             Job
		remoteLocation string
		err            error
		resp           *http.Response
	)
	type respDownload struct {
		Rtn int    `json:"rtn"`
		Msg string `json:"msg"`
	}
	//解析url
	if urlinfo, infohash, err := wky.UrlResolve(urlstr); err != nil {
		return err
	} else {
		log.Printf("URL信息:%+v", urlinfo)
		jb.Filesize = urlinfo.Size
		jb.Name = urlinfo.Name
		jb.Infohash = infohash
		jb.Url = urlstr
	}

	//构造下载路径 及 试算存储空间
	if usbinfos, err := wky.GetUSBInfos(); err != nil {
		return err
	} else {
		if len(usbinfos) <= 0 {
			return fmt.Errorf("无可用存储设备")
		} else {
			remoteLocation = usbinfos[0].Path + "/onecloud/tddownload"
		}
		//迁移到arm平台导致计算溢出
		//capacity, _ := strconv.Atoi(usbinfos[0].Capacity)
		//used, _ := strconv.Atoi(usbinfos[0].Used)
		//filesize, _ := strconv.Atoi(jb.Filesize)
		//free2Use := capacity - used - filesize
		//if free2Use <= 0 {
		//	log.Printf("剩余空间不足: %+v", free2Use)
		//	return fmt.Errorf("剩余空间不足: %+v", free2Use)
		//}
	}

	type Body struct {
		Path  string `json:"path"`
		Tasks []Job  `json:"tasks"`
	}

	var body Body
	body.Path = remoteLocation
	body.Tasks = append(body.Tasks, jb)

	query := map[string]string{
		"pid": wky.Peers.Devices[0].Peerid,
		"ct":  "31",
		"v":   "1",
	}
	u, _ := url.Parse(URL_CreateTasks)
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	jsonbytes, _ := json.Marshal(body)
	buffer := bytes.NewBuffer(jsonbytes)
	if resp, err = wky.client.Post(u.String(), "application/json", buffer); err != nil {
		log.Printf("创建下载任务（%+v）失败: %+v", URL_CreateTasks, err)
		return err
	}
	content, _ := ioutil.ReadAll(resp.Body)
	log.Printf("下载任务返回：%+v", string(content))
	var v respDownload
	if err := json.Unmarshal(content, &v); err != nil {
		log.Printf("解析返回值失败：%+v", err.Error())
		return err
	}
	if v.Rtn != 0 {
		return errors.New("Invalid response result")
	}
	return nil
}
