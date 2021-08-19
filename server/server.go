package main

import (
	"context"
	"log"
	"net"

	"github.com/AlexDiru/cat-world/catworldpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	catworldpb.UnimplementedCatWorldServiceServer
}

func (*server) Connect(ctx context.Context, req *catworldpb.ConnectRequest) (*catworldpb.ConnectResponse, error) {
	// No auth so it just works

	return &catworldpb.ConnectResponse{
		Success: true,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	defer lis.Close()

	grpcServer := grpc.NewServer()

	catworldpb.RegisterCatWorldServiceServer(grpcServer, &server{})

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
