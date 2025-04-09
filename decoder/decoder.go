package decoder

import (
	misc "cpu-emulator/utils"
	"fmt"
)

type Opcode struct {
	Code        byte
	Name        string
	Instruction uint8
	Condition   uint8
	LowNibble   uint8
	HighNibble  uint8
	Cycles      uint8
}

// conditions
const (
	NotZero uint8 = iota + 1
	Zero
	NoCarry
	Carry
	ParityOdd
	ParityEven
	Minus
	Positive
)

const (
	NOP uint8 = iota + 1
	RLC
	RRC
	RAL
	RAR
	RIM
	SHLD
	DAA
	LHLD
	CMA
	SIM
	STA
	STC
	LDA
	CMC
	RNZ
	JNZ
	JMP
	CNZ
	RZ
	RET
	JZ
	CZ
	CALL
	RNC
	JNC
	OUT
	CNC
	RC
	JC
	IN
	CC
	RPO
	JPO
	XTHL
	CPO
	RPE
	PCHL
	JPE
	XCHG
	CPE
	RP
	JP
	DI
	CP
	RM
	SPHL
	JM
	EI
	CM
	LXI
	STAX
	INX
	INR
	DCR
	MVI
	DAD
	LDAX
	DCX
	MOV
	ADD
	ADC
	SUB
	SBB
	ANA
	XRA
	ORA
	CMP
	POP
	PUSH
	RST
	ADI
	ACI
	SUI
	SBI
	ANI
	XRI
	ORI
	CPI
	HLT
)

const BDOS = 0x05

