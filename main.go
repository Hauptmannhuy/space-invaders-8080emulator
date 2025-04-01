package main

import (
	"cpu-emulator/machine"
	spaceinvaders "cpu-emulator/space-invaders"
	"flag"
	"fmt"
	"os"
)

const defaultPath string = "./roms/space-invaders.rom"

func main() {
	// setFlags()

	// args := os.Args[1:]
	// flag := os.Args[3]
	// cpu := machine.InitCpu()
	// buff, err := os.ReadFile(args[1])
	// if err != nil {
	// 	log.Panic(err)
	// }
	// cpu.LoadRom(buff)

	// if flag == "-rd" || flag == "-d" {
	// 	dbg := machine.InitDebugger(flag)
	// 	dbg.Debug(cpu)
	// }

	// cpu := machine.InitCpu()
	// buff, _ := os.ReadFile("./roms/space-invaders.rom")
	// cpu.LoadRom(buff)
	// dbg := machine.InitDebugger("-d")
	// dbg.Debug(cpu)

	cpu := machine.InitCpu()
	buff, _ := os.ReadFile("./roms/space-invaders.rom")
	cpu.LoadRom(buff)
	spaceinvaders.Start(cpu)
}

func setFlags() {
	rFlag := flag.String("r", "", "Path to the ROM file")
	drFlag := flag.Bool("rd", true, "debug with remote debugger")
	dFlag := flag.Bool("d", true, "default debug")

	flag.Parse()

	if *rFlag == "" {
		fmt.Println("Usage: go run . -r <path>")
		os.Exit(1)
	}

	if !*drFlag && !*dFlag {
		fmt.Println("Usage: go run . -r <path> (-rd | -r)")
		os.Exit(1)
	}
}
