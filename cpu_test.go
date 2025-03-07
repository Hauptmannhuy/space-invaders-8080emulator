package main

import (
	"testing"
)

func runCpuInstruction(name string, cpu *cpu, test *testing.T, fn func(*testing.T, *cpu)) {

	test.Run(name, func(*testing.T) {
		fn(test, cpu)
	})
	cpu.flags.ac = 0
	cpu.flags.p = 0
	cpu.flags.s = 0
	cpu.flags.cy = 0
	cpu.regs.a = 0
	cpu.regs.b = 0
	cpu.regs.c = 0
	cpu.regs.d = 0
	cpu.regs.h = 0
	cpu.regs.l = 0
	cpu.regs.pc = 0
	cpu.regs.sp = 0
	cpu.memory = &memory{}
}

func TestRar(t *testing.T) {

	cpuState := initCPU()
	runCpuInstruction("rotate right through carry", cpuState, t, func(t *testing.T, c *cpu) {
		expected := uint8(0x9)
		oldAccum := uint8(0x12)
		cpuState.regs.a = oldAccum
		cpuState.rar()
		if cpuState.regs.a != expected {
			t.Errorf("Value of reg A should be %d, not %d", expected, cpuState.regs.a)
		}
	})

}

func TestMov(t *testing.T) {
	cpuState := initCPU()
	cpuState.romBuffer = make([]byte, 10)
	var pc uint16
	var desiredVal uint8
	runCpuInstruction("MOV B,C", cpuState, t, func(t *testing.T, c *cpu) {
		desiredVal = 5
		cpuState.romBuffer[pc] = 0x41 // mov B,C
		cpuState.regs.c = desiredVal
		cpuState.step()
		if cpuState.regs.b != desiredVal {
			t.Errorf("Expected value should be %d, instead got %d", desiredVal, cpuState.regs.b)
		}
	})

	runCpuInstruction("MOV B,M", cpuState, t, func(t *testing.T, c *cpu) {
		desiredVal = 15
		cpuState.romBuffer[pc] = 0x46 // mov B,M
		cpuState.regs.h = 0x1
		cpuState.regs.l = 0x2
		addr := cpuState.regs.getPair("M")
		cpuState.memory.write(addr, desiredVal)
		cpuState.step()
		res := cpuState.regs.b
		if res != desiredVal {
			t.Errorf("Expected value should be %d, got %d", desiredVal, res)
		}
	})

	runCpuInstruction("MOV M,C", cpuState, t, func(t *testing.T, c *cpu) {
		desiredVal = 254
		cpuState.romBuffer[pc] = 0x71 // mov M,C
		cpuState.regs.c = desiredVal
		cpuState.step()
		res := cpuState.memory.read(cpuState.regs.getPair("M"))
		if res != desiredVal {
			t.Errorf("Expected value should be %d, got %d", desiredVal, res)
		}
	})
}
