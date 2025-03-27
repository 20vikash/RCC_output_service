package main

import (
	"log"
	"net"
	output "outputservice/grpc/server"

	"google.golang.org/grpc"
)

type Server struct {
	output.UnimplementedOutputServiceServer
}

func main() {
	lis, err := net.Listen("tcp", ":6971")
	if err != nil {
		log.Println(err)
	}

	s := grpc.NewServer()
	output.RegisterOutputServiceServer(s, &Server{})

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
