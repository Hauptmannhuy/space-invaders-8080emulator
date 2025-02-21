package main

import (
	"cpu-emulator/decoder"
	"fmt"
	"os"
)

func loadHex() []byte {
	buff, _ := os.ReadFile("space-invaders.rom")
	return buff
}

func disassebmle(buffer []byte, pc int) int {
	code := buffer[pc]
	opcodes := 1
	fmt.Printf("%04x ", pc)
	instruction := decoder.GetInstruction(code)
	register := decoder.GetRegister(instruction, code)
	switch instruction {
	case "NOP":
		fmt.Printf("NOP")
	case "RLC":
		fmt.Printf("RLC")
	case "RRC":
		fmt.Printf("RRC")
	case "RAL":
		fmt.Printf("RAL")
	case "RAR":
		fmt.Printf("RAR")
	case "RIM":
		fmt.Printf("RIM")
	case "SHLD":
		fmt.Printf("SHLD 0x%02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "DAA":
		fmt.Printf("DAA 0x%02x", code)
	case "LHD":
		fmt.Printf("LHD 0x%02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "CMA":
		fmt.Printf("CMA")
	case "SIM":
		fmt.Printf("SIM")
	case "STA":
		fmt.Printf("STA %02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "STC":
		fmt.Printf("STC")
	case "LDA":
		fmt.Printf("LDA %02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "CMC":
		fmt.Printf("CMC")
	case "RNZ":
		fmt.Printf("RNZ")
	case "JNZ":
		fmt.Printf("JNZ %02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "JMP":
		fmt.Printf("JMP %02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "CNZ":
		fmt.Printf("CNZ %02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "RZ":
		fmt.Printf("RZ")
	case "RET":
		fmt.Printf("RET")
	case "JZ":
		fmt.Printf("JZ 0x%02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "CZ":
		fmt.Printf("CZ 0x%02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "CALL":
		fmt.Printf("CALL 0x%02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "RNC":
		fmt.Printf("RNC")
	case "JNC":
		fmt.Printf("JNC 0x%02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "OUT":
		fmt.Printf("OUT 0x%02x", buffer[pc+1])
	case "CNC":
		fmt.Printf("CNC 0x%02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "RC":
		fmt.Printf("RC")
	case "JC":
		fmt.Printf("JC 0x%02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "IN":
		fmt.Printf("IN, 0x%02x", buffer[pc+1])
		opcodes = 2
	case "CC":
		fmt.Printf("CC 0x%02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "RPO":
		fmt.Printf("RPO")
	case "JPO":
		fmt.Printf("JPO 0x%02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "XHTL":
		fmt.Printf("XHTL")
	case "CPO":
		fmt.Printf("CPO 0x%02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "RPE":
		fmt.Printf("RPE")
	case "PCHL":
		fmt.Printf("PCHL")
	case "JPE":
		fmt.Printf("JPE 0x%02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "XCHG":
		fmt.Printf("XCHG")
	case "CPE":
		fmt.Printf("CPE 0x%02x%02x", buffer[pc+2], buffer[pc+1])
	case "RP":
		fmt.Printf("RP")
	case "JP":
		fmt.Printf("JP 0x%02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "DI":
		fmt.Printf("DI")
	case "CP":
		fmt.Printf("CP 0x%02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "RM":
		fmt.Printf("RM")
	case "SPHL":
		fmt.Printf("SPHL")
	case "JM":
		fmt.Printf("JM 0x%02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "EI":
		fmt.Printf("EI")
	case "CM":
		fmt.Printf("CM 0x%02x%02x", buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "LXI":
		fmt.Printf("LXI %s, 0x%02x%02x", register, buffer[pc+2], buffer[pc+1])
		opcodes = 3
	case "STAX":
		fmt.Printf("STAX %s", register)
	case "INX":
		fmt.Printf("INX %s", register)
	case "INR":
		fmt.Printf("INR %s", register)
	case "DCR":
		fmt.Printf("DCR %s", register)
	case "MVI":
		fmt.Printf("MVI %s, 0x%02x", register, buffer[pc+1])
		opcodes = 2
	case "DAD":
		fmt.Printf("DAD %s 0x%02x", register, code)
	case "LDAX":
		fmt.Printf("LDAX %s 0x%02x", register, code)
	case "DCX":
		fmt.Printf("DCX %s 0x%02x", register, code)
	case "MOV":
		fmt.Printf("MOV %s", register)
	case "ADD":
		fmt.Printf("ADD %s", register)
	case "ADC":
		fmt.Printf("ADC %s", register)
	case "SUB":
		fmt.Printf("SUB %s", register)
	case "SBB":
		fmt.Printf("SBB %s", register)
	case "ANA":
		fmt.Printf("ANA %s", register)
	case "XRA":
		fmt.Printf("XRA %s", register)
	case "ORA":
		fmt.Printf("ORA %s", register)
	case "CMP":
		fmt.Printf("CMP %s", register)
	case "POP":
		fmt.Printf("POP %s", register)
	case "PUSH":
		fmt.Printf("PUSH %s", register)
	case "RST":
		fmt.Printf("RST %s", register)
	default:
		if instruction == "ADI" || instruction == "ACI" || instruction == "SUI" || instruction == "SBI" || instruction == "ANI" || instruction == "XRI" || instruction == "ORI" || instruction == "CPI" {
			fmt.Printf("%s %02x", instruction, buffer[pc+1])
			opcodes = 2
		}
	}

	fmt.Printf("\n")
	return opcodes
}
