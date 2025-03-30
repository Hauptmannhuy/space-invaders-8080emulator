package main

// import (
// 	"cpu-emulator/machine"
// 	"fmt"
// 	"testing"
// )

// func runCpuInstruction(name string, cpu *machine.Cpu, test *testing.T, fn func(*testing.T, *machine.Cpu)) {

// 	test.Run(name, func(*testing.T) {
// 		fn(test, cpu)
// 	})
// 	cpu.ResetCpu()
// }

// func TestRar(t *testing.T) {

// 	cpuState := machine.InitCpu()
// 	runCpuInstruction("rotate right through carry", cpuState, t, func(t *testing.T, c *cpu) {
// 		expected := uint8(0x9)
// 		oldAccum := uint8(0x12)
// 		cpuState.regs.a = oldAccum
// 		cpuState.rar()
// 		if cpuState.regs.a != expected {
// 			t.Errorf("Value of reg A should be %d, not %d", expected, cpuState.regs.a)
// 		}
// 	})

// }

// func TestMov(t *testing.T) {
// 	cpuState := machine.InitCpu()
// 	var pc uint16
// 	var desiredVal uint8
// 	runCpuInstruction("MOV B,C", cpuState, t, func(t *testing.T, c *cpu) {
// 		desiredVal = 5
// 		cpuState.memory[pc] = 0x41 // mov B,C
// 		cpuState.regs.c = desiredVal
// 		cpuState.step()
// 		if cpuState.regs.b != desiredVal {
// 			t.Errorf("Expected value should be %d, instead got %d", desiredVal, cpuState.regs.b)
// 		}
// 	})

// 	runCpuInstruction("MOV B,M", cpuState, t, func(t *testing.T, c *cpu) {
// 		desiredVal = 15
// 		cpuState.memory[pc] = 0x46 // mov B,M
// 		cpuState.regs.h = 0x1
// 		cpuState.regs.l = 0x2
// 		addr := cpuState.getPair("M")
// 		cpuState.memory.write(addr, desiredVal)
// 		cpuState.step()
// 		res := cpuState.regs.b
// 		if res != desiredVal {
// 			t.Errorf("Expected value should be %d, got %d", desiredVal, res)
// 		}
// 	})

// 	runCpuInstruction("MOV M,C", cpuState, t, func(t *testing.T, c *cpu) {
// 		desiredVal = 254
// 		cpuState.memory[pc] = 0x71 // mov M,C
// 		cpuState.regs.c = desiredVal
// 		cpuState.step()
// 		res := cpuState.memory.read(cpuState.getPair("M"))
// 		if res != desiredVal {
// 			t.Errorf("Expected value should be %d, got %d", desiredVal, res)
// 		}
// 	})
// }

// func TestAdd(t *testing.T) {
// 	cpuState := machine.InitCpu()
// 	testInput := []byte{
// 		0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86,
// 		0x87,
// 	}
// 	copy(cpuState.memory[0:], testInput)
// 	regNames := []string{"B", "C", "D", "E", "H", "L", "M"}
// 	expected := []uint8{4, 5, 6, 7, 8, 9, 8}
// 	cpuState.regs.b = 4
// 	cpuState.regs.c = 1
// 	cpuState.regs.d = 1
// 	cpuState.regs.e = 1
// 	cpuState.regs.h = 1
// 	cpuState.regs.l = 1
// 	addr := cpuState.getPair("M")
// 	cpuState.memory.write(addr, 255)
// 	for i := 0; i < len(regNames); i++ {
// 		regName := regNames[i]
// 		opcodeName := fmt.Sprintf("ADD %s", regName)
// 		t.Run(opcodeName, func(t *testing.T) {
// 			cpuState.step()
// 			accumRegVal := cpuState.regs.a
// 			if accumRegVal != expected[i] {
// 				fmt.Println(i)
// 				t.Errorf("Expected value should be %d, got %d", expected[i], accumRegVal)
// 			}
// 		})
// 	}
// }

// func TestCPI(t *testing.T) {
// 	cpuState := machine.InitCpu()
// 	testInput := []byte{0xfe, 0x1, 0xfe, 129, 0xfe, 0}
// 	copy(cpuState.memory[0:], testInput)
// 	cpuState.regs.a = 0x2
// 	s := []uint8{0, 1, 0}
// 	z := []uint8{0, 0, 1}
// 	p := []uint8{0, 1, 0}
// 	var j int
// 	var i int
// 	for i < 4 {
// 		n := cpuState.memory[i]
// 		name := fmt.Sprintf("CPI %d", n)
// 		t.Run(name, func(t *testing.T) {
// 			cpuState.step()
// 			if cpuState.flags.z != z[j] {
// 				t.Errorf("Zero flag should be %d, got %d", z[j], cpuState.flags.z)
// 			}

// 			if cpuState.flags.s != s[j] {
// 				t.Errorf("Sign flag should be %d, got %d", s[j], cpuState.flags.s)
// 			}

// 			if cpuState.flags.p != p[j] {
// 				t.Errorf("Parity flag should be %d, got %d", p[j], cpuState.flags.p)
// 			}
// 			j += 1
// 			i += 2
// 		})
// 	}
// }
