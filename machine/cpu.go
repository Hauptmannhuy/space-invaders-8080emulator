package machine

import (
	"cpu-emulator/decoder"
	misc "cpu-emulator/utils"
	"fmt"
	"math/bits"
	"os"
)

type Cpu struct {
	memory           *Memory
	regs             *registers
	flags            *flags
	sp               uint16 // stack pointer
	pc               uint16 //  program counter
	interruptEnabled bool
	currentOp        *decoder.Opcode
}

func InitCpu() *Cpu {
	cpu := &Cpu{
		memory: &Memory{},
		regs:   &registers{},
		flags:  &flags{},
	}
	return cpu
}

func (cpu *Cpu) ResetCpu() {

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
	cpu.pc = 0
	cpu.sp = 0
	cpu.memory = &Memory{}
}

func (cpu *Cpu) executeInstruction() uint8 {

	switch cpu.currentOp.Instruction {
	// data transfer group
	case decoder.MOV:
		return cpu.mov()
	case decoder.LXI:
		return cpu.lxi()
	case decoder.MVI:
		return cpu.mvi()
	case decoder.LDA:
		return cpu.lda()
	case decoder.STA:
		return cpu.sta()
	case decoder.LHLD:
		return cpu.lhld()
	case decoder.SHLD:
		return cpu.shld()
	case decoder.STAX:
		return cpu.stax()
	case decoder.LDAX:
		return cpu.ldax()
	case decoder.XCHG:
		return cpu.xchg()

	// arithmetic group
	case decoder.ADD:
		return cpu.add()
	case decoder.ADI:
		return cpu.add()
	case decoder.ADC:
		return cpu.add()
	case decoder.SUB:
		return cpu.sub()
	case decoder.SBB:
		return cpu.sub()
	case decoder.SUI:
		return cpu.sub()
	case decoder.SBI:
		return cpu.sub()
	case decoder.ACI:
		return cpu.add()
	case decoder.DAD:
		return cpu.dad()
	case decoder.RAL:
		return cpu.ral()
	case decoder.RAR:
		return cpu.rar()
	case decoder.RLC:
		return cpu.rlc()
	case decoder.RRC:
		return cpu.rrc()
	case decoder.CMA:
		return cpu.cma()
	case decoder.STC:
		return cpu.stc()
	case decoder.CMC:
		return cpu.cmc()
	case decoder.DAA:
		return cpu.daa()
	case decoder.INR:
		return cpu.inr()
	case decoder.INX:
		return cpu.inx()
	case decoder.DCX:
		return cpu.dcx()
	case decoder.DCR:
		return cpu.dcr()

	// logic group
	case decoder.ANA:
		return cpu.ana()
	case decoder.ANI:
		return cpu.ana()
	case decoder.XRA:
		return cpu.xra()
	case decoder.XRI:
		return cpu.xra()
	case decoder.ORA:
		return cpu.ora()
	case decoder.ORI:
		return cpu.ora()
	case decoder.CMP:
		return cpu.cmp()
	case decoder.CPI:
		return cpu.cmp()

	// branch group
	case decoder.JMP:
		return cpu.jmp()
	case decoder.JC:
		return cpu.jmp()
	case decoder.JNC:
		return cpu.jmp()
	case decoder.JZ:
		return cpu.jmp()
	case decoder.JNZ:
		return cpu.jmp()
	case decoder.JM:
		return cpu.jmp()
	case decoder.JPE:
		return cpu.jmp()
	case decoder.JP:
		return cpu.jmp()
	case decoder.JPO:
		return cpu.jmp()
	case decoder.PCHL:
		return cpu.pchl()
	case decoder.CALL:
		return cpu.call()
	case decoder.CC:
		return cpu.call()
	case decoder.CZ:
		return cpu.call()
	case decoder.CNZ:
		return cpu.call()
	case decoder.CM:
		return cpu.call()
	case decoder.CPE:
		return cpu.call()
	case decoder.CPO:
		return cpu.call()
	case decoder.CNC:
		return cpu.call()
	case decoder.CP:
		return cpu.call()

	case decoder.RET:
		return cpu.ret()
	case decoder.RC:
		return cpu.ret()
	case decoder.RZ:
		return cpu.ret()
	case decoder.RNZ:
		return cpu.ret()
	case decoder.RM:
		return cpu.ret()
	case decoder.RP:
		return cpu.ret()
	case decoder.RPE:
		return cpu.ret()
	case decoder.RPO:
		return cpu.ret()
	case decoder.RNC:
		return cpu.ret()

	case decoder.RST:
		return cpu.rst()

	// stack i/o machine control group
	case decoder.PUSH:
		return cpu.push()
	case decoder.POP:
		return cpu.pop()
	case decoder.XTHL:
		return cpu.xthl()
	case decoder.SPHL:
		return cpu.sphl()
	case decoder.OUT:
		return cpu.out()
	case decoder.IN:
		return cpu.in()
	case decoder.DI:
		return cpu.di()
	case decoder.EI:
		return cpu.ei()
	case decoder.HLT:
		return cpu.hlt()
	case decoder.NOP:
		return cpu.nop()

	case decoder.RIM:
		return cpu.rim()

	//special

	case decoder.SIM:
		return cpu.sim()
	default:
		if cpu.currentOp.Instruction == 0 {
			return 1
		} else {
			panic(fmt.Sprintf("Instruction %d is not found, pc: x%02x", cpu.currentOp.Instruction, cpu.pc))
		}
	}
}

