package machine

const MemorySize = uint16(0xFFFF)
const (
	VRAMStart uint16 = 0x2400
	VRAMEnd   uint16 = 0x3FFF

	ROMstart uint16 = 0
	ROMend   uint16 = 0x1fff

	RAMend uint16 = 0x4000
)

// 64KB RAM
type Memory [MemorySize]byte

func (mem *Memory) write(addr uint16, val uint8) {
	if ROMend >= addr {
		// fmt.Printf("Error! unathorized write to ROM mem location %d\n", addr)
		return
	}

	if addr >= RAMend {
		// fmt.Printf("Error! unathorized write to not allocated RAM mem location %d\n", addr)
		return
	}

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
	copy(buffer, cpu.memory[VRAMStart:VRAMEnd])
	return buffer
}