func GetInstruction(memory []byte, pc uint16) *Opcode {
	var instruction uint8
	var name string
	var cycles uint8

	code := memory[pc]
	lowNibble := code & 0x0f
	highNibble := (code & 0xf0) >> 4

	switch code {
	case 0x00:
		instruction = NOP
		name = "NOP"
		cycles = 1
	case 0x07:
		instruction = RLC
		name = "RLC"
		cycles = 1
	case 0x0f:
		instruction = RRC
		name = "RRC"
		cycles = 1
	case 0x17:
		instruction = RAL
		name = "RAL"
		cycles = 1
	case 0x1f:
		instruction = RAR
		name = "RAR"
		cycles = 1
	case 0x20:
		instruction = RIM
		name = "RIM"
		cycles = 1
	case 0x22:
		instruction = SHLD
		name = fmt.Sprintf("SHLD 0x%x", misc.Make16bit(memory[pc+2], memory[pc+1]))
		cycles = 5
	case 0x27:
		instruction = DAA
		name = "DAA"
		cycles = 1
	case 0x2a:
		instruction = LHLD
		name = fmt.Sprintf("LHLD 0x%x", misc.Make16bit(memory[pc+2], memory[pc+1]))
		cycles = 5
	case 0x2f:
		instruction = CMA
		name = "CMA"
		cycles = 1
	case 0x30:
		instruction = SIM
		name = "SIM"
		cycles = 0
	case 0x32:
		instruction = STA
		name = fmt.Sprintf("STA 0x%x", misc.Make16bit(memory[pc+2], memory[pc+1]))
		cycles = 4
	case 0x37:
		instruction = STC
		name = "STC"
		cycles = 1
	case 0x3a:
		instruction = LDA
		name = fmt.Sprintf("LDA 0x%x", misc.Make16bit(memory[pc+2], memory[pc+1]))
		cycles = 4
	case 0x3f:
		instruction = CMC
		name = "CMC"
		cycles = 5
	case 0xc0:
		instruction = RNZ
		name = "RNZ"
		cycles = 3
	case 0xc2:
		instruction = JNZ
		name = "JNZ"
		cycles = 3
	case 0xc3:
		instruction = JMP
		name = "JMP"
		cycles = 3
	case 0xc4:
		instruction = CNZ
		name = "CNZ"
		cycles = 3
	case 0xc8:
		instruction = RZ
		name = "RZ"
		cycles = 5
	case 0xc9:
		instruction = RET
		name = "RET"
		cycles = 3
	case 0xca:
		instruction = JZ
		name = "JZ"
		cycles = 3
	case 0xcc:
		instruction = CZ
		name = "CZ"
		cycles = 3
	case 0xcd:
		instruction = CALL
		name = fmt.Sprintf("CALL 0x%x", misc.Make16bit(memory[pc+2], memory[pc+1]))
		cycles = 5
	case 0xd0:
		instruction = RNC
		name = "RNC"
		cycles = 3
	case 0xd2:
		instruction = JNC
		name = "JNC"
		cycles = 3
	case 0xd3:
		instruction = OUT
		name = "OUT"
		cycles = 3
	case 0xd4:
		instruction = CNC
		name = "CNC"
		cycles = 3
	case 0xd8:
		instruction = RC
		name = "RC"
		cycles = 5
	case 0xda:
		instruction = JC
		name = "JC"
		cycles = 3
	case 0xdb:
		instruction = IN
		name = "IN"
		cycles = 3
	case 0xdc:
		instruction = CC
		name = "CC"
		cycles = 3
	case 0xe0:
		instruction = RPO
		name = "RPO"
		cycles = 5
	case 0xe2:
		instruction = JPO
		name = "JPO"
		cycles = 3
	case 0xe3:
		instruction = XTHL
		name = "XTHL"
		cycles = 5
	case 0xe4:
		instruction = CPO
		name = "CPO"

	case 0xe8:
		instruction = RPE
		name = "RPE"
		cycles = 5
	case 0xe9:
		instruction = PCHL
		name = "PCHL"
		cycles = 5
	case 0xea:
		instruction = JPE
		name = "JPE"
		cycles = 3
	case 0xeb:
		instruction = XCHG
		name = "XCHG"
		cycles = 1
	case 0xec:
		instruction = CPE
		name = "CPE"
		cycles = 3
	case 0xf0:
		instruction = RP
		name = "RP"
		cycles = 5
	case 0xf2:
		instruction = JP
		name = "JP"
		cycles = 3
	case 0xf3:
		instruction = DI
		name = "DI"
		cycles = 1
	case 0xf4:
		instruction = CP
		name = "CP"
		cycles = 3
	case 0xf8:
		instruction = RM
		name = "RM"
		cycles = 5
	case 0xf9:
		instruction = SPHL
		name = "SPHL"
		cycles = 5
	case 0xfa:
		instruction = JM
		name = "JM"
		cycles = 3
	case 0xfb:
		instruction = EI
		name = "EI"
		cycles = 1
	case 0xfc:
		instruction = CM
		name = "CM"
		cycles = 3
	default:

		if code >= 0x01 && code <= 0x31 && code&0xF == 0x1 {
			instruction = LXI
			name = fmt.Sprintf("LXI %s, 0x%x", misc.RegPairToString(highNibble), misc.Make16bit(memory[pc+2], memory[pc+1]))
			cycles = 3
		} else if code >= 0x02 && code <= 0x12 && code&0xF == 0x2 {
			instruction = STAX
			name = fmt.Sprintf("STAX %s, 0x%x", misc.RegPairToString(highNibble), misc.Make16bit(memory[pc+2], memory[pc+1]))
			cycles = 2
		} else if code >= 0x03 && code <= 0x33 && code&0xF == 0x3 {
			name = fmt.Sprintf("INX %s", misc.RegPairToString(highNibble))
			instruction = INX
			cycles = 2
		} else if code >= 0x04 && code <= 0x3c && (code&0xF == 0x4 || code&0xf == 0xc) {
			name = fmt.Sprintf("INR %s", misc.RegToString(code>>3))
			instruction = INR
			cycles = 1
		} else if code >= 0x05 && code <= 0x3d && (code&0xf == 0xd || code&0xf == 0x5) {
			name = "DCR"
			name = fmt.Sprintf("DCR %s", misc.RegToString(code>>3))
			instruction = DCR
			cycles = 1
		} else if code >= 0x06 && code <= 0x3e && (code&0xf == 0x6 || code&0xf == 0xe) {
			name = fmt.Sprintf("MVI %s 0x%x", misc.RegToString(code>>3), memory[pc+1])
			instruction = MVI
		} else if code >= 0x09 && code <= 0x39 && code&0xf == 0x9 {
			name = fmt.Sprintf("DAD %s", misc.RegPairToString(highNibble))
			instruction = DAD
			cycles = 3
		} else if code >= 0x0a && code <= 0x1a && code&0xf == 0xa {
			name = fmt.Sprintf("LDAX %s", misc.RegPairToString(highNibble))
			instruction = LDAX
			cycles = 2
		} else if code >= 0x0b && code <= 0x3b && code&0xf == 0xb {
			name = fmt.Sprintf("DCX %s", misc.RegPairToString(highNibble))
			instruction = DCX
			cycles = 1
		} else if code >= 0x40 && 0x7f >= code && code != 0x76 {

			name = fmt.Sprintf("MOV %s, %s", misc.RegToString(code>>3), misc.RegToString(code&0b111))
			instruction = MOV
			if ((code>>3)&0xf0) == 0b110 || uint8(code) == 0b110 {
				cycles = 2
			} else {
				cycles = 1
			}
			cycles = 3
		} else if code >= 0x80 && code <= 0x87 {

			name = fmt.Sprintf("ADD %s", misc.RegToString(lowNibble))
			instruction = ADD
			if uint8(code) == 0b110 {
				cycles = 2
			} else {
				cycles = 1
			}

		} else if code >= 0x88 && code <= 0x8f {

			name = "ADC"
			name = fmt.Sprintf("ADC %s", misc.RegToString(lowNibble))
			if uint8(code) == 0b110 {
				cycles = 2
			} else {
				cycles = 1
			}

			instruction = ADC
		} else if code >= 0x90 && code <= 0x97 {

			name = fmt.Sprintf("SUB %s", misc.RegToString(lowNibble))
			instruction = SUB
			if uint8(code) == 0b110 {
				cycles = 2
			} else {
				cycles = 1
			}

		} else if code >= 0x98 && code <= 0x9f {

			name = fmt.Sprintf("SBB %s", misc.RegToString(lowNibble))
			instruction = SBB
			if uint8(code) == 0b110 {
				cycles = 2
			} else {
				cycles = 1
			}

		} else if code >= 0xa0 && code <= 0xa7 {

			name = fmt.Sprintf("ANA %s", misc.RegToString(lowNibble))
			instruction = ANA
			cycles = 2
		} else if code >= 0xa8 && code <= 0xaf {
			name = fmt.Sprintf("XRA %s", misc.RegToString(lowNibble))
			instruction = XRA
			cycles = 2
		} else if code >= 0xb0 && code <= 0xb7 {
			name = fmt.Sprintf("ORA %s", misc.RegToString(lowNibble))
			instruction = ORA
			cycles = 2
		} else if code >= 0xb8 && code <= 0xbf {
			name = fmt.Sprintf("CMP %s", misc.RegToString(lowNibble))
			instruction = CMP
			cycles = 2
		} else if code >= 0xc1 && code <= 0xf1 && code&0xf == 0x1 {
			if code == 0xf1 {
				name = "POP PSW"
			} else {
				name = fmt.Sprintf("POP %s", misc.RegPairToString(highNibble))
			}
			cycles = 3
			instruction = POP
		} else if code >= 0xc5 && code <= 0xf5 && code&0xf == 0x5 {
			if code == 0xf5 {
				name = "PUSH PSW"
			} else {
				name = fmt.Sprintf("PUSH %s", misc.RegPairToString(highNibble))
			}
			cycles = 3
			instruction = PUSH
		} else if code >= 0xc6 && code <= 0xfe && (code&0xf == 0x6 || code&0xf == 0xe) {

			instructs := map[byte]uint8{
				0xc6: ADI, 0xce: ACI, 0xd6: SUI, 0xde: SBI, 0xe6: ANI,
				0xee: XRI, 0xf6: ORI, 0xfe: CPI,
			}
			names := map[byte]string{
				0xc6: "ADI", 0xce: "ACI", 0xd6: "SUI", 0xde: "SBI", 0xe6: "ANI",
				0xee: "XRI", 0xf6: "ORI", 0xfe: "CPI",
			}
			cycles = 2
			name = fmt.Sprintf("%s 0x%0x", names[code], memory[pc+1])
			instruction = instructs[code]
		} else if code >= 0xc7 && code <= 0xff && (code&0xf == 0x7 || code&0xf == 0xf) {
			name = "RST"
			instruction = RST
			cycles = 3
		}
	}
	opcode := &Opcode{
		Code:        code,
		Name:        name,
		Instruction: instruction,
		LowNibble:   lowNibble,
		HighNibble:  highNibble,
		Cycles:      cycles,
	}
	if (code != 0xd9 && code != 0xcb) && (code >= 0xc0 && code <= 0xff) && (string(opcode.Name[0]) == "J" || string(opcode.Name[0]) == "R" || string(opcode.Name[0]) == "C") {
		setConditionOpcode(opcode)
	}
	return opcode
}

func setConditionOpcode(opcode *Opcode) {
	condition := strConditionToByte(string(opcode.Name[1:]))
	if condition == NotZero || condition == Zero || condition == NoCarry || condition == Carry || condition == ParityOdd || condition == ParityEven || condition == Minus || condition == Positive {
		opcode.Condition = condition
		var name string
		switch string(opcode.Name[0]) {
		case "J":
			name = "JMP"
		case "R":
			name = "RET"
		case "C":
			name = "CALL"
		}
		opcode.Name = name + " " + opcode.Name[1:]
	}
}

func strConditionToByte(cond string) byte {
	switch cond {
	case "NZ":
		return NotZero
	case "Z":
		return Zero
	case "NC":
		return NoCarry
	case "C":
		return Carry
	case "PO":
		return ParityOdd
	case "PE":
		return ParityEven
	case "P":
		return Positive
	case "M":
		return Minus
	default:
		return 0
	}
}
