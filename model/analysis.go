package model

import (
	"encoding/json"

	"goanalysis/pkg/errno"
)

// 页面浏览日志结构
// vtype: 0,1,2 window.performance.navigation.type
type PageViewLog struct {
	PId     string `json:"i"` //页面id
	WebName string //网站客户端名称
	Url     string `json:"u"` //网页链接
	Title   string `json:"t"` //页面标题
	Status  int    `json:"s"` //页面状态码

	MId    string `json:"m"` //机器识别码
	UA     string //useragen 代理
	AutoUA bool   `json:"w"`  //是否是模拟的ua
	Screen string `json:"sc"` //屏幕
	Proto  string `json:"p"`  //http 版本（http 1.1 , h2, h3）

	Lang string //语言
	IP   string //ip
	
	Continent string //所属洲
	Country   string //国家
	Provinces string //省或州
	City      string //城市
	Location  string //坐标（纬度，经度）

	Sessionid string //会话识别码
	Depth     int    //流量深度,只有在viewtype为正常访问才加深度
	Referer   string `json:"r"` //上一页面
	Viewtype  string `json:"v"` //浏览类型(需要对客户端发来的window.performance.navigation.type值进行语义化)

	Time int64 //时间
}

// 事件日志结构
type Eventlog struct {
	Id      string //事件标识
	EType   string `json:"e"` //事件类型
	WebName string //网页名称
	Url     string `json:"u"` //当前网页链接地址
	PId     string `json:"i"` //pageview id
	JsonDb  string `json:"c"` //事件内容
	Time    int64  //时间
}

// PageViewLog 转化成json字符串
func (model PageViewLog) ConvertJson() (string, error) {
	jsons, errs := json.Marshal(model)
	if errs != nil {
		return "", errno.ModelError
	}
	return string(jsons), nil
}

// Eventlog 转化成json字符串
func (model Eventlog) ConvertJson() (string, error) {
	jsons, errs := json.Marshal(model)
	if errs != nil {
		return "", errno.ModelError
	}
	return string(jsons), nil
}
