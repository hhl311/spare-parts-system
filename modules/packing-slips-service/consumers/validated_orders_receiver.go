package consumers

import (
	"bytes"
	"encoding/json"
	"github.com/AntoineAube/spare-parts-system/modules/business-structures"
	"github.com/streadway/amqp"
	"log"
)

type ValidatedOrdersReceiver struct {
	ChannelName    string
	BusLocation    string
	BusCredentials string
}

func (receiver *ValidatedOrdersReceiver) LaunchAcknowledgment(consumer func(models.Order)) {
	conn, err := amqp.Dial("amqp://" + receiver.BusCredentials + "@" + receiver.BusLocation + "/")

	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		err := conn.Close()

		if err != nil {
			log.Println(err)
		}
	}()

	ch, err := conn.Channel()
	defer func() {
		err := ch.Close()

		if err != nil {
			log.Println(err)
		}
	}()

	q, err := ch.QueueDeclare(
		receiver.ChannelName,
		true,
		false,
		false,
		false,
		nil,
	)

	err = ch.Qos(
		1,
		0,
		false,
	)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	log.Println("Listening bus")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			order, err := deserialize(d.Body)

			if err == nil {
				consumer(order)

				if err := d.Ack(false); err != nil {
					log.Println(err)
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func deserialize(b []byte) (models.Order, error) {
	var msg models.Order
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}
