package config
 
import (
	"github.com/spf13/viper"
	"time"
	"os"
	"log"
)
 
// LogInfo 初始化日志配置
func LogInfo() {
	file := "./logs/" + time.Now().Format("2006-01-02") + ".log"
	logFile, _ := os.OpenFile(file,os.O_RDWR| os.O_CREATE| os.O_APPEND, 0755)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(logFile)
}
 
// Init 读取初始化配置文件
func Init() error {
	//初始化配置
	if err := Config();err != nil{
		return err
	}
 
	//初始化日志
	LogInfo()
	return nil
}
 
// Config viper解析配置文件
func Config() error{
	viper.AddConfigPath("conf")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig();err != nil{
		return err
	}
	return nil
}
