package asm

import . "bootstrap/defs"

import "strings"
import "strconv"
import "io"
import "bufio"
import "encoding/base64"

//import "fmt"

type EmitContext struct {
	offset       int
	memory       []uint8
	labels       []labelDef
	lateResolves []lateResolve
}

type labelDef struct {
	name string
	pos  int
}

type lateResolve struct {
	name string
	pos  int
	op   uint8
}

func NewDefaultEmitContext() *EmitContext {
	return &EmitContext{
		offset: 0,
		memory: make([]uint8, 16),
		labels: make([]labelDef, 0),
	}
}

func NewEmitContext(offset int, memsz int) *EmitContext {
	return &EmitContext {
		offset: offset,
		memory: make([]uint8, memsz),
		labels: make([]labelDef, 0),
	}
}

func (ec *EmitContext) Memory() []uint8 {
	return ec.memory
}

func (ec *EmitContext) Resolve() {
	for _, lr := range ec.lateResolves {
		found := false
		for _, v := range ec.labels {
			if v.name == lr.name {
				found = true

				ec.offset = lr.pos
				EmitLDC(ec, lr.op, uint32(v.pos))

				break
			}
		}

		if !found {
			panic("Label `" + lr.name + "` not found!")
		}
	}
}

func EmitBytes(ec *EmitContext, data []byte) {
	jx := 0
	for _, v := range data {
		ec.memory[ec.offset+jx] = v
		jx++
	}
	ec.offset += jx
}

func EmitLDC(ec *EmitContext, op uint8, val uint32) {
	for ec.offset+4 >= len(ec.memory) {
		new_mem := make([]uint8, len(ec.memory)+16)
		copy(new_mem, ec.memory)
		ec.memory = new_mem
	}

	ec.memory[ec.offset+0] = op
	ec.memory[ec.offset+1] = uint8((val >> 24) & 0xFF)
	ec.memory[ec.offset+2] = uint8((val >> 16) & 0xFF)
	ec.memory[ec.offset+3] = uint8((val >> 8) & 0xFF)
	ec.memory[ec.offset+4] = uint8(val & 0xFF)
	ec.offset += 5
}

func EmitOP(ec *EmitContext, op uint8, dst uint8, src uint8) {
	for ec.offset+1 >= len(ec.memory) {
		new_mem := make([]uint8, len(ec.memory)+16)
		copy(new_mem, ec.memory)
		ec.memory = new_mem
	}

	ec.memory[ec.offset+0] = op
	ec.memory[ec.offset+1] = (dst << 4) | src
	ec.offset += 2
}

func Op2Str(o uint8) string {
	switch o {
	case OP_LDCA:
		return "ldca"
	case OP_LDCB:
		return "ldcb"
	case OP_LDCC:
		return "ldcc"
	case OP_ADD:
		return "add"
	case OP_SUB:
		return "sub"
	case OP_MUL:
		return "mul"
	case OP_DIV:
		return "div"
	case OP_INC:
		return "inc"
	case OP_DEC:
		return "dec"
	case OP_MOV:
		return "mov"
	case OP_MOD:
		return "mod"
	case OP_SADD:
		return "sadd"
	case OP_SSUB:
		return "ssub"
	case OP_SMUL:
		return "smul"
	case OP_SDIV:
		return "sdiv"
	case OP_XCHG:
		return "xchg"
	case OP_SHR:
		return "shr"
	case OP_SHL:
		return "shl"
	case OP_JMP:
		return "jmp"
	case OP_JNZ:
		return "jnz"
	case OP_JIZ:
		return "jiz"
	case OP_CALL:
		return "call"
	case OP_RET:
		return "ret"
	case OP_LDB:
		return "ldb"
	case OP_STB:
		return "stb"
	case OP_LDH:
		return "ldh"
	case OP_STH:
		return "sth"
	case OP_LDW:
		return "ldw"
	case OP_STW:
		return "stw"
	case OP_PUSH:
		return "push"
	case OP_POP:
		return "pop"
	case OP_FAIL:
		return "fail"
	case OP_OUT:
		return "out"
	case OP_HLT:
		return "hlt"
	case OP_NOP:
		return "nop"
	case OP_BYT:
		return "byt"
	case OP_JNE:
		return "jne"
	default:
		panic("can't handle this")
	}
}

func Reg2Str(r uint8) string {
	switch r {
	case REG_A:
		return "ra"
	case REG_B:
		return "rb"
	case REG_C:
		return "rc"
	case REG_D:
		return "rd"
	case REG_E:
		return "re"
	case REG_F:
		return "rf"
	case REG_G:
		return "rg"
	case REG_H:
		return "rh"
	case REG_I:
		return "ri"
	case REG_J:
		return "rj"
	case REG_K:
		return "rk"
	case REG_L:
		return "rl"
	case REG_M:
		return "rm"
	case REG_N:
		return "rn"
	case REG_O:
		return "ro"
	case REG_IP:
		return "rip"
	default:
		panic("can't handle this")
	}
}

