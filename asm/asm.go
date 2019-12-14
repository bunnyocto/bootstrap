package asm

import . "bootstrap/defs"

import "strings"
import "strconv"
import "io"
import "bufio"
import "encoding/base64"
import "fmt"

type EmitContext struct {
	offset        int
	memory        []uint8
	labels        map[string]labelDef
	lateResolves  []lateResolve
	xLabels       map[string]labelDef
	xLateResolves []lateResolve
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
		offset:  0,
		memory:  make([]uint8, 16),
		labels:  make(map[string]labelDef),
		xLabels: make(map[string]labelDef),
	}
}

func NewEmitContext(offset int, memsz int) *EmitContext {
	return &EmitContext{
		offset:  offset,
		memory:  make([]uint8, memsz),
		labels:  make(map[string]labelDef),
		xLabels: make(map[string]labelDef),
	}
}

func (ec *EmitContext) Memory() []uint8 {
	return ec.memory
}

func (ec *EmitContext) Resolve() {
	for _, lr := range ec.lateResolves {
		v, found := ec.labels[lr.name]

		if !found {
			panic("Label `" + lr.name + "` not found!")
		}

		ec.offset = lr.pos
		EmitLDC(ec, lr.op, uint32(v.pos))
	}

	remXLateResolves := make([]lateResolve, 0)

	for _, xlr := range ec.xLateResolves {
		v, found := ec.xLabels[xlr.name]

		if !found {
			remXLateResolves = append(remXLateResolves, xlr)
		}

		ec.offset = xlr.pos
		EmitLDC(ec, xlr.op, uint32(v.pos))
	}

	ec.lateResolves = make([]lateResolve, 0)
	ec.labels = make(map[string]labelDef)
	ec.xLateResolves = remXLateResolves
}

func (ec *EmitContext) GetXUnresolved() []string {
	strs := make([]string, 0)

	for _, xlr := range ec.xLateResolves {
		strs = append(strs, xlr.name)
	}

	return strs
}

func EmitBytes(ec *EmitContext, data []byte) error {
	for ec.offset+len(data) >= len(ec.memory) {
		new_mem := make([]uint8, len(ec.memory)+16)
		copy(new_mem, ec.memory)
		ec.memory = new_mem
	}

	jx := 0
	for _, v := range data {
		ec.memory[ec.offset+jx] = v
		jx++
	}
	ec.offset += jx

	return nil
}

func EmitLDC(ec *EmitContext, op uint8, val uint32) error {
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

	return nil
}

