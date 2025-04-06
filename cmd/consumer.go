package main

import (
	"encoding/json"
	"log"
	"outputservice/internal/utility/consch"
)

type Mq struct {
	Code     string
	Jid      string
	Language string
}

func (app *Application) listenToQueue() {
	q, err := app.Mq.QueueDeclare(
		"code", // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		log.Println(err)
	}

	msgs, err := app.Mq.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Println(err)
	}

	go func() {
		r := &Mq{}

		for d := range msgs {
			json.Unmarshal(d.Body, r)

			log.Printf("Received a message: %s", d.Body)
			err := d.Ack(false)

			consch.AddToPythonQueue(r.Language, r.Code, r.Jid)

			if err != nil {
				log.Println("Failed to ack the message")
			}
		}
	}()
}
