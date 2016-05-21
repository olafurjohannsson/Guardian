package QueueSender

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

// Queue ...
type Queue struct {
	TimeStamp  time.Time
	Id         string
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// Item ...
type Item struct {
	Data      interface{}
	ID        string
	TimeStamp time.Time
}

// Init a queue
func Init(items chan Item, key interface{}) {

	go func() {
		fmt.Println("go func() in sender.go")

		//m := make(map[int]Item)

		// register interface
		gob.Register(key)
		//counter := 0
		for {
			select {

			// item popped off the queue
			case item := <-items:
				fmt.Printf("Item received %s\n", item)

				// m[counter] = item
				// counter = counter + 1

				// if len(m) >= 1 {
				// 	for i := 0; i < counter; i++ {
				// 		Publish(q, m[i])
				// 	}
				// 	counter = 0
				// 	m = make(map[int]Item)
				// }

				// publish item to queue

				Publish(item)

				// if err != nil {
				// 	fmt.Printf("publish err %s\n", err)
				//}
			}
		}
	}()

}

// Publish an item to amqp
func Publish(key interface{}) error {
	exchange := "test-exchange"
	routingKey := "test-key"
	exchangeType := "direct"

	var url string
	url = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}

	// get channel
	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	if err := channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	data, err := json.Marshal(key)

	if err != nil {
		fmt.Printf("err marshalling json: %s\n", err)
	}

	if err := channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            data,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}
