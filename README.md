# golang追剧工具
简单的追剧工具，根据config.json中的配置爬区电视剧最新信息，然后邮件推送

# 使用说明
1. config.json中配置需要追的电视剧  type:暂时只支持 btjia,hao6v
2. 修改config.json抓取电视信息
3. 增加crontab定时任务

# 持续开发
sitelib资源网站库，新开发网站库步骤
1. 实现GetHtml,ParseHtml方法
2. 在go-teleplay.go 105行注册新资源

# 路由器K3编译
```
  linux/unix:
    env GOOS=linux GOARCH=arm go build go-teleplay.go
  windows:
    set GOOS=linux
    set GOARCH=arm
    go build go-teleplay.go
```
# 路由器K2p编译:
```
  windows:
    set GOOS=linux
    set GOARCH=mipsle
    go build -ldflags "-w -s" go-teleplay.go
```
# 使用upx压缩
  `upx.exe -9 go-teleplay`

# 部署
 上传执行文件和配置文件文档结构如图   
├── go-teleplay  
├── config.json  
└── tmp  
增加定时任务   
  crontab -e   
  0,30 9-23 * * *  /root/go-teleplay >>/tmp/go-teleplay.log 2>&1 #9点到23点 每30分钟执行一次
    

# 更新
* 2020年11月12日.更新玩客云直接推下载


# 感谢
* 玩客云python接口[https://github.com/allanpk716/WanKeYunApi](https://github.com/allanpk716/WanKeYunApi)
* [https://github.com/Adyzng/go-jd](https://github.com/Adyzng/go-jd)对cookie的封装
* [https://github.com/jsfaint/sentinel](https://github.com/jsfaint/sentinel)