package main

import "fmt"

/// write an assembler
/// make this emulator take an input file as the prog

var prog = []word{	0x7c01, 0x0030, 0x7de1, 0x1000, 0x0020, 0x7803, 0x1000, 0xc00d,
					0x7dc1, 0x001a, 0x8861, 0x7c01, 0x2000, 0x2161, 0x2000, 0x8463,
					0x806d, 0x7dc1, 0x000d, 0x9031, 0x7c10, 0x0018, 0x7dc1, 0x001d,
					0x9037, 0x61c1, 0x0000, 0x0000, 0x0000, 0x0000}

type word uint16
type hword uint8

var lit = []word{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
				 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F}

type Cpu struct {
	r [8]word
	pc, sp word
	o word///
	mem [0x10000]word
	running bool
	cycles int
}


func newCpu() *Cpu {
	s := Cpu{}
	s.sp = 0xFFFF
	s.running = false
	return &s
}


func (s *Cpu) start() {
	s.running = true
	s.run()
}

func (s *Cpu) stop() {
	s.running = false
}

func (s *Cpu) error(err string) {
	fmt.Println(err)
	s.stop()
}

var instr word

func peekNextWord(pc word) word {
	return prog[pc]
}


var back int

func (s *Cpu) fetchNextWord() word {
	back++
	s.cycles++
	pc := s.pc
	s.pc++
	return prog[pc]
}

func (s *Cpu) fetch() {
	pc := s.pc
	s.pc++
	instr = prog[pc]
}

func (s *Cpu) skip() {
	s.fetch()
	s.decode()
}

var opcode byte
var vo word // opcode value for non-basic instructions
var va, vb word // a, b values
var vma, vmb word // a, b memory values
var op1, op2 *word

func (s *Cpu) decode() {
	opcode = byte(instr & 0x000F)
	if opcode == 0 {
		vo = (instr & 0x03F0) >> 4
		va = (instr & 0xFC00) >> 10

		op1 = s.decode_arg(1)
	} else {
		va = (instr & 0x03F0) >> 4
		vb = (instr & 0xFC00) >> 10

//	dbg("va =", va)///
//	dbg("vb =", vb)///
		op1 = s.decode_arg(1)
//	dbg("*op1 =", *op1)///
		op2 = s.decode_arg(2)
//	dbg("*op2 =", *op2)///
	}
}

func (s *Cpu) decode_arg(i int) *word {
	var c, val word
	var mc *word

	if i == 1 {
		c = va
		mc = &vma
	} else {
		c = vb
		mc = &vmb
	}

	switch c {
	case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07:
		return &s.r[c]
	case 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F:
		return &s.mem[s.r[c-0x08]]
	case 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17:
		*mc = s.fetchNextWord()
		return &s.mem[*mc + s.r[c-0x10]]
	case 0x18:
		val = s.stackPop()
		return &val
	case 0x19:
		val = s.stackPeek()
		return &val
	case 0x1A:///
	case 0x1B:///
	case 0x1C:
		return &s.pc
	case 0x1D:///
	case 0x1E:
		*mc = s.fetchNextWord()
		return &s.mem[val]
	case 0x1F:
		val = s.fetchNextWord()
		return &val
	case 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F,
	     0x30, 0X31, 0X32, 0X33, 0X34, 0X35, 0X36, 0X37, 0X38, 0X39, 0X3A, 0X3B, 0X3C, 0X3D, 0X3E, 0X3F:///
		return &lit[c-0x20]
	}
	return nil
}

func disassemble() string {
	i := disassemble_op()
	if opcode != 0x0 {
		i += " " + disassemble_arg(1) + ", " + disassemble_arg(2)
	} else {
		i += " " + disassemble_arg(1)
	}
	return i
}

func disassemble_op() string {
	switch opcode {
	case 0x0:
		switch vo {
		case 0x00:
		case 0x01:
			return "JSR"
		case 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		     0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		     0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F,
		     0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F:
		}
	case 0x1:
		return "SET"
	case 0x2:
		return "ADD"
	case 0x3:
		return "SUB"
	case 0x4:
		return "MUL"
	case 0x5:
		return "DIV"
	case 0x6:
		return "MOD"
	case 0x7:
		return "SHL"
	case 0x8:
		return "SHR"
	case 0x9:
		return "AND"
	case 0xA:
		return "BOR"
	case 0xB:
		return "XOR"
	case 0xC:
		return "IFE"
	case 0xD:
		return "IFN"
	case 0xE:
		return "IFG"
	case 0xF:
		return "IFB"
	}
	return ""
}

