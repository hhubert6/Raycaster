package main

import (
	"image/color"
	"log"
	"math"

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
	TILE_SIZE     = 10
	FOV           = math.Pi / 2
)

var COLOR_GREY = color.RGBA{200, 200, 200, 255}

type Game struct {
	gridMap     [][]int
	playerPos   vec.Vec2
	playerAngle float64
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dx, dy := math.Cos(g.playerAngle), math.Sin(g.playerAngle)
		g.playerPos[0] += 0.1 * dx
		g.playerPos[1] += 0.1 * dy
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		dx, dy := math.Cos(g.playerAngle), math.Sin(g.playerAngle)
		g.playerPos[0] -= 0.1 * dx
		g.playerPos[1] -= 0.1 * dy
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dx, dy := math.Cos(g.playerAngle-math.Pi/2), math.Sin(g.playerAngle-math.Pi/2)
		g.playerPos[0] += 0.1 * dx
		g.playerPos[1] += 0.1 * dy
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dx, dy := math.Cos(g.playerAngle+math.Pi/2), math.Sin(g.playerAngle+math.Pi/2)
		g.playerPos[0] += 0.1 * dx
		g.playerPos[1] += 0.1 * dy
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.playerAngle -= 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.playerAngle += 0.02
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
	angle := g.playerAngle - FOV/2

	numOfRays := 160
	offsetStep := FOV / float64(numOfRays)

	for i := 0; i < numOfRays; i++ {
		angle += offsetStep
		r := ray.NewRayFromAngle(g.playerPos, angle)
		intersection, wall, solidFound := r.Cast(g.gridMap)
		if solidFound {
			x := i * (SCREEN_WIDTH / numOfRays)
			dist := intersection.Copy()
			dist.Sub(g.playerPos)
			height := SCREEN_HEIGHT / dist.Length()
			y := SCREEN_HEIGHT/2 - height/2

			v := uint8(255 - 10*dist.Length() - float64(wall*50))
			clr := color.RGBA{v, v, v, 255}
			vector.DrawFilledRect(screen, float32(x), float32(y), float32(SCREEN_WIDTH/numOfRays), float32(height), clr, false)
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
	playerDisplayX, playerDisplayY := getDisplayPos(g.playerPos)
	dst := g.playerPos.Copy()
	playerDir := vec.Vec2{math.Cos(g.playerAngle), math.Sin(g.playerAngle)}
	playerDir.Scale(5)
	dst.Add(playerDir)
	dstX, dstY := getDisplayPos(dst)
	vector.StrokeLine(screen, playerDisplayX, playerDisplayY, dstX, dstY, 1, color.White, false)
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
		gridMap:     make([][]int, MAP_HEIGHT),
		playerPos:   vec.Vec2{float64(MAP_WIDTH / 2), float64(MAP_HEIGHT / 2)},
		playerAngle: 0,
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
