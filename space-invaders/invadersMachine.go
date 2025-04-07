package spacegameMachine

import (
	"cpu-emulator/machine"
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	fire       = "Enter"
	left       = "ArrowLeft"
	right      = "ArrowRight"
	insertCoin = "S"
	start      = "W"
)

const (
	width  = 224
	height = 256
)

type spaceInvadersMachine struct {
	cpu *machine.Cpu

	lastInterrupt  time.Time
	whichInterrupt int
}

type gameIO struct {
	shift0      uint8 //LSB of Space Invader's external shift hardware
	shift1      uint8 //MSB
	shiftOffset uint8 //offset for external shift hardware
	gameMachine *spaceInvadersMachine
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

func handleInput() string {
	var keyBuff []ebiten.Key
	pressed := inpututil.AppendPressedKeys(keyBuff)
	if len(pressed) > 0 {
		return pressed[0].String()
	}

	return ""
}

func handleRelease() string {
	var keyBuff []ebiten.Key
	released := inpututil.AppendJustReleasedKeys(keyBuff)
	if len(released) > 0 {
		return released[0].String()
	}
	return ""
}

func (gameMachine *spaceInvadersMachine) Update() error {
	gameMachine.cpu.Step()

	if gameMachine.lastInterrupt.IsZero() {
		gameMachine.lastInterrupt = time.Now()
	}

	if time.Since(gameMachine.lastInterrupt) >= time.Second/60 {
		if gameMachine.cpu.InterruptEnabled {
			fmt.Println("interrupt")
			gameMachine.cpu.GenerateInterrupt(gameMachine.whichInterrupt)
		}

		gameMachine.lastInterrupt = time.Now()

		if gameMachine.whichInterrupt == 1 {
			gameMachine.whichInterrupt = 2
		} else {
			gameMachine.whichInterrupt = 1
		}
	}

	input := handleInput()
	if input != "" {
		fmt.Println(gameMachine.cpu.Ports)
		gameMachine.machineKeyPressed(&gameMachine.cpu.Ports, input)
	}
	input = handleRelease()
	if input != "" {
		gameMachine.machineKeyReleased(&gameMachine.cpu.Ports, input)
	}

	time.Sleep(gameMachine.cpu.CycleRun)
	return nil
}

func (gameMachine *spaceInvadersMachine) Draw(screen *ebiten.Image) {
	buffer := gameMachine.cpu.CopyFrameBuffer()
	bitmap := make([]byte, width*height*4)
	for x := 0; x < 224; x++ {
		for y := 0; y < 256; y += 8 {
			p := buffer[(x*(256/8))+y/8]
			offset := (255-y)*(224*4) + (x * 4)
			for i := 0; i < 8; i++ {
				if p&0x1 == 1 {
					bitmap[offset+0] = 255
					bitmap[offset+1] = 255
					bitmap[offset+2] = 255
				} else {
					bitmap[offset+0] = 0
					bitmap[offset+1] = 0
					bitmap[offset+2] = 0
				}
				bitmap[offset+3] = 255
				p <<= 1
			}
		}
	}
	screen.WritePixels(bitmap)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%f", ebiten.ActualFPS()))
}

func (gameMachine *spaceInvadersMachine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
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

func Start(cpu *machine.Cpu) {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Space Invaders")
	gameMachine := initEmulation(cpu)
	if err := ebiten.RunGame(gameMachine); err != nil {
		log.Fatal(err)
	}
}
