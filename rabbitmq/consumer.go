package logiarabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/yusologia/go-core/v2/model"
	"github.com/yusologia/go-core/v2/pkg"
	"log"
	"os"
	"strings"
	"time"
)

type RabbitMQConsumerInterface interface {
	Consume(message any) error
}

type RabbitMQConsumeOpt struct {
	Exchange string
	Queue    string
	Consumer RabbitMQConsumerInterface
}

type rabbitMQBody struct {
	MessageId uint `json:"messageId"`
	Data      any  `json:"data"`
}

func Consume(connection string, options []RabbitMQConsumeOpt) {
	if connection == "" || (connection != RABBITMQ_CONNECTION_GLOBAL && connection != RABBITMQ_CONNECTION_LOCAL) {
		log.Panicf("Please choose connection %s or %s", RABBITMQ_CONNECTION_GLOBAL, RABBITMQ_CONNECTION_LOCAL)
	}

	for _, opt := range options {
		if (opt.Exchange == "" && opt.Queue == "") || (opt.Exchange != "" && opt.Queue != "") {
			log.Panicf("Please select one of them: Exhange or Queue!!")
		}
	}

	mqConnection, ok := RabbitMQConnectionCache[connection]
	if !ok {
		if len(RabbitMQConnectionCache) == 0 {
			RabbitMQConnectionCache = make(map[string]logiamodel.RabbitMQConnection)
		}

		mqConnQuery := RabbitMQSQL.Where("connection = ?", connection)
		if connection == RABBITMQ_CONNECTION_LOCAL {
			mqConnQuery = mqConnQuery.Where("service = ?", os.Getenv("SERVICE"))
		}

		err := mqConnQuery.First(&mqConnection).Error
		if err != nil || mqConnection.ID == 0 {
			log.Panicf("Data connection does not exists: %s", err)
		}

		RabbitMQConnectionCache[connection] = mqConnection
	}

	connConf := RabbitMQConf.Connection[connection]
	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", connConf.Username, connConf.Password, connConf.Host, connConf.Port))
	if err != nil {
		log.Panicf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	var forever chan struct{}

	for _, opt := range options {
		if opt.Exchange != "" {
			fanoutConsumer(ch, mqConnection, opt)
		} else if opt.Queue != "" {
			directConsumer(ch, mqConnection, opt)
		}
	}

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func fanoutConsumer(ch *amqp091.Channel, connection logiamodel.RabbitMQConnection, opt RabbitMQConsumeOpt) {
	err := ch.ExchangeDeclare(
		opt.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panicf("Failed to declare exchange %s: %s", opt.Exchange, err)
	}

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Panicf("Failed to declare a queue: %s", err)
	}

	err = ch.QueueBind(
		q.Name,
		"",
		opt.Exchange,
		false,
		nil,
	)
	if err != nil {
		log.Panicf("Failed to bind queue %s to exchange %s: %s", q.Name, opt.Exchange, err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panicf("Failed to register a consumer: %s", err)
	}

	go func() {
		for d := range msgs {
			process(connection, opt, d.Body)
		}
	}()
}

func directConsumer(ch *amqp091.Channel, connection logiamodel.RabbitMQConnection, opt RabbitMQConsumeOpt) {
	q, err := ch.QueueDeclare(
		opt.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panicf("Failed to declare a queue: %s", err)
	}

	err = ch.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		log.Panicf("Failed to set QoS: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panicf("Failed to register a consumer: %s", err)
	}

	go func() {
		for d := range msgs {
			process(connection, opt, d.Body)
		}
	}()
}

func process(connection logiamodel.RabbitMQConnection, opt RabbitMQConsumeOpt, body []byte) {
	var consumerKey string
	if opt.Queue != "" {
		consumerKey = opt.Queue
	} else if opt.Exchange != "" {
		consumerKey = opt.Exchange
	} else {
		consumerKey = "CONSUMER_DOES_NOT_EXISTS"
	}

	log.Printf("CONSUMING: %s %s", printMessage(consumerKey), time.DateTime)

	var mqBody rabbitMQBody
	err := json.Unmarshal(body, &mqBody)
	if err != nil {
		logiapkg.LogError(fmt.Sprintf("Error unmarshalling: %s", err))
		return
	}

	var message logiamodel.RabbitMQMessage
	err = RabbitMQSQL.First(&message, mqBody.MessageId).Error
	if err != nil {
		failed(connection, opt, mqBody, fmt.Sprintf("Get message data: %s", err))
		return
	}

	err = opt.Consumer.Consume(mqBody.Data)
	if err != nil {
		failed(connection, opt, mqBody, fmt.Sprintf("Consume message is failed: %s", err))
		return
	}

	finish(message)

	log.Printf("%-10s %s %s", "SUCCESS:", printMessage(consumerKey), time.DateTime)
}

func finish(message logiamodel.RabbitMQMessage) {
	message.Finished = true

	err := RabbitMQSQL.Save(&message).Error
	if err != nil {
		logiapkg.LogError(fmt.Sprintf("Update message status is failed: %s", err))
	}
}

func failed(connection logiamodel.RabbitMQConnection, opt RabbitMQConsumeOpt, mqBody rabbitMQBody, message string) {
	logiapkg.LogError(message)

	payload, _ := json.Marshal(mqBody.Data)

	var messageFailed logiamodel.RabbitMQMessageFailed
	messageFailed.ConnectionId = connection.ID
	messageFailed.MessageId = mqBody.MessageId
	messageFailed.Service = os.Getenv("SERVICE")
	messageFailed.Exchange = opt.Exchange
	messageFailed.Queue = opt.Queue
	messageFailed.Payload = payload
	messageFailed.Exception = map[string]interface{}{"message": message, "trace": ""}

	err := RabbitMQSQL.Create(&messageFailed).Error
	if err != nil {
		logiapkg.LogError(fmt.Sprintf("Save message failed failed: %s", err))
	}
}

func printMessage(message string) string {
	paddedStr := fmt.Sprintf("%-60s", message)
	return strings.ReplaceAll(paddedStr, " ", ".")
}
