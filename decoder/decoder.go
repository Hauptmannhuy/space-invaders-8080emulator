package decoder

var lxi = map[byte]string{
	0x01: "B",
	0x11: "D",
	0x21: "H",
	0x31: "SP",
}
var stax = map[byte]string{
	0x02: "B",
	0x12: "D"}
var inx = map[byte]string{
	0x03: "B",
	0x13: "D",
	0x23: "H",
	0x33: "SP",
}
var inr = map[byte]string{
	0x04: "B",
	0x0c: "C",
	0x14: "D",
	0x1c: "E",
	0x24: "H",
	0x2c: "L",
	0x34: "M",
	0x3c: "A",
}
var dcr = map[byte]string{
	0x05: "B", 0x0d: "C", 0x15: "D", 0x1d: "E",
	0x25: "H", 0x2d: "L", 0x35: "M", 0x3d: "A",
}
var mvi = map[byte]string{
	0x06: "B", 0x0e: "C", 0x16: "D", 0x1e: "E",
	0x26: "H", 0x2e: "L", 0x36: "M", 0x3e: "A",
}
var dad = map[byte]string{
	0x9:  "B",
	0x19: "D",
	0x29: "H",
	0x39: "SP",
}
var dcx = map[byte]string{
	0x0b: "B",
	0x1b: "D",
	0x2b: "H",
	0x3b: "SP",
}
var mov = map[byte]string{
	0x40: "B,B", 0x41: "B,C", 0x42: "B,D", 0x43: "B,E",
	0x44: "B,H", 0x45: "B,L", 0x46: "B,M", 0x47: "B,A",

	0x48: "C,B", 0x49: "C,C", 0x4a: "C,D", 0x4b: "C,E",
	0x4c: "C,H", 0x4d: "C,L", 0x4e: "C,M", 0x4f: "C,A",

	0x50: "D,B", 0x51: "D,C", 0x52: "D,D", 0x53: "D,E",
	0x54: "D,H", 0x55: "D,L", 0x56: "D,M", 0x57: "D,A",

	0x58: "E,B", 0x59: "E,C", 0x5a: "E,D", 0x5b: "E,E",
	0x5c: "E,H", 0x5d: "E,L", 0x5e: "E,M", 0x5f: "E,A",

	0x60: "H,B", 0x61: "H,C", 0x62: "H,D", 0x63: "H,E",
	0x64: "H,H", 0x65: "H,L", 0x66: "H,M", 0x67: "H,A",

	0x68: "L,B", 0x69: "L,C", 0x6a: "L,D", 0x6b: "L,E",
	0x6c: "L,H", 0x6d: "L,L", 0x6e: "L,E", 0x6f: "L,A",

	0x70: "M,B", 0x71: "M,C", 0x72: "M,D", 0x73: "M,E",
	0x74: "M,H", 0x75: "M,L", 0x77: "M,A",

	0x78: "A,B", 0x79: "A,C", 0x7a: "A,D", 0x7b: "A,E",
	0x7c: "A,H", 0x7d: "A,L", 0x7e: "A,M", 0x7f: "A,A",
}

var add = map[byte]string{
	0x80: "B", 0x81: "C", 0x82: "D", 0x83: "E",
	0x84: "H", 0x85: "L", 0x86: "M", 0x87: "A",
}

var adc = map[byte]string{
	0x88: "B", 0x89: "C", 0x8a: "D", 0x8b: "E",
	0x8c: "H", 0x8d: "L", 0x8e: "M", 0x8f: "A",
}

var sub = map[byte]string{
	0x90: "B", 0x91: "C", 0x92: "D", 0x93: "E",
	0x94: "H", 0x95: "L", 0x96: "M", 0x97: "A",
}

var sbb = map[byte]string{
	0x98: "B", 0x99: "C", 0x9a: "D", 0x9b: "E",
	0x9c: "H", 0x9d: "L", 0x9e: "M", 0x9f: "A",
}

var ana = map[byte]string{
	0xa0: "B", 0xa1: "C", 0xa2: "D", 0xa3: "E",
	0xa4: "H", 0xa5: "L", 0xa6: "M", 0xa7: "A",
}

var xra = map[byte]string{
	0xa8: "B", 0xa9: "C", 0xaa: "D", 0xab: "E",
	0xac: "H", 0xad: "L", 0xae: "M", 0xaf: "A",
}

var ora = map[byte]string{
	0xb0: "B", 0xb1: "C", 0xb2: "D", 0xb3: "E",
	0xb4: "H", 0xb5: "L", 0xb6: "M", 0xb7: "A",
}

var cmp = map[byte]string{
	0xb8: "B", 0xb9: "C", 0xba: "D", 0xbb: "E",
	0xbc: "H", 0xbd: "L", 0xbe: "M", 0xbf: "A",
}

var pop = map[byte]string{
	0xc1: "B", 0xd1: "D", 0xe1: "H", 0xf1: "PSW",
}

var push = map[byte]string{
	0xc5: "B", 0xd5: "D", 0xe5: "H", 0xf5: "PSW",
}

var rst = map[byte]string{
	0xc7: "1", 0xcf: "1", 0xd7: "2", 0xdf: "3",
	0xe7: "4", 0xef: "5", 0xf7: "6", 0xff: "7",
}

