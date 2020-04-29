package conf

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func StartListener(appName string, amqpServer string, exchangeName string) {
	err := NewConsumer(amqpServer, exchangeName, "topic", "springCloudBus", exchangeName, appName)
	if err != nil {
		log.Fatalf("%s", err)
	}

	log.Printf("running forever")
	select {} // Yet another way to stop a Goroutine from finishing...
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func NewConsumer(amqpURI, exchange, exchangeType, queue, key, ctag string) error {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error

	log.Printf("dialing %s", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}

	log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring Exchange (%s)", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	log.Printf("declared Exchange, declaring Queue (%s)", queue)
	state, err := c.channel.QueueDeclare(
		queue, // name of the queue
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}

	log.Printf("declared Queue (%d messages, %d consumers), binding to Exchange (key '%s')",
		state.Messages, state.Consumers, key)

	if err = c.channel.QueueBind(
		queue,    // name of the queue
		key,      // bindingKey
		exchange, // sourceExchange
		false,    // noWait
		nil,      // arguments
	); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag '%s')", c.tag)
	deliveries, err := c.channel.Consume(
		queue, // name
		c.tag, // consumerTag,
		false, // noAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Consume: %s", err)
	}

	go handle(deliveries, c.done)

	return nil
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

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf(
			"got %dB consumer: [%v] delivery: [%v] routingkey: [%v] %s",
			len(d.Body),
			d.ConsumerTag,
			d.DeliveryTag,
			d.RoutingKey,
			d.Body,
		)
		handleRefreshEvent(d.Body, d.ConsumerTag)
		d.Ack(false)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}

// {"type":"RefreshRemoteApplicationEvent","timestamp":1494514362123,"originService":"config-server:docker:8888","destinationService":"xxxaccoun:**","id":"53e61c71-cbae-4b6d-84bb-d0dcc0aeb4dc"}
type UpdateToken struct {
	Type               string `json:"type"`
	Timestamp          int    `json:"timestamp"`
	OriginService      string `json:"originService"`
	DestinationService string `json:"destinationService"`
	Id                 string `json:"id"`
}
