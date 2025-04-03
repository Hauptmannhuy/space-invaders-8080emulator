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
	width  = 256
	height = 224
)

type spaceInvadersMachine struct {
	cpu   *machine.Cpu
	timer time.Time
}

type gameIO struct{}

func (gIO *gameIO) InPort(cpu *machine.Cpu) uint8 {
	var accum uint8
	pc := cpu.GetPC()
	port := cpu.GetMemoryAt(pc + 1)
	portVal := cpu.Ports[port]

	switch port {
	case 1:
		accum = portVal
	case 2:

	case 3:

	}
	fmt.Println("performing specific IN")
	return accum
}

func (gameMachine *spaceInvadersMachine) sendToPort(ports *machine.Ports, input string) {
	switch input {
	case fire:
		ports[0] = 0x04
	case left:
		ports[0] = 0b00100000
	case right:
		ports[0] = 0b01000000
	case insertCoin:
		ports[0] = 0
	case start:
		ports[0] = 0x02
	}
}

func handleInput() string {
	var keyBuff []ebiten.Key
	key := inpututil.AppendPressedKeys(keyBuff)
	if len(key) > 0 {
		return key[0].String()
	}
	return ""
}

func (gameMachine *spaceInvadersMachine) Update() error {
	input := handleInput()
	if input != "" {
		gameMachine.sendToPort(&gameMachine.cpu.Ports, input)
	}
	// elapsed := time.Since(gameMachine.cpu.LastInterrupt)

	// if gameMachine.cpu.ExternalSourceInt() {
	// 	gameMachine.cpu.HandleExternalInt()
	// }

	// if elapsed > 1/60 {
	// 	gameMachine.cpu.GenerateInterrupt()
	// }
	gameMachine.cpu.Step()

	return nil
}

func (gameMachine *spaceInvadersMachine) Draw(screen *ebiten.Image) {
	buffer := gameMachine.cpu.CopyFrameBuffer()
	bitmap := make([]byte, width*height*4)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			p := buffer[x+y]
			for i := 0; i < 8; i++ {
				index := (y*width + x) * 4

				if p&0x1 == 1 {
					bitmap[index+0] = 255
					bitmap[index+1] = 255
					bitmap[index+2] = 255
					bitmap[index+3] = 255
				} else {
					bitmap[index+0] = 0
					bitmap[index+1] = 0
					bitmap[index+2] = 0
					bitmap[index+3] = 255
				}

			}
		}
	}
	screen.WritePixels(bitmap)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%f", ebiten.ActualFPS()))
}

func (gameMachine *spaceInvadersMachine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func Start(cpu *machine.Cpu) {
	gameMachine := &spaceInvadersMachine{
		cpu: cpu,
	}
	cpu.IO_handler = &gameIO{}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(gameMachine); err != nil {
		log.Fatal(err)
	}
}
