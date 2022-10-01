package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Object interface {
	X() float64
	Y() float64
	Image() *ebiten.Image
	Update(*World, []ebiten.Key)
	Active() bool
}

type World struct {
	keys          []ebiten.Key
	playerObjects []Object
	wallObjects   []Object
	switchObjects []*Switch
	keyList       []*Key
	players       []*Player
}

func NewWorld() *World {
	world := &World{}
	player := &Player{x: 16 * 5, y: 640 - 32, newX: 16 * 5, newY: 640 - 32, first: true, active: true, letGoOfPickup: true}
	world.players = append(world.players, player)
	world.playerObjects = append(world.playerObjects, player)
	for i := 0; i < 60; i++ {
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * float64(i), y: 640 - 16})
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * float64(i), y: 0})
	}
	for i := 1; i < 39; i++ {
		world.wallObjects = append(world.wallObjects, &Wall{x: 0, y: 16 * float64(i)})
		world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16, y: 16 * float64(i)})
	}

	// Make the map

	// Home chamber
	for i := 1; i < 59; i++ {
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * float64(i), y: 640 - 16*7})
	}

	// Top chamber
	for i := 1; i < 47; i++ {
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * float64(i), y: 16 * 8})
	}
	for i := 0; i < 7; i++ {
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 46, y: 16 + float64(i)*16})
	}

	tele0D := &Teleporter{x: 16 * 3, y: 16 * 3, spriteIndex: 0}
	tele0S := &Teleporter{x: 16 * 3, y: 640 - 32, destination: tele0D, spriteIndex: 0}
	world.wallObjects = append(world.wallObjects, tele0D)
	world.wallObjects = append(world.wallObjects, tele0S)

	key0 := &Key{x: 960 - 16*20, y: 16 * 4, originX: 960 - 16*20, originY: 16 * 4, active: true, spriteIndex: 0}
	world.keyList = append(world.keyList, key0)
	world.wallObjects = append(world.wallObjects, key0)
	world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16*20, y: 16 * 5})
	world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16*20, y: 16 * 6})
	world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16*20, y: 16 * 7})
	world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16*19, y: 16 * 7})
	world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16*19, y: 16 * 6})
	world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16*18, y: 16 * 7})
	world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16*21, y: 16 * 7})
	world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16*22, y: 16 * 7})
	world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16*23, y: 16 * 7})
	world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16*26, y: 16 * 7})
	world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16*27, y: 16 * 7})
	world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16*28, y: 16 * 7})

	for i := 0; i < 3; i++ {
		world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16*35 - 16*float64(i), y: 16 * 7})
		world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16*35 - 16*float64(i), y: 16 * 6})
	}
	world.wallObjects = append(world.wallObjects, &Wall{x: 960 - 16*35 - 16*3, y: 16 * 7})

	swtch0 := &Switch{x: 16 * 3, y: 16 * 7, key: key0, spriteIndex: 0, resets: false}
	world.switchObjects = append(world.switchObjects, swtch0)
	world.wallObjects = append(world.wallObjects, swtch0)

	door0a := &Wall{x: 16 * 4, y: 640 - 32 - 16*3, spriteIndex: 0, wallSwitch: swtch0}
	world.wallObjects = append(world.wallObjects, door0a)
	door0b := &Wall{x: 16 * 5, y: 640 - 32 - 16*3, spriteIndex: 0, wallSwitch: swtch0}
	world.wallObjects = append(world.wallObjects, door0b)
	door0c := &Wall{x: 16 * 6, y: 640 - 32 - 16*3, spriteIndex: 0, wallSwitch: swtch0}
	world.wallObjects = append(world.wallObjects, door0c)
	door0d := &Wall{x: 16 * 4, y: 640 - 32 - 16*4, spriteIndex: 0, wallSwitch: swtch0}
	world.wallObjects = append(world.wallObjects, door0d)
	door0e := &Wall{x: 16 * 6, y: 640 - 32 - 16*4, spriteIndex: 0, wallSwitch: swtch0}
	world.wallObjects = append(world.wallObjects, door0e)

	// Second chamber
	tele1D := &Teleporter{x: 16 * 23, y: 16*8 + 16*13, spriteIndex: 1}
	world.wallObjects = append(world.wallObjects, tele1D)
	tele1S := &Teleporter{x: 16 * 5, y: 640 - 32 - 16*4, spriteIndex: 1, destination: tele1D}
	world.wallObjects = append(world.wallObjects, tele1S)

	for i := 1; i < 46; i++ {
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * float64(i), y: 16*8 + 16*14})
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * float64(i), y: 16*8 + 16*15})
	}
	for i := 0; i < 24; i++ {
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 46, y: 16*9 + float64(i)*16})
	}

	// world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 22, y: 16*7 + 16*12})
	// world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 23, y: 16*7 + 16*12})
	// world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 24, y: 16*7 + 16*12})
	// world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 25, y: 16*7 + 16*12})
	// world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 21, y: 16*7 + 16*12})

	for i := 0; i < 13; i++ {
		for j := 0; j < i; j++ {
			world.wallObjects = append(world.wallObjects, &Wall{x: 16*(23+10) + float64(i)*16, y: 16*8 + 16*13 - float64(j)*16})
			world.wallObjects = append(world.wallObjects, &Wall{x: 16*(23-10) - float64(i)*16, y: 16*8 + 16*13 - float64(j)*16})
		}
	}

	key1 := &Key{x: 16, y: 16 * 9, originX: 16, originY: 16 * 9, active: true, spriteIndex: 1}
	world.keyList = append(world.keyList, key1)
	world.wallObjects = append(world.wallObjects, key1)

	key2 := &Key{x: 16 * 45, y: 16 * 9, originX: 16 * 45, originY: 16 * 9, active: true, spriteIndex: 2}
	world.keyList = append(world.keyList, key2)
	world.wallObjects = append(world.wallObjects, key2)

	swtch1 := &Switch{x: 16 * 42, y: 16 * 12, key: key1, spriteIndex: 1, resets: false}
	world.switchObjects = append(world.switchObjects, swtch1)
	world.wallObjects = append(world.wallObjects, swtch1)

	swtch2 := &Switch{x: 16 * 3, y: 16 * 12, key: key2, spriteIndex: 2, resets: false}
	world.switchObjects = append(world.switchObjects, swtch2)
	world.wallObjects = append(world.wallObjects, swtch2)

	// Side chamber
	tele2D := &Teleporter{x: 16 * 50, y: 16*8 + 16*24, spriteIndex: 2}
	world.wallObjects = append(world.wallObjects, tele2D)
	tele2S := &Teleporter{x: 16 * 8, y: 640 - 32 - 16*4, spriteIndex: 2, destination: tele2D}
	world.wallObjects = append(world.wallObjects, tele2S)

	for i := 0; i < 10; i++ {
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 50, y: 16*8 + 16*22 - 16*3*float64(i)})
		// world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 53, y: 16*8 + 16*22 - 16*3*float64(i)})
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 55, y: 16*8 + 16*22 - 16*3*float64(i)})
	}
	for i := 0; i < 10; i++ {
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 47, y: 16*8 + 16*21 - 16*3*float64(i)})
		// world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 48, y: 16*8 + 16*21 - 16*3*float64(i)})
		// world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 47, y: 16*8 + 16*22 - 16*3*float64(i)})
		// world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 53, y: 16*8 + 16*22 - 16*3*float64(i)})
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 58, y: 16*8 + 16*21 - 16*3*float64(i)})
		// world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 58, y: 16*8 + 16*22 - 16*3*float64(i)})
		// world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 57, y: 16*8 + 16*21 - 16*3*float64(i)})
	}

	swtch3 := &Switch{x: 16 * 50, y: 16*8 + 16*22 - 16*3*9 - 16, spriteIndex: 3, resets: false}
	world.switchObjects = append(world.switchObjects, swtch3)
	world.wallObjects = append(world.wallObjects, swtch3)
	swtch4 := &Switch{x: 16 * 55, y: 16*8 + 16*22 - 16*3*9 - 16, spriteIndex: 4, resets: false}
	world.switchObjects = append(world.switchObjects, swtch4)
	world.wallObjects = append(world.wallObjects, swtch4)

	door1a := &Wall{x: 16 * 7, y: 640 - 32 - 16*3, spriteIndex: 1, wallSwitch: swtch1}
	world.wallObjects = append(world.wallObjects, door1a)
	door1b := &Wall{x: 16 * 8, y: 640 - 32 - 16*3, spriteIndex: 1, wallSwitch: swtch1}
	world.wallObjects = append(world.wallObjects, door1b)
	door1c := &Wall{x: 16 * 9, y: 640 - 32 - 16*3, spriteIndex: 1, wallSwitch: swtch1}
	world.wallObjects = append(world.wallObjects, door1c)
	door1d := &Wall{x: 16 * 7, y: 640 - 32 - 16*4, spriteIndex: 1, wallSwitch: swtch1}
	world.wallObjects = append(world.wallObjects, door1d)
	door1e := &Wall{x: 16 * 9, y: 640 - 32 - 16*4, spriteIndex: 1, wallSwitch: swtch1}
	world.wallObjects = append(world.wallObjects, door1e)

	door2a := &Wall{x: 16 * 10, y: 640 - 32 - 16*3, spriteIndex: 2, wallSwitch: swtch2}
	world.wallObjects = append(world.wallObjects, door2a)
	door2b := &Wall{x: 16 * 11, y: 640 - 32 - 16*3, spriteIndex: 2, wallSwitch: swtch2}
	world.wallObjects = append(world.wallObjects, door2b)
	door2c := &Wall{x: 16 * 12, y: 640 - 32 - 16*3, spriteIndex: 2, wallSwitch: swtch2}
	world.wallObjects = append(world.wallObjects, door2c)
	door2d := &Wall{x: 16 * 10, y: 640 - 32 - 16*4, spriteIndex: 2, wallSwitch: swtch2}
	world.wallObjects = append(world.wallObjects, door2d)
	door2e := &Wall{x: 16 * 12, y: 640 - 32 - 16*4, spriteIndex: 2, wallSwitch: swtch2}
	world.wallObjects = append(world.wallObjects, door2e)

	// Bottom chamber
	tele3D := &Teleporter{x: 16 * 42, y: 16*6 + 16*24, spriteIndex: 3}
	world.wallObjects = append(world.wallObjects, tele3D)
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 42, y: 16*7 + 16*25})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 41, y: 16*7 + 16*25})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 43, y: 16*7 + 16*25})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 41, y: 16*6 + 16*25})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 42, y: 16*6 + 16*25})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 43, y: 16*6 + 16*25})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 40, y: 16*7 + 16*25})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 40, y: 16*6 + 16*25})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 44, y: 16*7 + 16*25})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 44, y: 16*6 + 16*25})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 40, y: 16*7 + 16*25})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 40, y: 16*6 + 16*25})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 45, y: 16*7 + 16*25})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 45, y: 16*6 + 16*25})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 39, y: 16*7 + 16*25})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 39, y: 16*6 + 16*25})
	tele3S := &Teleporter{x: 16 * 11, y: 640 - 32 - 16*4, spriteIndex: 3, destination: tele3D}
	world.wallObjects = append(world.wallObjects, tele3S)

	for i := 6; i < 35; i++ {
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * float64(i), y: 16*7 + 16*14 + 7*16})
	}
	for i := 0; i < 4; i++ {
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 6, y: 16*7 + 16*14 + 7*16 + 16*(1+float64(i))})
	}
	for i := 1; i < 4; i++ {
		for j := 0; j < 3; j++ {
			world.wallObjects = append(world.wallObjects, &Wall{x: 16*6 + 16*float64(i)*5, y: 16*7 + 16*14 + 7*16 + 16*(1+float64(j))})
			world.wallObjects = append(world.wallObjects, &Wall{x: 16*6 + 16*float64(i)*5, y: 16*2 + 16*14 + 7*16 + 16*(1+float64(j))})
		}
	}

	swtch5 := &Switch{x: 16 * 18, y: 16*6 + 16*26, spriteIndex: 5, resets: false}
	world.switchObjects = append(world.switchObjects, swtch5)
	world.wallObjects = append(world.wallObjects, swtch5)

	swtch6 := &Switch{x: 16 * 18, y: 16*6 + 16*21, spriteIndex: 6, resets: false}
	world.switchObjects = append(world.switchObjects, swtch6)
	world.wallObjects = append(world.wallObjects, swtch6)

	swtch7 := &Switch{x: 16 * 13, y: 16*6 + 16*26, spriteIndex: 7, resets: false}
	world.switchObjects = append(world.switchObjects, swtch7)
	world.wallObjects = append(world.wallObjects, swtch7)

	swtch8 := &Switch{x: 16 * 13, y: 16*6 + 16*21, spriteIndex: 8, resets: false}
	world.switchObjects = append(world.switchObjects, swtch8)
	world.wallObjects = append(world.wallObjects, swtch8)

	swtch9 := &Switch{x: 16 * 8, y: 16*6 + 16*26, spriteIndex: 9, resets: false}
	world.switchObjects = append(world.switchObjects, swtch9)
	world.wallObjects = append(world.wallObjects, swtch9)

	door4 := &Wall{x: 16 * 21, y: 16*6 + 16*21, spriteIndex: 5, wallSwitch: swtch5}
	world.wallObjects = append(world.wallObjects, door4)
	door5 := &Wall{x: 16 * 16, y: 16*6 + 16*26, spriteIndex: 6, wallSwitch: swtch6}
	world.wallObjects = append(world.wallObjects, door5)
	door6 := &Wall{x: 16 * 16, y: 16*6 + 16*21, spriteIndex: 7, wallSwitch: swtch7}
	world.wallObjects = append(world.wallObjects, door6)
	door7 := &Wall{x: 16 * 11, y: 16*6 + 16*26, spriteIndex: 8, wallSwitch: swtch8}
	world.wallObjects = append(world.wallObjects, door7)
	door8 := &Wall{x: 16 * 11, y: 16*6 + 16*21, spriteIndex: 9, wallSwitch: swtch9}
	world.wallObjects = append(world.wallObjects, door8)

	// world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 3, y: 640 - 16*3})

	// keyA := &Key{x: 16 * 7, y: 640 - 16*2, originX: 16 * 7, originY: 640 - 16*2, active: true}
	// world.wallObjects = append(world.wallObjects, keyA)
	// world.keyList = append(world.keyList, keyA)
	// switchA := &Switch{x: 16 * 5, y: 640 - 16*2, key: keyA}
	// world.switchObjects = append(world.switchObjects, switchA)

	// world.wallObjects = append(world.wallObjects, switchA)
	// world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 6, y: 640 - 16*3, wallSwitch: switchA})

	// teleB := &Teleporter{x: 16 * 8, y: 16 * 1}
	// teleA := &Teleporter{x: 16 * 8, y: 640 - 16*2, destination: teleB}
	// world.wallObjects = append(world.wallObjects, teleA)
	// world.wallObjects = append(world.wallObjects, teleB)

	return world
}

