package game

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Object interface {
	X() float64
	Y() float64
	Image() *ebiten.Image
	Update(*World, []ebiten.Key)
	Active() bool
	UpdateAnyway()
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

	initSounds()

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

	swtch0 := &Switch{x: 16 * 3, y: 16 * 7, key: key0, spriteIndex: 0, resets: false, cover: 3}
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

	swtch1 := &Switch{x: 16 * 42, y: 16 * 12, key: key1, spriteIndex: 1, resets: false, cover: 4}
	world.switchObjects = append(world.switchObjects, swtch1)
	world.wallObjects = append(world.wallObjects, swtch1)

	swtch2 := &Switch{x: 16 * 4, y: 16 * 12, key: key2, spriteIndex: 2, resets: false, cover: 2}
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

	swtch10 := &Switch{x: 16 * 3, y: 16*6 + 16*25, spriteIndex: 10, resets: false}
	world.switchObjects = append(world.switchObjects, swtch10)
	world.wallObjects = append(world.wallObjects, swtch10)
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 3, y: 16*6 + 16*26})

	swtch11 := &Switch{x: 16 * 54, y: 16*6 + 16*30, spriteIndex: 11, resets: false, final: true}
	world.switchObjects = append(world.switchObjects, swtch11)
	world.wallObjects = append(world.wallObjects, swtch11)
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 54, y: 16*7 + 16*30})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 54, y: 16*8 + 16*30})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 53, y: 16*8 + 16*30})
	world.wallObjects = append(world.wallObjects, &Wall{x: 16 * 55, y: 16*8 + 16*30})

	// Final Doors
	for i := 0; i < 3; i++ {
		for j := 0; j < 5; j++ {
			door9 := &Wall{x: 16*31 + 16*float64(i), y: 16*6 + 16*28 + 16*float64(j), spriteIndex: 10, wallSwitch: swtch10}
			world.wallObjects = append(world.wallObjects, door9)
		}
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 5; j++ {
			door10 := &Wall{x: 16*35 + 16*float64(i), y: 16*6 + 16*28 + 16*float64(j), spriteIndex: 3, wallSwitch: swtch3}
			world.wallObjects = append(world.wallObjects, door10)
		}
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 5; j++ {
			door11 := &Wall{x: 16*39 + 16*float64(i), y: 16*6 + 16*28 + 16*float64(j), spriteIndex: 4, wallSwitch: swtch4}
			world.wallObjects = append(world.wallObjects, door11)
		}
	}

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

var interluding = false
var interludeCounter = 0

var shadowCounter = 0

var starting = false
var canStart = false
var gameover = true
var goCause = "start"