func (cpu *Cpu) Step() {
	cpu.currentOp = getOpcode(&cpu.memory[cpu.pc])
	n := cpu.executeInstruction()
	cpu.pc += uint16(n)
}

func getOpcode(code *byte) *decoder.Opcode {
	return decoder.GetInstruction(code)
}

func (cpu *Cpu) nop() uint8 {
	return 1
}

func (cpu *Cpu) lxi() uint8 {
	msb := cpu.memory.read(cpu.pc + 2)
	lsb := cpu.memory.read(cpu.pc + 1)
	cpu.updatePairRegs(cpu.currentOp.HighNibble, msb, lsb)
	return 3
}

func (cpu *Cpu) stax() uint8 {
	reg := cpu.currentOp.HighNibble
	addr := cpu.getPair(reg)
	accumVal := cpu.regs.a
	cpu.memory.write(addr, accumVal)
	return 1
}

func (cpu *Cpu) inx() uint8 {
	reg := cpu.currentOp.HighNibble
	bits := cpu.getPair(reg) + 1
	b1 := uint8(bits >> 8)
	b2 := uint8(bits & 0xff)
	cpu.updatePairRegs(reg, b1, b2)
	return 1
}

func (cpu *Cpu) inr() uint8 {
	reg := cpu.currentOp.Code >> 3
	regVal := cpu.GetReg(reg)
	res := regVal + 1
	cpu.updateReg(reg, res)
	cpu.setAux((regVal&0x0F == 0x0F))
	cpu.updateFlags(res, 0)
	return 1
}

func (cpu *Cpu) dcr() uint8 {
	reg := cpu.currentOp.Code >> 3
	regVal := cpu.GetReg(reg)
	res := regVal - 1
	cpu.updateReg(reg, res)
	cpu.setAux((regVal&0x0F == 0x00))
	cpu.updateFlags(res, 0)
	return 1
}

func (cpu *Cpu) mvi() uint8 {
	immediate := cpu.memory[cpu.pc+1]
	cpu.updateReg((cpu.currentOp.Code >> 3), immediate)
	return 2
}

func (cpu *Cpu) dad() uint8 {
	regPair := cpu.currentOp.HighNibble
	hl := cpu.getPair(HL_REG)
	value := cpu.getPair(regPair)

	res := hl + value
	sum := uint32(hl) + uint32(value)

	if sum > 65535 {
		cpu.flags.cy = 1
	} else {
		cpu.flags.cy = 0
	}

	cpu.updatePairRegs(HL_REG, uint8(res>>8), uint8(res))

	return 1
}

func (cpu *Cpu) ldax() uint8 {
	reg := cpu.currentOp.HighNibble
	addr := cpu.getPair(reg)
	memVal := cpu.memory.read(addr)
	cpu.updateReg(A_REG, memVal)
	return 1
}

func (cpu *Cpu) dcx() uint8 {
	reg := cpu.currentOp.HighNibble
	pairVal := cpu.getPair(reg) - 1
	cpu.updatePairRegs(reg, uint8(pairVal>>8), uint8(pairVal&0xff))
	return 1
}

func (cpu *Cpu) rlc() uint8 {
	accum := cpu.regs.a

	res := bits.RotateLeft8(accum, 1)
	cpu.regs.a = res

	if (res & 0x01) == 0x01 {
		cpu.flags.cy = 1
	} else {
		cpu.flags.cy = 0
	}

	return 1
}

func (cpu *Cpu) rrc() uint8 {
	accum := cpu.regs.a

	res := bits.RotateLeft8(accum, -1)
	cpu.regs.a = res

	if res&0x80 == 0x80 {
		cpu.flags.cy = 1
	} else {
		cpu.flags.cy = 0
	}

	return 1
}

