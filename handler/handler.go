package handler

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"goanalysis/pkg/errno"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 返回json 格式
func SendResponse(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)

	//总是返回http状态ok
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})

}

// 返回html 格式
func SendResponseHtml(c *gin.Context, err error, data string) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	//总是返回http状态ok
	c.String(http.StatusOK, data)
}

// http请求
func HttpRequest(api string, json string, method string) (string, error) {
	jsonStr := []byte(json)
	req, err := http.NewRequest(method, api, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json") //使用json格式传参

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errno.ApiServerError
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if !(resp.StatusCode == 200) {
		return "", errno.ApiServerError
	}
	return string(body), nil
}
