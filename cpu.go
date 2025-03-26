package main

import (
	"cpu-emulator/decoder"
	"fmt"
	"log"
	"math/bits"
	"os"
)

type cpu struct {
	memory           *memory
	regs             *registers
	flags            *flags
	sp               uint16 // stack pointer
	pc               uint16 //  program counter
	interruptEnabled bool
	currentOp        *decoder.Opcode
}

type registers struct {
	b, c, d, e, l, h, a uint8
}

// flags Z (zero), S (sign), P (parity), CY (carry), CA (auxillary  carry)
type flags struct {
	s, z, ac, p, cy uint8
}

func initCPU() *cpu {

	cpu := &cpu{
		memory: &memory{},
		regs:   &registers{},
		flags:  &flags{},
	}
	return cpu
}

func (cpu *cpu) executeInstruction() uint8 {

	switch cpu.currentOp.Instruction {
	// data transfer group
	case "MOV":
		return cpu.mov()
	case "LXI":
		return cpu.lxi()
	case "MVI":
		return cpu.mvi()
	case "LDA":
		return cpu.lda()
	case "STA":
		return cpu.sta()
	case "LHLD":
		return cpu.lhld()
	case "SHLD":
		return cpu.shld()
	case "STAX":
		return cpu.stax()
	case "LDAX":
		return cpu.ldax()
	case "XCHG":
		return cpu.xchg()

	// arithmetic group
	case "ADD":
		return cpu.add()
	case "ADI":
		return cpu.add()
	case "ADC":
		return cpu.add()
	case "SUB":
		return cpu.sub()
	case "SBB":
		return cpu.sub()
	case "SUI":
		return cpu.sub()
	case "SBI":
		return cpu.sub()
	case "ACI":
		return cpu.add()
	case "DAD":
		return cpu.dad()
	case "RAL":
		return cpu.ral()
	case "RAR":
		return cpu.rar()
	case "RLC":
		return cpu.rlc()
	case "RRC":
		return cpu.rrc()
	case "CMA":
		return cpu.cma()
	case "STC":
		return cpu.stc()
	case "CMC":
		return cpu.cmc()
	case "DAA":
		return cpu.daa()
	case "INR":
		return cpu.inr()
	case "INX":
		return cpu.inx()
	case "DCX":
		return cpu.dcx()
	case "DCR":
		return cpu.dcr()

	// logic group
	case "ANA":
		return cpu.ana()
	case "ANI":
		return cpu.ana()
	case "XRA":
		return cpu.xra()
	case "XRI":
		return cpu.xra()
	case "ORA":
		return cpu.ora()
	case "ORI":
		return cpu.ora()
	case "CMP":
		return cpu.cmp()
	case "CPI":
		return cpu.cmp()

	// branch group
	case "JMP":
		return cpu.jmp()
	case "PCHL":
		return cpu.pchl()
	case "CALL":
		return cpu.call()
	case "RET":
		return cpu.ret()
	case "RST":
		return cpu.rst()

	// stack i/o machine control group
	case "PUSH":
		return cpu.push()
	case "POP":
		return cpu.pop()
	case "XTHL":
		return cpu.xthl()
	case "SPHL":
		return cpu.sphl()
	case "OUT":
		return cpu.out()
	case "IN":
		return cpu.in()
	case "DI":
		return cpu.di()
	case "EI":
		return cpu.ei()
	case "HLT":
		return cpu.hlt()
	case "NOP":
		return cpu.nop()

	case "RIM":
		return cpu.rim()

	//special

	case "SIM":
		return cpu.sim()
	default:
		if cpu.currentOp.Instruction == "" {
			return 1
		} else {
			panic(fmt.Sprintf("Instruction %s is not found, pc: x%02x", cpu.currentOp.Instruction, cpu.pc))
		}
	}
}

func (cpu *cpu) step() {
	cpu.currentOp = getOpcode(cpu.memory[cpu.pc])
	n := cpu.executeInstruction()
	cpu.pc += uint16(n)
}

func getOpcode(code byte) *decoder.Opcode {
	return decoder.GetInstruction(code)
}

func (cpu *cpu) updatePairRegs(reg string, msb, lsb uint8) {
	switch reg {
	case "B":
		cpu.regs.b = msb
		cpu.regs.c = lsb
	case "D":
		cpu.regs.d = msb
		cpu.regs.e = lsb
	case "H":
		cpu.regs.h = msb
		cpu.regs.l = lsb
	case "SP":
		cpu.sp = make16bit(msb, lsb)
	}
}

