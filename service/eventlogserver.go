package service

import (
	"encoding/json"
	. "goanalysis/handler"
	. "goanalysis/model"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func EnventLogDo(c *gin.Context, producer Producer) {
	elog, err := GetEventLogDB(c)

	if err == nil {
		b, _ := json.Marshal(elog)
		producer.Send("eventlog", string(b))
	}
	SendResponse(c, nil, nil)
}

// 整理事件日志（pageview）
func GetEventLogDB(c *gin.Context) (Eventlog, error) {
	elog := Eventlog{}

	if err := c.ShouldBindJSON(&elog); err != nil {
		log.Println("GetEventLogDB ShouldBind error:", err)
		return elog, err
	}

	logid, _ := uuid.NewRandom()

	elog.Id = logid.String()
	elog.WebName = c.Param("webname")
	elog.Time = time.Now().Unix()

	return elog, nil
}
