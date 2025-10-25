# GoAnalysis 埋点服务

## 项目简介
GoAnalysis 是一个基于 Go 语言开发的埋点数据收集与分析服务，支持通过 HTTP/HTTPS 接口接收前端或后端埋点数据，并将数据写入 Kafka 进行后续分析处理。

后续，我们将围绕 Go、.NET、JavaScript/Node.js 构建的多语言微服务，以及基于 Elasticsearch 的搜索场景，系统拆解微服务架构与大中型搜索方案的设计与落地。干货持续更新，敬请关注「代码扳手」微信公众号。：

![wx.jpg](wx.jpg)

## 主要功能
- 埋点数据收集接口（支持 HTTPS）
- 数据写入 Kafka
- 支持多种运行模式（debug/release）
- TLS/SSL 配置支持

## 埋点接口说明

### 1. 埋点数据上报接口
- **URL**: `/eventlog`  
- **方法**: `POST`
- **Content-Type**: `application/json`
- **请求示例**:

```json
{
  "event": "page_view",
  "user_id": "123456",
  "timestamp": 1719820800,
  "properties": {
    "page": "/home",
    "referrer": "https://www.example.com"
  }
}
```

- **参数说明**：
  - `event`：事件名称，如 `page_view`、`click` 等
  - `user_id`：用户唯一标识
  - `timestamp`：事件发生时间（Unix 时间戳，秒）
  - `properties`：事件属性，键值对

- **返回示例**：
```json
{
  "code": 0,
  "msg": "success"
}
```

### 2. 其他接口
如需扩展其他埋点接口，请参考 `handler/handler.go` 和 `router/router.go` 文件。

## 配置说明
配置文件位于 `conf/config.yaml`，主要参数如下：
- `server.addr`：监听地址及端口
- `server.tls`：是否启用 TLS
- `server.tlspem`/`tlskey`：证书路径
- `kafka.addr`：Kafka 地址
- `kafka.topic`：Kafka 主题

## 运行方式

1. 安装依赖：
   ```sh
   go mod tidy
   ```
2. 启动服务：
   ```sh
   go run analysis.go
   ```

## 目录结构
- `analysis.go`         主程序入口
- `conf/config.yaml`    配置文件
- `handler/`            业务处理逻辑
- `router/`             路由与中间件
- `service/`            服务层
- `model/`              数据模型
- `pkg/errno/`          错误码定义
- `logs/`               日志文件

## 联系方式
如有问题请联系维护者。 邮箱：26377227@qq.com
