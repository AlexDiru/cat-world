package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/AlexDiru/cat-world/catworldpb"
	"github.com/AlexDiru/cat-world/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	catworldpb.UnimplementedCatWorldServiceServer
}

var cat = common.Cat{
	Location: common.Location{
		X: 200,
		Y: 300,
	},
}

var cats = []common.Cat{cat}
var gameWorld = common.GameWorld{
	Cats: cats,
}

var r = rand.New(rand.NewSource(99))

func SimulateGameWorld() {
	time.Sleep(1 * time.Second)

	randomLocation := common.Location{
		X: int(r.Float64() * 600),
		Y: int(r.Float64() * 480),
	}

	gameWorld.Cats[0].Location = common.Location{
		X: randomLocation.X,
		Y: randomLocation.Y,
	}
}

func (*server) Connect(ctx context.Context, req *catworldpb.ConnectRequest) (*catworldpb.ConnectResponse, error) {
	// No auth so it just works

	return &catworldpb.ConnectResponse{
		Success: true,
	}, nil
}

func (*server) GetGameState(ctx context.Context, req *catworldpb.GetGameStateRequest) (*catworldpb.GetGameStateResponse, error) {

	location := gameWorld.Cats[0].Location

	locations := []*catworldpb.GetGameStateResponse_Location{
		{
			X: int32(location.X),
			Y: int32(location.Y),
		},
	}

	return &catworldpb.GetGameStateResponse{
		CatLocations: locations,
	}, nil
}

func main() {

	go func() {
		for {
			SimulateGameWorld()
		}
	}()

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