var recordCounter = 0
var respawnCounter = 0
var respawnShadowCounter = 0
var respawning = false

var shadowCounter = 0

var gameover = false

func UpdateWorld(world *World) {
	if gameover {
		return
	}
	world.keys = inpututil.AppendPressedKeys(world.keys[:0])

	for _, object := range world.playerObjects {
		if object.Active() {
			object.Update(world, world.keys)
		}
	}
	for _, object := range world.wallObjects {
		if object.Active() {
			object.Update(world, world.keys)
		}
	}
	for _, player := range world.players {
		if player.first || player.keyIndex < len(player.keyRecord) {
			player.x = player.newX
			player.y = player.newY
		}
	}

	recordCounter++
	if recordCounter%600 == 0 {
		keyRecordCopy := make([][]ebiten.Key, len(keyRecording))
		copy(keyRecordCopy, keyRecording)
		keyRecording = [][]ebiten.Key{}

		player := &Player{x: 16 * 5, y: 640 - 32, newX: 16 * 5, newY: 640 - 32, first: false, keyRecord: keyRecordCopy, active: true, letGoOfPickup: true}
		world.players = append(world.players, player)
		world.playerObjects = append(world.playerObjects, player)
		shadowCounter++
		respawning = true
		respawnShadowCounter = 0
		respawnCounter = 0
		for _, player := range world.players {
			player.x = -32
			player.newX = -32
			player.y = -32
			player.newY = -32
			player.active = false
			player.carrying = nil
			player.letGoOfPickup = true
		}
		for _, swtch := range world.switchObjects {
			if swtch.resets {
				swtch.activated = false
			}
		}
		for _, key := range world.keyList {
			key.x = key.originX
			key.y = key.originY
			key.active = true
		}
	}

	if respawning {
		if respawning && respawnCounter%60 == 0 {
			if respawnShadowCounter == len(world.players)-1 {
				world.players[0].x = 16 * 5
				world.players[0].newX = 16 * 5
				world.players[0].y = 640 - 32
				world.players[0].newY = 640 - 32
				world.players[0].jumped = false
				world.players[0].active = true
				respawning = false
				recordCounter = 0
			} else {
				world.players[1+respawnShadowCounter].x = 16 * 5
				world.players[1+respawnShadowCounter].newX = 16 * 5
				world.players[1+respawnShadowCounter].y = 640 - 32
				world.players[1+respawnShadowCounter].newY = 640 - 32
				world.players[1+respawnShadowCounter].jumped = false
				world.players[1+respawnShadowCounter].keyIndex = 0
				world.players[1+respawnShadowCounter].active = true
				respawnShadowCounter++
			}
		}
		respawnCounter++
	}
}

func DrawObjects(screen *ebiten.Image, world *World) {
	screen.Fill(color.White)
	for _, object := range append(world.playerObjects, world.wallObjects...) {
		if object.Active() {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(object.X(), object.Y())
			screen.DrawImage(object.Image(), op)
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%d", recordCounter))
}
