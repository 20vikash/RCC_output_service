package main

import (
	"log"
	"net"
	output "outputservice/grpc/server"

	"google.golang.org/grpc"
)

type Application struct {
	output.UnimplementedOutputServiceServer
	Port string
}

func main() {
	app := &Application{
		Port: ":6971",
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
