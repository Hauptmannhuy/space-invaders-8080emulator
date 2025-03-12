package main

import (
	"cpu-emulator/decoder"
	"log"
	"math/bits"
	"os"
)

type cpu struct {
	memory         *memory
	regs           *registers
	flags          *flags
	instructionSet map[string]func() int
	romBuffer      []byte
	opcode         *decoder.Opcode
}

type registers struct {
	b, c, d, e, l, h, a uint8
	sp                  uint16 // stack pointer
	pc                  uint16 //  program counter
}

// flags Z (zero), S (sign), P (parity), CY (carry), CA (auxillary  carry)
type flags struct {
	z, s, p, cy, ac uint8
}

func initCPU() *cpu {
	cpu := &cpu{
		memory: &memory{},
		regs:   &registers{},
		flags:  &flags{},
	}
	cpu.instructionSet = initOpcodeSet(cpu)
	return cpu
}

func initOpcodeSet(cpu *cpu) map[string]func() int {
	return map[string]func() int{
		"NOP": cpu.nop,
		"LXI": cpu.lxi,
		// "STAX":  func() {},
		// "INX":   func() {},
		// "INR B": func() {},
		"MOV": cpu.mov,
		"ADD": cpu.srcRegOperationSet,
		"SUB": cpu.srcRegOperationSet,
		"CMP": cpu.srcRegOperationSet,
		"CPI": cpu.immediateOperationSet,
	}
}

func (cpu *cpu) step() {
	cpu.opcode = getOpcode(cpu.romBuffer[cpu.regs.pc])
	handleFunc := cpu.instructionSet[cpu.opcode.Instruction]
	n := handleFunc()
	cpu.regs.pc += uint16(n)
}

func getOpcode(code byte) *decoder.Opcode {
	return decoder.GetInstruction(code)
}

func (regs *registers) writePairRegs(reg string, b1, b2 uint8) {
	switch reg {
	case "B":
		regs.b = b1
		regs.c = b2
	case "D":
		regs.d = b1
		regs.e = b2
	case "H":
		regs.h = b1
		regs.l = b2
	case "SP":
		regs.sp = make16bit(b1, b2)
	}
}

func (regs *registers) getPair(dest string) uint16 {
	switch dest {
	case "B":
		return make16bit(regs.b, regs.c)
	case "D":
		return make16bit(regs.d, regs.e)
	default:
		if dest == "HL" || dest == "M" || dest == "H" {
			return make16bit(regs.h, regs.l)
		}
		log.Fatal()
	}
	return 0
}

func (regs *registers) getReg(reg string) uint8 {
	switch reg {
	case "B":
		return regs.b
	case "C":
		return regs.c
	case "D":
		return regs.d
	case "E":
		return regs.e
	case "H":
		return regs.h
	case "L":
		return regs.l
	case "A":
		return regs.a
	}
	return 0
}

func (regs *registers) updateReg(reg string, val uint8) {
	switch reg {
	case "B":
		regs.b = val
	case "C":
		regs.c = val
	case "D":
		regs.d = val
	case "E":
		regs.e = val
	case "H":
		regs.h = val
	case "L":
		regs.l = val
	case "A":
		regs.a = val
	}
}

func (cpu *cpu) checkConditionFlag() bool {

	switch cpu.opcode.Condition {
	case decoder.Minus:
		{
			if cpu.flags.s == 1 {
				return true
			}
			return false
		}
	case decoder.Positive:
		{
			if cpu.flags.s == 0 {
				return true
			}
			return false
		}
	case decoder.Carry:
		{
			if cpu.flags.cy == 1 {
				return true
			}
			return false
		}
	case decoder.NoCarry:
		{
			if cpu.flags.cy == 0 {
				return true
			}
			return false
		}
	case decoder.NotZero:
		{
			if cpu.flags.z == 1 {
				return true
			}
			return false
		}
	case decoder.Zero:
		{
			if cpu.flags.z == 0 {
				return true
			}
			return false
		}
	case decoder.ParityOdd:
		{
			if cpu.flags.p == 1 {
				return true
			}
			return false
		}
	case decoder.ParityEven:
		{
			if cpu.flags.p == 0 {
				return true
			}
			return false
		}
	default:
		log.Fatal("Screwed up condition flag")
		return false
	}
}