func disassemble_arg(i int) string {
	var c, mc, op word

	if i == 1 {
		c = va
		mc = vma
		op = *op1
	} else {
		c = vb
		mc = vmb
		op = *op2
	}
	switch c {
	case 0x00:
		return "A"
	case 0x01:
		return "B"
	case 0x02:
		return "C"
	case 0x03:
		return "X"
	case 0x04:
		return "Y"
	case 0x05:
		return "Z"
	case 0x06:
		return "I"
	case 0x07:
		return "J"
	case 0x08:
		return "[A]"
	case 0x09:
		return "[B]"
	case 0x0A:
		return "[C]"
	case 0x0B:
		return "[X]"
	case 0x0C:
		return "[Y]"
	case 0x0D:
		return "[Z]"
	case 0x0E:
		return "[I]"
	case 0x0F:
		return "[J]"
	case 0x10:
		return fmt.Sprintf("[0x%04X + A]", mc)
	case 0x11:
		return fmt.Sprintf("[0x%04X + B]", mc)
	case 0x12:
		return fmt.Sprintf("[0x%04X + C]", mc)
	case 0x13:
		return fmt.Sprintf("[0x%04X + X]", mc)
	case 0x14:
		return fmt.Sprintf("[0x%04X + Y]", mc)
	case 0x15:
		return fmt.Sprintf("[0x%04X + Z]", mc)
	case 0x16:
		return fmt.Sprintf("[0x%04X + I]", mc)
	case 0x17:
		return fmt.Sprintf("[0x%04X + J]", mc)
	case 0x18:
		return "POP"
	case 0x19:
		return "PEEK"
	case 0x1A:
		return "PUSH"
	case 0x1B:
		return "SP"
	case 0x1C:
		return "PC"
	case 0x1D:
		return "O"
	case 0x1E:
		return fmt.Sprintf("[0x%04X]", mc)
	case 0x1F:
		return fmt.Sprintf("0x%04X", op)
	case 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F,
	     0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F:
		return fmt.Sprintf("0x%02X", c-0x20)
	}
	return ""
}


func (s *Cpu) eval() {
	switch opcode {
	case 0x0:
		switch vo {
		case 0x00:
		case 0x01:
			s.jsr()
		case 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		     0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		     0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F,
		     0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F:
		}
	case 0x1:
		s.set()
	case 0x2:
		s.add()
	case 0x3:
		s.sub()
	case 0x4:
		s.mul()
	case 0x5:
		s.div()
	case 0x6:
		s.mod()
	case 0x7:
		s.shl()
	case 0x8:
		s.shr()
	case 0x9:
		s.and()
	case 0xA:
		s.bor()
	case 0xB:
		s.xor()
	case 0xC:
		s.ife()
	case 0xD:
		s.ifn()
	case 0xE:
		s.ifg()
	case 0xF:
		s.ifb()
	default:
		s.error("eval(): wrong opcode")
	}
}

func (s *Cpu) set() {
	*op1 = *op2
	s.cycles++
}

func (s *Cpu) add() {
	*op1 += *op2
	s.cycles += 2
	if
}

func (s *Cpu) sub() {
	*op1 -= *op2
	s.cycles += 2
}

func (s *Cpu) mul() {
	*op1 *= *op2
	s.cycles += 2
}

func (s *Cpu) div() {
	*op1 /= *op2
	s.cycles += 3
}

func (s *Cpu) mod() {
	*op1 %= *op2
	s.cycles += 3
}

func (s *Cpu) shl() {
	*op1 <<= *op2
	s.cycles += 2
}

func (s *Cpu) shr() {
	*op1 >>= *op2
	s.cycles += 2
}

func (s *Cpu) and() {
	*op1 &= *op2
	s.cycles++
}

func (s *Cpu) bor() {
	*op1 |= *op2
	s.cycles++
}

func (s *Cpu) xor() {
	*op1 ^= *op2
	s.cycles++
}

func (s *Cpu) ife() {
	s.cycles += 2
	if *op1 == *op2 {
		s.skip()
		s.cycles++
	}

}

func (s *Cpu) ifn() {
	s.cycles += 2
	if *op1 == *op2 {
		s.skip()
		s.cycles++
	}
}

func (s *Cpu) ifg() {
	s.cycles += 2
	if *op1 <= *op2 {
		s.skip()
		s.cycles++
	}
}

func (s *Cpu) ifb() {
	s.cycles += 2
	if *op1 & *op2 == 0 {
		s.skip()
		s.cycles++
	}
}

func (s *Cpu) jsr() {
	s.cycles += 2
	ret := s.pc
	s.pc = *op1
	*op1 = ret
	s.stackPush()
}

func (s *Cpu) stackPush() {
	s.sp--
	s.mem[s.sp] = *op1
}

func (s *Cpu) stackPop() word {
	s.sp++
	return s.mem[s.sp-1]
}

func (s *Cpu) stackPeek() word {
	return s.mem[s.sp]
}

func dumpHeader() {
	fmt.Println("PC   SP   O    A    B    C    X    Y    Z    I    J    instr Instruction")
	fmt.Println("---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ----- -----------")
}

func (s Cpu) dumpState() {
	fmt.Printf("%04X %04X %04X ", s.pc-word(back+1), s.sp, s.o)
	back=0
	for i := 0; i < 8; i++ {
		fmt.Printf("%04X ", s.r[i])
	}
	fmt.Printf("%04X  ", instr)
	fmt.Printf("%s", disassemble())
	fmt.Println()
}

func (s *Cpu) run() {
	dumpHeader()
	for s.running {
		s.fetch()
		s.decode()
		s.dumpState()
		s.eval()
	}
}

var debug bool ////make it a flag
func main() {
	debug = true
	s := newCpu()
	s.start()
}

func dbg(s string, w word) {
		fmt.Printf("%s %04X\n", s, w)
}
