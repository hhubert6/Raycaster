package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	ray "github.com/hhubert6/Raycaster/internal/raycaster"
	vec "github.com/hhubert6/Raycaster/internal/vector"
)

const (
	SCREEN_WIDTH  = 960
	SCREEN_HEIGHT = 640
	MAP_WIDTH     = 24
	MAP_HEIGHT    = 16
	TILE_SIZE     = 40
)

var COLOR_GREY = color.RGBA{200, 200, 200, 255}

type Game struct {
	gridMap   [][]int
	playerPos vec.Vec2
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.playerPos[1] -= 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.playerPos[1] += 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.playerPos[0] -= 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.playerPos[0] += 0.1
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		mouseX, mouseY := ebiten.CursorPosition()
		g.gridMap[mouseY/TILE_SIZE][mouseX/TILE_SIZE] = 1
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y := range MAP_HEIGHT {
		for x := range MAP_WIDTH {
			disX, disY := float32(x*TILE_SIZE), float32(y*TILE_SIZE)
			if g.gridMap[y][x] == 1 {
				vector.DrawFilledRect(screen, disX, disY, TILE_SIZE, TILE_SIZE, color.RGBA{50, 50, 200, 255}, false)
			}
			vector.StrokeRect(screen, disX, disY, TILE_SIZE, TILE_SIZE, 1, COLOR_GREY, false)
		}
	}

	mouseX, mouseY := ebiten.CursorPosition()

	playerDisplayX, playerDisplayY := getDisplayPos(g.playerPos)
	vector.StrokeLine(screen, playerDisplayX, playerDisplayY, float32(mouseX), float32(mouseY), 1, color.White, true)

	r := ray.NewRayFromTarget(g.playerPos, vec.Vec2{float64(mouseX) / TILE_SIZE, float64(mouseY) / TILE_SIZE})
	intersection, solidFound := r.Cast(g.gridMap)
	if solidFound {
		interX, interY := getDisplayPos(intersection)
		vector.StrokeCircle(screen, interX, interY, TILE_SIZE/5, 1, color.RGBA{200, 200, 0, 255}, true)
	}

	vector.DrawFilledCircle(screen, float32(mouseX), float32(mouseY), TILE_SIZE/4, color.RGBA{0, 200, 0, 255}, true)
	vector.DrawFilledCircle(screen, playerDisplayX, playerDisplayY, TILE_SIZE/4, color.RGBA{255, 100, 100, 255}, true)

}

func getDisplayPos(v vec.Vec2) (float32, float32) {
	return float32(v[0] * TILE_SIZE), float32(v[1] * TILE_SIZE)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func main() {
	game := &Game{
		gridMap:   make([][]int, MAP_HEIGHT),
		playerPos: vec.Vec2{float64(MAP_WIDTH / 2), float64(MAP_HEIGHT / 2)},
	}

	for i := range game.gridMap {
		game.gridMap[i] = make([]int, MAP_WIDTH)
	}

	ebiten.SetWindowSize(960, 640)
	ebiten.SetWindowTitle("Raycaster")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