func (cpu *cpu) readROM(n int) byte {
	pc := cpu.regs.pc
	return cpu.romBuffer[pc+uint16(n)]
}

func (cpu *cpu) nop() int {
	return 1
}

func (cpu *cpu) lxi() int {
	b1 := cpu.readROM(2)
	b2 := cpu.readROM(1)
	cpu.regs.writePairRegs(cpu.opcode.Operand, b1, b2)
	return 3
}

func (cpu *cpu) stax() int {
	reg := cpu.opcode.Operand
	addr := cpu.regs.getPair(reg)
	accumVal := cpu.regs.a
	cpu.memory.write(addr, accumVal)
	return 1
}

func (cpu *cpu) inx() int {
	reg := cpu.opcode.Operand
	bits := cpu.regs.getPair(reg) + 1
	b1 := uint8(bits >> 8)
	b2 := uint8(bits & 0xff)
	cpu.regs.writePairRegs(reg, b1, b2)
	return 1
}

func (cpu *cpu) inr() int {
	reg := cpu.opcode.Operand
	var val uint8
	if reg == "M" {
		addr := cpu.regs.getPair(reg)
		val = cpu.memory.read(addr)
		val += 1
		cpu.memory.write(addr, val)
	} else {
		val = cpu.regs.getReg(reg) + 1
		cpu.regs.updateReg(reg, val)
	}
	cpu.flags.updateZSPAC(val)
	return 1
}

func (cpu *cpu) dcr() int {
	reg := cpu.opcode.Operand
	var val uint8
	if reg == "M" {
		addr := cpu.regs.getPair(reg)
		val = cpu.memory.read(addr)
		val -= 1
		cpu.memory.write(addr, val)
	} else {
		val = cpu.regs.getReg(reg) - 1
		cpu.regs.updateReg(reg, val)
	}
	cpu.flags.updateZSPAC(val)
	return 1
}

func (cpu *cpu) mvi() int {
	byte2 := cpu.romBuffer[cpu.regs.pc+1]
	if cpu.opcode.Operand == "M" {
		addr := cpu.regs.getPair(cpu.opcode.Operand)
		cpu.memory.write(addr, byte2)
	} else {
		cpu.regs.updateReg(cpu.opcode.Operand, byte2)
	}
	return 2
}

func (cpu *cpu) rlc() int {
	reg := "A"
	regVal := cpu.regs.getReg(reg)
	prev7Bit := regVal & 0xf
	regVal = regVal << 1
	regVal = regVal | prev7Bit
	cpu.regs.updateReg(reg, regVal)
	cpu.flags.cy = prev7Bit
	return 1
}

func (cpu *cpu) dad() int {
	destReg := cpu.opcode.Operand
	hl := cpu.regs.getPair("HL")
	value := cpu.regs.getPair(destReg)
	sum := hl + value
	if sum < 0xfff {
		cpu.flags.cy = 1
	} else {
		cpu.flags.cy = 0
	}
	return 1
}

func (cpu *cpu) ldax() int {
	reg := cpu.opcode.Operand
	pairVal := cpu.regs.getPair(reg)
	memVal := cpu.memory.read(pairVal)
	cpu.regs.updateReg(reg, memVal)
	return 1
}

func (cpu *cpu) dcx() int {
	pairVal := cpu.regs.getPair(cpu.opcode.Operand) - 1
	cpu.regs.writePairRegs(cpu.opcode.Operand, uint8(pairVal>>8), uint8(pairVal&0xff))
	return 1
}

