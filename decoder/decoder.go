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

func GetInstruction(memory *byte) *Opcode {
	var instruction uint8
	var name string
	lowNibble := *memory & 0x0f
	highNibble := (*memory & 0xf0) >> 4
	code := *memory

	switch code {
	case 0x00:
		instruction = NOP
		name = "NOP"
	case 0x07:
		instruction = RLC
		name = "RLC"
	case 0x0f:
		instruction = RRC
		name = "RRC"
	case 0x17:
		instruction = RAL
		name = "RAL"
	case 0x1f:
		instruction = RAR
		name = "RAR"
	case 0x20:
		instruction = RIM
		name = "RIM"
	case 0x22:
		instruction = SHLD
		name = fmt.Sprintf("SHLD 0x%x", misc.Make16bit(*memory+2, *memory+1))
	case 0x27:
		instruction = DAA
		name = "DAA"
	case 0x2a:
		instruction = LHLD
		name = fmt.Sprintf("LHLD 0x%x", misc.Make16bit(*memory+2, *memory+1))
	case 0x2f:
		instruction = CMA
		name = "CMA"
	case 0x30:
		instruction = SIM
		name = "SIM"
	case 0x32:
		instruction = STA
		name = fmt.Sprintf("STA 0x%x", misc.Make16bit(*memory+2, *memory+1))
	case 0x37:
		instruction = STC
		name = "STC"
	case 0x3a:
		instruction = LDA
		name = fmt.Sprintf("LDA 0x%x", misc.Make16bit(*memory+2, *memory+1))
	case 0x3f:
		instruction = CMC
		name = "CMC"
	case 0xc0:
		instruction = RNZ
		name = "RNZ"
	case 0xc2:
		instruction = JNZ
		name = "JNZ"
	case 0xc3:
		instruction = JMP
		name = "JMP"
	case 0xc4:
		instruction = CNZ
		name = "CNZ"
	case 0xc8:
		instruction = RZ
		name = "RZ"
	case 0xc9:
		instruction = RET
		name = "RET"
	case 0xca:
		instruction = JZ
		name = "JZ"
	case 0xcc:
		instruction = CZ
		name = "CZ"
	case 0xcd:
		instruction = CALL
		name = fmt.Sprintf("CALL 0x%x", misc.Make16bit(*memory+2, *memory+1))
	case 0xd0:
		instruction = RNC
		name = "RNC"
	case 0xd2:
		instruction = JNC
		name = "JNC"
	case 0xd3:
		instruction = OUT
		name = "OUT"
	case 0xd4:
		instruction = CNC
		name = "CNC"
	case 0xd8:
		instruction = RC
		name = "RC"
	case 0xda:
		instruction = JC
		name = "JC"
	case 0xdb:
		instruction = IN
		name = "IN"
	case 0xdc:
		instruction = CC
		name = "CC"
	case 0xe0:
		instruction = RPO
		name = "RPO"
	case 0xe2:
		instruction = JPO
		name = "JPO"
	case 0xe3:
		instruction = XTHL
		name = "XTHL"
	case 0xe4:
		instruction = CPO
		name = "CPO"
	case 0xe8:
		instruction = RPE
		name = "RPE"
	case 0xe9:
		instruction = PCHL
		name = "PCHL"
	case 0xea:
		instruction = JPE
		name = "JPE"
	case 0xeb:
		instruction = XCHG
		name = "XCHG"
	case 0xec:
		instruction = CPE
		name = "CPE"
	case 0xf0:
		instruction = RP
		name = "RP"
	case 0xf2:
		instruction = JP
		name = "JP"
	case 0xf3:
		instruction = DI
		name = "DI"
	case 0xf4:
		instruction = CP
		name = "CP"
	case 0xf8:
		instruction = RM
		name = "RM"
	case 0xf9:
		instruction = SPHL
		name = "SPHL"
	case 0xfa:
		instruction = JM
		name = "JM"
	case 0xfb:
		instruction = EI
		name = "EI"
	case 0xfc:
		instruction = CM
		name = "CM"
	default:

		if code >= 0x01 && code <= 0x31 && code&0xF == 0x1 {
			instruction = LXI
			name = fmt.Sprintf("LXI %s, 0x%x", misc.RegPairToString(highNibble), misc.Make16bit(*memory+2, *memory+1))

		} else if code >= 0x02 && code <= 0x12 && code&0xF == 0x2 {
			instruction = STAX
			name = fmt.Sprintf("STAX %s, 0x%x", misc.RegPairToString(highNibble), misc.Make16bit(*memory+2, *memory+1))
		} else if code >= 0x03 && code <= 0x33 && code&0xF == 0x3 {
			name = fmt.Sprintf("INX %s", misc.RegPairToString(highNibble))
			instruction = INX
		} else if code >= 0x04 && code <= 0x3c && (code&0xF == 0x4 || code&0xf == 0xc) {
			name = fmt.Sprintf("INR %s", misc.RegToString(code>>3))
			instruction = INR
		} else if code >= 0x05 && code <= 0x3d && (code&0xf == 0xd || code&0xf == 0x5) {
			name = "DCR"
			name = fmt.Sprintf("DCR %s", misc.RegToString(code>>3))
			instruction = DCR
		} else if code >= 0x06 && code <= 0x3e && (code&0xf == 0x6 || code&0xf == 0xe) {
			name = fmt.Sprintf("MVI %s 0x%x", misc.RegToString(code>>3), *memory+1)
			instruction = MVI
		} else if code >= 0x09 && code <= 0x39 && code&0xf == 0x9 {
			name = fmt.Sprintf("DAD %s", misc.RegPairToString(highNibble))
			instruction = DAD
		} else if code >= 0x0a && code <= 0x1a && code&0xf == 0xa {
			name = fmt.Sprintf("LDAX %s", misc.RegPairToString(highNibble))
			instruction = LDAX
		} else if code >= 0x0b && code <= 0x3b && code&0xf == 0xb {
			name = fmt.Sprintf("DCX %s", misc.RegPairToString(highNibble))
			instruction = DCX
		} else if code >= 0x40 && 0x7f >= code && code != 0x76 {
			name = fmt.Sprintf("MOV %s, %s", misc.RegToString(code>>3), misc.RegToString(code&0b111))
			instruction = MOV
		} else if code >= 0x80 && code <= 0x87 {
			name = fmt.Sprintf("ADD %s", misc.RegToString(lowNibble))
			instruction = ADD
		} else if code >= 0x88 && code <= 0x8f {
			name = "ADC"
			name = fmt.Sprintf("ADC %s", misc.RegToString(lowNibble))
			instruction = ADC
		} else if code >= 0x90 && code <= 0x97 {
			name = fmt.Sprintf("SUB %s", misc.RegToString(lowNibble))
			instruction = SUB
		} else if code >= 0x98 && code <= 0x9f {
			name = fmt.Sprintf("SBB %s", misc.RegToString(lowNibble))
			instruction = SBB
		} else if code >= 0xa0 && code <= 0xa7 {
			name = fmt.Sprintf("ANA %s", misc.RegToString(lowNibble))
			instruction = ANA
		} else if code >= 0xa8 && code <= 0xaf {
			name = fmt.Sprintf("XRA %s", misc.RegToString(lowNibble))
			instruction = XRA
		} else if code >= 0xb0 && code <= 0xb7 {
			name = fmt.Sprintf("ORA %s", misc.RegToString(lowNibble))
			instruction = ORA
		} else if code >= 0xb8 && code <= 0xbf {
			name = fmt.Sprintf("CMP %s", misc.RegToString(lowNibble))
			instruction = CMP
		} else if code >= 0xc1 && code <= 0xf1 && code&0xf == 0x1 {
			if code == 0xf1 {
				name = "POP PSW"
			} else {
				name = fmt.Sprintf("POP %s", misc.RegPairToString(highNibble))
			}
			instruction = POP
		} else if code >= 0xc5 && code <= 0xf5 && code&0xf == 0x5 {
			if code == 0xf5 {
				name = "PUSH PSW"
			} else {
				name = fmt.Sprintf("PUSH %s", misc.RegPairToString(highNibble))
			}
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

			name = fmt.Sprintf("%s 0x%0x", names[code], *memory+1)
			instruction = instructs[code]
		} else if code >= 0xc7 && code <= 0xff && (code&0xf == 0x7 || code&0xf == 0xf) {
			name = "RST"
			instruction = RST
		}
	}
	opcode := &Opcode{
		Code:        code,
		Name:        name,
		Instruction: instruction,
		LowNibble:   lowNibble,
		HighNibble:  highNibble,
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
		return Minus
	case "M":
		return Positive
	default:
		return 0
	}
}
