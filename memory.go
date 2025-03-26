package main

const MemorySize = uint16(65535)

// 64KB RAM
type memory [MemorySize]byte

func (mem *memory) write(addr uint16, val uint8) {
	mem[addr] = val
}

func (mem *memory) read(addr uint16) uint8 {
	return mem[addr]
}

func (memory *memory) loadRom(buff []byte) {
	copy(memory[:], buff)
}

func make16bit(hi, lo uint8) uint16 {
	return (uint16(hi) << 8) | uint16(lo)
}
