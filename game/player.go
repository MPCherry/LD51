package game

import (
	"fmt"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var spritePlayer, _, _ = ebitenutil.NewImageFromFile("resources/sprites/player3.png")
var spritePlayerLeft, _, _ = ebitenutil.NewImageFromFile("resources/sprites/player3left.png")
var spritePlayerJump, _, _ = ebitenutil.NewImageFromFile("resources/sprites/playerjump.png")
var spritePlayerJumpLeft, _, _ = ebitenutil.NewImageFromFile("resources/sprites/playerjumpleft.png")
var spriteGhost, _, _ = ebitenutil.NewImageFromFile("resources/sprites/ghost3.png")
var spriteGhostLeft, _, _ = ebitenutil.NewImageFromFile("resources/sprites/ghost3left.png")
var spriteGhostJump, _, _ = ebitenutil.NewImageFromFile("resources/sprites/ghostjump.png")
var spriteGhostJumpLeft, _, _ = ebitenutil.NewImageFromFile("resources/sprites/ghostjumpleft.png")
var keyRecording = [][]ebiten.Key{}
var xRecording = []float64{}
var yRecording = []float64{}

const (
	runSpeed   = 2
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
	xRecord       []float64
	yRecord       []float64
	keyIndex      int
	active        bool
	carrying      *Key
	letGoOfPickup bool
	frame         int
	frameCounter  int
	facing        bool
}

func (p *Player) X() float64 {
	return p.x
}

func (p *Player) Y() float64 {
	return p.y
}

func (p *Player) Image() *ebiten.Image {
	if p.first {
		if p.facing {
			if p.jumped {
				return spritePlayerJumpLeft.SubImage(image.Rect(16, 0, 0, 16)).(*ebiten.Image)
			}
			return spritePlayerLeft.SubImage(image.Rect(16, 16*p.frame, 0, 16+16*p.frame)).(*ebiten.Image)
		}
		if p.jumped {
			return spritePlayerJump.SubImage(image.Rect(16, 0, 0, 16)).(*ebiten.Image)
		}
		return spritePlayer.SubImage(image.Rect(16, 16*p.frame, 0, 16+16*p.frame)).(*ebiten.Image)
	} else {
		if p.facing {
			if p.jumped {
				return spriteGhostJumpLeft.SubImage(image.Rect(16, 0, 0, 16)).(*ebiten.Image)
			}
			return spriteGhostLeft.SubImage(image.Rect(16, 16*p.frame, 0, 16+16*p.frame)).(*ebiten.Image)
		}
		if p.jumped {
			return spriteGhostJump.SubImage(image.Rect(16, 0, 0, 16)).(*ebiten.Image)
		}
		return spriteGhost.SubImage(image.Rect(16, 16*p.frame, 0, 16+16*p.frame)).(*ebiten.Image)
	}
}

func (p *Player) Active() bool {
	return p.active
}

func (p *Player) Update(world *World, keys []ebiten.Key) {
	p.frameCounter++
	if p.frameCounter == 10 {
		p.frame++
		p.frameCounter = 0
	}
	if p.frame == 4 {
		p.frame = 0
	}
	if p.first {
		keyCopy := make([]ebiten.Key, len(keys))
		copy(keyCopy, keys)
		keyRecording = append(keyRecording, keyCopy)
	} else {
		if p.keyIndex < len(p.keyRecord) {
			keys = p.keyRecord[p.keyIndex]
		} else {
			keys = []ebiten.Key{}
		}
		p.keyIndex++
	}

	sawDown := false
	for _, k := range keys {
		switch k {
		case ebiten.KeyLeft:
			p.newX = p.x - runSpeed
			p.facing = true
		case ebiten.KeyRight:
			p.newX = p.x + runSpeed
			p.facing = false
		case ebiten.KeyUp:
			if !p.jumped {
				p.verticalSpeed = -jumpHeight
				p.jumped = true
				if p.first {
					jump.Rewind()
					jump.Play()
				}
			}
		case ebiten.KeyDown:
			if p.letGoOfPickup {
				if p.carrying != nil && !p.jumped {
					p.carrying = nil
					p.letGoOfPickup = false
					if p.first {
						keyDown.Rewind()
						keyDown.Play()
					}
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
						if p.first {
							keyUp.Rewind()
							keyUp.Play()
						}
					}
				}
			}
		}
	}

	if p.carrying != nil {
		p.carrying.x = p.x
		p.carrying.y = p.y
	}

	if p.first {
		xRecording = append(xRecording, p.newX)
		yRecording = append(yRecording, p.newY)
	}

	if !p.first {
		// if p.keyIndex < 600 {
		// 	if math.Abs(p.newX-p.xRecord[p.keyIndex]) > 64 {
		// 		p.newX = p.xRecord[p.keyIndex]
		// 	}
		// 	if math.Abs(p.newY-p.yRecord[p.keyIndex]) > 64 {
		// 		p.newY = p.yRecord[p.keyIndex]
		// 	}
		// }

		for _, shadow := range world.players {
			if p == shadow {
				continue
			}

			if math.Abs(shadow.y-p.newY) < 16 && math.Abs(shadow.x-p.newX) < 16 {
				gameover = true
				shaking = true
				shakeCounter = 120
				goCause = "collision"
				fmt.Println("gameover, collision")
				lostSound.Rewind()
				lostSound.Play()
			}
		}
	}

}

func (w *Player) UpdateAnyway() {

}
