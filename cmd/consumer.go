package main

import "log"

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
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			err := d.Ack(false)

			if err != nil {
				log.Println("Failed to ack the message")
			}
		}
	}()
}