func UpdateWorld(world *World) {
	for _, object := range world.wallObjects {
		object.UpdateAnyway()
	}
	if shaking {
		shakeCounter--
		if shakeCounter == 0 {
			shaking = false
		}
	}
	if gameover {
		stopMusic()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	world.keys = inpututil.AppendPressedKeys(world.keys[:0])
	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		canStart = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) && !interluding {
		stopMusic()
		canStart = false
		fmt.Println("Resetting game over")
		gameover = false
		starting = true
		bottomCoverDraw = true
		middleCoverDraw = true
		rightCoverDraw = true

		for _, swtch := range world.switchObjects {
			swtch.activated = false
		}
		for _, key := range world.keyList {
			key.x = key.originX
			key.y = key.originY
			key.active = true
		}
		respawning = false
		respawnShadowCounter = 0
		respawnCounter = 0
		keyRecording = [][]ebiten.Key{}
		xRecording = []float64{}
		yRecording = []float64{}
		recordCounter = 0

		for _, player := range world.players {
			if !player.first {
				player.active = false
			}
		}
		player := world.players[0]
		world.players = []*Player{}
		world.players = append(world.players, player)
		world.players[0].x = 16 * 5
		world.players[0].newX = 16 * 5
		world.players[0].y = 640 - 32
		world.players[0].newY = 640 - 32
		world.players[0].jumped = false
		world.players[0].verticalSpeed = 0
		world.players[0].active = true
		player.carrying = nil
		player.letGoOfPickup = true

	}
	if gameover {
		return
	}
	if starting {
		if len(world.keys) > 0 && canStart {
			fmt.Println("Starting game")
			starting = false
			startSound.Rewind()
			startSound.Play()
			playMusic()
			// world.wallObjects = append(world.wallObjects, &Spawn{x: 16 * 5, y: 640 - 32, spriteIndex: 0, life: 12, active: true})
		}
		return
	}

	for _, object := range world.playerObjects {
		if object.Active() && !interluding {
			object.Update(world, world.keys)
		}
	}
	for _, object := range world.wallObjects {
		if object.Active() && !interluding {
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
	if recordCounter%600 == 0 && !respawning {
		interluding = true
		interludeCounter = 0
		timeUp.Rewind()
		timeUp.Play()
		for _, player := range world.players {
			world.wallObjects = append(world.wallObjects, &Spawn{x: player.x, y: player.y, spriteIndex: 1, life: 120, active: true})
		}
	}
	if interluding {
		interludeCounter++
		if interludeCounter == 120 {
			interluding = false
			if world.players[0].carrying != nil {
				gameover = true
				shaking = true
				shakeCounter = 120
				goCause = "key"
				fmt.Println("gameover, carrying key")
				lostSound.Rewind()
				lostSound.Play()
				return
			}
			keyRecordCopy := make([][]ebiten.Key, len(keyRecording))
			xRecordCopy := make([]float64, len(xRecording))
			yRecordCopy := make([]float64, len(yRecording))
			copy(keyRecordCopy, keyRecording)
			copy(xRecordCopy, xRecording)
			copy(yRecordCopy, yRecording)
			keyRecording = [][]ebiten.Key{}
			xRecording = []float64{}
			yRecording = []float64{}

			player := &Player{x: 16 * 5, y: 640 - 32, newX: 16 * 5, newY: 640 - 32, first: false, keyRecord: keyRecordCopy, xRecord: xRecordCopy, yRecord: yRecordCopy, active: true, letGoOfPickup: true}
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
	}

	if respawning {
		if respawning && respawnCounter%60 == 0 {
			if respawnShadowCounter == len(world.players)-1 {
				world.players[0].x = 16 * 5
				world.players[0].newX = 16 * 5
				world.players[0].y = 640 - 32
				world.players[0].newY = 640 - 32
				world.players[0].jumped = false
				world.players[0].verticalSpeed = 0
				world.players[0].active = true
				world.wallObjects = append(world.wallObjects, &Spawn{x: 16 * 5, y: 640 - 32, spriteIndex: 0, life: 12, active: true})
				respawning = false
				recordCounter = 0
				playerSpawn.Rewind()
				playerSpawn.Play()
				shaking = true
				shakeCounter = 5
				playMusic()
			} else {
				world.players[1+respawnShadowCounter].x = 16 * 5
				world.players[1+respawnShadowCounter].newX = 16 * 5
				world.players[1+respawnShadowCounter].y = 640 - 32
				world.players[1+respawnShadowCounter].newY = 640 - 32
				world.players[1+respawnShadowCounter].jumped = false
				world.players[1+respawnShadowCounter].keyIndex = 0
				world.players[1+respawnShadowCounter].verticalSpeed = 0
				world.players[1+respawnShadowCounter].active = true
				world.wallObjects = append(world.wallObjects, &Spawn{x: 16 * 5, y: 640 - 32, spriteIndex: 1, life: 12, active: true})
				shaking = true
				shakeCounter = 5
				respawnShadowCounter++
				ghostSpawn.Rewind()
				ghostSpawn.Play()
			}
		}
		respawnCounter++
	}

}

var startScreen, _, _ = ebitenutil.NewImageFromFile("resources/start.png")
var collisionScreen, _, _ = ebitenutil.NewImageFromFile("resources/collision.png")
var keyScreen, _, _ = ebitenutil.NewImageFromFile("resources/key.png")
var winScreen, _, _ = ebitenutil.NewImageFromFile("resources/win.png")

var bottomCover, _, _ = ebitenutil.NewImageFromFile("resources/bottom_cover.png")
var bottomCoverDraw = true
var middleCover, _, _ = ebitenutil.NewImageFromFile("resources/middle_cover.png")
var middleCoverDraw = true
var rightCover, _, _ = ebitenutil.NewImageFromFile("resources/righ_cover.png")
var rightCoverDraw = true

var bg, _, _ = ebitenutil.NewImageFromFile("resources/bg.png")
var fg, _, _ = ebitenutil.NewImageFromFile("resources/fg.png")

var shaking = true
var shakeCounter = 30

func DrawObjects(screen *ebiten.Image, world *World) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(bg, op)
	offsetX := rand.Intn(4)
	offsetY := rand.Intn(4)
	for _, object := range append(world.playerObjects, world.wallObjects...) {
		if object.Active() {
			op := &ebiten.DrawImageOptions{}
			if shaking {
				op.GeoM.Translate(object.X()+float64(offsetX), object.Y()+float64(offsetY))
			} else {
				op.GeoM.Translate(object.X(), object.Y())
			}
			screen.DrawImage(object.Image(), op)
		}
	}

	op = &ebiten.DrawImageOptions{}
	screen.DrawImage(fg, op)
	if gameover {
		switch goCause {
		case "start":
			op := &ebiten.DrawImageOptions{}
			screen.DrawImage(startScreen, op)
		case "collision":
			op := &ebiten.DrawImageOptions{}
			screen.DrawImage(collisionScreen, op)
		case "key":
			op := &ebiten.DrawImageOptions{}
			screen.DrawImage(keyScreen, op)
		case "win":
			op := &ebiten.DrawImageOptions{}
			screen.DrawImage(winScreen, op)
		}
	} else {
		op := &ebiten.DrawImageOptions{}
		if bottomCoverDraw {
			screen.DrawImage(bottomCover, op)
		}
		if middleCoverDraw {
			screen.DrawImage(middleCover, op)
		}
		if rightCoverDraw {
			screen.DrawImage(rightCover, op)
		}
	}
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("%d", len(world.players)))
}

var audioContext = audio.NewContext(48000)
var jump *audio.Player
var keyUp *audio.Player
var keyDown *audio.Player
var teleportSound *audio.Player
var playerSpawn *audio.Player
var ghostSpawn *audio.Player
var startSound *audio.Player
var timeUp *audio.Player
var switchSound *audio.Player
var winSound *audio.Player
var lostSound *audio.Player
var song1 *audio.Player
var song2 *audio.Player
var song3 *audio.Player
var song4 *audio.Player

func initSounds() {
	dat, _ := os.ReadFile("resources/sounds/jump.wav")
	d, _ := wav.DecodeWithoutResampling(bytes.NewReader(dat))
	jump, _ = audioContext.NewPlayer(d)

	dat, _ = os.ReadFile("resources/sounds/keypickup.wav")
	d, _ = wav.DecodeWithoutResampling(bytes.NewReader(dat))
	keyUp, _ = audioContext.NewPlayer(d)

	dat, _ = os.ReadFile("resources/sounds/key_down.wav")
	d, _ = wav.DecodeWithoutResampling(bytes.NewReader(dat))
	keyDown, _ = audioContext.NewPlayer(d)

	dat, _ = os.ReadFile("resources/sounds/teleport.wav")
	d, _ = wav.DecodeWithoutResampling(bytes.NewReader(dat))
	teleportSound, _ = audioContext.NewPlayer(d)

	dat, _ = os.ReadFile("resources/sounds/playerspawn.wav")
	d, _ = wav.DecodeWithoutResampling(bytes.NewReader(dat))
	playerSpawn, _ = audioContext.NewPlayer(d)

	dat, _ = os.ReadFile("resources/sounds/ghostspawn.wav")
	d, _ = wav.DecodeWithoutResampling(bytes.NewReader(dat))
	ghostSpawn, _ = audioContext.NewPlayer(d)

	dat, _ = os.ReadFile("resources/sounds/start.wav")
	d, _ = wav.DecodeWithoutResampling(bytes.NewReader(dat))
	startSound, _ = audioContext.NewPlayer(d)

	dat, _ = os.ReadFile("resources/sounds/timesup.wav")
	d, _ = wav.DecodeWithoutResampling(bytes.NewReader(dat))
	timeUp, _ = audioContext.NewPlayer(d)

	dat, _ = os.ReadFile("resources/sounds/switch.wav")
	d, _ = wav.DecodeWithoutResampling(bytes.NewReader(dat))
	switchSound, _ = audioContext.NewPlayer(d)

	dat, _ = os.ReadFile("resources/sounds/win.wav")
	d, _ = wav.DecodeWithoutResampling(bytes.NewReader(dat))
	winSound, _ = audioContext.NewPlayer(d)

	dat, _ = os.ReadFile("resources/sounds/lost.wav")
	d, _ = wav.DecodeWithoutResampling(bytes.NewReader(dat))
	lostSound, _ = audioContext.NewPlayer(d)

	dat, _ = os.ReadFile("resources/sounds/song1.wav")
	d, _ = wav.DecodeWithoutResampling(bytes.NewReader(dat))
	song1, _ = audioContext.NewPlayer(d)

	dat, _ = os.ReadFile("resources/sounds/song2.wav")
	d, _ = wav.DecodeWithoutResampling(bytes.NewReader(dat))
	song2, _ = audioContext.NewPlayer(d)

	dat, _ = os.ReadFile("resources/sounds/song3.wav")
	d, _ = wav.DecodeWithoutResampling(bytes.NewReader(dat))
	song3, _ = audioContext.NewPlayer(d)

	dat, _ = os.ReadFile("resources/sounds/song4.wav")
	d, _ = wav.DecodeWithoutResampling(bytes.NewReader(dat))
	song4, _ = audioContext.NewPlayer(d)

	song1.SetVolume(0.5)
	song2.SetVolume(0.5)
	song3.SetVolume(0.5)
	song4.SetVolume(0.5)
	jump.SetVolume(0.5)
}

var lastPlayed = 0

func playMusic() {
	song := rand.Intn(4)
	for song == lastPlayed {
		song = rand.Intn(4)
	}
	lastPlayed = song
	switch song {
	case 0:
		song1.Rewind()
		song1.Play()
	case 1:
		song2.Rewind()
		song2.Play()
	case 2:
		song3.Rewind()
		song3.Play()
	case 3:
		song4.Rewind()
		song4.Play()
	}
}

func stopMusic() {
	song1.Pause()
	song2.Pause()
	song3.Pause()
	song4.Pause()
}
