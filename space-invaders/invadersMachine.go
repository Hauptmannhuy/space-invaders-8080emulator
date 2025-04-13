package spacegameMachine

import (
	"cpu-emulator/machine"
	"fmt"
	"log"
	"sync"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	fire       = sdl.K_SPACE
	left       = sdl.K_a
	right      = sdl.K_d
	insertCoin = sdl.K_s
	start      = sdl.K_w
)

const (
	RGB_ON  uint32 = 0xFFFFFFFF
	RGB_OFF uint32 = 0x000000FF
)
const (
	width  = 224
	height = 256
)

type spaceInvadersMachine struct {
	cpu *machine.Cpu

	bitmap []byte

	cyclesRan          uint64
	lastInterruptCycle uint64
	whichInterrupt     int

	pause     uint8
	syncPause *sync.WaitGroup
}

type gameIO struct {
	shift0      uint8 //LSB of Space Invader's external shift hardware
	shift1      uint8 //MSB
	shiftOffset uint8 //offset for external shift hardware
}

func Main(cpu *machine.Cpu) {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("space invaders", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
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

	go keyboardUpdate(gameMachine, &running)
	go gameMachine.internalUpdate()

	for running {
		if gameMachine.pause == 2 {
			gameMachine.syncPause.Wait()
		}

		renderer.Clear()
		updateTexture(texture, gameMachine)
		renderer.Copy(texture, nil, nil)
		renderer.Present()
		sdl.Delay(16)
	}
}

func initEmulation(cpu *machine.Cpu) *spaceInvadersMachine {
	gameMachine := &spaceInvadersMachine{
		cpu:            cpu,
		whichInterrupt: 1,
		bitmap:         make([]byte, width*height*4),
		syncPause:      &sync.WaitGroup{},
	}
	cpu.InterruptEnabled = true
	cpu.IO_handler = &gameIO{}
	return gameMachine
}

func (gameMachine *spaceInvadersMachine) internalUpdate() {
	// ticker := time.NewTicker(time.Second)
	// defer ticker.Stop()

	// go func(gm *spaceInvadersMachine) {
	// 	for {
	// 		<-ticker.C
	// 		fmt.Println("cycles run", gameMachine.cyclesRan)
	// 		fmt.Println("elapsed 1 second")
	// 		ticker.Reset(time.Second)

	// 	}
	// }(gameMachine)

	for {
		if gameMachine.pause == 2 {
			fmt.Println("pause")
			gameMachine.syncPause.Wait()
		}
		gameMachine.cpu.Step()
		op := gameMachine.cpu.GetCurrentOP()
		cycles := op.Cycles

		if gameMachine.cyclesRan-gameMachine.lastInterruptCycle > 33_333 {
			if gameMachine.cpu.InterruptEnabled {

				gameMachine.cpu.GenerateInterrupt(gameMachine.whichInterrupt)

				// fmt.Println("interrupt")
				gameMachine.lastInterruptCycle = gameMachine.cyclesRan
				gameMachine.whichInterrupt ^= 3
			}
		}

		gameMachine.cyclesRan += uint64(cycles)
	}
}

func keyboardUpdate(gameMachine *spaceInvadersMachine, running *bool) {
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch ev := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				*running = false
				break
			case *sdl.KeyboardEvent:
				gameMachine.handleKey(ev)

				if ev.Keysym.Sym == sdl.K_p {
					if ev.State == sdl.PRESSED {
						if gameMachine.pause == 0 {
							gameMachine.pause = 2
							gameMachine.syncPause.Add(1)
						} else if gameMachine.pause == 2 {
							gameMachine.syncPause.Done()
							gameMachine.pause = 0
						}
						break
					}
				}

			}
		}
	}

}

func updateTexture(texture *sdl.Texture, gameMachine *spaceInvadersMachine) {
	buffer := gameMachine.cpu.CopyFrameBuffer()

	for x := 0; x < 224; x++ {
		for y := 0; y < 256; y += 8 {
			p := buffer[(x*(256/8))+y/8]
			offset := (255-y)*(224*4) + (x * 4)
			ptr := (*uint32)(unsafe.Pointer(&gameMachine.bitmap[offset]))

			for i := 0; i < 8; i++ {
				if p&0x1 == 1 {
					*ptr = RGB_ON
				} else {
					*ptr = RGB_OFF
				}

				ptr = (*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) - 224*4))
				p >>= 1
			}
		}
	}

	pixels, _, err := texture.Lock(nil)

	if err != nil {
		log.Fatal(err)
	}

	copy(pixels, gameMachine.bitmap)
	texture.Unlock()
}

func (io *gameIO) InPort(cpu *machine.Cpu) uint8 {

	shift0 := io.shift0
	shift1 := io.shift1
	shiftOffset := io.shiftOffset
	var accum uint8
	pc := cpu.GetPC()
	port := cpu.GetMemoryAt(pc + 1)
	portVal := cpu.Ports[port]

	switch port {
	case 0:
		return 1
	case 1:
		accum = portVal
	case 3:
		v := (uint16(shift1) << 8) | uint16(shift0)
		accum = uint8((v >> (8 - shiftOffset)) & 0xff)
	}
	return accum
}

func (io *gameIO) OutPort(cpu *machine.Cpu) {

	port := cpu.GetMemoryAt(cpu.GetPC() + 1)
	accum := cpu.GetAccumulator()
	switch port {
	case 2:
		io.shiftOffset = accum & 0x7
	case 4:
		io.shift0 = io.shift1
		io.shift1 = accum
	}
}

func getSetBit(input sdl.Keycode) uint8 {
	switch input {
	case fire:
		return 0x10
	case left:
		return 0x20
	case right:
		return 0x40
	case insertCoin:
		return 0x01
	case start:
		return 0x04
	default:
		fmt.Println("Unknown input")
		return 0
	}
}

func getClearBit(input sdl.Keycode) uint8 {
	switch input {
	case fire:
		return 0xef
	case left:
		return 0xdf
	case right:
		return 0xbf
	case insertCoin:
		return 0
	case start:
		return 0xfb
	default:
		fmt.Println("Unknown input")
		return 0
	}
}

func (gameMachine *spaceInvadersMachine) handleKey(ev *sdl.KeyboardEvent) {
	var result uint8
	port := &gameMachine.cpu.Ports[1]

	if ev.State == sdl.PRESSED {
		result = *port | getSetBit(ev.Keysym.Sym)
	} else {
		result = *port & getClearBit(ev.Keysym.Sym)
	}

	*port = result
}
