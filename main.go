package main

import (
	"cpu-emulator/decoder"
	"log"
	"time"
)

func main() {
	start()
}

func start() {
	cpu := initCPU()
	cpu.romBuffer = loadHex()
	for cpu.regs.pc < uint16(len(cpu.romBuffer)) {
		pc := cpu.regs.pc
		time.Sleep(150 * time.Millisecond)
		cpu.step()
		disassebmle(cpu.romBuffer, int(pc))
		debugCpuState(cpu)
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
			if currOp.Instruction != "" {
				log.Fatalf("Instruction %s is not implemented, code 0x%02x", currOp.Instruction, code)
			}
		}
		n := disassebmle(rom, int(pc))
		pc += n
		// time.Sleep(150 * time.Millisecond)
	}
}
