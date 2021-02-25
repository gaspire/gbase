package base

import (
	"os"

	pb "github.com/gaspire/gbase/jobrequest"
	"github.com/gaspire/gbase/util"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// IJob is an interface for utron Jobs
type IJob interface {
	Handle(req *pb.JobRequest) error
}

// RabbitMq 消息类型
type RabbitMq struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
	Exchange   string
}

// Close close connection
func (me *RabbitMq) Close() {
	me.Connection.Close()
	me.Channel.Close()
}

// NewRabbitMq 初始化rabbitmq
func NewRabbitMq(url, queue string) *RabbitMq {
	conn, err := amqp.Dial(url)
	util.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		queue,    // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name, // queue name
		q.Name, // routing key
		q.Name, // exchange
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")

	return &RabbitMq{conn, ch, q, q.Name}
}

// Consume 处理消息
func (me *RabbitMq) Consume(job IJob, async bool) {
	defer func() {
		me.Close()
	}()
	var autoAck = true
	//异步消费
	if async == true {
		autoAck = false
	}

	msgs, err := me.Channel.Consume(
		me.Queue.Name, // queue
		"",            // consumer
		autoAck,       // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if flag := util.FailOnError(err, "Failed to register a consumer"); flag {
		return
	}

	log.WithFields(log.Fields{
		"exchange": me.Exchange,
		"queue":    me.Queue.Name,
	}).Info(" [*] Waiting for msgs. To exit press CTRL+C")

	//处理消息
	for d := range msgs {
		//protobuf解码
		preq := &pb.JobRequest{}

		err := proto.Unmarshal(d.Body, preq)
		util.FailOnError(err, "proto unmarshaling error")
		if async == true {
			go func(d amqp.Delivery) {
				job.Handle(preq)
				d.Ack(false)
			}(d)
		} else {
			job.Handle(preq)
		}
	}
}

// AmqpPublish 发布消息
func AmqpPublish(queue string, data *pb.JobRequest) (err error) {
	conn, err := amqp.Dial(os.Getenv("MQ_URL"))
	if ok := util.FailOnError(err, "Failed to connect to RabbitMQ"); ok {
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if ok := util.FailOnError(err, "Failed to open a channel"); ok {
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if ok := util.FailOnError(err, "Failed to declare a queue"); ok {
		return
	}

	body, _ := proto.Marshal(data)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	util.FailOnError(err, "Failed to publish")
	return
}
