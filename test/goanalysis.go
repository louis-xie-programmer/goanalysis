package testanalysis

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// 创建过程数据的存储容器
type Queue struct {
	Maxsize int
	array   [5]EventDbLog //数组
	head    int           // 队首
	tail    int           // 队尾
}

// 添加数据到容器
func (s *Queue) Push(value EventDbLog) error {
	if s.IsFull() {
		log.Println("queue full")
		return errors.New("queue full")
	}
	s.array[s.tail] = value
	s.tail = (s.tail + 1) % s.Maxsize
	return nil
}

// 出队列
func (s *Queue) Pop() (val EventDbLog, err error) {
	if s.IsEmpty() {
		log.Println("queue empty")
		return val, errors.New("queue empty")
	}

	val = s.array[s.head]
	s.head = (s.head + 1) % s.Maxsize
	return val, nil
}

func (s *Queue) Shows() {
	size := s.Size()
	if size == 0 {
		log.Println("queue is empty")
	}
	tempHead := s.head
	for i := 0; i < size; i++ {
		log.Println(s.array[tempHead])
		tempHead = (tempHead + 1) % s.Maxsize
	}
}

// 是否是空
func (s *Queue) IsEmpty() bool {
	return s.tail == s.head
}

// 是否满
func (s *Queue) IsFull() bool {
	return (s.tail+1)%s.Maxsize == s.head
}

// 统计队列中的个数
func (s *Queue) Size() int {
	return (s.tail + s.Maxsize - s.head) % s.Maxsize
}

type EventDbLog struct {
	LogContent string
	Type       string
}

// 页面浏览日志结构
// vtype: 0,1,2 window.performance.navigation.type
type PageViewlog struct {
	Id      string
	WebName string
	Host    string
	Url     string
	UA      string
	Lang    string
	IP      string
	Sid     string
	Index   int
	Ref     string
	Vtype   int
	Title   string
	Uuid    string
	ScreenH int
	ScreenW int
	T       int64
}

// 事件日志结构
type Eventlog struct {
	Id      string
	WebName string
	PVId    string //pageview id
	EType   string
	Url     string
	JsonDb  string
	T       int64
}

// 整理日志字段（pageview）
func GetPageViewLogDB(c *gin.Context) EventDbLog {
	elog := PageViewlog{}

	if err := c.ShouldBindJSON(&elog); err != nil {
		log.Fatalln("GetPageViewLogDB ShouldBind error:", err)
	}

	//elog.Url = UrlBase64Decode(elog.Url)
	elog.WebName = c.Param("webname")
	elog.Host = c.Request.Host
	elog.UA = c.Request.Header["User-Agent"][0]
	elog.Lang = c.Request.Header["Accept-Language"][0]
	elog.IP = GetRealIP(c)
	elog.Sid = GetSession(c, "sessionID").(string) //会话标记
	elog.Index = GetSession(c, "index").(int)      //访问顺序标记
	elog.T = time.Now().UnixMilli()
	//elog.Ref = UrlBase64Decode(elog.Ref)
	//elog.Title = UrlBase64Decode(elog.Title)

	b, _ := json.Marshal(elog)

	log := EventDbLog{
		LogContent: string(b),
		Type:       "pageview",
	}

	return log
}

// 整理事件日志（pageview）
func GetEventLogDB(c *gin.Context) EventDbLog {
	elog := Eventlog{}

	if err := c.ShouldBindJSON(&elog); err != nil {
		log.Fatalln("GetEventLogDB ShouldBind error:", err)
	}

	logid, _ := uuid.NewRandom()

	elog.Id = logid.String()
	elog.WebName = c.Param("webname")
	elog.T = time.Now().UnixMilli()

	b, _ := json.Marshal(elog)

	log := EventDbLog{
		LogContent: string(b),
		Type:       "eventlog",
	}

	return log
}

// 获取真实ip
func GetRealIP(c *gin.Context) string {
	ip := c.Request.Header.Get("x-Real-IP")
	if ip == "" {
		ip = c.Request.Header.Get("x-Forwarded-For")
	}
	if ip == "" {
		ip = c.Request.RemoteAddr
	}
	return ip
}

// urlbase64转义解析
// func UrlBase64Decode(input string) string {
// 	input = strings.ReplaceAll(input, "_", "/")
// 	input = strings.ReplaceAll(input, "-", "=")
// 	output, _ := base64.StdEncoding.DecodeString(input)
// 	return string(output)
// }

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

