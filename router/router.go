package router

import (
	"goanalysis/router/middleware"
	"goanalysis/service"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

//初始化路由
func InitRouter(g *gin.Engine){
	middlewares := []gin.HandlerFunc{}
	//中间件
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(middlewares...)
	g.Use(gin.Logger())

	//注入会话存储
	store := cookie.NewStore([]byte("secret"))
	g.Use(sessions.Sessions("mysession", store))
 
	//404处理
	g.NoRoute(func(c *gin.Context){
		c.String(http.StatusNotFound,"该路径不存在")
	})

	g.GET("/",service.Index)//主页

	//页面浏览组
	pView := g.Group("/api/v1/pageview")
	//页面事件组
	pEvent := g.Group("/api/v1/event")

	producer := service.NewProducer()

	pView.POST(":webname", func(ctx *gin.Context) {
		service.PageViewDo(ctx,producer)
	})

	pEvent.POST("analytics/:webname",func(ctx *gin.Context) {
		service.EnventLogDo(ctx,producer)
	})
}