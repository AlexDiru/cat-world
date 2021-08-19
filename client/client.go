package main

import (
	"context"
	"fmt"
	"log"
	"math"

	_ "image/png"

	"github.com/AlexDiru/cat-world/catworldpb"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"google.golang.org/grpc"
)

type Game struct {
	inited bool
	op     ebiten.DrawImageOptions
}

type Sprite struct {
	imageWidth  int
	imageHeight int
	x           int
	y           int
	vx          int
	vy          int
	angle       int
}

type RenderableCat struct {
	sprite *Sprite
	image  *ebiten.Image
}

func (r *RenderableCat) Draw(g *Game, screen *ebiten.Image) {
	w, h := r.image.Size()
	s := r.sprite
	g.op.GeoM.Reset()
	g.op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	g.op.GeoM.Rotate(2 * math.Pi * float64(s.angle) / 360)
	g.op.GeoM.Translate(float64(w)/2, float64(h)/2)
	g.op.GeoM.Translate(float64(s.x), float64(s.y))
	screen.DrawImage(r.image, &g.op)
}

func CreateCatSprite() *RenderableCat {
	return &RenderableCat{
		sprite: &Sprite{
			imageWidth:  32,
			imageHeight: 32,
			x:           0,
			y:           0,
			vx:          100,
			vy:          100,
			angle:       0,
		},
		image: LoadImage(),
	}
}

func LoadImage() *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile("assets/cat.png") //image.Decode(bytes.NewReader(images.Ebiten_png))
	if err != nil {
		fmt.Println("Failed2")
		log.Fatal(err)
	}
	origEbitenImage := ebiten.NewImageFromImage(img)

	w, h := origEbitenImage.Size()
	ebitenImage := ebiten.NewImage(w, h)

	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 0.5)
	ebitenImage.DrawImage(origEbitenImage, op)

	return ebitenImage
}

func (g *Game) Init() {
	defer func() {
		g.inited = true
	}()

}

func (g *Game) Draw(screen *ebiten.Image) {
	sprite := CreateCatSprite()
	sprite.Draw(g, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func (g *Game) Update() error {
	//if !g.inited {
	g.Init()
	//}

	return nil
}

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer conn.Close()

	client := catworldpb.NewCatWorldServiceClient(conn)

	ConnectToServer(client, "alex")

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Cat World")
	game := &Game{}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

func ConnectToServer(client catworldpb.CatWorldServiceClient, username string) {
	req := &catworldpb.ConnectRequest{
		Username: username,
	}

	res, _ := client.Connect(context.Background(), req)

	if res.GetSuccess() {
		fmt.Println("Successfully connected")
	}
}