func (cpu *cpu) getPair(dest string) uint16 {
	switch dest {
	case "B":
		return make16bit(cpu.regs.b, cpu.regs.c)
	case "D":
		return make16bit(cpu.regs.d, cpu.regs.e)
	default:
		if dest == "HL" || dest == "M" || dest == "H" {
			return make16bit(cpu.regs.h, cpu.regs.l)
		}
		log.Fatalf("Unknown pair %s", dest)
	}
	return 0
}

func (cpu *cpu) getReg(reg string) uint8 {
	switch reg {
	case "B":
		return cpu.regs.b
	case "C":
		return cpu.regs.c
	case "D":
		return cpu.regs.d
	case "E":
		return cpu.regs.e
	case "H":
		return cpu.regs.h
	case "L":
		return cpu.regs.l
	case "A":
		return cpu.regs.a
	case "M":
		addr := cpu.getPair(reg)
		return cpu.memory.read(addr)
	}
	return 0
}

func (cpu *cpu) updateReg(reg string, val uint8) {
	switch reg {
	case "B":
		cpu.regs.b = val
	case "C":
		cpu.regs.c = val
	case "D":
		cpu.regs.d = val
	case "E":
		cpu.regs.e = val
	case "H":
		cpu.regs.h = val
	case "L":
		cpu.regs.l = val
	case "A":
		cpu.regs.a = val
	}
}

func (cpu *cpu) checkConditionFlag() bool {

	switch cpu.currentOp.Condition {
	case decoder.Minus:
		{
			return cpu.flags.s == 1
		}
	case decoder.Positive:
		{
			return cpu.flags.s == 0
		}
	case decoder.Carry:
		{
			return cpu.flags.cy == 1
		}
	case decoder.NoCarry:
		{
			return cpu.flags.cy == 0
		}
	case decoder.NotZero:
		{
			return cpu.flags.z == 0
		}
	case decoder.Zero:
		{
			return cpu.flags.z == 1
		}
	case decoder.ParityOdd:
		{
			return cpu.flags.p == 0
		}
	case decoder.ParityEven:
		{
			return cpu.flags.p == 1
		}
	default:
		log.Fatal("Error! Unknown condition flag")
		return false
	}
}

func (cpu *cpu) nop() uint8 {
	return 1
}

func (cpu *cpu) lxi() uint8 {
	msb := cpu.memory.read(cpu.pc + 2)
	lsb := cpu.memory.read(cpu.pc + 1)
	cpu.updatePairRegs(cpu.currentOp.Register, msb, lsb)
	return 3
}

func (cpu *cpu) stax() uint8 {
	reg := cpu.currentOp.Register
	addr := cpu.getPair(reg)
	accumVal := cpu.regs.a
	cpu.memory.write(addr, accumVal)
	return 1
}

func (cpu *cpu) inx() uint8 {
	reg := cpu.currentOp.Register
	bits := cpu.getPair(reg) + 1
	b1 := uint8(bits >> 8)
	b2 := uint8(bits & 0xff)
	cpu.updatePairRegs(reg, b1, b2)
	return 1
}

func (cpu *cpu) inr() uint8 {
	reg := cpu.currentOp.Register
	var val uint8
	if reg == "M" {
		addr := cpu.getPair(reg)
		val = cpu.memory.read(addr)
		val += 1
		cpu.memory.write(addr, val)
	} else {
		val = cpu.getReg(reg) + 1
		cpu.updateReg(reg, val)
	}
	cpu.updateFlags(val)
	return 1
}

func (cpu *cpu) dcr() uint8 {
	reg := cpu.currentOp.Register
	var val uint8
	if reg == "M" {
		addr := cpu.getPair(reg)
		val = cpu.memory.read(addr)
		val -= 1
		cpu.memory.write(addr, val)
	} else {
		val = cpu.getReg(reg) - 1
		cpu.updateReg(reg, val)
	}
	cpu.updateFlags(val)
	return 1
}

func (cpu *cpu) mvi() uint8 {
	byte2 := cpu.memory[cpu.pc+1]
	if cpu.currentOp.Register == "M" {
		addr := cpu.getPair(cpu.currentOp.Register)
		cpu.memory.write(addr, byte2)
	} else {
		cpu.updateReg(cpu.currentOp.Register, byte2)
	}
	return 2
}

func (cpu *cpu) rlc() uint8 {
	reg := "A"
	regVal := cpu.getReg(reg)
	prev7Bit := regVal & 0xf
	regVal = regVal << 1
	regVal = regVal | prev7Bit
	cpu.updateReg(reg, regVal)
	cpu.flags.cy = prev7Bit
	return 1
}