// 自定义头部信息（包含sid）
func SetHeaderDb(c *gin.Context) {
	sessionID := GetSession(c, "sessionID")
	c.Writer.Header().Set("sid", sessionID.(string))
}

// 设置浏览序号
func SetPageViewIndex(c *gin.Context, isnew bool) {
	s := sessions.Default(c)
	sindex := s.Get("index")
	if sindex == nil || isnew {
		s.Set("index", 1)
		s.Save()
	} else {
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

// kafka 写入
// func producerHandler(kafkaconn *kafka.Conn, key string, value string) {
// 	msg := kafka.Message{
// 		Key:   []byte(key),
// 		Value: []byte(value),
// 	}

// 	_, err := kafkaconn.WriteMessages(msg)

// 	if err != nil {
// 		log.Println("kafka WriteMessages error:" + err.Error())

// 		kafkaconn, _ = kafka.DialLeader(context.Background(), "tcp", "192.168.127.160:9092", "weblogs", 0)
// 		_, err := kafkaconn.WriteMessages(msg)

// 		if err != nil {
// 			log.Println("kafka WriteMessages 2 error:" + err.Error())
// 		}
// 	}
// }

type Client struct {
	host     string      // kafka地址,一般是:localhost:9092
	instance *kafka.Conn //
}

// Connect 连接kafka服务
func (client *Client) Connect(topic string) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", client.host, topic, 0)
	if err != nil {
		log.Println("failed to dial learder: " + err.Error())
	}

	//_ = conn.SetWriteDeadline(time.Now().Add(120 * time.Second))

	client.instance = conn
}

func (producer Producer) Send(key string, value string) {

	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}

	_, err := producer.client.instance.WriteMessages(msg)

	if err != nil {
		log.Println("failed to write message: " + err.Error())
		producer.client.Connect("weblogs")
		_, err1 := producer.client.instance.WriteMessages(msg)
		if err1 != nil {
			log.Println("failed to write 2 message: " + err1.Error())
		}
	}
}

// NewClient 实例化
func NewClient(host string) *Client {
	return &Client{host: host}
}

// Producer 生产者
type Producer struct {
	client *Client
}

// 入口函数
func main() {
	router := gin.Default()

	//跨域处理(线上处理时一定要填写正确的授权orgins)
	corsconfig := cors.DefaultConfig()
	//config.AllowOrigins = []string{"www.laiys.com"}
	corsconfig.AllowAllOrigins = true
	router.Use(cors.New(corsconfig))

	//注入会话存储
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	//kafka 连接
	kafkaClient := NewClient("192.168.127.160:9092")

	kafkaClient.Connect("weblogs")

	// kafkaconn, kerr := kafka.DialLeader(context.Background(), "tcp", "192.168.127.160:9092", "weblogs", 0)

	// if kerr != nil {
	// 	log.Fatal("kafka conn error:", kerr)
	// }

	// defer kafkaconn.Close()

	//kafkaWriter := getKafkaWriter("192.168.127.160:9092", "weblogs")

	//defer kafkaWriter.Close()

	//页面浏览组
	pView := router.Group("/api/v1/pageview")
	//页面事件组
	pEvent := router.Group("/api/v1/event")

	producer := Producer{
		client: kafkaClient,
	}

	var queues = Queue{}

	go func() {
		for {
			item, err := queues.Pop()
			if err == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			b, _ := json.Marshal(item)
			producer.Send("pageview", string(b))
		}
	}()

	//访问记录收集接口 webname: 客户端名称
	pView.POST(":webname",
		func(c *gin.Context) {
			//初始化会话
			isnew := InitSession(c)
			//设置会话访问序号
			SetPageViewIndex(c, isnew)

			elog := GetPageViewLogDB(c)

			err := queues.Push(elog)
			if err != nil {
				log.Println("queues push error:" + elog.LogContent)
			}

			//b, _ := json.Marshal(elog)
			//producer.Send("pageview", string(b))
			c.Status(http.StatusOK)
		},
	)

	//页面性能数据收集接口(lcp收集)
	pEvent.POST("analytics/:webname",
		func(c *gin.Context) {

			elog := GetEventLogDB(c)

			err := queues.Push(elog)
			if err != nil {
				log.Println("queues push error:" + elog.LogContent)
			}

			//b, _ := json.Marshal(elog)

			//producer.Send("eventlog", string(b))

			c.Status(http.StatusOK)
		},
	)

	//router.Run(":6001")
	router.RunTLS("www.laiys.com:443", "www.laiys.com.pem", "www.laiys.com.key")
}
