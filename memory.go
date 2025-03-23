package main

import "os"

const MemorySize = uint16(65535)

// 64KB RAM
type memory [MemorySize]byte

func (mem *memory) write(addr uint16, val uint8) {
	mem[addr] = val
}

func (mem *memory) read(addr uint16) uint8 {
	return mem[addr]
}

func loadRom() *memory {
	buff, _ := os.ReadFile("./roms/cpudiag.bin")
	memory := memory{}
	copy(memory[:], buff)
	return &memory
}

func make16bit(hi, lo uint8) uint16 {
	return (uint16(hi) << 8) | uint16(lo)
}