func (cpu *cpu) dad() uint8 {
	destReg := cpu.currentOp.Register
	hl := cpu.getPair("HL")
	value := cpu.getPair(destReg)
	sum := hl + value
	if sum < 0xfff {
		cpu.flags.cy = 1
	} else {
		cpu.flags.cy = 0
	}
	return 1
}

func (cpu *cpu) ldax() uint8 {
	reg := cpu.currentOp.Register
	pairVal := cpu.getPair(reg)
	memVal := cpu.memory.read(pairVal)
	cpu.updateReg(reg, memVal)
	return 1
}

func (cpu *cpu) dcx() uint8 {
	pairVal := cpu.getPair(cpu.currentOp.Register) - 1
	cpu.updatePairRegs(cpu.currentOp.Register, uint8(pairVal>>8), uint8(pairVal&0xff))
	return 1
}

func (cpu *cpu) rrc() uint8 {
	regVal := cpu.getReg(cpu.currentOp.Register)
	shifted := regVal >> 1
	prev0Bit := regVal & 1
	shifted = (shifted & 0x7f) | (prev0Bit << 7)
	cpu.updateReg(cpu.currentOp.Register, shifted)
	cpu.flags.cy = prev0Bit
	return 1
}

// A = A << 1; bit 0 = prev CY; CY = prev bit 7
func (cpu *cpu) ral() uint8 {
	regVal := cpu.getReg(cpu.currentOp.Register)
	newVal := regVal << 1
	newVal = newVal | cpu.flags.cy
	cpu.flags.cy = regVal >> 7 & 1
	cpu.updateReg(cpu.currentOp.Register, newVal)
	return 1
}

// A = A >> 1; bit 7 = prev bit 7; CY = prev bit 0
func (cpu *cpu) rar() uint8 {
	regVal := cpu.regs.a
	prevBit7 := regVal & (1 << 7)
	newVal := (prevBit7 | regVal>>1)
	cpu.regs.a = newVal
	cpu.flags.cy = regVal & 1
	return 1
}

// some I/O thing
func (cpu *cpu) rim() uint8 {
	return 1
}

func (cpu *cpu) shld() uint8 {
	addr := make16bit(cpu.memory.read(cpu.pc+2), cpu.memory.read(cpu.pc+1))
	cpu.memory.write(addr, cpu.regs.l)
	cpu.memory.write(addr+1, cpu.regs.h)
	return 3
}

func (cpu *cpu) lhld() uint8 {
	addr := make16bit(cpu.memory.read(cpu.pc+2), cpu.memory.read(cpu.pc+1))
	l := cpu.memory.read(addr)
	h := cpu.memory.read(addr + 1)
	cpu.regs.l = l
	cpu.regs.h = h
	return 3
}

func (cpu *cpu) cma() uint8 {
	cpu.regs.a = ^cpu.regs.a
	return 1
}

// some useless command
func (cpu *cpu) daa() uint8 {
	return 1
}

// special
func (cpu *cpu) sim() uint8 {
	return 1
}

func (cpu *cpu) sta() uint8 {
	addr := make16bit(cpu.memory.read(cpu.pc+2), cpu.memory.read(cpu.pc+1))
	cpu.memory.write(addr, cpu.regs.a)
	return 3
}

func (cpu *cpu) stc() uint8 {
	cpu.flags.cy = 1
	return 1
}

func (cpu *cpu) lda() uint8 {
	addr := make16bit(cpu.memory.read(cpu.pc+2), cpu.memory.read(cpu.pc+1))
	cpu.regs.a = cpu.memory.read(addr)
	return 3
}

func (cpu *cpu) cmc() uint8 {
	cpu.flags.cy = ^cpu.flags.cy
	return 1
}

func (cpu *cpu) mov() uint8 {
	destReg := string(cpu.currentOp.Register[0])
	sourceReg := string(cpu.currentOp.Register[len(cpu.currentOp.Register)-1])
	if destReg == "M" {
		addr := cpu.getPair(destReg)
		cpu.memory.write(addr, cpu.getReg(sourceReg))
		return 1
	} else if sourceReg == "M" {
		addr := cpu.getPair(sourceReg)
		memVal := cpu.memory.read(addr)
		cpu.updateReg(destReg, memVal)
		return 1
	}

	cpu.updateReg(destReg, cpu.getReg(sourceReg))
	return 1
}

// halt
func (cpu *cpu) hlt() uint8 {
	os.Exit(0)
	return 0
}

