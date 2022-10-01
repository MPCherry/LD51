package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var spritePlayer, _, _ = ebitenutil.NewImageFromFile("resources/small_black_square.png")
var keyRecording = [][]ebiten.Key{}

const (
	runSpeed   = 2.2
	jumpHeight = 6
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
	carrying      *Key
	letGoOfPickup bool
}

func (p *Player) X() float64 {
	return p.x
}

func (p *Player) Y() float64 {
	return p.y
}

func (p *Player) Image() *ebiten.Image {
	if p.first {
		return spritePlayer.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image)
	} else {
		return spritePlayer.SubImage(image.Rect(16, 0, 32, 16)).(*ebiten.Image)
	}
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

	sawDown := false
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
		case ebiten.KeyDown:
			if p.letGoOfPickup {
				if p.carrying != nil && !p.jumped {
					p.carrying = nil
					p.letGoOfPickup = false
				}
			}
			sawDown = true
		}
	}
	if !sawDown {
		p.letGoOfPickup = true
	}

	p.newY = p.y + p.verticalSpeed
	if p.verticalSpeed < 3 {
		p.verticalSpeed += 0.3
	}

	for _, key := range world.keyList {
		if math.Abs(key.y-p.newY) < 16 && math.Abs(key.x-p.newX) < 16 {
			if p.carrying == nil {
				for _, k := range keys {
					if k == ebiten.KeyDown && p.letGoOfPickup {
						p.carrying = key
						p.letGoOfPickup = false
					}
				}
			}
		}
	}

	if p.carrying != nil {
		p.carrying.x = p.x
		p.carrying.y = p.y
	}

	for _, shadow := range world.players {
		if p == shadow {
			continue
		}

		if math.Abs(shadow.y-p.newY) < 16 && math.Abs(shadow.x-p.newX) < 16 {
			gameover = true
		}
	}
}
