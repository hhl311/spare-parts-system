package consumers

import (
	"../../../business-structures"
	"bytes"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type ValidatedOrdersSender struct {
	ChannelName    string
	BusLocation    string
	BusCredentials string
}

func (sender *ValidatedOrdersSender) Send(order models.Order) (err error) {
	conn, err := amqp.Dial("amqp://" + sender.BusCredentials + "@" + sender.BusLocation + "/")
	defer func() {
		err := conn.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		log.Fatal(err)
		return
	}

	ch, err := conn.Channel()
	defer func() {
		err := ch.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		log.Fatal(err)
		return
	}

	q, err := ch.QueueDeclare(
		sender.ChannelName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
		return
	}

	body, _ := serialize(order)
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/octet-stream",
			Body:         []byte(body),
		})

	return
}

func serialize(msg models.Order) ([]byte, error) {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	err := encoder.Encode(msg)
	return b.Bytes(), err
}
