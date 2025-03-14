package main

import (
	"cpu-emulator/decoder"
	"log"
)

func main() {
	debug()
}

func start() {
	cpu := initCPU()
	cpu.romBuffer = loadHex()
	for cpu.regs.pc < uint16(len(cpu.romBuffer)) {
		pc := cpu.regs.pc
		disassebmle(cpu.romBuffer, int(pc))
		// time.Sleep(150 * time.Millisecond)
		cpu.step()
	}
}

func debug() {
	cpu := initCPU()
	cpuInstructions := initOpcodeSet(cpu)
	rom := loadHex()
	pc := 0
	for pc < len(rom) {
		code := rom[pc]
		currOp := decoder.GetInstruction(code)

		_, ok := cpuInstructions[currOp.Instruction]
		if !ok {
			log.Fatalf("Instruction %s is not implemented", currOp.Instruction)
		}
		n := disassebmle(rom, int(pc))
		pc += n
		// time.Sleep(150 * time.Millisecond)
	}
}
