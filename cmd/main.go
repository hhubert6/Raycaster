package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	entities "github.com/hhubert6/Raycaster/internal/entities"
	ray "github.com/hhubert6/Raycaster/internal/raycaster"
	vec "github.com/hhubert6/Raycaster/internal/vector"
)

const (
	SCREEN_WIDTH  = 960
	SCREEN_HEIGHT = 640
	MAP_WIDTH     = 24
	MAP_HEIGHT    = 16
	TILE_SIZE     = 10
	FOV           = math.Pi / 3
	RESOLUTION    = 2
	NUM_OF_RAYS   = SCREEN_WIDTH / RESOLUTION
)

var (
	COLOR_GREY  = color.RGBA{200, 200, 200, 255}
	SCREEN_DIST = (SCREEN_WIDTH / 2) / math.Tan(FOV/2)
	imgTest     image.Image
)

type Game struct {
	gridMap [][]int
	player  entities.Player
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Move(entities.FORWARD, g.gridMap)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Move(entities.BACKWARD, g.gridMap)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.Move(entities.LEFT, g.gridMap)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.Move(entities.RIGHT, g.gridMap)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.Rotate(-0.02)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.Rotate(0.02)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		mouseX, mouseY := ebiten.CursorPosition()
		x, y := mouseX/TILE_SIZE, mouseY/TILE_SIZE
		if 0 <= y && y < MAP_HEIGHT && 0 <= x && x < MAP_WIDTH {
			g.gridMap[y][x] = 1
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	angle := g.player.Angle - FOV/2
	offsetStep := FOV / float64(NUM_OF_RAYS)

	for i := 0; i < NUM_OF_RAYS; i++ {
		angle += offsetStep
		r := ray.NewRayFromAngle(g.player.Pos, angle)
		offset, distance, _, solidFound := r.Cast(g.gridMap)
		if solidFound {
			height := SCREEN_DIST / (distance * math.Cos(g.player.Angle-angle))

			x := float32(i * RESOLUTION)
			startY := float32(SCREEN_HEIGHT/2 - height/2)

			dy := float32(height / 64)

			for i := 0; i < 64; i += 1 {
				clr := imgTest.At(int(64*offset), i)
				vector.DrawFilledRect(screen, x, startY+float32(i)*dy, RESOLUTION, dy, clr, false)
			}
		}
	}

	for y := range MAP_HEIGHT {
		for x := range MAP_WIDTH {
			disX, disY := float32(x*TILE_SIZE), float32(y*TILE_SIZE)
			if g.gridMap[y][x] == 1 {
				vector.DrawFilledRect(screen, disX, disY, TILE_SIZE, TILE_SIZE, color.RGBA{50, 50, 200, 255}, false)
			}
			vector.StrokeRect(screen, disX, disY, TILE_SIZE, TILE_SIZE, 1, COLOR_GREY, false)
		}
	}
	dst := g.player.Pos.Copy().Add(g.player.Dir().Scale(5))
	playerDisplayX, playerDisplayY := getDisplayPos(g.player.Pos)
	dstX, dstY := getDisplayPos(*dst)
	vector.StrokeLine(screen, playerDisplayX, playerDisplayY, dstX, dstY, 1, color.White, false)
	vector.DrawFilledCircle(screen, playerDisplayX, playerDisplayY, TILE_SIZE/4, color.RGBA{255, 100, 100, 255}, true)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f", ebiten.ActualFPS()), SCREEN_WIDTH-50, 10)
}

func getDisplayPos(v vec.Vec2) (float32, float32) {
	return float32(v[0] * TILE_SIZE), float32(v[1] * TILE_SIZE)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func main() {
	game := &Game{
		gridMap: make([][]int, MAP_HEIGHT),
	}

	for i := range game.gridMap {
		game.gridMap[i] = make([]int, MAP_WIDTH)
	}

	game.player.Pos = vec.Vec2{float64(MAP_WIDTH / 2), float64(MAP_HEIGHT / 2)}
	// img, _, err := ebitenutil.NewImageFromFile("assets/redbrick1.png")
	img, err := getImageFromFilePath("assets/redbrick1.png")
	if err != nil {
		log.Fatal(err)
	}
	imgTest = img

	ebiten.SetWindowSize(960, 640)
	ebiten.SetWindowTitle("Raycaster")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
