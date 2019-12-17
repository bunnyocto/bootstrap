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
	origin        int
	saves         []uint8
	macros        map[string]string
	vars          map[string]uint8
	labelCounter  uint32
	nestingInfos  []nestingInfo
}

type nestingInfo struct {
	counter uint32
	kind    string
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
		offset:       0,
		memory:       make([]uint8, 16),
		labels:       make(map[string]labelDef),
		xLabels:      make(map[string]labelDef),
		macros:       make(map[string]string),
		vars:         make(map[string]uint8),
		labelCounter: 0,
	}
}

func NewEmitContext(offset int, memsz int) *EmitContext {
	return &EmitContext{
		offset:       offset,
		memory:       make([]uint8, memsz),
		labels:       make(map[string]labelDef),
		xLabels:      make(map[string]labelDef),
		macros:       make(map[string]string),
		vars:         make(map[string]uint8),
		labelCounter: 0,
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

		old_offset := ec.offset
		ec.offset = lr.pos
		EmitLDC(ec, lr.op, uint32(v.pos))
		ec.offset = old_offset
	}

	remXLateResolves := make([]lateResolve, 0)

	for _, xlr := range ec.xLateResolves {
		v, found := ec.xLabels[xlr.name]

		if !found {
			remXLateResolves = append(remXLateResolves, xlr)
		}

		old_offset := ec.offset
		ec.offset = xlr.pos
		EmitLDC(ec, xlr.op, uint32(v.pos))
		ec.offset = old_offset
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

func EmitLabel(ec *EmitContext, name string) error {
	pos := ec.offset
	ld := labelDef{
		name: name,
		pos:  pos + ec.origin,
	}

	_, found := ec.labels[name]

	if found {
		return fmt.Errorf("Duplicate label: %q!", name)
	}

	ec.labels[name] = ld

	return nil
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
	case OP_LDN:
		return "ldn"
	case OP_LDP:
		return "ldp"
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

func Str2Reg(i string, ec *EmitContext) (uint8, error) {
	switch i {
	case "ra", "r1":
		return REG_A, nil
	case "rb", "r2":
		return REG_B, nil
	case "rc", "r3":
		return REG_C, nil
	case "rd", "r4":
		return REG_D, nil
	case "re", "r5":
		return REG_E, nil
	case "rf", "r6":
		return REG_F, nil
	case "rg", "r7":
		return REG_G, nil
	case "rh", "r8":
		return REG_H, nil
	case "ri", "r9":
		return REG_I, nil
	case "rj", "r10":
		return REG_J, nil
	case "rk", "r11":
		return REG_K, nil
	case "rl", "r12":
		return REG_L, nil
	case "rm", "r13":
		return REG_M, nil
	case "rn", "r14":
		return REG_N, nil
	case "ro", "r15":
		return REG_O, nil
	case "rip", "r0":
		return REG_IP, nil
	}

	if ec != nil {
		r, found := ec.vars[i]

		if found {
			return r, nil
		}
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
	case "ldp":
		return OP_LDP, nil
	case "ldn":
		return OP_LDN, nil
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
	fmt.Printf("ASM [%08x] %q\n", ec.offset, i)

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
		case ".defm":
			if len(flds) < 3 {
				return fmt.Errorf("Invalid line: %q! Need more arguments.", i)
			}

			s := ""
			for i := 2; i < len(flds); i++ {
				s += flds[i] + " "
			}
			ec.macros[flds[1]] = s
		case ".m":
			if len(flds) < 2 {
				return fmt.Errorf("Invalid line: %q! Need more arguments.", i)
			}

			macro, found := ec.macros[flds[1]]

			if !found {
				return fmt.Errorf("Invalid line: %q!Unknown macro: %q!", i, flds[1])
			}

			for i := 2; i < len(flds); i++ {
				macro = strings.Replace(macro, fmt.Sprintf("#%d", i-2), flds[i], -1)
			}

			lns := strings.Split(macro, "\\")

			for _, ln := range lns {
				err := Asm(ec, ln)

				if err != nil {
					return err
				}
			}

			return nil
		}
	}

	if len(flds) == 1 {
		switch flds[0] {
		case ".end":
			if len(ec.nestingInfos) < 1 {
				return fmt.Errorf("Invalid line: %q! Nothing to end.", i)
			}

			ni := ec.nestingInfos[len(ec.nestingInfos)-1]
			ec.nestingInfos = ec.nestingInfos[:len(ec.nestingInfos)-1]

			switch ni.kind {
			case "while":
				err := AsmLns(ec, []string{
					fmt.Sprintf(".adrc __L_%d_begin", ni.counter),
					"jmp rc",
				})

				if err != nil {
					return err
				}
			case "if":
			}

			return EmitLabel(ec, fmt.Sprintf("__L_%d_end", ni.counter))
		case ".eof":
			ec.saves = nil
		case ".ret":
			fmt.Println("LEN: %d", len(ec.saves))
			for i := len(ec.saves) - 1; i >= 0; i-- {
				err := EmitOP(ec, OP_POP, ec.saves[i], REG_O)

				if err != nil {
					return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
				}
			}

			return EmitOP(ec, OP_RET, REG_O, 0)
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
		case ".save":
			r, err := Str2Reg(flds[1], ec)

			if err != nil {
				return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
			}

			ec.saves = append(ec.saves, r)

			return EmitOP(ec, OP_PUSH, REG_O, r)
		case ".call":
			r, err := Str2Reg(flds[1], ec)

			if err != nil {
				return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
			}

			return EmitOP(ec, OP_CALL, r, REG_O)
		case ".origin":
			v, _ := strconv.ParseInt(flds[1], 0, 32)

			ec.origin = int(v)
		case ".alloc":
			v, _ := strconv.ParseInt(flds[1], 0, 32)
			ec.offset += int(v)
		case ".l":
			return EmitLabel(ec, flds[1])
		case ".x":
			name := flds[1]
			pos := ec.offset
			ld := labelDef{
				name: name,
				pos:  pos + ec.origin,
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
				EmitLDC(ec, OP_LDCA, uint32(v.pos+ec.origin))
			case ".lrb":
				EmitLDC(ec, OP_LDCB, uint32(v.pos+ec.origin))
			case ".lrc":
				EmitLDC(ec, OP_LDCC, uint32(v.pos+ec.origin))
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

			reg, err := Str2Reg(flds[1], ec)

			if err != nil {
				return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
			}

			return EmitOP(ec, op, reg, 0)
		case ".whilenz":
			ec.labelCounter++
			err := EmitLabel(ec, fmt.Sprintf("__L_%d_begin", ec.labelCounter))

			if err != nil {
				return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
			}

			err = AsmLns(ec, []string{
				fmt.Sprintf(".adrc __L_%d_end", ec.labelCounter),
				fmt.Sprintf("jiz rc %s", flds[1]),
			})

			if err != nil {
				return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
			}

			ni := nestingInfo{
				kind:    "while",
				counter: ec.labelCounter,
			}

			ec.nestingInfos = append(ec.nestingInfos, ni)

			return nil
		default:
			return fmt.Errorf("Invalid line: %q!", i)
		}
	} else if len(flds) == 3 {
		switch flds[0] {
		case ".var":
			name := flds[1]
			treg := flds[2]

			reg, err := Str2Reg(treg, ec)

			if err != nil {
				return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
			}

			ec.vars[name] = reg

			return nil
		}

		op, err := Str2Op(flds[0])

		if err != nil {
			return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
		}

		dst, err := Str2Reg(flds[1], ec)

		if err != nil {
			return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
		}

		switch op {
		case OP_LDN, OP_LDP:
			imm, _ := strconv.ParseUint(flds[2], 0, 4)

			return EmitOP(ec, op, dst, uint8(imm))
		}

		src, err := Str2Reg(flds[2], ec)

		if err != nil {
			return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
		}

		return EmitOP(ec, op, dst, src)
	} else if len(flds) == 4 {
		switch flds[0] {
		case ".if":
			ec.labelCounter++
			err := EmitLabel(ec, fmt.Sprintf("__L_%d_begin", ec.labelCounter))

			if err != nil {
				return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
			}

			op := ""

			switch flds[1] {
			case "eq":
				op = "jne"
			case "ne":
				op = "jeq"
			default:
				return fmt.Errorf("Invalid line: %q! Invalid argument %q!", i, flds[1])
			}

			err = AsmLns(ec, []string{
				fmt.Sprintf(".adrc __L_%d_end", ec.labelCounter),
				fmt.Sprintf("%s %s %s", op, flds[2], flds[3]),
			})

			if err != nil {
				return fmt.Errorf("Invalid line: %q! %s", i, err.Error())
			}

			ni := nestingInfo{
				kind:    "if",
				counter: ec.labelCounter,
			}

			ec.nestingInfos = append(ec.nestingInfos, ni)

			return nil
		}
	}

	return nil
}
