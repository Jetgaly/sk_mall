package main

// import (
// 	"fmt"
// 	utils "sk_mall/utils/RabbitMQ"
// 	"time"

// 	amqp "github.com/rabbitmq/amqp091-go"
// )

// func main() {
// 	r, err := utils.NewRMQ("amqp://admin:admin@localhost:5672/", 3, 6, 3)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		panic("err")
// 	}
// 	c, _ := r.Get()
// 	// q, e1 := c.QueueDeclare(
// 	// 	"hello", // name
// 	// 	false,   // durable
// 	// 	false,   // delete when unused
// 	// 	false,   // exclusive
// 	// 	false,   // no-wait
// 	// 	nil,     // arguments
// 	// )
// 	// 在 channel 上开启发布确认
// 	e3 := c.Confirm(false) // false 表示不开启事务模式
// 	if e3 != nil {
// 		panic(e3.Error())
// 	}
// 	// 确认和未确认的消息通道
// 	confirms := c.NotifyPublish(make(chan amqp.Confirmation, 1))
// 	c.ExchangeDeclare(
// 		"ex1",
// 		"direct",
// 		false,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	e1 := c.QueueBind(
// 		"hello",    // 队列名
// 		"hellokey", // 绑定键
// 		"ex1",      // 交换机
// 		false,      // no-wait
// 		nil,        // args
// 	)
// 	if e1 != nil {
// 		fmt.Println(e1.Error())
// 		panic("err")
// 	}
// 	body := "Hello World2222!"
// 	// 4.将消息发布到声明的队列
// 	err = c.Publish(
// 		"ex1",      // exchange
// 		"hellokey", // routing key
// 		false,      // mandatory
// 		false,      // immediate
// 		amqp.Publishing{
// 			ContentType: "text/plain",
// 			Body:        []byte(body),
// 		})
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	} else {
// 		select {
// 		case confirm := <-confirms:
// 			if confirm.Ack {
// 				fmt.Println("Message confirmed")
// 			} else {
// 				fmt.Println("Message failed")
// 			}
// 		case <-time.After(5 * time.Second):
// 			fmt.Println("Confirm timeout")
// 		}
// 		fmt.Println("success")
// 	}
// }

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// func main() {
// 	conn, _ := amqp.Dial("amqp://admin:admin@localhost:5672/")
// 	ch, _ := conn.Channel()
// 	ch.Confirm(false) //用手动确认模式
// 	//When noWait is true, the client will not wait for a response. A channel exception could occur if the server does not support this method.

// 	confirm := ch.NotifyPublish(make(chan amqp.Confirmation, 2))
// 	confirm2 := ch.NotifyPublish(make(chan amqp.Confirmation, 2))
// 	ch.QueueDeclare(
// 		"confirmtest",
// 		false,
// 		false,
// 		false,
// 		false,
// 		nil)
// 	ch.PublishWithContext(
// 		context.Background(),
// 		"",
// 		"confirmtest",
// 		false,
// 		false,
// 		amqp.Publishing{
// 			ContentType: "text/plain",
// 			Body:        []byte("hello,rmq"),
// 		})

// 	select {
// 	case cf := <-confirm:
// 		if cf.Ack {
// 			fmt.Println("confirm1")
// 		}
// 	case <-time.After(5 * time.Second):
// 		fmt.Println("timeout")
// 	}

// 	select {
// 	case cf := <-confirm2:
// 		if cf.Ack {
// 			fmt.Println("confirm2")
// 		}
// 	case <-time.After(5 * time.Second):
// 		fmt.Println("timeout2")
// 	}
// }

func main() {
	conn, _ := amqp.Dial("amqp://admin:admin@localhost:5672/")
	ch, _ := conn.Channel()
	ch.Confirm(false)

	// 旧监听器：收到一个后不读（一直阻塞）
	confirmOld := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	// 并发发布多个消息，让旧监听器 buffer 撑满

	ch.PublishWithContext(context.Background(),
		"", "confirmtest", false, false,
		amqp.Publishing{ContentType: "text/plain", Body: []byte(fmt.Sprintf("msg-%d", 1))})

	// 新的监听器

	var confirmNew [10000]chan amqp.Confirmation
	for i := 0; i < 10000; i++ {
		confirmNew[i] = ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	}

	// 新发布一条，理论上 confirmNew 应该收到
	ch.PublishWithContext(context.Background(),
		"", "confirmtest", false, false,
		amqp.Publishing{ContentType: "text/plain", Body: []byte("msg-new")})

	// 看新监听器是否阻塞
	select {
	case <-confirmNew[1000]:
		fmt.Println("confirmNew received ACK ✅")
	case <-time.After(5 * time.Second):
		fmt.Println("confirmNew timeout ❌")
	}

	// 读旧监听器
	select {
	case <-confirmOld:
		fmt.Println("confirmOld finally consumed")
	default:
		fmt.Println("confirmOld still blocked")
	}

	//fmt.Println(confirmNew, confirmOld)
}
