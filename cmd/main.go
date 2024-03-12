package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	vec "github.com/hhubert6/Raycaster/internal"
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
	gridMap   [MAP_HEIGHT][MAP_WIDTH]int
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

	// raycasting
	rayStart := g.playerPos.Copy()
	rayDir := vec.Vec2{float64(mouseX) / TILE_SIZE, float64(mouseY) / TILE_SIZE}
	rayDir.Sub(g.playerPos)
	rayDir.Normalize()
	vRayUnitStepSize := vec.Vec2{math.Hypot(1, rayDir[1]/rayDir[0]), math.Hypot(1, rayDir[0]/rayDir[1])}

	vMapCheck := vec.Vec2{math.Floor(rayStart[0]), math.Floor(rayStart[1])}
	vRayLength1D := vec.Vec2{}

	vStep := vec.Vec2{}

	if rayDir[0] < 0 {
		vStep[0] = -1
		vRayLength1D[0] = (g.playerPos[0] - vMapCheck[0]) * vRayUnitStepSize[0]
	} else {
		vStep[0] = 1
		vRayLength1D[0] = (vMapCheck[0] + 1 - g.playerPos[0]) * vRayUnitStepSize[0]
	}
	if rayDir[1] < 0 {
		vStep[1] = -1
		vRayLength1D[1] = (g.playerPos[1] - vMapCheck[1]) * vRayUnitStepSize[1]
	} else {
		vStep[1] = 1
		vRayLength1D[1] = (vMapCheck[1] + 1 - g.playerPos[1]) * vRayUnitStepSize[1]
	}

	solidFound := false
	distance := 0.0
	maxDistance := 100.0

	for !solidFound && distance < maxDistance {
		if vRayLength1D[0] < vRayLength1D[1] {
			vMapCheck[0] += vStep[0]
			distance = vRayLength1D[0]
			vRayLength1D[0] += vRayUnitStepSize[0]
		} else {
			vMapCheck[1] += vStep[1]
			distance = vRayLength1D[1]
			vRayLength1D[1] += vRayUnitStepSize[1]
		}

		if vMapCheck[0] >= 0 && vMapCheck[0] < MAP_WIDTH && vMapCheck[1] >= 0 && vMapCheck[1] < MAP_HEIGHT {
			if g.gridMap[int(vMapCheck[1])][int(vMapCheck[0])] == 1 {
				solidFound = true
			}
		}
	}

	if solidFound {
		intersection := rayStart.Copy()
		vDist := rayDir.Copy()
		vDist.Scale(distance)
		intersection.Add(vDist)
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
	game := &Game{}

	game.playerPos[0] = float64(MAP_WIDTH / 2)
	game.playerPos[1] = float64(MAP_HEIGHT / 2)

	ebiten.SetWindowSize(960, 640)
	ebiten.SetWindowTitle("Raycaster")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
