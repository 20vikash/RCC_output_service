package main

import (
	"log"
	"net"
	"outputservice/internal/env"
	"outputservice/internal/mq"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type Application struct {
	Port string
	Mq   *amqp.Channel
}

func main() {
	mq := &mq.MQ{
		User: env.GetMqUser(),
		Pass: env.GetMqPassword(),
		Port: "5672",
	}

	con := mq.ConnectToMq()

	app := &Application{
		Port: ":6971",
		Mq:   mq.CreateChannel(con),
	}

	lis, err := net.Listen("tcp", app.Port)
	if err != nil {
		log.Println(err)
	}

	s := grpc.NewServer()

	log.Println("Yes")
	app.listenToQueue()

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