func GetInstruction(code byte) string {
	var instruction string
	switch code {
	case 0x00:
		instruction = "NOP"
	case 0x07:
		instruction = "RLC"
	case 0x0f:
		instruction = "RRC"
	case 0x17:
		instruction = "RAL"
	case 0x1f:
		instruction = "RAR"
	case 0x20:
		instruction = "RIM"
	case 0x22:
		instruction = "SHLD"
	case 0x27:
		instruction = "DAA"
	case 0x2a:
		instruction = "LHD"
	case 0x2f:
		instruction = "CMA"
	case 0x30:
		instruction = "SIM"
	case 0x32:
		instruction = "STA"
	case 0x37:
		instruction = "STC"
	case 0x3a:
		instruction = "LDA"
	case 0x3f:
		instruction = "CMC"
	case 0xc0:
		instruction = "RNZ"
	case 0xc2:
		instruction = "JNZ"
	case 0xc3:
		instruction = "JMP"
	case 0xc4:
		instruction = "CNZ"
	case 0xc8:
		instruction = "RZ"
	case 0xc9:
		instruction = "RET"
	case 0xca:
		instruction = "JZ"
	case 0xcc:
		instruction = "CZ"
	case 0xcd:
		instruction = "CALL"
	case 0xd0:
		instruction = "RNC"
	case 0xd2:
		instruction = "JNC"
	case 0xd3:
		instruction = "OUT"
	case 0xd4:
		instruction = "CNC"
	case 0xd8:
		instruction = "RC"
	case 0xda:
		instruction = "JC"
	case 0xdb:
		instruction = "IN"
	case 0xdc:
		instruction = "CC"
	case 0x0e:
		instruction = "RPO"
	case 0xe2:
		instruction = "JPO"
	case 0xe3:
		instruction = "XHTL"
	case 0xe4:
		instruction = "CPO"
	case 0xe8:
		instruction = "RPE"
	case 0xe9:
		instruction = "PCHL"
	case 0xea:
		instruction = "JPE"
	case 0xeb:
		instruction = "XCHG"
	case 0xec:
		instruction = "CPE"
	case 0xf0:
		instruction = "RP"
	case 0xf2:
		instruction = "JP"
	case 0xf3:
		instruction = "DI"
	case 0xf4:
		instruction = "CP"
	case 0xf8:
		instruction = "RM"
	case 0xf9:
		instruction = "SPHL"
	case 0xfa:
		instruction = "JM"
	case 0xfb:
		instruction = "EI"
	case 0xfc:
		instruction = "CM"
	default:
		if code >= 0x01 && code <= 0x31 && code&0xF == 0x1 {
			instruction = "LXI"
		} else if code >= 0x02 && code <= 0x12 && code&0xF == 0x2 {
			instruction = "STAX"
		} else if code >= 0x03 && code <= 0x33 && code&0xF == 0x3 {
			instruction = "INX"
		} else if code >= 0x04 && code <= 0x3c && (code&0xF == 0x4 || code&0xf == 0xc) {
			instruction = "INR"
		} else if code >= 0x05 && code <= 0x3d && (code&0xf == 0xd || code&0xf == 0x5) {
			instruction = "DCR"
		} else if code >= 0x06 && code <= 0x3e && (code&0xf == 0x6 || code&0xf == 0xe) {
			instruction = "MVI"
		} else if code >= 0x09 && code <= 0x39 && code&0xf == 0x9 {
			instruction = "DAD"
		} else if code >= 0x0a && 0x1a <= code && code&0xf == 0xa {
			instruction = "LDAX"
		} else if code >= 0x0b && code <= 0x3b && code&0xf == 0xb {
			instruction = "DCX"
		} else if code >= 0x40 && 0x7f >= code && code != 0x76 {
			instruction = "MOV"
		} else if code >= 0x80 && code <= 0x87 {
			instruction = "ADD"
		} else if code >= 0x88 && code <= 0x8f {
			instruction = "ADC"
		} else if code >= 0x90 && code <= 0x97 {
			instruction = "SUB"
		} else if code >= 0x98 && code <= 0x9f {
			instruction = "SBB"
		} else if code >= 0xa0 && code <= 0xa7 {
			instruction = "ANA"
		} else if code >= 0xa8 && code <= 0xaf {
			instruction = "XRA"
		} else if code >= 0xb0 && code <= 0xb7 {
			instruction = "ORA"
		} else if code >= 0xb8 && code <= 0xbf {
			instruction = "CMP"
		} else if code >= 0xc1 && code <= 0xf1 && code&0xf == 0x1 {
			instruction = "POP"
		} else if code >= 0xc5 && code <= 0xf5 && code&0xf == 0x5 {
			instruction = "PUSH"
		} else if code >= 0xc6 && code <= 0xfe && (code&0xf == 0x6 || code&0xf == 0xe) {
			instructs := map[byte]string{
				0xc6: "ADI", 0xce: "ACI", 0xd6: "SUI", 0xde: "SBI", 0xe6: "ANI",
				0xee: "XRI", 0xf6: "ORI", 0xfe: "CPI",
			}
			instruct := instructs[code]
			instruction = instruct

		} else if code >= 0xc7 && code <= 0xff && (code&0xf == 0x7 || code&0xf == 0xf) {
			instruction = "RST"
		}
	}
	return instruction
}

func GetDestination(instruction string, code byte) string {

	regs := map[string]map[byte]string{
		"LXI":  lxi,
		"STAX": stax,
		"INX":  inx,
		"INR":  inr,
		"DCR":  dcr,
		"MVI":  mvi,
		"DAD":  dad,
		"DCX":  dcx,
		"MOV":  mov,
		"ADD":  add,
		"ADC":  adc,
		"SUB":  sub,
		"SBB":  sbb,
		"ANA":  ana,
		"XRA":  xra,
		"ORA":  ora,
		"CMP":  cmp,
		"POP":  pop,
		"PUSH": push,
		"RST":  rst,
	}
	return regs[instruction][code]
}
