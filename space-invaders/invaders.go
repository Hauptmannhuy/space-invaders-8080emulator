package spaceinvaders

import (
	"cpu-emulator/machine"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	width  = 256
	height = 224
)

type Game struct {
	cpu *machine.Cpu
}

func (g *Game) Update() error {
	g.cpu.Step()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	buffer := g.cpu.CopyFrameBuffer()
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func Start(cpu *machine.Cpu) {

	game := &Game{
		cpu: cpu,
	}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