func (cpu *cpu) cmp() uint8 {
	var operand uint8

	if cpu.currentOp.Instruction == "CPI" {
		operand = cpu.memory.read(cpu.pc + 1)
	} else {
		operand = cpu.getReg(cpu.currentOp.Register)
	}

	res, carry := overflowingSub(cpu.regs.a, operand, 0)
	cpu.updateFlags(res, carry)

	if cpu.currentOp.Instruction == "CPI" {
		return 2
	}

	return 1
}

func (cpu *cpu) xra() uint8 {
	var res uint8
	var operand uint8

	if cpu.currentOp.Instruction == "XRI" {
		operand = cpu.memory.read(cpu.pc + 1)
	} else {
		operand = cpu.getReg(cpu.currentOp.Register)
	}

	res = cpu.regs.a ^ operand
	cpu.regs.a = res

	cpu.updateFlags(res, 0)

	if cpu.currentOp.Instruction == "XRI" {
		return 2
	}

	return 1
}

func (cpu *cpu) ana() uint8 {
	var operand uint8
	var res uint8
	prevAccum := cpu.regs.a

	if cpu.currentOp.Instruction == "ANI" {
		operand = cpu.memory.read(cpu.pc + 1)
	} else {
		operand = cpu.getReg(cpu.currentOp.Register)
	}

	res = prevAccum & operand
	cpu.regs.a = res
	cpu.updateFlags(res, 0)

	if cpu.currentOp.Instruction == "ANI" {
		return 2
	}

	return 1
}

func (cpu *cpu) add() uint8 {
	var res uint8
	var operand uint8
	var carry uint8
	prevAccum := cpu.regs.a

	if cpu.currentOp.Instruction == "ADI" || cpu.currentOp.Instruction == "ACI" {
		operand = cpu.memory.read(cpu.pc + 1)
	} else {
		operand = cpu.getReg(cpu.currentOp.Register)
	}

	switch cpu.currentOp.Instruction {
	case "ADD":
		res, carry = overflowingAdd(prevAccum, operand, 0)
	case "ADC":
		res, carry = overflowingAdd(prevAccum, operand, cpu.flags.cy)
	case "ADI":
		res, carry = overflowingAdd(prevAccum, operand, 0)
	case "ACI":
		res, carry = overflowingAdd(prevAccum, operand, cpu.flags.cy)
	}
	cpu.regs.a = res
	cpu.updateFlags(res, carry)

	if cpu.currentOp.Instruction == "ADI" || cpu.currentOp.Instruction == "ACI" {
		return 2
	}

	return 1
}

func (cpu *cpu) sub() uint8 {
	var res uint8
	var carry uint8
	var operand uint8
	prevAccum := cpu.regs.a

	if cpu.currentOp.Instruction == "SUI" || cpu.currentOp.Instruction == "SBI" {
		operand = cpu.memory.read(cpu.pc + 1)
	} else {
		operand = cpu.getReg(cpu.currentOp.Register)
	}

	switch cpu.currentOp.Instruction {
	case "SUB":
		res, carry = overflowingSub(prevAccum, operand, 0)
	case "SBB":
		res, carry = overflowingSub(prevAccum, operand, cpu.flags.cy)
	case "SUI":
		res, carry = overflowingSub(prevAccum, operand, 0)
	case "SBI":
		res, carry = overflowingSub(prevAccum, operand, cpu.flags.cy)
	}

	cpu.regs.a = res
	cpu.updateFlags(res, carry)

	if cpu.currentOp.Instruction == "SUI" || cpu.currentOp.Instruction == "SBI" {
		return 2
	}

	return 1
}

func (cpu *cpu) ora() uint8 {
	var res uint8
	var operand uint8

	if cpu.currentOp.Instruction == "ORI" {
		operand = cpu.memory.read(cpu.pc + 1)
	} else {
		operand = cpu.getReg(cpu.currentOp.Register)
	}

	res = cpu.regs.a | operand

	cpu.regs.a = res
	cpu.updateFlags(res, 0)

	if cpu.currentOp.Instruction == "ORI" {
		return 2
	}

	return 1
}