func (cpu *cpu) rrc() int {
	regVal := cpu.regs.getReg(cpu.opcode.Operand)
	shifted := regVal >> 1
	prev0Bit := regVal & 1
	shifted = (shifted & 0x7f) | (prev0Bit << 7)
	cpu.regs.updateReg(cpu.opcode.Operand, shifted)
	cpu.flags.cy = prev0Bit
	return 1
}

// A = A << 1; bit 0 = prev CY; CY = prev bit 7
func (cpu *cpu) ral() int {
	regVal := cpu.regs.getReg(cpu.opcode.Operand)
	newVal := regVal << 1
	newVal = newVal | cpu.flags.cy
	cpu.flags.cy = regVal >> 7 & 1
	cpu.regs.updateReg(cpu.opcode.Operand, newVal)
	return 1
}

// A = A >> 1; bit 7 = prev bit 7; CY = prev bit 0
func (cpu *cpu) rar() int {
	regVal := cpu.regs.a
	prevBit7 := regVal & (1 << 7)
	newVal := (prevBit7 | regVal>>1)
	cpu.regs.a = newVal
	cpu.flags.cy = regVal & 1
	return 1
}

// some I/O thing
func (cpu *cpu) rim() int {
	return 1
}

func (cpu *cpu) shld() int {
	addr := make16bit(cpu.readROM(2), cpu.readROM(1))
	cpu.memory.write(addr, cpu.regs.l)
	cpu.memory.write(addr+1, cpu.regs.h)
	return 3
}

func (cpu *cpu) lhld() int {
	addr := make16bit(cpu.readROM(2), cpu.readROM(1))
	l := cpu.memory.read(addr)
	h := cpu.memory.read(addr + 1)
	cpu.regs.l = l
	cpu.regs.h = h
	return 3
}

func (cpu *cpu) cma() int {
	cpu.regs.a = ^cpu.regs.a
	return 1
}

// some useless command
func (cpu *cpu) daa() int {
	return 1
}

// special
func (cpu *cpu) sim() int {
	return 1
}

func (cpu *cpu) sta() int {
	addr := make16bit(cpu.readROM(2), cpu.readROM(1))
	cpu.memory.write(addr, cpu.regs.a)
	return 3
}

func (cpu *cpu) stc() int {
	cpu.flags.cy = 1
	return 1
}

func (cpu *cpu) lda() int {
	addr := make16bit(cpu.readROM(2), cpu.readROM(1))
	cpu.regs.a = cpu.memory.read(addr)
	return 3
}

func (cpu *cpu) cmc() int {
	cpu.flags.cy = ^cpu.flags.cy
	return 1
}

func (cpu *cpu) mov() int {
	destReg := string(cpu.opcode.Operand[0])
	sourceReg := string(cpu.opcode.Operand[len(cpu.opcode.Operand)-1])
	if destReg == "M" {
		addr := cpu.regs.getPair(destReg)
		cpu.memory.write(addr, cpu.regs.getReg(sourceReg))
		return 1
	} else if sourceReg == "M" {
		addr := cpu.regs.getPair(sourceReg)
		memVal := cpu.memory.read(addr)
		cpu.regs.updateReg(destReg, memVal)
		return 1
	}

	cpu.regs.updateReg(destReg, cpu.regs.getReg(sourceReg))
	return 1
}

// halt
func (cpu *cpu) hlt() int {
	os.Exit(0)
	return 0
}

func (flags *flags) updateCY(prev, new uint8) {

	if prev > new {
		flags.cy = 1
	} else {
		flags.cy = 0
	}
}

