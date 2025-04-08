package spacegameMachine

import (
	"cpu-emulator/machine"
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	fire       = "Enter"
	left       = "ArrowLeft"
	right      = "ArrowRight"
	insertCoin = "S"
	start      = "W"
)

const (
	RGB_ON  uint32 = 0xFFFFFFFF // Белый
	RGB_OFF uint32 = 0x000000FF // Чёрный (A=255, но R=G=B=0)
)
const (
	width  = 224
	height = 256
)

type spaceInvadersMachine struct {
	cpu *machine.Cpu

	timer          time.Time
	lastInterrupt  time.Time
	whichInterrupt int
}

type gameIO struct {
	shift0      uint8 //LSB of Space Invader's external shift hardware
	shift1      uint8 //MSB
	shiftOffset uint8 //offset for external shift hardware
	gameMachine *spaceInvadersMachine
}

func Main(cpu *machine.Cpu) {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("space invaders", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		panic(err)
	}

	gameMachine := initEmulation(cpu)

	loop(window, texture, gameMachine)
}

func loop(window *sdl.Window, texture *sdl.Texture, gameMachine *spaceInvadersMachine) {

	renderer, _ := window.GetRenderer()
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent: // NOTE: Please use `*sdl.QuitEvent` for `v0.4.x` (current version).
				println("Quit")
				running = false
				break
			case *sdl.KeyboardEvent:

			}
		}

		renderer.Clear()
		gameMachine.cpu.Step()
		updateTexture(texture, gameMachine)
		gameMachine.internalUpdate()
		renderer.Copy(texture, nil, nil)
		renderer.Present()
		// sdl.Delay(16)
	}
}

func initEmulation(cpu *machine.Cpu) *spaceInvadersMachine {
	gameMachine := &spaceInvadersMachine{
		cpu:            cpu,
		whichInterrupt: 1,
	}
	cpu.InterruptEnabled = true
	cpu.IO_handler = &gameIO{}
	return gameMachine
}

func (gameMachine *spaceInvadersMachine) internalUpdate() {

	if gameMachine.lastInterrupt.IsZero() {
		gameMachine.lastInterrupt = time.Now()
		gameMachine.timer = time.Now()
	}

	if gameMachine.timer.Sub(gameMachine.lastInterrupt) >= time.Second/60 {
		if gameMachine.cpu.InterruptEnabled {
			err := gameMachine.cpu.GenerateInterrupt(gameMachine.whichInterrupt)
			if err != nil {
				time.Sleep(gameMachine.cpu.CycleRun)
				return
			}
			fmt.Println("interrupt")
			gameMachine.lastInterrupt = time.Now()
			if gameMachine.whichInterrupt == 1 {
				gameMachine.whichInterrupt = 2
			} else {
				gameMachine.whichInterrupt = 1
			}
		}
	}
	gameMachine.timer = time.Now()
}

func updateTexture(texture *sdl.Texture, gameMachine *spaceInvadersMachine) {
	buffer := gameMachine.cpu.CopyFrameBuffer()
	bitmap := make([]byte, width*height*4)
	for x := 0; x < 224; x++ {
		for y := 0; y < 256; y += 8 {
			p := buffer[(x*(256/8))+y/8]
			offset := y*(224*4) + x*4 // Начинаем сверху

			ptr := (*uint32)(unsafe.Pointer(&bitmap[offset]))
			for i := 0; i < 8; i++ {
				if p&0x1 == 1 {
					os.Exit(2)
					*ptr = RGB_ON
				} else {
					*ptr = RGB_OFF
				}
				ptr = (*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + 224*4))
				p >>= 1
			}
		}
	}
	pixels, _, err := texture.Lock(nil)
	if err != nil {
		log.Fatal(err)
	}
	copy(pixels, bitmap)
	texture.Unlock()
}

func (gIO *gameIO) InPort(cpu *machine.Cpu) uint8 {
	shift0 := gIO.shift0
	shift1 := gIO.shift1
	shiftOffset := gIO.shiftOffset
	var accum uint8
	pc := cpu.GetPC()
	port := cpu.GetMemoryAt(pc + 1)
	portVal := cpu.Ports[port]

	switch port {
	case 1:
		accum = portVal
	case 2:

	case 3:
		v := (uint16((shift1<<8)|shift0) >> (8 - uint16(shiftOffset)) & 0xff)
		accum = uint8(v)
	}
	return accum
}

func (gIO *gameIO) OutPort(cpu *machine.Cpu) {
	port := cpu.Ports[cpu.GetMemoryAt(cpu.GetPC()+1)]
	val := cpu.Ports[port]
	switch port {
	case 2:
		gIO.shiftOffset = val & 0x7
	case 4:
		gIO.shift0 = gIO.shift1
		gIO.shift1 = val
	}
}

func (gameMachine *spaceInvadersMachine) machineKeyPressed(ports *machine.Ports, input string) {
	switch input {
	case fire:
		ports[1] &= 0xf7
	case left:
		ports[1] &= 0xdf
	case right:
		ports[1] &= 0xbf
	case insertCoin:
		ports[1] &= 0
	case start:
		ports[1] &= 0xfd
	}
}

func (gameMachine *spaceInvadersMachine) machineKeyReleased(ports *machine.Ports, input string) {
	switch input {
	case fire:
		ports[1] |= 0x04
	case left:
		ports[1] |= 0b0100000
	case right:
		ports[1] |= 0b1000000
	case insertCoin:
		ports[1] |= 0
	case start:
		ports[1] |= 0x02
	}
}

// func handleInput() string {
// 	var keyBuff []ebiten.Key
// 	pressed := inpututil.AppendPressedKeys(keyBuff)
// 	if len(pressed) > 0 {
// 		return pressed[0].String()
// 	}

// 	return ""
// }

// func handleRelease() string {
// 	var keyBuff []ebiten.Key
// 	released := inpututil.AppendJustReleasedKeys(keyBuff)
// 	if len(released) > 0 {
// 		return released[0].String()
// 	}
// 	return ""
// }
