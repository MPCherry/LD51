package resources

import (
	_ "embed"
)

var (
	//go:embed ghostspawn.wav
	GhostSpawn_wav []byte

	//go:embed jump.wav
	Jump_wav []byte

	//go:embed key_down.wav
	KeyDown_wav []byte

	//go:embed keypickup.wav
	KeyUp_wav []byte

	//go:embed land.wav
	Land_wav []byte

	//go:embed lost.wav
	Lost_wav []byte

	//go:embed playerspawn.wav
	PlayerSpawn_wav []byte

	//go:embed song1.wav
	Song1_wav []byte

	//go:embed song2.wav
	Song2_wav []byte

	//go:embed song3.wav
	Song3_wav []byte

	//go:embed song4.wav
	Song4_wav []byte

	//go:embed sound.wav
	Sound_wav []byte

	//go:embed start.wav
	Start_wav []byte

	//go:embed switch.wav
	Switch_wav []byte

	//go:embed teleport.wav
	Teleport_wav []byte

	//go:embed timesup.wav
	TimesUp_wav []byte

	//go:embed win.wav
	Win_wav []byte
)