func (cpu *cpu) call() uint8 {
	if make16bit(cpu.memory.read(cpu.pc+2), cpu.memory.read(cpu.pc+1)) == decoder.BDOS {
		if cpu.regs.c == 0x9 {
			addr := cpu.getPair("D")
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
		if cpu.currentOp.Condition == "" || cpu.checkConditionFlag() {
			lsb := uint8(cpu.memory[cpu.pc+1])
			msb := uint8(cpu.memory[cpu.pc+2])
			addr := make16bit(msb, lsb)

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

func (cpu *cpu) jmp() uint8 {
	if cpu.currentOp.Condition == "" || cpu.checkConditionFlag() {
		cpu.pc = make16bit(cpu.memory.read(cpu.pc+2), cpu.memory.read(cpu.pc+1))

		return 0
	}

	return 3
}

func (cpu *cpu) ret() uint8 {
	if cpu.currentOp.Condition == "" || cpu.checkConditionFlag() {
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

func (cpu *cpu) push() uint8 {
	sp := cpu.sp
	switch cpu.currentOp.Register {
	case "B":
		cpu.memory.write(sp-2, cpu.regs.c)
		cpu.memory.write(sp-1, cpu.regs.b)
	case "D":
		cpu.memory.write(sp-2, cpu.regs.e)
		cpu.memory.write(sp-1, cpu.regs.d)
	case "H":
		cpu.memory.write(sp-2, cpu.regs.l)
		cpu.memory.write(sp-1, cpu.regs.h)
	case "PSW":
		cpu.memory.write(sp-2, cpu.flags.s|cpu.flags.z<<1|cpu.flags.ac<<2|cpu.flags.p<<3|cpu.flags.cy<<4)
		cpu.memory.write(sp-1, cpu.regs.a)
	}
	cpu.sp += 2
	return 1
}

func (cpu *cpu) pop() uint8 {
	sp := cpu.sp
	switch cpu.currentOp.Register {
	case "B":
		cpu.regs.c = cpu.memory.read(sp)
		cpu.regs.b = cpu.memory.read(sp + 1)
	case "D":
		cpu.regs.e = cpu.memory.read(sp)
		cpu.regs.d = cpu.memory.read(sp + 1)
	case "H":
		cpu.regs.l = cpu.memory.read(sp)
		cpu.regs.h = cpu.memory.read(sp + 1)
	case "PSW":
		psw := cpu.memory.read(sp)
		cpu.flags.cy = psw & 0x1
		cpu.flags.p = (psw >> 2) & 0x1
		cpu.flags.ac = (psw >> 4) & 0x1
		cpu.flags.z = (psw >> 6) & 0x1
		cpu.flags.s = (psw >> 7) & 0x1
		cpu.regs.a = cpu.memory.read(sp + 1)
	}
	cpu.sp += 2
	return 1
}

func (cpu *cpu) sphl() uint8 {
	cpu.sp = cpu.getPair("HL")
	return 1
}

func (cpu *cpu) xthl() uint8 {
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

func (cpu *cpu) xchg() uint8 {
	h, d, l, e := cpu.regs.h, cpu.regs.d, cpu.regs.l, cpu.regs.e
	cpu.regs.h = d
	cpu.regs.d = h
	cpu.regs.l = e
	cpu.regs.e = l
	return 1
}

func (cpu *cpu) rst() uint8 {
	resetAddr := cpu.currentOp.Code & 0b00111000
	cpu.memory.write(cpu.sp-1, uint8(cpu.pc&0x0f))
	cpu.memory.write(cpu.sp-2, uint8((cpu.pc>>8)&0x0f))
	cpu.sp -= 2
	cpu.pc = uint16(resetAddr - 1)

	return 0
}

func (cpu *cpu) in() uint8 {
	return 2
}

func (cpu *cpu) out() uint8 {
	return 2
}

func (cpu *cpu) ei() uint8 {
	cpu.interruptEnabled = true
	return 1
}

func (cpu *cpu) di() uint8 {
	cpu.interruptEnabled = false
	return 1
}

func (cpu *cpu) pchl() uint8 {
	cpu.pc = make16bit(cpu.regs.h, cpu.regs.l)
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

func (cpu *cpu) updateFlags(val uint8, carry ...uint8) {
	if val == 0 {
		cpu.flags.z = 1
	} else {
		cpu.flags.z = 0
	}

	if val >= 128 {
		cpu.flags.s = 1
	} else {
		cpu.flags.s = 0
	}

	if carry != nil {
		cpu.flags.cy = carry[0]
	}

	cpu.flags.p = parity(val)
}

func overflowingSub(x, y, cy uint8) (uint8, uint8) {
	var carry uint8
	if int16(x)-int16(y)-int16(cy) < 0 {
		carry = 1
	}
	return x - y - cy, carry
}

func overflowingAdd(x, y, cy uint8) (uint8, uint8) {
	var carry uint8
	if uint16(x)+uint16(y)+uint16(cy) >= 255 {
		carry = 1
	}
	return x + y + cy, carry
}
