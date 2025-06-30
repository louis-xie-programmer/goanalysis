package service

import (
	. "goanalysis/handler"

	"github.com/gin-gonic/gin"
)

// 首页
func Index(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>hello world</title>
</head>
<body>
    hello world
</body>
</html>
`
	SendResponseHtml(c, nil, html)
}
