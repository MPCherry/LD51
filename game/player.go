package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var spritePlayer, _, _ = ebitenutil.NewImageFromFile("resources/small_black_square.png")
var keyRecording = [][]ebiten.Key{}

const (
	runSpeed   = 2
	jumpHeight = 5
)

type Player struct {
	first         bool
	x             float64
	newX          float64
	y             float64
	newY          float64
	verticalSpeed float64
	jumped        bool
	keyRecord     [][]ebiten.Key
	keyIndex      int
	active        bool
}

func (p *Player) X() float64 {
	return p.x
}

func (p *Player) Y() float64 {
	return p.y
}

func (p *Player) Image() *ebiten.Image {
	return spritePlayer
}

func (p *Player) Active() bool {
	return p.active
}

func (p *Player) Update(world *World, keys []ebiten.Key) {
	if p.first {
		keyCopy := make([]ebiten.Key, len(keys))
		copy(keyCopy, keys)
		keyRecording = append(keyRecording, keyCopy)
	} else {
		p.keyIndex++
		if p.keyIndex < len(p.keyRecord) {
			keys = p.keyRecord[p.keyIndex]
		} else {
			keys = []ebiten.Key{}
		}
	}

	for _, k := range keys {
		switch k {
		case ebiten.KeyLeft:
			p.newX = p.x - runSpeed
		case ebiten.KeyRight:
			p.newX = p.x + runSpeed
		case ebiten.KeyUp:
			if !p.jumped {
				p.verticalSpeed = -jumpHeight
				p.jumped = true
			}
		}
	}

	p.newY = p.y + p.verticalSpeed
	if p.verticalSpeed < 3 {
		p.verticalSpeed += 0.2
	}
}
