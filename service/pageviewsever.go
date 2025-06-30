package service

import (
	"crypto/sha256"
	"fmt"
	. "goanalysis/handler"
	. "goanalysis/model"
	"log"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func PageViewDo(c *gin.Context, producer Producer) {
	elog,err := GetPageViewLogDB(c)

	if err == nil {
		b, _ := elog.ConvertJson()
	  producer.Send("pageview", string(b))
	}

	SendResponse(c, nil, nil)
}

// 整理日志字段（pageview）
func GetPageViewLogDB(c *gin.Context) (PageViewLog, error) {
	//初始化会话
	isnew := InitSession(c)

	elog := PageViewLog{}

	if err := c.ShouldBindJSON(&elog); err != nil {
		log.Println("GetPageViewLogDB ShouldBind error:", err.Error())

		return elog, err
	}

	//深度，过滤掉刷新页面和历史记录页面
	isview := true

	switch elog.Viewtype {
	case "0":
		elog.Viewtype = "normal"
	case "1":
		elog.Viewtype = "reload"
		isview = false
	case "2":
		elog.Viewtype = "back"
		isview = false
	case "255":
		elog.Viewtype = "other"
	default:
		elog.Viewtype = "unknown"
	}

	//设置会话访问序号
	SetPageViewIndex(c, isnew, isview)

	elog.WebName = c.Param("webname")
	elog.UA = c.Request.Header["User-Agent"][0]
	elog.Lang = c.Request.Header["Accept-Language"][0]
	elog.IP = GetRealIP(c)                              //GetRealIP(c)
	elog.Sessionid = GetSession(c, "sessionID").(string) //会话标记
	elog.Depth = GetSession(c, "index").(int)            //访问顺序标记
	elog.Time = time.Now().Unix()

	//通过cf的安全头部获取定位信息 cf-ipcontinent
	elog.Continent = c.Request.Header["cf-ipcontinent"][0]
	elog.Country = c.Request.Header["cf-ipcountry"][0]
	elog.Provinces = c.Request.Header["cf-region"][0]
	elog.City = c.Request.Header["cf-ipcity"][0]
	elog.Location = c.Request.Header["cf-iplongitude"][0] + "," + c.Request.Header["cf-iplatitude"][0]

	return elog, nil
}

// 获取真实ip
func GetRealIP(c *gin.Context) string {
	ip := c.Request.Header.Get("x-Forwarded-For")
	if ip == "" {
		ip = c.ClientIP()
	}
	return ip
}

// 初始化会话,如果是新的会话则放回true
func InitSession(c *gin.Context) bool {
	s := sessions.Default(c)
	sessionID := s.Get("sessionID")
	if sessionID == nil {
		str := strconv.FormatInt(time.Now().UnixMilli(), 10)
		nums := sha256.Sum256([]byte(str))
		sid := fmt.Sprintf("%x", nums)
		SetSession(c, "sessionID", sid)
		return true
	}
	return false
}

// 设置浏览序号
func SetPageViewIndex(c *gin.Context, isnew bool, isview bool) {
	s := sessions.Default(c)
	sindex := s.Get("index")
	if sindex == nil || isnew {
		s.Set("index", 1)
		s.Save()
	} else if isview {
		sindexnum := sindex.(int)
		sindexnum += 1
		s.Set("index", sindexnum)
		s.Save()
	}
}

// 设置会话
func SetSession(c *gin.Context, key string, value string) {
	s := sessions.Default(c)
	s.Set(key, value)
	s.Save()
}

// 获取会话值
func GetSession(c *gin.Context, sname string) any {
	s := sessions.Default(c)
	val := s.Get(sname)
	if val == nil {
		return nil
	} else {
		return val
	}
}
