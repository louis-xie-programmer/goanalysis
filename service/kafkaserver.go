package service

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type Client struct {
	host     string      // kafka地址
	instance *kafka.Conn //
}

// Connect 连接kafka服务
func (client *Client) Connect(topic string) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", client.host, topic, 0)
	if err != nil {
		log.Println("failed to dial learder: " + err.Error())
	}

	client.instance = conn
}

func (producer Producer) Send(key string, value string) {

	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}

	_, err := producer.client.instance.WriteMessages(msg)

	if err != nil {
		log.Println("reconn kafka error: " + err.Error())
		producer.client.Connect(viper.GetString("common.kafka.topic"))
		_, err1 := producer.client.instance.WriteMessages(msg)
		if err1 != nil {
			log.Println("failed to write 2 message: " + err1.Error())
		}
	}
}

func NewProducer() Producer {
	//kafka 连接
	kafkaClient := NewClient(viper.GetString("common.kafka.addr"))

	kafkaClient.Connect(viper.GetString("common.kafka.topic"))

	producer := Producer{
		client: kafkaClient,
	}

	return producer
}

// NewClient 实例化
func NewClient(host string) *Client {
	return &Client{host: host}
}

// Producer 生产者
type Producer struct {
	client *Client
}