func EmitOP(ec *EmitContext, op uint8, dst uint8, src uint8) error {
	for ec.offset+1 >= len(ec.memory) {
		new_mem := make([]uint8, len(ec.memory)+16)
		copy(new_mem, ec.memory)
		ec.memory = new_mem
	}

	ec.memory[ec.offset+0] = op
	ec.memory[ec.offset+1] = (dst << 4) | src
	ec.offset += 2

	return nil
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
	case OP_JEQ:
		return "jeq"
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

func Str2Reg(i string) (uint8, error) {
	switch i {
	case "ra", "A":
		return REG_A, nil
	case "rb", "B":
		return REG_B, nil
	case "rc", "C":
		return REG_C, nil
	case "rd", "D":
		return REG_D, nil
	case "re", "E":
		return REG_E, nil
	case "rf", "F":
		return REG_F, nil
	case "rg", "G":
		return REG_G, nil
	case "rh", "H":
		return REG_H, nil
	case "ri", "I":
		return REG_I, nil
	case "rj", "J":
		return REG_J, nil
	case "rk", "K":
		return REG_K, nil
	case "rl", "L":
		return REG_L, nil
	case "rm", "M":
		return REG_M, nil
	case "rn", "N":
		return REG_N, nil
	case "ro", "O":
		return REG_O, nil
	case "rip", "IP":
		return REG_IP, nil
	}

	return 0, fmt.Errorf("Unknown register %q!", i)
}

func Str2Op(i string) (uint8, error) {
	switch i {
	case "add":
		return OP_ADD, nil
	case "sub":
		return OP_SUB, nil
	case "div":
		return OP_DIV, nil
	case "mul":
		return OP_MUL, nil
	case "nop":
		return OP_NOP, nil
	case "hlt":
		return OP_HLT, nil
	case "ldcc", "ldc.c":
		return OP_LDCC, nil
	case "ldcb", "ldc.b":
		return OP_LDCB, nil
	case "ldca", "ldc.a":
		return OP_LDCA, nil
	case "dec":
		return OP_DEC, nil
	case "inc":
		return OP_INC, nil
	case "out":
		return OP_OUT, nil
	case "jmp":
		return OP_JMP, nil
	case "ldb":
		return OP_LDB, nil
	case "mov":
		return OP_MOV, nil
	case "jnz":
		return OP_JNZ, nil
	case "fail":
		return OP_FAIL, nil
	case "jiz":
		return OP_JIZ, nil
	case "pop":
		return OP_POP, nil
	case "push":
		return OP_PUSH, nil
	case "shl":
		return OP_SHL, nil
	case "shr":
		return OP_SHR, nil
	case "stb":
		return OP_STB, nil
	case "call":
		return OP_CALL, nil
	case "ret":
		return OP_RET, nil
	case "byt":
		return OP_BYT, nil
	case "jne":
		return OP_JNE, nil
	case "jeq":
		return OP_JEQ, nil
	}

	return 0, fmt.Errorf("Unknown op %q!", i)
}

func AsmReader(ec *EmitContext, rd io.Reader) error {
	sc := bufio.NewScanner(rd)

	for sc.Scan() {
		ln := sc.Text()
		err := Asm(ec, ln)

		if err != nil {
			return err
		}
	}

	return nil
}

func AsmLns(ec *EmitContext, i []string) error {
	for _, v := range i {
		err := Asm(ec, v)

		if err != nil {
			return err
		}
	}

	return nil
}

func Asm(ec *EmitContext, i string) error {
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
		return nil
	}

	if len(flds) >= 1 {
		switch flds[0] {
		case ".c":
			return nil
		}
	}

	if len(flds) == 1 {
		switch flds[0] {
		case ".c":
			return nil
		case "nop", "hlt", "fail":
			op, err := Str2Op(flds[0])

			if err != nil {
				return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
			}

			return EmitOP(ec, op, 0, 0)
		default:
			return fmt.Errorf("Invalid line: %q! Unknown instruction or wrong number of arguments!", i)
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

			_, found := ec.labels[name]

			if found {
				return fmt.Errorf("Invalid line: %q! Duplicate label: %q!", i, name)
			}

			ec.labels[name] = ld
		case ".x":
			name := flds[1]
			pos := ec.offset
			ld := labelDef{
				name: name,
				pos:  pos,
			}

			_, found := ec.xLabels[name]

			if found {
				return fmt.Errorf("Invalid line: %q! Duplicate label: %q!", i, name)
			}

			ec.xLabels[name] = ld
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

		case ".xadra", ".xadrb", ".xadrc":
			name := flds[1]
			lr := lateResolve{
				name: name,
				pos:  ec.offset,
			}

			switch flds[0] {
			case ".xadra":
				lr.op = OP_LDCA
			case ".xadrb":
				lr.op = OP_LDCB
			case ".xadrc":
				lr.op = OP_LDCC
			default:
				panic("can't handle this")
			}

			ec.xLateResolves = append(ec.xLateResolves, lr)

			ec.offset += 5

		case ".lra", ".lrb", ".lrc":
			name := flds[1]
			v, found := ec.labels[name]

			if !found {
				return fmt.Errorf("Invalid line: %q! Unknown label: %q!", i, name)
			}

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

		case "strdata":
			strval, _ := base64.StdEncoding.DecodeString(flds[1])
			return EmitBytes(ec, strval)
		case "ldca", "ldcb", "ldcc":
			val, _ := strconv.ParseUint(flds[1], 0, 32)
			op, err := Str2Op(flds[0])

			if err != nil {
				return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
			}

			return EmitLDC(ec, op, uint32(val))
		case "inc", "dec", "jmp", "ret", "byt":
			op, err := Str2Op(flds[0])

			if err != nil {
				return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
			}

			reg, err := Str2Reg(flds[1])

			if err != nil {
				return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
			}

			return EmitOP(ec, op, reg, 0)
		default:
			return fmt.Errorf("Invalid line: %q!", i)
		}
	} else if len(flds) == 3 {
		op, err := Str2Op(flds[0])

		if err != nil {
			return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
		}

		dst, err := Str2Reg(flds[1])

		if err != nil {
			return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
		}

		src, err := Str2Reg(flds[2])

		if err != nil {
			return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
		}

		return EmitOP(ec, op, dst, src)
	}

	return nil
}
