package main

import (
	_ "embed"
	_ "image/png"

	"log"
	"test/game"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	windowWidth  = 600
	windowHeight = 600
	gameWidth    = 600
	gameHeight   = 600
)

type Game struct {
	world *game.World
}

func newGame() *Game {
	return &Game{
		world: game.NewWorld(),
	}
}

func (g *Game) Update() error {
	game.UpdateWorld(g.world)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	game.DrawObjects(screen, g.world)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return gameWidth, gameHeight
}

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetTPS(120)
	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
