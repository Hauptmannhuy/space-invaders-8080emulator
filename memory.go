package main

// 64KB RAM
type memory [65536]byte

func (mem *memory) write() {

}

func (mem *memory) read() {

}

func make16BitAddr(b1, b2 uint8) uint16 {
	return (uint16(b1) << 8) | uint16(b2)
}
