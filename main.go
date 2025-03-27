package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const defaultPath string = "./roms/space-invaders.rom"

func main() {
	setFlags()

	args := os.Args[1:]
	flag := os.Args[3]
	cpu := initCPU()
	buff, err := os.ReadFile(args[1])
	if err != nil {
		log.Panic(err)
	}
	cpu.memory.loadRom(buff)

	if flag == "-rd" || flag == "-d" {
		dbg := initDebugger(flag)
		dbg.debug(cpu)
	}

	// buff, _ := os.ReadFile("./roms/cpudiag.bin")
	// cpu.memory.loadRom(buff)
	// dbg := initDebugger("-d")
	// dbg.debug(cpu)
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
