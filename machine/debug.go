package machine

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type proceeder interface {
	next()
}

type debugger struct {
	instructionExec proceeder
	advanceOP       *int
}

type defaultProceeder struct{}

type remoteProceeder struct {
	conn       net.Conn
	ptrAdvance *int
}

func (rmp *remoteProceeder) next() {
	n, s := getInput()
	fmt.Println(n)
	*rmp.ptrAdvance = n
	rmp.conn.Write([]byte(s))
}

func (defp defaultProceeder) next() {
	// time.Sleep(500 * time.Millisecond)
}

func (dbg debugger) Debug(cpu *Cpu) {
	cpu.pc = 0x100
	for cpu.pc < uint16(MemorySize) {
		disassebmle(cpu)
		debugCpuState(cpu)
		if *dbg.advanceOP == 0 {
			dbg.advance()
		} else {
			*dbg.advanceOP -= 1
		}
		cpu.Step()
	}
}

func (dbg debugger) advance() {
	dbg.instructionExec.next()
}

// func debugInstructions() {
// 	cpu := InitCpu()
// 	pc := 0
// 	for pc < len(cpu.memory) {
// 		cpu.executeInstruction()
// 		n := disassebmle(cpu.memory, int(pc))
// 		pc += n
// 		time.Sleep(150 * time.Millisecond)
// 	}
// }

func InitDebugger(flag string) debugger {
	dbg := debugger{advanceOP: new(int)}

	if flag == "-rd" {
		dbg.instructionExec = &remoteProceeder{
			conn:       connRemoteDbg(),
			ptrAdvance: dbg.advanceOP,
		}
	} else {
		dbg.instructionExec = defaultProceeder{}
	}

	return dbg
}

func getInput() (int, string) {
	var s string
	buffer := make([]byte, 128)
	for {

		fmt.Print()
		r := bufio.NewReader(os.Stdin)
		n, _ := r.Read(buffer)
		s = strings.TrimSpace(string(buffer[0:n]))

		if len(s) > 1 {

			n, err := strconv.Atoi(s[2:])
			if err != nil {
				log.Fatal(err)
			}

			return n - 1, s

		} else if s == "q" {

			os.Exit(1)

		} else if s == "s" {
			return 0, s
		}
	}
}

func connRemoteDbg() net.Conn {
	for {
		conn, err := net.Dial("tcp", "127.0.0.1:8080")
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			fmt.Println(err)
		} else {
			return conn
		}
	}
}

func disassebmle(cpuState *Cpu) {
	pc := cpuState.pc
	opcode := getOpcode(&cpuState.memory[pc])
	fmt.Printf("%02x ", pc)
	fmt.Printf("Instruction: %s ", opcode.Name)
	fmt.Printf("          Immediate: 0x%02x\n", cpuState.memory[pc+1])
}

func debugCpuState(cpu *Cpu) {
	fmt.Printf(" A: %d\n B: %d\n C: %d\n D: %d\n E: %d\n  H: %d\n L: %d\n", cpu.regs.a, cpu.regs.b, cpu.regs.c, cpu.regs.d, cpu.regs.e, cpu.regs.h, cpu.regs.l)
	fmt.Printf(" SP: %02x\n PC: %02x \n", cpu.sp, cpu.pc)
	fmt.Printf(" flags: sign - %d, zero - %d, aux carry - %d, parity - %d, carry - %d  \n", cpu.flags.s, cpu.flags.z, cpu.flags.ac, cpu.flags.p, cpu.flags.cy)
	fmt.Println("=========================================================================")

}
