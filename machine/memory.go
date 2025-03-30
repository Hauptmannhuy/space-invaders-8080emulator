package machine

const MemorySize = uint16(65535)

// 64KB RAM
type Memory [MemorySize]byte

func (mem *Memory) write(addr uint16, val uint8) {
	mem[addr] = val
}

func (mem *Memory) read(addr uint16) uint8 {
	return mem[addr]
}

func (cpu *Cpu) LoadRom(buff []byte) {
	copy(cpu.memory[:], buff)
}
