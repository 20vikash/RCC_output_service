package consch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MqNode struct {
	Jid    string
	Output string
}

type flaskResponse struct {
	Output   string `json:"output"`
	Error    string `json:"error"`
	ExitCode int    `json:"exit_code"`
}

func execPython(node *conNode, containerNumber int) {
	log.Println("Executing")

	defer PyDoneExec(containerNumber)

	containerHost := fmt.Sprintf("http://rcc-python_runner-%v:8000/run", containerNumber)

	payload := map[string]string{"code": node.Code}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(containerHost, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("HTTP error:", err)
		pushToMq(context.Background(), &MqNode{
			Jid:    node.Jid,
			Output: "Failed to contact Python runner",
		}, ch)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var res flaskResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Println("Failed to decode response:", err)
		pushToMq(context.Background(), &MqNode{
			Jid:    node.Jid,
			Output: "Failed to parse Python runner output",
		}, ch)
		return
	}

	output := res.Output
	if res.Error != "" {
		output = res.Error
	}

	mq := &MqNode{
		Jid:    node.Jid,
		Output: output,
	}

	pushToMq(context.Background(), mq, ch)
	log.Println("Output:", mq.Output)
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
