package machine

import (
	"cpu-emulator/decoder"
	"log"
	"math/bits"
)

// flags Z (zero), S (sign), P (parity), CY (carry), CA (auxillary  carry)
type flags struct {
	s, z, ac, p, cy uint8
}

func parity(val uint8) uint8 {
	count := bits.OnesCount8(val)
	if count%2 == 0 {
		return 1
	} else {
		return 0
	}
}

func (cpu *Cpu) updateFlags(val uint8, carry ...uint8) {
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

func (cpu *Cpu) setAux(expression bool) {
	if expression {
		cpu.flags.ac = 1
	} else {
		cpu.flags.ac = 0
	}
}

func (cpu *Cpu) checkConditionFlag() bool {

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
