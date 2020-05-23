package sitelib

//配置文件中电影数据结构
type Teleplay struct {
	Name     string
	Url      string
	Type     string
	Download []string
}

// 定义一个接口，里面有两个方法
type Siteinfunc interface {
	GetHtml(url string) (string, error)
	ParseHtml(html string, tel Teleplay) (map[string]string, error)
}

// 定义一个类，来存放我们的插件
type Sites struct {
	Sitelist map[string]Siteinfunc
}

// 初始化插件
func (p *Sites) Init() {
	p.Sitelist = make(map[string]Siteinfunc)
}

// 注册插件
func (p *Sites) Register(name string, Site Siteinfunc) {
	p.Sitelist[name] = Site
}
