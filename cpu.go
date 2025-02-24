package main

import (
	"cpu-emulator/decoder"
	"log"
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
		regs.sp = make16BitAddr(b1, b2)
	}
}

func (regs *registers) getPair(reg string) uint16 {
	switch reg {
	case "B":
		return make16BitAddr(regs.b, regs.c)
	case "D":
		return make16BitAddr(regs.d, regs.e)
	default:
		log.Fatal()
	}
	return 0
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

}
