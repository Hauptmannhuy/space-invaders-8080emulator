package main

import (
	"cpu-emulator/decoder"
	"log"
	"math/bits"
)

type cpu struct {
	memory         *memory
	regs           *registers
	flags          *flags
	instructionSet map[string]func() uint16
	romBuffer      []byte
	opcode         *opcode
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

type opcode struct {
	code       byte
	opcodeName string
	operand    string
}

func initCPU() *cpu {
	return &cpu{
		memory: &memory{},
		regs:   &registers{},
		flags:  &flags{},
	}
}

func initOpcodeSet(cpu *cpu) map[string]func() int {
	return map[string]func() int{
		"NOP": cpu.nop,
		"LXI": cpu.lxi,
		// "STAX":  func() {},
		// "INX":   func() {},
		// "INR B": func() {},
	}
}

func (cpu *cpu) step() {
	cpu.opcode = getOpcode(cpu.romBuffer[cpu.regs.pc])
	handleFunc := cpu.instructionSet[cpu.opcode.opcodeName]
	n := handleFunc()
	cpu.regs.pc += n
}

func getOpcode(code byte) *opcode {
	instruction := decoder.GetInstruction(code)
	operand := decoder.GetDestination(instruction, code)

	return &opcode{
		opcodeName: instruction,
		operand:    operand,
		code:       code,
	}
}

func (regs *registers) writeRegs(reg string, b1, b2 uint8) {
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
	cpu.regs.writeRegs(cpu.opcode.operand, b1, b2)
	return 3
}

func (cpu *cpu) stax() int {
	reg := cpu.opcode.operand
	addr := cpu.regs.getPair(reg)
	accumVal := cpu.regs.a
	cpu.memory.write(addr, accumVal)
	return 1
}

func (cpu *cpu) inx() int {
	reg := cpu.opcode.operand
	bits := cpu.regs.getPair(reg) + 1
	b1 := uint8(bits >> 8)
	b2 := uint8(bits & 0xff)
	cpu.regs.writeRegs(reg, b1, b2)
	return 1
}

func (cpu *cpu) inr() int {
	reg := cpu.opcode.operand
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
	reg := cpu.opcode.operand
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
	if cpu.opcode.operand == "M" {
		addr := cpu.regs.getPair(cpu.opcode.operand)
		cpu.memory.write(addr, byte2)
	} else {
		cpu.regs.updateReg(cpu.opcode.operand, byte2)
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
	destReg := cpu.opcode.operand
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

func (flags *flags) updateZSPAC(val uint8) {
	if val == 0 {
		flags.z = 1
	} else {
		flags.z = 0
	}

	if val == 128 {
		flags.s = 1
	} else {
		flags.s = 0
	}
	flags.p = parity(val)
}

func parity(val uint8) uint8 {
	count := bits.OnesCount8(val)
	if count%2 == 0 {
		return 1
	} else {
		return 0
	}
}
