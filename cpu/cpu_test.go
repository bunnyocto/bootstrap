package cpu

import "testing"
import . "bootstrap/defs"
import "bootstrap/asm"

func TestDummy(t *testing.T) {
}

func emit_ldc(memory []uint8, op uint8, val uint32) []uint8 {
	memory[0] = op
	memory[1] = uint8((val >> 24) & 0xFF)
	memory[2] = uint8((val >> 16) & 0xFF)
	memory[3] = uint8((val >> 8) & 0xFF)
	memory[4] = uint8(val & 0xFF)
	return memory[5:]
}

func emit_op(memory []uint8, op uint8, dst uint8, src uint8) []uint8 {
	memory[0] = op
	memory[1] = (dst << 4) | src
	return memory[2:]
}

func TestLdca(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			"ldca 0x33445566",
			"hlt",
		})

	Execute(regs, ec.Memory())

	if regs[REG_A] != 0x33445566 {
		t.Fatalf("Expected %x but got %x!", 0x33445566, regs[REG_A])
	}
}

func TestLdcb(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			"ldcb 0x33445566",
			"hlt",
		})

	Execute(regs, ec.Memory())

	if regs[REG_B] != 0x33445566 {
		t.Fatalf("Expected %x but got %x!", 0x33445566, regs[REG_B])
	}
}

func TestLdcc(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			"ldcc 0x33445566",
			"hlt",
		})

	Execute(regs, ec.Memory())

	if regs[REG_C] != 0x33445566 {
		t.Fatalf("Expected %x but got %x!", 0x33445566, regs[REG_C])
	}
}

func TestAdd(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			"ldcb 0x12345678",
			"ldcc 0x23451209",
			"add rb rc",
			"hlt",
		})

	Execute(regs, ec.Memory())

	if regs[REG_B] != (0x12345678 + 0x23451209) {
		t.Fatalf("Expected %x but got %x1", (0x12345678 + 0x23451209), regs[REG_B])
	}
}

func TestSub(t *testing.T) {
	regs := make([]uint32, 0x10)
	memory := make([]uint8, 4096)
	base := memory
	memory = memory[1024:]

	memory = emit_ldc(memory, OP_LDCB, 0x12345678)
	memory = emit_ldc(memory, OP_LDCC, 0x23451209)
	memory = emit_op(memory, OP_SUB, REG_C, REG_B)
	memory = emit_op(memory, OP_HLT, 0, 0)

	Execute(regs, base)

	if regs[REG_C] != (0x23451209 - 0x12345678) {
		t.Fatalf("Expected %x but got %x1", (0x23451209 - 0x23451209), regs[REG_C])
	}
}

func TestMul(t *testing.T) {
	regs := make([]uint32, 0x10)
	memory := make([]uint8, 0x4096)
	base := memory
	memory = memory[1024:]

	memory = emit_ldc(memory, OP_LDCB, 0x1234)
	memory = emit_ldc(memory, OP_LDCC, 0x23)
	memory = emit_op(memory, OP_MUL, REG_C, REG_B)
	memory = emit_op(memory, OP_HLT, 0, 0)

	Execute(regs, base)

	exp := uint32(0x1234) * uint32(0x23)

	if regs[REG_C] != exp {
		t.Fatalf("Expected %x but got %x1", exp, regs[REG_C])
	}
}

func TestDiv(t *testing.T) {
	regs := make([]uint32, 0x10)
	memory := make([]uint8, 4096)
	base := memory
	memory = memory[1024:]

	memory = emit_ldc(memory, OP_LDCB, 0x1234)
	memory = emit_ldc(memory, OP_LDCC, 0x23)
	memory = emit_op(memory, OP_DIV, REG_B, REG_C)
	memory = emit_op(memory, OP_HLT, 0, 0)

	Execute(regs, base)

	exp := uint32(0x1234) / uint32(0x23)

	if regs[REG_B] != exp {
		t.Fatalf("Expected %x but got %x!", exp, regs[REG_B])
	}
}

func TestInc(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			"ldca 0x00",
			"inc ra",
			"inc ra",
			"hlt",
		})

	Execute(regs, ec.Memory())

	if regs[REG_A] != 2 {
		t.Fatalf("Expected %x but got %x!", 2, regs[REG_A])
	}
}

func TestInc2(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			"ldca 0xFFFFFFFF",
			"inc ra",
			"inc ra",
			"hlt",
		})

	Execute(regs, ec.Memory())

	if regs[REG_A] != 1 {
		t.Fatalf("Expected %x but got %x!", 1, regs[REG_A])
	}
}

func TestLra(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			"ldcc 0x01",
			".l test",
			".lra test",
			"hlt",
		})

	Execute(regs, ec.Memory())

	if regs[REG_A] != 0x405 {
		t.Fatalf("Expected %x but got %x!", 0x405, regs[REG_A])
	}
}

func TestLrb(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			"ldcc 0x01",
			".l test",
			".lrb test",
			"hlt",
		})

	Execute(regs, ec.Memory())

	if regs[REG_B] != 0x405 {
		t.Fatalf("Expected %x but got %x!", 0x405, regs[REG_B])
	}
}

func TestLrc(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			"ldcc 0x01",
			".l test",
			".lrc test",
			"hlt",
		})

	Execute(regs, ec.Memory())

	if regs[REG_C] != 0x405 {
		t.Fatalf("Expected %x but got %x!", 0x405, regs[REG_C])
	}
}

