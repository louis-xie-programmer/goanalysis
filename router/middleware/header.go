package middleware
 
import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)
 
//无缓存头部中间件 ，
//要来防止客户端获取已经缓存的响应信息
func NoCache(c *gin.Context){
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}
 
//选项中间件
//要来给预请求 终止并退出中间件 ,链接并结束请求
func Options(c *gin.Context){
	if c.Request.Method != "OPTIONS"{
		c.Next()
	}else{
		c.Header("Access-Control-Allow-Origin","*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(200)
	}
}
 
//安全中间件
//要来保障数据安全的头部
func Secure(c *gin.Context){
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")
	if c.Request.TLS != nil {
		c.Header("Strict-Transport-Security", "max-age=31536000")
	}
 
	//也可以考虑添加一个安全代理的头部
	//c.Header("Content-Security-Policy", "script-src 'self' https://cdnjs.cloudflare.com")
}
