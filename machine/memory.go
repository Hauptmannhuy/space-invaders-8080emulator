package machine

const MemorySize = uint16(65535)
const (
	VRAMStart uint16 = 0x2400
	VRAMEnd   uint16 = 0x3FFF
)

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

func (cpu *Cpu) GetMemoryAt(offset uint16) uint8 {
	return cpu.memory.read(offset)
}

func (cpu *Cpu) CopyFrameBuffer() []byte {
	buffer := make([]byte, 7168)
	copy(cpu.memory[VRAMStart:VRAMEnd], buffer)
	return buffer
}
