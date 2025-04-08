package main

import (
	"cpu-emulator/machine"
	spacegameMachine "cpu-emulator/space-invaders"
	"flag"
	"fmt"
	"log"
	"os"
)

const defaultPath string = "./roms/space-invaders.rom"

func main() {
	setFlags()
	var buffer []byte
	cpu := machine.InitCpu()

	args := os.Args[1:]

	if len(args) > 1 {
		flag := os.Args[3]
		fileBuffer, err := os.ReadFile(args[1])
		buffer = fileBuffer
		if err != nil {
			log.Panic(err)
		}

		if flag == "-rd" || flag == "-d" {
			dbg := machine.InitDebugger(flag)
			dbg.Debug(cpu)
			return
		}

	} else {
		fileBuffer, err := os.ReadFile(defaultPath)
		if err != nil {
			log.Panic(err)
		}
		buffer = fileBuffer
	}

	cpu.LoadRom(buffer)

	spacegameMachine.Main(cpu)
}

func setFlags() {
	usageText := "Usage: go run . -r <path> (-rd | -r)"

	rFlag := flag.String("r", "", "path to the ROM file")
	drFlag := flag.Bool("rd", true, "debug with remote debugger")
	dFlag := flag.Bool("d", true, "default debug")
	pFlag := flag.Bool("p", true, "play space invaders")
	flag.Parse()

	if *pFlag {
		if *rFlag != "" {
			fmt.Println(usageText)
			os.Exit(1)
		} else if !*drFlag && !*dFlag {
			fmt.Println(usageText)
			os.Exit(1)
		}
	}
}
