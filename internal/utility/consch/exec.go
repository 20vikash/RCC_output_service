package consch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MqNode struct {
	Jid    string
	Output string
}

func execPython(node *conNode, containerNumber int) {
	defer PyDoneExec(containerNumber)

	containerName := fmt.Sprintf("rcc-python_runner-%v", containerNumber)

	code := node.Code

	cmd := exec.Command("docker", "exec", containerName, "python3", "-c", code)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		mq := &MqNode{
			Jid:    node.Jid,
			Output: stderr.String(),
		}

		pushToMq(context.Background(), mq, ch)
		return
	}

	mq := &MqNode{
		Jid:    node.Jid,
		Output: stdout.String(),
	}

	pushToMq(context.Background(), mq, ch)
}

func pushToMq(ctx context.Context, mq *MqNode, channel *amqp.Channel) {
	mqJson, err := json.Marshal(mq)
	if err != nil {
		log.Println(err)
	}

	q, err := channel.QueueDeclare(
		"result", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Println(err)
	}

	err = channel.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         mqJson,
		})
	if err != nil {
		log.Println(err)
	}

	log.Println(" [x] Sent")
}
