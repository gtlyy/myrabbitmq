package myrabbitmq

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试用的主类
var rabbit *RabbitMqClass

// 测试：Init()。先单独看看连接和创建channel是否出错。
func TestInit(t *testing.T) {
	rabbit1 := &RabbitMqClass{}

	err := rabbit1.Init("guest", "guest", "localhost", "5672")
	if err != nil {
		t.Errorf("Failed to initialize RabbitMQ: %s", err)
	}
	if rabbit1.Conn == nil {
		t.Error("Failed to establish connection to RabbitMQ")
	}
	if rabbit1.Ch == nil {
		t.Error("Failed to open a channel")
	}

	// Clean up
	rabbit1.Close()
}

// 测试入口
func TestMain(m *testing.M) {
	fmt.Println("=================begin")

	rabbit = &RabbitMqClass{}
	rabbit.Init("userdw01", "pq328hu7", "127.0.0.1", "5672")
	// 运行测试函数
	m.Run()

	fmt.Println("=================end")
	rabbit.Close()
}

// 测试：CreateQueue()
func TestCreateQueue(t *testing.T) {
	queue_name := "test_queue" // "" 会生成一个唯一的随机queue
	err := rabbit.CreateQueue(queue_name)
	assert.True(t, err == nil)
	t.Log(queue_name)
}

// 测试：ReceiveQueue()，通过queue接收消息。
func TestReceiveQueue(t *testing.T) {
	queue_name := "hello_queue"
	err := rabbit.CreateQueue(queue_name)
	assert.True(t, err == nil)

	msgs, err := rabbit.ReceiveQueue()
	assert.True(t, err == nil)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {}
}

// 测试：通用函数Send()
func TestSendQueue(t *testing.T) {
	exchange := ""
	routing_key := "hello_queue" // queue_name
	msg := "Hello queue!"
	err := rabbit.Send(exchange, routing_key, msg)
	assert.True(t, err == nil)
}

// 测试：接受消息的通用写法1
func TestReceive1(t *testing.T) {
	exchange := "amq.topic"
	exchange_type := "topic"
	routing_key := "key1"
	queue_name := "queue1"
	err := rabbit.Setup(exchange, exchange_type, routing_key, queue_name)
	assert.True(t, err == nil)

	msgs, err := rabbit.ReceiveQueue()
	assert.True(t, err == nil)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {}
}

// 测试：配合上面
func TestSend1(t *testing.T) {
	exchange := "amq.topic"
	routing_key := "key1"
	msg := "exchange=amq.topic, routing_key=key1, queue_name=queue1"
	err := rabbit.Send(exchange, routing_key, msg)
	assert.True(t, err == nil)
}

// 测试：接受消息的通用写法2
func TestReceive2(t *testing.T) {
	// 第一步：
	exchange := "ex2"
	exchange_type := "topic"
	routing_key := "key2"
	queue_name := "queue2"
	err := rabbit.Setup(exchange, exchange_type, routing_key, queue_name)
	assert.True(t, err == nil)

	// 第二步：
	msgs, err := rabbit.ReceiveQueue()
	assert.True(t, err == nil)

	//第三步：
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {}
}

// 测试：配合上面
func TestSend2(t *testing.T) {
	exchange := "ex2"
	routing_key := "key2"
	msg := "exchange=ex2, routing_key=key2, queue_name=queue2"
	err := rabbit.Send(exchange, routing_key, msg)
	assert.True(t, err == nil)
}

// 测试：接受消息的通用写法3，实际是 queue
func TestReceive3(t *testing.T) {
	// 第一步：
	exchange := ""
	exchange_type := ""
	routing_key := "key3"
	queue_name := "queue3"
	err := rabbit.Setup(exchange, exchange_type, routing_key, queue_name)
	assert.True(t, err == nil)

	// 第二步：
	msgs, err := rabbit.ReceiveQueue()
	assert.True(t, err == nil)

	//第三步：
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {}
}

// 测试：配合上面
func TestSend3(t *testing.T) {
	exchange := ""
	routing_key := "queue3" // queue的时候，routing_key = queue_name
	msg := "exchange=, routing_key=queue3"
	err := rabbit.Send(exchange, routing_key, msg)
	assert.True(t, err == nil)
}

// 测试：send：  ~/code/echarts/stomp1.html
func TestSendStomp1(t *testing.T) {
	exchange := "amq.topic"
	key := "echo"
	msg := "hello world."
	err := rabbit.Send(exchange, key, msg)
	assert.True(t, err == nil)

	err = rabbit.Send(exchange, "test", "testing...")
	assert.True(t, err == nil)
}

// 测试：receive：  ~/code/echarts/stomp2.html
func TestReceive4(t *testing.T) {
	// 第一步：
	exchange := "amq.topic"
	exchange_type := "topic"
	routing_key := "test"
	queue_name := ""
	err := rabbit.Setup(exchange, exchange_type, routing_key, queue_name)
	assert.True(t, err == nil)

	// 第二步：
	msgs, err := rabbit.ReceiveQueue()
	assert.True(t, err == nil)

	//第三步：
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {}
}

// 测试：send：  ~/code/echarts/stomp3.html
func TestSendStomp3(t *testing.T) {
	exchange := "amq.topic"
	key := "account"
	msg := `{"usdt":` + `"1000"` + `,"coin":` + `"2"` + `}`
	err := rabbit.Send(exchange, key, msg)
	assert.True(t, err == nil)

	err = rabbit.Send(exchange, "test", "testing...")
	assert.True(t, err == nil)

	err = rabbit.Send(exchange, "trade", "enable_open_button")
	assert.True(t, err == nil)

	err = rabbit.Send(exchange, "trade", "enable_close_button")
	assert.True(t, err == nil)
}

// 测试：receive：  ~/code/echarts/stomp3.html
func TestReceiveTradeReply(t *testing.T) {
	// 第一步：
	exchange := "amq.topic"
	exchange_type := "topic"
	routing_key := "tradeReply"
	queue_name := ""
	err := rabbit.Setup(exchange, exchange_type, routing_key, queue_name)
	assert.True(t, err == nil)

	// 第二步：
	msgs, err := rabbit.ReceiveQueue()
	assert.True(t, err == nil)

	//第三步：
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {}
}

// 测试：send：  ~/code/echarts/lines-async-simple.html
func TestSendLine(t *testing.T) {
	exchange := "amq.topic"
	routing_key := "priceAB"
	msg := "0.09"
	err := rabbit.Send(exchange, routing_key, msg)
	assert.True(t, err == nil)
}
