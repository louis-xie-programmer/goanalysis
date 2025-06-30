package main

import (
	"goanalysis/config"
	"goanalysis/router"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	//设置gin模式
	gin.SetMode(viper.GetString("common.server.runmode"))

	//创建一个gin引擎
	g := gin.New()

	router.InitRouter(g)

	istls := viper.GetBool("common.server.tls")

	if istls {
		log.Printf("开始监听服务器地址: %s\n", "https://"+viper.GetString("common.server.addr"))
		g.RunTLS(viper.GetString("common.server.addr"), viper.GetString("common.server.tlspem"), viper.GetString("common.server.tlskey"))
	} else {
		log.Printf("开始监听服务器地址: %s\n", "http://"+viper.GetString("common.server.addr"))
		g.Run(viper.GetString("common.server.addr"))
	}
}
