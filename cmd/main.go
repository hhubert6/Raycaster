package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	SCREEN_WIDTH  = 960
	SCREEN_HEIGHT = 640
	MAP_WIDTH     = 30
	MAP_HEIGHT    = 20
	TILE_SIZE     = 32
)

var COLOR_GREY = color.RGBA{200, 200, 200, 255}

type Game struct {
	gridMap   [MAP_WIDTH * MAP_HEIGHT]int
	playerPos [2]float32
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.playerPos[1] -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.playerPos[1] += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.playerPos[0] -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.playerPos[0] += 2
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y := range MAP_HEIGHT {
		for x := range MAP_WIDTH {
			disX, disY := float32(x*TILE_SIZE), float32(y*TILE_SIZE)
			if g.gridMap[y*MAP_WIDTH+x] == 1 {
				vector.DrawFilledRect(screen, disX, disY, TILE_SIZE, TILE_SIZE, color.White, false)
			}
			vector.StrokeRect(screen, disX, disY, TILE_SIZE, TILE_SIZE, 1, COLOR_GREY, false)
		}
	}

	mouseX, mouseY := ebiten.CursorPosition()
	vector.DrawFilledCircle(screen, float32(mouseX), float32(mouseY), TILE_SIZE/4, color.RGBA{200, 200, 0, 255}, false)

	vector.DrawFilledCircle(screen, g.playerPos[0], g.playerPos[1], TILE_SIZE/4, color.RGBA{255, 100, 100, 255}, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func main() {
	game := &Game{}

	game.gridMap[MAP_WIDTH+1] = 1
	game.playerPos[0] = SCREEN_HEIGHT / 2
	game.playerPos[1] = SCREEN_WIDTH / 2

	ebiten.SetWindowSize(960, 640)
	ebiten.SetWindowTitle("Raycaster")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