// A = A << 1; bit 0 = prev CY; CY = prev bit 7
func (cpu *Cpu) ral() uint8 {
	accum := cpu.regs.a
	carry := cpu.flags.cy

	if (accum & 0x80) == 0x80 {
		cpu.flags.cy = 1
	} else {
		cpu.flags.cy = 0
	}

	accum = accum << 1
	accum = accum | carry

	cpu.regs.a = accum

	return 1
}

// A = A >> 1; bit 7 = prev bit 7; CY = prev bit 0
func (cpu *Cpu) rar() uint8 {
	accum := cpu.regs.a
	carry := cpu.flags.cy

	if (accum & 0x01) == 0x01 {
		cpu.flags.cy = 1
	} else {
		cpu.flags.cy = 0
	}

	accum = accum >> 1
	accum = accum | (carry << 7)

	cpu.regs.a = accum
	return 1
}

// some I/O thing
func (cpu *Cpu) rim() uint8 {
	return 1
}

func (cpu *Cpu) shld() uint8 {
	addr := misc.Make16bit(cpu.memory.read(cpu.pc+2), cpu.memory.read(cpu.pc+1))
	cpu.memory.write(addr, cpu.regs.l)
	cpu.memory.write(addr+1, cpu.regs.h)
	return 3
}

func (cpu *Cpu) lhld() uint8 {
	addr := misc.Make16bit(cpu.memory.read(cpu.pc+2), cpu.memory.read(cpu.pc+1))
	l := cpu.memory.read(addr)
	h := cpu.memory.read(addr + 1)
	cpu.regs.l = l
	cpu.regs.h = h
	return 3
}

func (cpu *Cpu) cma() uint8 {
	cpu.regs.a = ^cpu.regs.a
	return 1
}

// some useless command
func (cpu *Cpu) daa() uint8 {
	accum := cpu.regs.a

	if accum&0xf > 9 || cpu.flags.ac == 1 {
		accum += 0x06
		if (accum & 0x0F) < 0x09 {
			cpu.flags.ac = 1
		}
	}

	if ((accum & 0xF0) > 0x90) || cpu.flags.cy == 1 {
		accum, cpu.flags.cy = misc.OverflowingAdd(accum, 0x60, 0)
	}

	cpu.regs.a = accum
	// cpu.updateFlags(cpu.regs.a)

	return 1
}

// special
func (cpu *Cpu) sim() uint8 {
	return 1
}

func (cpu *Cpu) sta() uint8 {
	addr := misc.Make16bit(cpu.memory.read(cpu.pc+2), cpu.memory.read(cpu.pc+1))
	cpu.memory.write(addr, cpu.regs.a)
	return 3
}

func (cpu *Cpu) stc() uint8 {
	cpu.flags.cy = 1
	return 1
}

func (cpu *Cpu) lda() uint8 {
	addr := misc.Make16bit(cpu.memory.read(cpu.pc+2), cpu.memory.read(cpu.pc+1))
	cpu.regs.a = cpu.memory.read(addr)
	return 3
}

func (cpu *Cpu) cmc() uint8 {
	cpu.flags.cy = 1 ^ cpu.flags.cy
	return 1
}

func (cpu *Cpu) mov() uint8 {
	sReg := cpu.currentOp.LowNibble
	dReg := (cpu.currentOp.Code >> 3) & 0b111
	sRegVal := cpu.GetReg(sReg)
	cpu.updateReg(dReg, sRegVal)
	return 1
}

// halt
func (cpu *Cpu) hlt() uint8 {
	os.Exit(0)
	return 0
}

func (cpu *Cpu) cmp() uint8 {
	var operand uint8

	if cpu.currentOp.Instruction == decoder.CPI {
		operand = cpu.memory.read(cpu.pc + 1)
	} else {
		operand = cpu.GetReg(cpu.currentOp.LowNibble)
	}

	res, carry := misc.OverflowingSub(cpu.regs.a, operand, 0)
	cpu.updateFlags(res, carry)
	cpu.setAux((cpu.regs.a & 0x0F) < (operand & 0x0F))

	if cpu.currentOp.Instruction == decoder.CPI {
		return 2
	}

	return 1
}

func (cpu *Cpu) xra() uint8 {
	var res uint8
	var operand uint8

	if cpu.currentOp.Instruction == decoder.XRI {
		operand = cpu.memory.read(cpu.pc + 1)
	} else {
		operand = cpu.GetReg(cpu.currentOp.LowNibble)
	}

	res = cpu.regs.a ^ operand
	cpu.regs.a = res

	cpu.updateFlags(res, 0)

	if cpu.currentOp.Instruction == decoder.XRI {
		return 2
	}

	return 1
}

