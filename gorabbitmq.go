package myrabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func OnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// 面向对象编程： ====================================================  Start
type RabbitMqClass struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
	// exchange    string
	// routing_key string
	Queue amqp.Queue
}

// 函数：初始化，包括连接、生成channel
func (rabbit *RabbitMqClass) Init(user string, passwd string, url string, port string) (err error) {
	amqp_url := "amqp://" + user + ":" + passwd + "@" + url + ":" + port + "/"
	rabbit.Conn, err = amqp.Dial(amqp_url)
	OnError(err, "Failed to connect to RabbitMQ")
	rabbit.Ch, err = rabbit.Conn.Channel()
	OnError(err, "Failed to open a channel")
	return
}

// 函数：创建 Exchange
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
	OnError(err, "In CreateExchange().")
	return err
}

// 函数：创建Queue
func (rabbit *RabbitMqClass) CreateQueue(queue_name string) error {
	// log.Println("In CreateQueue1()")
	var err error
	rabbit.Queue, err = rabbit.Ch.QueueDeclare(
		queue_name, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	// log.Println(rabbit.Queue.Name)
	OnError(err, "In CreateQueue().")
	return err
}

// 函数：创建Queue
// durable：持久化
func (rabbit *RabbitMqClass) CreateQueue2(queue_name string, durable bool) error {
	log.Println("In CreateQueue2()")
	var err error
	rabbit.Queue, err = rabbit.Ch.QueueDeclare(
		queue_name, // name
		durable,    // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	log.Println("Create Queue:", rabbit.Queue.Name)
	OnError(err, "In CreateQueue2().")
	return err
}

// 函数：创建Queue ， 返回名字的版本
func (rabbit *RabbitMqClass) CreateQueue3() (string, error) {
	// log.Println("In CreateQueue3()")
	q, err := rabbit.Ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	// log.Println("Create Queue 3:", q.Name)
	OnError(err, "In CreateQueue3().")
	return q.Name, err
}

// 函数：配置  ;  实际上是 Bind
func (rabbit *RabbitMqClass) Setup(exchange, exchange_type, routing_key, queue_name string) (err error) {
	if queue_name == "" || rabbit.Queue.Name != queue_name {
		err = rabbit.CreateQueue(queue_name)
		OnError(err, "In Setup():Fail to create queue.")
	}

	// 不是默认的，或者 ""，就新建一个。
	if exchange != "" && exchange != "amq.topic" && exchange != "amq.direct" && exchange != "amq.headers" &&
		exchange != "amq.fanout" && exchange != "amq.match" && exchange != "amq.rabbitmq.trace" {
		err = rabbit.CreateExchange(exchange, exchange_type)
		OnError(err, "In Setup():Fail to create exchange.")
	}

	if exchange != "" {
		// 绑定
		err = rabbit.Ch.QueueBind(
			rabbit.Queue.Name, // queue name
			routing_key,       // routing key
			exchange,          // exchange
			false,             // no-wait
			nil,               // args
		)
		OnError(err, "In Setup(): Failed to bind a queue")
	}

	return err
}

// 函数：Bind
func (rabbit *RabbitMqClass) Bind(exchange, exchange_type, routing_key, queue_name string) (err error) {
	if exchange != "" {
		// 绑定
		err = rabbit.Ch.QueueBind(
			queue_name,  // queue name
			routing_key, // routing key
			exchange,    // exchange
			false,       // no-wait
			nil,         // args
		)
		OnError(err, "In Bind(): Failed to bind a queue")
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
	OnError(err, "In ReceiveQueue().")
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
	OnError(err, "In ReceiveQueue().")
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
	// log.Println("In Send(): Sending msg", msg)
	OnError(err, "In Send().")
	return
}

// 函数：关闭连接和channel
func (rabbit *RabbitMqClass) Close() {
	rabbit.Conn.Close()
	rabbit.Ch.Close()
}

// 面向对象编程： ====================================================  End

// func Connect(user string, passwd string, url string, port string) (conn *amqp.Connection, err error) {
// 	amqp_url := "amqp://" + user + ":" + passwd + "@" + url + ":" + port + "/"
// 	conn, err = amqp.Dial(amqp_url)
// 	// OnError(err, "Failed to connect to RabbitMQ")
// 	// defer conn.Close()
// 	return
// }

// func Send(conn *amqp.Connection, ch *amqp.Channel, topic string, routing_key string, msg string) error {
// 	// bug: 每次开一个channel ，用完就会出错：Exception (504) Reason: "channel id space exhausted"
// 	// ch, err := conn.Channel()
// 	// OnError(err, "Failed to open a channel")
// 	// defer ch.Close()
// 	err := ch.Publish(
// 		topic,       // exchange: "amq.topic"
// 		routing_key, // routing key
// 		false,       // mandatory
// 		false,       // immediate
// 		amqp.Publishing{
// 			ContentType: "text/plain",
// 			Body:        []byte(msg),
// 		})
// 	OnError(err, "Failed to publish a message")
// 	log.Printf(" [x] Sent %s", msg)
// 	return err
// }

// func bodyFrom(args []string) string {
// 	var s string
// 	if (len(args) < 3) || os.Args[2] == "" {
// 		// s = "hello22"
// 		s = `{"xx":["2022-05-03", "2022-04-04"], "yy":[[25,45,42,35], [31,23,44,67]]}`
// 	} else {
// 		s = strings.Join(args[2:], " ")
// 	}
// 	return s
// }

// func severityFrom(args []string) string {
// 	var s string
// 	if (len(args) < 2) || os.Args[1] == "" {
// 		s = "echo"
// 	} else {
// 		s = os.Args[1]
// 	}
// 	return s
// }
