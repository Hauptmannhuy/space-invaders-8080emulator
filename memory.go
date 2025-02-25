package main

// 64KB RAM
type memory [65536]byte

func (mem *memory) write(addr uint16, val uint8) {
	mem[addr] = val
}

func (mem *memory) read(addr uint16) uint8 {
	return mem[addr]
}

func make16bit(b1, b2 uint8) uint16 {
	return (uint16(b1) << 8) | uint16(b2)
}