func TestDec(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			"ldca 0x04",
			"dec ra",
			"dec ra",
			"hlt",
		})

	Execute(regs, ec.Memory())

	if regs[REG_A] != 2 {
		t.Fatalf("Expected %x but got %x!", 2, regs[REG_A])
	}
}

func TestDec2(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			"ldca 0x01",
			"dec ra",
			"dec ra",
			"hlt",
		})

	Execute(regs, ec.Memory())

	if regs[REG_A] != 0xFFFFFFFF {
		t.Fatalf("Expected %x but got %x!", 0xFFFFFFFF, regs[REG_A])
	}
}

func TestAdra(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			".adra test",
			"nop",
			".l test",
			"hlt",
		})
	ec.Resolve()

	Execute(regs, ec.Memory())

	if regs[REG_A] != 0x407 {
		t.Fatalf("Expected %x but got %x!", 0x407, regs[REG_A])
	}
}

func TestAdrb(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			".adrb test",
			"nop",
			".l test",
			"hlt",
		})
	ec.Resolve()

	Execute(regs, ec.Memory())

	if regs[REG_B] != 0x407 {
		t.Fatalf("Expected %x but got %x!", 0x407, regs[REG_B])
	}
}

func TestAdrc(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			".adrc test",
			"nop",
			".l test",
			"hlt",
		})
	ec.Resolve()

	Execute(regs, ec.Memory())

	if regs[REG_C] != 0x407 {
		t.Fatalf("Expected %x but got %x!", 0x407, regs[REG_C])
	}
}

func TestJmp(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			".adrc test",
			"jmp rc",
			"hlt",
			".l test",
			"inc ra",
			"hlt",
		})
	ec.Resolve()
	
	Execute(regs, ec.Memory())

	if regs[REG_A] != 1 {
		t.Fatalf("Expected %x but got %x!", 1, regs[REG_A])
	}
}

func TestJnz(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			".adrc test",
			"ldca 0x01",
			"jnz rc ra",
			"ldcb 0x4321",
			"hlt",
			".l test",
			"ldcb 0x1234",
			"hlt",
		})
	ec.Resolve()
	
	Execute(regs, ec.Memory())

	if regs[REG_B] != 0x1234 {
		t.Fatalf("Expected %x but got %x!", 0x1234, regs[REG_B])
	}
}

func TestJiz(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			".adrc test",
			"ldca 0x01",
			"jiz rc ra",
			"ldcb 0x4321",
			"hlt",
			".l test",
			"ldcb 0x1234",
			"hlt",
		})
	ec.Resolve()
	
	Execute(regs, ec.Memory())

	if regs[REG_B] != 0x4321 {
		t.Fatalf("Expected %x but got %x!", 0x4321, regs[REG_B])
	}
}

func TestMov(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			"ldca 0x1234",
			"mov re ra",
			"hlt",
		})
	ec.Resolve()
	
	Execute(regs, ec.Memory())

	if regs[REG_A] != 0x1234 {
		t.Fatalf("Expected %x but got %x!", 0x1234, regs[REG_A])
	}
}

func TestPush(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			"ldca 2048",
			"ldcb 0x12345678",
			"push ra rb",
			"hlt",
		})
	ec.Resolve()
	
	Execute(regs, ec.Memory())

	if regs[REG_A] != 2052 {
		t.Fatalf("Expected %x but got %x!", 2052, regs[REG_A])
	}
	
	if ec.Memory()[2048] != 0x12 {
		t.Fatalf("Expected %x but got %x!", 0x12, ec.Memory()[2048])
	}
	
	if ec.Memory()[2049] != 0x34 {
		t.Fatalf("Expected %x but got %x!", 0x34, ec.Memory()[2049])
	}
	
	if ec.Memory()[2050] != 0x56 {
		t.Fatalf("Expected %x but got %x!", 0x56, ec.Memory()[2050])
	}
	
	if ec.Memory()[2051] != 0x78 {
		t.Fatalf("Expected %x but got %x!", 0x78, ec.Memory()[2051])
	}
}

func TestPop(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			"ldca 2052",
			"pop rb ra",
			"hlt",
		})
	ec.Resolve()
	
	mem := ec.Memory()
	mem[2051] = 0x78
	mem[2050] = 0x56
	mem[2049] = 0x34
	mem[2048] = 0x12
	
	Execute(regs, ec.Memory())
	
	if regs[REG_A] != 2048 {
		t.Fatalf("Expected %x but got %x!", 2048, regs[REG_A])
	}
	
	if regs[REG_B] != 0x12345678 {
		t.Fatalf("Expected %x but got %x!", 0x12345678, regs[REG_B])
	}
}

func TestShl(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			"ldca 0x1234",
			"ldcb 0x2",
			"shl ra rb",
			"hlt",
		})
	ec.Resolve()
	
	Execute(regs, ec.Memory())

	if regs[REG_A] != 0x48d0 {
		t.Fatalf("Expected %x but got %x!", 0x48d0, regs[REG_A])
	}
}

func TestShr(t *testing.T) {
	regs := make([]uint32, 0x10)
	ec := asm.NewDefaultEmitContext()

	asm.AsmLns(ec,
		[]string{
			"ldca 0x48d0",
			"ldcb 0x2",
			"shr ra rb",
			"hlt",
		})
	ec.Resolve()
	
	Execute(regs, ec.Memory())

	if regs[REG_A] != 0x1234 {
		t.Fatalf("Expected %x but got %x!", 0x1234, regs[REG_A])
	}
}