package main

import (
	"context"
	"fmt"
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

var cats = []common.Cat{}
var gameWorld = common.GameWorld{
	Cats: cats,
}

var r = rand.New(rand.NewSource(99))

func SimulateGameWorld() {
	time.Sleep(1 * time.Second)
	fmt.Printf("Simulating Game World with %v cats\n", len(gameWorld.Cats))

	gameWorld.Mutex.Lock()
	defer gameWorld.Mutex.Unlock()

	for catIndex := range gameWorld.Cats {
		randomLocation := GenerateRandomLocation()
		gameWorld.Cats[catIndex].Location = common.Location{
			X: randomLocation.X,
			Y: randomLocation.Y,
		}
	}
}

func GenerateRandomLocation() common.Location {
	return common.Location{
		X: int(r.Float64() * 600),
		Y: int(r.Float64() * 480),
	}
}

func (*server) Connect(ctx context.Context, req *catworldpb.ConnectRequest) (*catworldpb.ConnectResponse, error) {
	gameWorld.Mutex.Lock()
	defer gameWorld.Mutex.Unlock()

	// No auth so it just works

	// Add a cat
	newCat := common.Cat{
		Location: GenerateRandomLocation(),
	}
	gameWorld.Cats = append(gameWorld.Cats, newCat)

	return &catworldpb.ConnectResponse{
		Success: true,
	}, nil
}

func (*server) GetGameState(ctx context.Context, req *catworldpb.GetGameStateRequest) (*catworldpb.GetGameStateResponse, error) {
	gameWorld.Mutex.Lock()
	defer gameWorld.Mutex.Unlock()

	locations := make([]*catworldpb.GetGameStateResponse_Location, len(gameWorld.Cats))

	for catIndex := range gameWorld.Cats {
		cat := gameWorld.Cats[catIndex]

		locations[catIndex] = &catworldpb.GetGameStateResponse_Location{
			X: int32(cat.Location.X),
			Y: int32(cat.Location.Y),
		}
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
