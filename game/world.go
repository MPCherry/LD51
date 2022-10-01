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
	for i := 1; i < 59; i++ {
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * float64(i), y: 640 - 16*7})
		world.wallObjects = append(world.wallObjects, &Wall{x: 16 * float64(i), y: 16 * 8})
	}

	tele0D := &Teleporter{x: 16 * 3, y: 16 * 3, spriteIndex: 0}
	tele0S := &Teleporter{x: 16 * 3, y: 640 - 32, destination: tele0D, spriteIndex: 0}
	world.wallObjects = append(world.wallObjects, tele0D)
	world.wallObjects = append(world.wallObjects, tele0S)

	key0 := &Key{x: 960 - 16*20, y: 16 * 7, originX: 960 - 16*20, originY: 16 * 7, active: true, spriteIndex: 0}
	world.keyList = append(world.keyList, key0)
	world.wallObjects = append(world.wallObjects, key0)

	swtch0 := &Switch{x: 16 * 3, y: 16 * 7, key: key0, spriteIndex: 0}
	world.switchObjects = append(world.switchObjects, swtch0)
	world.wallObjects = append(world.wallObjects, swtch0)

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

func UpdateWorld(world *World) {
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
			swtch.activated = false
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
