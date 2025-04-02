package main

import (
	"log"
	"net"
	output "outputservice/grpc/server"
	"outputservice/internal/env"
	"outputservice/internal/mq"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type Application struct {
	output.UnimplementedOutputServiceServer
	Port string
	Mq   *amqp.Channel
}

func main() {
	mq := &mq.MQ{
		User: env.GetMqUser(),
		Pass: env.GetMqPassword(),
		Port: "5672",
	}

	app := &Application{
		Port: ":6971",
		Mq:   mq.ConnectToMq(),
	}

	lis, err := net.Listen("tcp", app.Port)
	if err != nil {
		log.Println(err)
	}

	s := grpc.NewServer()
	output.RegisterOutputServiceServer(s, app)

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