func (cpu *Cpu) ana() uint8 {
	var operand uint8
	var res uint8
	prevAccum := cpu.regs.a

	if cpu.currentOp.Instruction == decoder.ANI {
		operand = cpu.memory.read(cpu.pc + 1)
	} else {
		operand = cpu.GetReg(cpu.currentOp.LowNibble)
	}

	res = prevAccum & operand
	cpu.regs.a = res
	cpu.updateFlags(res, 0)

	if cpu.currentOp.Instruction == decoder.ANI {
		return 2
	}

	return 1
}

func (cpu *Cpu) add() uint8 {
	var res uint8
	var operand uint8
	var carry uint8
	prevAccum := cpu.regs.a

	if cpu.currentOp.Instruction == decoder.ADI || cpu.currentOp.Instruction == decoder.ACI {
		operand = cpu.memory.read(cpu.pc + 1)
	} else {
		operand = cpu.GetReg(cpu.currentOp.LowNibble)
	}

	switch cpu.currentOp.Instruction {
	case decoder.ADD:
		res, carry = misc.OverflowingAdd(prevAccum, operand, 0)
	case decoder.ADC:
		res, carry = misc.OverflowingAdd(prevAccum, operand, cpu.flags.cy)
	case decoder.ADI:
		res, carry = misc.OverflowingAdd(prevAccum, operand, 0)
	case decoder.ACI:
		res, carry = misc.OverflowingAdd(prevAccum, operand, cpu.flags.cy)
	}

	cpu.regs.a = res
	cpu.updateFlags(res, carry)
	cpu.setAux((prevAccum&0x0F)+(operand&0x0F) > 0x0F)

	if cpu.currentOp.Instruction == decoder.ADI || cpu.currentOp.Instruction == decoder.ACI {
		return 2
	}

	return 1
}

func (cpu *Cpu) sub() uint8 {
	var res uint8
	var carry uint8
	var operand uint8
	prevAccum := cpu.regs.a

	if cpu.currentOp.Instruction == decoder.SUI || cpu.currentOp.Instruction == decoder.SBI {
		operand = cpu.memory.read(cpu.pc + 1)
	} else {
		operand = cpu.GetReg(cpu.currentOp.LowNibble)
	}

	switch cpu.currentOp.Instruction {
	case decoder.SUB:
		res, carry = misc.OverflowingSub(prevAccum, operand, 0)
	case decoder.SBB:
		res, carry = misc.OverflowingSub(prevAccum, operand, cpu.flags.cy)
	case decoder.SUI:
		res, carry = misc.OverflowingSub(prevAccum, operand, 0)
	case decoder.SBI:
		res, carry = misc.OverflowingSub(prevAccum, operand, cpu.flags.cy)
	}

	cpu.regs.a = res
	cpu.updateFlags(res, carry)
	cpu.setAux((prevAccum & 0x0F) < (operand & 0x0F))
	if cpu.currentOp.Instruction == decoder.SUI || cpu.currentOp.Instruction == decoder.SBI {
		return 2
	}

	return 1
}

func (cpu *Cpu) ora() uint8 {
	var res uint8
	var operand uint8

	if cpu.currentOp.Instruction == decoder.ORI {
		operand = cpu.memory.read(cpu.pc + 1)
	} else {
		operand = cpu.GetReg(cpu.currentOp.LowNibble)
	}

	res = cpu.regs.a | operand

	cpu.regs.a = res
	cpu.updateFlags(res, 0)

	if cpu.currentOp.Instruction == decoder.ORI {
		return 2
	}

	return 1
}

func (cpu *Cpu) call() uint8 {
	if misc.Make16bit(cpu.memory.read(cpu.pc+2), cpu.memory.read(cpu.pc+1)) == decoder.BDOS {
		if cpu.regs.c == 0x9 {
			addr := cpu.getPair(DE_REG)
			var msg []byte
			for {
				char := cpu.memory.read(addr)
				msg = append(msg, char)
				addr++
				if string(char) == "$" {
					break
				}
			}
			fmt.Printf("OUTPUT MESSAGE: %s\n", msg)
			os.Exit(1)
		}
	} else {
		if cpu.currentOp.Condition == 0 || cpu.checkConditionFlag() {
			lsb := uint8(cpu.memory[cpu.pc+1])
			msb := uint8(cpu.memory[cpu.pc+2])
			addr := misc.Make16bit(msb, lsb)

			nextAddr := cpu.pc + 3
			lsbNextAddr := uint8(nextAddr & 0x00FF)
			msbNextAddr := uint8((nextAddr & 0xFF00) >> 8)

			cpu.memory.write(cpu.sp-1, msbNextAddr)
			cpu.memory.write(cpu.sp-2, lsbNextAddr)

			cpu.sp = cpu.sp - 2
			cpu.pc = addr

			return 0
		}
		return 3
	}
	return 3
}

