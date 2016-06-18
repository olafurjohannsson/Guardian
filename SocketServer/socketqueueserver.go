package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/googollee/go-socket.io"
	"github.com/streadway/amqp"
)

var (
	uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	exchange     = flag.String("exchange", "test-exchange", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	queue        = flag.String("queue", "test-queue", "Ephemeral AMQP queue name")
	bindingKey   = flag.String("key", "test-key", "AMQP binding key")
	consumerTag  = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
	lifetime     = flag.Duration("lifetime", 5*time.Second, "lifetime of process before shutdown (0s=infinite)")
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

type ConnectedClient struct {
	ID        string
	UserAgent string
	Socket    socketio.Socket
}

type CrossOriginServer struct{}

var m map[string]ConnectedClient

func init() {
	flag.Parse()
}

func NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Consumer, error) {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error

	log.Printf("dialing %q", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	go func() {
		fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring Exchange (%q)", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	log.Printf("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, key)

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	go handle(deliveries, c.done)

	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

// Item ...
type Item struct {
	Data      interface{}
	ID        string
	TimeStamp time.Time
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {

		var item Item
		json.Unmarshal(d.Body, &item)

		MessageConnectedClients("update", string(d.Body))

		log.Printf(
			"got %dB delivery: [%v] %s",
			len(d.Body),
			d.DeliveryTag,
			item.TimeStamp.String(),
		)
		d.Ack(false)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}

func SocketIOIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the socket.io server written in Go"))
}

func MessageConnectedClients(channel string, msg string) {
	//fmt.Printf("Messaging connected clients on channel %s with msg %s\n", channel, msg.TimeStamp.String())
	if m != nil {
		for c := range m {
			var client ConnectedClient
			client = m[c]

			if client.Socket != nil {
				client.Socket.Emit(channel, msg)
			}
		}
	}
}

func (s *CrossOriginServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	if m == nil {
		m = make(map[string]ConnectedClient)
	}

	// Allowed headers
	allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

	// Check Origin header
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Credentials", "true") // Without this, error: Credentials flag is 'true', but the 'Access-Control-Allow-Credentials' header is ''. It must be 'true' to allow credentials.
		rw.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers", allowHeaders)
		rw.Header().Set("X-SocketIO-Server-Version", "0.2.2")
	}

	// OPTIONS verb not supported
	if req.Method == "OPTIONS" {
		return
	}

	// WebSockets and Polling are supported
	transport := []string{"websocket", "polling"}

	// Create new server instance
	server, err := socketio.NewServer(transport)
	if err != nil {
		log.Fatal(err)
	}

	// Event handler for each connection
	server.On("connection", func(so socketio.Socket) {
		fmt.Println("connection: " + so.Id())

		// create new connected client
		m[so.Id()] = ConnectedClient{
			ID:        so.Id(),
			UserAgent: so.Request().UserAgent(),
			Socket:    so,
		}
		
		go func() {
			for {
				
				so.Emit("update", time.Now())
				time.Sleep(time.Second * 5);
			}
		}()

		// Client
		so.On("update", func(msg string) {

			//MessageConnectedClients("update", msg)
		})

		so.On("disconnection", func(so socketio.Socket) {
			fmt.Println("disconntion: " + so.Id())
			delete(m, so.Id())
		})

	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	mux := http.NewServeMux()
	mux.Handle("/socket.io/", server)
	mux.HandleFunc("/", SocketIOIndex)

	mux.ServeHTTP(rw, req)
}

func main() {

/*
	log.Println("Starting up RabbitMQ consumer)")
	c, err := NewConsumer(*uri, *exchange, *exchangeType, *queue, *bindingKey, *consumerTag)
	if err != nil {
		log.Fatalf("%s", err)
	}
*/

	// Start SocketServer
	log.Println("Listening on :5000 for socket.io")
	http.ListenAndServe(":5000", &CrossOriginServer{})

/*
	log.Printf("Shutting down")
	if err := c.Shutdown(); err != nil {
		log.Fatalf("error during shutdown: %s", err)
	}
*/
}
