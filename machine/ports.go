package machine

type Ports [8]uint8

type IO_handler interface {
	InPort(cpu *Cpu) uint8
	OutPort(cpu *Cpu)
}