func (cpu *Cpu) jmp() uint8 {
	if cpu.currentOp.Condition == 0 || cpu.checkConditionFlag() {
		cpu.pc = misc.Make16bit(cpu.memory.read(cpu.pc+2), cpu.memory.read(cpu.pc+1))

		return 0
	}

	return 3
}

func (cpu *Cpu) ret() uint8 {
	if cpu.currentOp.Condition == 0 || cpu.checkConditionFlag() {
		var addr uint16

		lsb := cpu.memory.read(cpu.sp)
		msb := cpu.memory.read(cpu.sp + 1)
		addr = uint16(uint16(lsb) | uint16(msb)<<8)
		cpu.sp += 2
		cpu.pc = addr
		return 0
	}

	return 1
}

func (cpu *Cpu) push() uint8 {
	sp := cpu.sp
	reg := cpu.currentOp.HighNibble

	pairVal := cpu.getPair(reg)
	lsb := uint8(pairVal)
	msb := uint8((pairVal & 0xff00) >> 8)

	if reg&0b11 == SP_REG {
		cpu.memory.write(sp-1, cpu.regs.a)
		cpu.memory.write(sp-2, cpu.flags.cy|cpu.flags.p<<2|cpu.flags.ac<<4|cpu.flags.z<<6|cpu.flags.s<<7)
	} else {
		cpu.memory.write(sp-1, msb)
		cpu.memory.write(sp-2, lsb)
	}
	cpu.sp -= 2

	return 1
}

func (cpu *Cpu) pop() uint8 {
	sp := cpu.sp
	reg := cpu.currentOp.HighNibble
	msb := cpu.memory.read(sp + 1)
	lsb := cpu.memory.read(sp)

	if reg&0b11 == SP_REG {
		psw := cpu.memory.read(sp)
		cpu.flags.cy = psw & 0x1
		cpu.flags.p = (psw >> 2) & 0x1
		cpu.flags.ac = (psw >> 4) & 0x1
		cpu.flags.z = (psw >> 6) & 0x1
		cpu.flags.s = (psw >> 7) & 0x1
		cpu.regs.a = cpu.memory.read(sp + 1)
	} else {
		cpu.updatePairRegs(reg, msb, lsb)
	}

	cpu.sp += 2
	return 1
}

func (cpu *Cpu) sphl() uint8 {
	cpu.sp = cpu.getPair(HL_REG)
	return 1
}

func (cpu *Cpu) xthl() uint8 {
	lVal := cpu.regs.l
	hVal := cpu.regs.h
	sp1Addr := cpu.sp
	sp2Addr := cpu.sp + 1
	cpu.regs.l = cpu.memory.read(cpu.sp)
	cpu.regs.h = cpu.memory.read(cpu.sp + 1)
	cpu.memory.write(sp1Addr, lVal)
	cpu.memory.write(sp2Addr, hVal)
	return 1
}

func (cpu *Cpu) xchg() uint8 {
	h, d, l, e := cpu.regs.h, cpu.regs.d, cpu.regs.l, cpu.regs.e
	cpu.regs.h = d
	cpu.regs.d = h
	cpu.regs.l = e
	cpu.regs.e = l
	return 1
}

func (cpu *Cpu) rst() uint8 {
	resetAddr := cpu.currentOp.Code & 0b00111000
	cpu.memory.write(cpu.sp-1, uint8(cpu.pc&0x0f))
	cpu.memory.write(cpu.sp-2, uint8((cpu.pc>>8)&0x0f))
	cpu.sp -= 2
	cpu.pc = uint16(resetAddr - 1)

	return 0
}

func (cpu *Cpu) in() uint8 {
	return 2
}

func (cpu *Cpu) out() uint8 {
	return 2
}

func (cpu *Cpu) ei() uint8 {
	cpu.interruptEnabled = true
	return 1
}

func (cpu *Cpu) di() uint8 {
	cpu.interruptEnabled = false
	return 1
}

func (cpu *Cpu) pchl() uint8 {
	cpu.pc = misc.Make16bit(cpu.regs.h, cpu.regs.l)
	return 0
}
