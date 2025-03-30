package machine

import (
	misc "cpu-emulator/utils"
	"log"
)

const (
	A_REG   uint8 = 0b111
	B_REG   uint8 = 0b000
	C_REG   uint8 = 0b001
	D_REG   uint8 = 0b010
	E_REG   uint8 = 0b011
	H_REG   uint8 = 0b100
	L_REG   uint8 = 0b101
	MEM_REG uint8 = 0b110
)

const (
	BC_REG uint8 = 0b00
	DE_REG uint8 = 0b01
	HL_REG uint8 = 0b10
)

const (
	SP_REG uint8 = 0b11
	PC_REG uint8 = 0b100
)

const PSW uint8 = 0b1010

type registers struct {
	b, c, d, e, h, l, a uint8
}

func (cpu *Cpu) updatePairRegs(register uint8, msb, lsb uint8) {
	switch register & 0b11 {
	case BC_REG:
		cpu.regs.b = msb
		cpu.regs.c = lsb
	case DE_REG:
		cpu.regs.d = msb
		cpu.regs.e = lsb
	case HL_REG:
		cpu.regs.h = msb
		cpu.regs.l = lsb
	case SP_REG:
		cpu.sp = misc.Make16bit(msb, lsb)
	default:
		log.Fatalf("Invalid register pair %d", register)
	}
}

func (cpu *Cpu) getPair(dest uint8) uint16 {
	switch dest & 0b11 {
	case BC_REG:
		return misc.Make16bit(cpu.regs.b, cpu.regs.c)
	case DE_REG:
		return misc.Make16bit(cpu.regs.d, cpu.regs.e)
	case HL_REG:
		return misc.Make16bit(cpu.regs.h, cpu.regs.l)
	case SP_REG:
		return cpu.sp
	default:
		log.Fatalf("Unknown pair %s", dest)
	}
	return 0
}

func (cpu *Cpu) GetReg(reg uint8) uint8 {
	switch reg & 0b111 {
	case B_REG:
		return cpu.regs.b
	case C_REG:
		return cpu.regs.c
	case D_REG:
		return cpu.regs.d
	case E_REG:
		return cpu.regs.e
	case H_REG:
		return cpu.regs.h
	case L_REG:
		return cpu.regs.l
	case A_REG:
		return cpu.regs.a
	case MEM_REG:
		addr := cpu.getPair(reg)
		return cpu.memory.read(addr)
	}
	return 0
}

func (cpu *Cpu) updateReg(reg uint8, val uint8) {
	switch reg & 0b111 {
	case B_REG:
		cpu.regs.b = val
	case C_REG:
		cpu.regs.c = val
	case D_REG:
		cpu.regs.d = val
	case E_REG:
		cpu.regs.e = val
	case H_REG:
		cpu.regs.h = val
	case L_REG:
		cpu.regs.l = val
	case A_REG:
		cpu.regs.a = val
	case MEM_REG:
		addr := cpu.getPair(HL_REG)
		cpu.memory.write(addr, val)
	default:
		log.Fatalf("Invalid register %d", reg)
	}
}