func (cpu *cpu) srcRegOperationSet() int {
	var OperandSelector uint8
	var newAccum uint8
	accumulator := cpu.regs.a

	if cpu.opcode.Operand == "M" {
		OperandSelector = cpu.memory.read(cpu.regs.getPair("M"))
	} else {
		OperandSelector = cpu.regs.getReg(cpu.opcode.Operand)
	}

	switch cpu.opcode.Instruction {
	case "ADD":
		newAccum = accumulator + OperandSelector
	case "ADC":
		newAccum = accumulator + OperandSelector + cpu.flags.cy
	case "SUB":
		newAccum = accumulator - OperandSelector
	case "SBB":
		newAccum = accumulator - OperandSelector - cpu.flags.cy
	case "ANA":
		newAccum = accumulator & OperandSelector
	case "XRA":
		newAccum = accumulator ^ OperandSelector
	case "ORA":
		newAccum = accumulator | OperandSelector
	case "CMP":
		newAccum = OperandSelector
	}

	cpu.flags.updateZSPAC(newAccum)
	cpu.flags.updateCY(accumulator, newAccum)

	if cpu.opcode.Instruction != "CMP" {
		cpu.regs.a = newAccum
	}

	return 1
}

func (cpu *cpu) immediateOperationSet() int {
	var newAccum uint8
	prevAccum := cpu.regs.a

	switch cpu.opcode.Instruction {
	case "ADI":
		newAccum = prevAccum + cpu.readROM(1)
	case "ACI":
		newAccum = prevAccum + cpu.flags.cy + cpu.readROM(1)
	case "SUI":
		newAccum = prevAccum - cpu.readROM(1)
	case "SBI":
		newAccum = prevAccum - cpu.readROM(1)
	case "ANI":
		newAccum = prevAccum & cpu.readROM(1)
	case "XRI":
		newAccum = prevAccum ^ cpu.readROM(1)
	case "ORI":
		newAccum = prevAccum | cpu.readROM(1)
	case "CPI":
		newAccum = cpu.readROM(1)
	}

	cpu.flags.updateZSPAC(newAccum)
	cpu.flags.updateCY(prevAccum, newAccum)

	if cpu.opcode.Instruction != "CPI" {
		cpu.regs.a = newAccum
	}
	return 2
}

func (cpu *cpu) call() int {
	if cpu.opcode.Condition == "" || cpu.checkConditionFlag() {
		sp := cpu.regs.sp
		hiByte := cpu.regs.pc >> 8
		loByte := uint8(cpu.regs.pc)

		cpu.memory.write(sp-1, uint8(hiByte))
		cpu.memory.write(sp-2, loByte)

		cpu.regs.sp = cpu.regs.sp + 2
		cpu.regs.pc = make16bit(cpu.romBuffer[cpu.regs.pc+2], cpu.romBuffer[cpu.regs.pc+1])
	}

	return 3
}

func (cpu *cpu) jmp() int {
	if cpu.checkConditionFlag() || cpu.opcode.Condition == "" {
		cpu.regs.pc = make16bit(cpu.readROM(2), cpu.readROM(1))
	}

	return 3
}

func (cpu *cpu) ret() int {
	if cpu.opcode.Condition == "" || cpu.checkConditionFlag() {
		var newPC uint16

		lowByte := cpu.memory.read(cpu.regs.sp)
		highByte := cpu.memory.read(cpu.regs.sp + 1)
		newPC = (uint16(highByte) << 8) | uint16(lowByte)

		cpu.regs.sp += cpu.regs.sp + 2
		cpu.regs.pc = newPC
	}

	return 1
}

func (flags *flags) updateZSPAC(val uint8) {
	if val == 0 {
		flags.z = 1
	} else {
		flags.z = 0
	}

	if val >= 128 {
		flags.s = 1
	} else {
		flags.s = 0
	}
	flags.p = parity(val)
}

func (cpu *cpu) pchl() int {
	cpu.regs.pc = make16bit(cpu.regs.h, cpu.regs.l)
	return 1
}

func parity(val uint8) uint8 {
	count := bits.OnesCount8(val)
	if count%2 == 0 {
		return 1
	} else {
		return 0
	}
}
