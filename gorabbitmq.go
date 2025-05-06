package myrabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 错误处理函数
func IfError(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %s\n", err, msg)
	}
}

// 面向对象编程： ====================================================  Start
type RabbitMqClass struct {
	Conn  *amqp.Connection
	Ch    *amqp.Channel
	Queue amqp.Queue
}

// 函数：初始化，包括连接、生成channel
func (rabbit *RabbitMqClass) Init(user string, passwd string, url string, port string) (err error) {
	amqp_url := "amqp://" + user + ":" + passwd + "@" + url + ":" + port + "/"
	rabbit.Conn, err = amqp.Dial(amqp_url)
	IfError("Failed to connect to RabbitMQ.", err)
	rabbit.Ch, err = rabbit.Conn.Channel()
	IfError("Failed to open a channel.", err)
	return
}

// 函数：配置  ;  实际上关键是 Bind
func (rabbit *RabbitMqClass) Setup(exchange, exchange_type, routing_key, queue_name string) (err error) {
	if queue_name == "" || rabbit.Queue.Name != queue_name {
		err = rabbit.CreateQueue(queue_name)
		IfError("In Setup(): Fail to create queue.", err)
	}

	// 不是默认的，或者 ""，就新建一个。
	if exchange != "" && exchange != "amq.topic" && exchange != "amq.direct" && exchange != "amq.headers" &&
		exchange != "amq.fanout" && exchange != "amq.match" && exchange != "amq.rabbitmq.trace" {
		err = rabbit.CreateExchange(exchange, exchange_type)
		IfError("In Setup(): Fail to create exchange.", err)
	}

	if exchange != "" {
		err = rabbit.Ch.QueueBind(
			rabbit.Queue.Name, // queue name
			routing_key,       // routing key
			exchange,          // exchange
			false,             // no-wait
			nil,               // args
		)
		IfError("In Setup(): Failed to bind a queue.", err)
	}

	return err
}

// 函数：快速创建 Exchange
func (rabbit *RabbitMqClass) CreateExchange(xname, xtype string) (err error) {
	log.Println("In CreateExchange().")
	err = rabbit.Ch.ExchangeDeclare(
		xname, // name 自定义的名称
		xtype, // type 类型："topic" "direct" "fanout" "headers"
		false, // durable 持久性
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	IfError("In CreateExchange().", err)
	return err
}

// 函数：快速创建Queue
func (rabbit *RabbitMqClass) CreateQueue(queue_name string) error {
	var err error
	rabbit.Queue, err = rabbit.Ch.QueueDeclare(
		queue_name, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	IfError("In CreateQueue().", err)
	return err
}

// 函数：快速创建Queue。增加 durable （持久化）
func (rabbit *RabbitMqClass) CreateQueueDurable(queue_name string, durable bool) error {
	var err error
	rabbit.Queue, err = rabbit.Ch.QueueDeclare(
		queue_name, // name
		durable,    // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	IfError("In CreateQueueDurable().", err)
	return err
}

// 函数：不指定名字快速创建 Queue ，并返回随机 Queue 名字
func (rabbit *RabbitMqClass) CreateQueueReturnName() (string, error) {
	q, err := rabbit.Ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	IfError("In CreateQueueReturnName().", err)
	return q.Name, err
}

// 函数：Bind
func (rabbit *RabbitMqClass) Bind(exchange, exchange_type, routing_key, queue_name string) (err error) {
	if exchange != "" {
		err = rabbit.Ch.QueueBind(
			queue_name,  // queue name
			routing_key, // routing key
			exchange,    // exchange
			false,       // no-wait
			nil,         // args
		)
		IfError("In Bind(): Failed to bind a queue.", err)
	}
	return err
}

// 函数：接收消息：queue
func (rabbit *RabbitMqClass) ReceiveQueue() (<-chan amqp.Delivery, error) {
	msgs, err := rabbit.Ch.Consume(
		rabbit.Queue.Name, // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	IfError("In ReceiveQueue().", err)
	return msgs, err
}

// 函数：接收消息，具备参数。。。
func (rabbit *RabbitMqClass) Receive(queue_name string, auto_ack bool) (<-chan amqp.Delivery, error) {
	msgs, err := rabbit.Ch.Consume(
		queue_name, // queue
		"",         // consumer
		auto_ack,   // auto-ack: true
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	IfError("In ReceiveQueue().", err)
	return msgs, err
}

// 函数：发送信息：根据 exchange, routing_key
func (rabbit *RabbitMqClass) Send(exchange, routing_key, msg string) (err error) {
	err = rabbit.Ch.Publish(
		exchange,    // exchange: "amq.topic"  默认的
		routing_key, // routing key: "echo" or other...
		false,       // mandatory: 命令的, 托管的 ?
		false,       // immediate 立即的
		amqp.Publishing{
			Expiration:  "5000", // timeout : 5s
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	IfError("In Send().", err)
	return
}

// 函数：关闭连接和channel
func (rabbit *RabbitMqClass) Close() {
	rabbit.Conn.Close()
	rabbit.Ch.Close()
}

// 面向对象编程： ====================================================  End
