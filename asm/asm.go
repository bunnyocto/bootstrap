package asm

import . "bootstrap/defs"

import "strings"
import "strconv"
import "io"
import "bufio"
import "encoding/base64"
//import "fmt"

type EmitContext struct {
	offset int
	memory []uint8
	labels []labelDef
	lateResolves []lateResolve
}

type labelDef struct {
	name string
	pos int
}

type lateResolve struct {
	name string
	pos int
	op uint8
}

func NewDefaultEmitContext() *EmitContext {
	return &EmitContext{
		offset: 1024,
		memory: make([]uint8, 2048),
		labels: make([]labelDef, 0),
	}
}

func (ec *EmitContext) Memory() []uint8 {
	return ec.memory
}

func (ec *EmitContext) Resolve() {
	for _,lr := range(ec.lateResolves) {
		found := false
		for _,v := range(ec.labels) {
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
	ec.memory[ec.offset+0] = op
	ec.memory[ec.offset+1] = uint8((val >> 24) & 0xFF)
	ec.memory[ec.offset+2] = uint8((val >> 16) & 0xFF)
	ec.memory[ec.offset+3] = uint8((val >> 8) & 0xFF)
	ec.memory[ec.offset+4] = uint8(val & 0xFF)
	ec.offset += 5
}

func EmitOP(ec *EmitContext, op uint8, dst uint8, src uint8) {
	ec.memory[ec.offset+0] = op
	ec.memory[ec.offset+1] = (dst << 4) | src
	ec.offset += 2
}

func Str2Reg(i string) uint8 {
	switch i {
	case "ra":
		return REG_A
	case "rb":
		return REG_B
	case "rc":
		return REG_C
	case "rd":
		return REG_D
	case "re":
		return REG_E
	}

	panic("no such reg")
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
	case "ldcc":
		return OP_LDCC
	case "ldcb":
		return OP_LDCB
	case "ldca":
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
	for _,v := range(rawflds) {
		if !strings.HasPrefix(v,"%:") {
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
		case "nop","hlt","fail":
			EmitOP(ec, Str2Op(flds[0]), 0, 0)
		default:
			panic("can't handle this")
		}
	} else if len(flds) == 2 {
		switch flds[0] {
		case ".l":
			name := flds[1]
			pos := ec.offset
			ld := labelDef {
				name: name,
				pos: pos,
			}
			ec.labels = append(ec.labels, ld)
		case ".adra",".adrb",".adrc":
			name := flds[1]
			lr := lateResolve {
				name: name,
				pos: ec.offset,
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
			
		case ".lra",".lrb",".lrc":
			name := flds[1]
			found := false
			for _,v := range(ec.labels) {
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
		case "inc", "dec", "jmp":
			EmitOP(ec, Str2Op(flds[0]), Str2Reg(flds[1]), 0)
		default:
			panic("can't handle this")
		}
	} else if len(flds) == 3 {
		EmitOP(ec, Str2Op(flds[0]), Str2Reg(flds[1]), Str2Reg(flds[2]))
	}

	return
}