func Str2Reg(i string) uint8 {
	switch i {
	case "ra", "A":
		return REG_A
	case "rb", "B":
		return REG_B
	case "rc", "C":
		return REG_C
	case "rd", "D":
		return REG_D
	case "re", "E":
		return REG_E
	case "rf", "F":
		return REG_F
	case "rg", "G":
		return REG_G
	case "rh", "H":
		return REG_H
	case "ri", "I":
		return REG_I
	case "rj", "J":
		return REG_J
	case "rk", "K":
		return REG_K
	case "rl", "L":
		return REG_L
	case "rm", "M":
		return REG_M
	case "rn", "N":
		return REG_N
	case "ro", "O":
		return REG_O
	case "rip", "IP":
		return REG_IP
	}

	panic("no such reg: " + i)
}

func Str2Op(i string) uint8 {
	switch i {
	case "add":
		return OP_ADD
	case "sub":
		return OP_SUB
	case "div":
		return OP_DIV
	case "mul":
		return OP_MUL
	case "nop":
		return OP_NOP
	case "hlt":
		return OP_HLT
	case "ldcc", "ldc.c":
		return OP_LDCC
	case "ldcb", "ldc.b":
		return OP_LDCB
	case "ldca", "ldc.a":
		return OP_LDCA
	case "dec":
		return OP_DEC
	case "inc":
		return OP_INC
	case "out":
		return OP_OUT
	case "jmp":
		return OP_JMP
	case "ldb":
		return OP_LDB
	case "mov":
		return OP_MOV
	case "jnz":
		return OP_JNZ
	case "fail":
		return OP_FAIL
	case "jiz":
		return OP_JIZ
	case "pop":
		return OP_POP
	case "push":
		return OP_PUSH
	case "shl":
		return OP_SHL
	case "shr":
		return OP_SHR
	case "stb":
		return OP_STB
	case "call":
		return OP_CALL
	case "ret":
		return OP_RET
	case "byt":
		return OP_BYT
	case "jne":
		return OP_JNE
	}

	panic("no such op: " + i)
}

func AsmReader(ec *EmitContext, rd io.Reader) {
	sc := bufio.NewScanner(rd)

	for sc.Scan() {
		ln := sc.Text()
		Asm(ec, ln)
	}
}

func AsmLns(ec *EmitContext, i []string) {
	for _, v := range i {
		Asm(ec, v)
	}
}

func Asm(ec *EmitContext, i string) {
	rawflds := strings.Fields(i)
	flds := make([]string, 0)
	for _, v := range rawflds {
		if !strings.HasPrefix(v, "%:") {
			flds = append(flds, v)
		} else {
			break
		}
	}

	if len(flds) == 0 {
		return
	}

	if len(flds) >= 1 {
		switch flds[0] {
		case ".c":
			return
		}
	}

	if len(flds) == 1 {
		switch flds[0] {
		case ".c":
			return
		case "nop", "hlt", "fail":
			EmitOP(ec, Str2Op(flds[0]), 0, 0)
		default:
			panic("can't handle this: " + flds[0])
		}
	} else if len(flds) == 2 {
		switch flds[0] {
		case ".l":
			name := flds[1]
			pos := ec.offset
			ld := labelDef{
				name: name,
				pos:  pos,
			}
			ec.labels = append(ec.labels, ld)
		case ".adra", ".adrb", ".adrc":
			name := flds[1]
			lr := lateResolve{
				name: name,
				pos:  ec.offset,
			}

			switch flds[0] {
			case ".adra":
				lr.op = OP_LDCA
			case ".adrb":
				lr.op = OP_LDCB
			case ".adrc":
				lr.op = OP_LDCC
			default:
				panic("can't handle this")
			}

			ec.lateResolves = append(ec.lateResolves, lr)

			ec.offset += 5

		case ".lra", ".lrb", ".lrc":
			name := flds[1]
			found := false
			for _, v := range ec.labels {
				if v.name == name {
					found = true
					switch flds[0] {
					case ".lra":
						EmitLDC(ec, OP_LDCA, uint32(v.pos))
					case ".lrb":
						EmitLDC(ec, OP_LDCB, uint32(v.pos))
					case ".lrc":
						EmitLDC(ec, OP_LDCC, uint32(v.pos))
					default:
						panic("can't handle this")
					}
					break
				}
			}

			if !found {
				panic("Label `" + name + "` not found!")
			}
		case "strdata":
			strval, _ := base64.StdEncoding.DecodeString(flds[1])
			EmitBytes(ec, strval)
		case "ldca", "ldcb", "ldcc":
			val, _ := strconv.ParseUint(flds[1], 0, 32)
			EmitLDC(ec, Str2Op(flds[0]), uint32(val))
		case "inc", "dec", "jmp", "ret", "byt":
			EmitOP(ec, Str2Op(flds[0]), Str2Reg(flds[1]), 0)
		default:
			panic("can't handle this: " + flds[0])
		}
	} else if len(flds) == 3 {
		EmitOP(ec, Str2Op(flds[0]), Str2Reg(flds[1]), Str2Reg(flds[2]))
	}

	return
}
