package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"go-teleplay/sitelib"
	"go-teleplay/wkycore"
	"io"
	"log"
	"os"
	"time"
)

//配置文件结构
type Config struct {
	ServerChanKey string
	Phone         string
	Pass          string
	Teleplays     []struct {
		Status   int
		Name     string
		Url      string
		Type     string
		Download []string
	}
}

//读取配置文件
func readJson(filePath string) (result string) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		result = ""
	}
	buf := bufio.NewReader(file)
	for {
		s, err := buf.ReadString('\n')
		result += s
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return ""
			}
		}
	}
	return result
}

//写入配置文件
func writeJson(filePath string, jsonString string) bool {
	file, err := os.Create(filePath)
	if err != nil {
		return false
	}
	defer file.Close()
	n, err1 := io.WriteString(file, jsonString) //写入文件(字符串)
	if err1 != nil || len(jsonString) != n {
		return false
	}
	return true
}

// 启动参数解析
func param() (string, string) {
	var configFile string
	var logFile string

	flag.StringVar(&configFile, "f", "config.json", "配置文件路径")
	flag.StringVar(&logFile, "l", "go-teleplay.log", "日志文件路径")
	flag.Parse()
	return configFile, logFile
}

var g_Configfile, g_Logfile string

func init() {
	g_Configfile, g_Logfile = param()

	fmt.Printf("启动参数 config_path=[%v] log_path=[%v]\n",
		g_Configfile, g_Logfile)

	// 设置日志
	logFile, err := os.OpenFile(g_Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(logFile)
}

func main() {
	cv := Config{}
	result := readJson(g_Configfile)
	err := json.Unmarshal([]byte(result), &cv)
	if err != nil {
		log.Println("ERROR: %+v", err.Error())
		return
	}

	//玩客云客户端初始化
	wkyClit := wkycore.NewWkyCore(cv.Phone, cv.Pass)
	defer wkyClit.Release()
	if !wkyClit.LoginEx() {
		log.Fatal("玩客云客户端登录失败")
		return
	}

	//站点插件初始化
	plugin := new(sitelib.Sites)
	plugin.Init()
	hao6v := new(sitelib.Hao6v)
	btjia := new(sitelib.Btjia)
	plugin.Register("hao6v", hao6v)
	plugin.Register("btjia", btjia)

	var i int
	var site sitelib.Siteinfunc
	for i = 0; i < len(cv.Teleplays); i++ {
		film := cv.Teleplays[i]
		// 加入监控启用状态
		if film.Status == 0 {
			log.Printf("[%v]监控未启用", film.Name)
			continue
		}
		site = plugin.Sitelist[film.Type]
		html, err := site.GetHtml(film.Url) //下载数据
		if err != nil {
			log.Println("下载数据失败")
			continue
		}
		ret, err := site.ParseHtml(html, film) //解析网页返回tab内容
		if err != nil {
			log.Println("解析网页失败")
			continue
		}

		// 发送邮件
		for k, v := range ret {
			log.Printf("k[%v][%v]\n", k, v)
			if err := wkyClit.CreateTasks(v); err != nil {
				log.Printf("请求下载失败:%+v", err.Error())
				retCode, retMsg := sitelib.SendToServerChan(k, v, cv.ServerChanKey)
				if retCode != 0 {
					log.Println(k, v)
					log.Printf("Send mail error!%v", retMsg)
				} else {
					log.Println("Send mail success!")
					film.Download = append(film.Download, k)
				}
			} else {
				log.Print("请求下载成功")
				film.Download = append(film.Download, k)
			}
		}
		cv.Teleplays[i] = film
		time.Sleep(5000)
	}

	bytes, err := json.MarshalIndent(cv, "", " ")
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	writeJson(g_Configfile, string(bytes))

}
