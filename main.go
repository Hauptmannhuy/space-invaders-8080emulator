package main

import (
	"time"
)

func main() {
	debugLogic()
}

func debugLogic() {
	cpu := initCPU()
	cpu.pc = 0x100
	for cpu.pc < uint16(MemorySize) {
		pc := cpu.pc
		// time.Sleep(200 * time.Millisecond)
		disassebmle(cpu.memory, int(pc))
		debugCpuState(cpu)
		nextStep()
		cpu.step()
	}
}

func debugInstruction() {
	cpu := initCPU()
	pc := 0
	for pc < len(cpu.memory) {
		cpu.executeInstruction()
		n := disassebmle(cpu.memory, int(pc))
		pc += n
		time.Sleep(150 * time.Millisecond)
	}
}
