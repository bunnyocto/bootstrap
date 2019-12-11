package cpu

import . "bootstrap/defs"
import "fmt"

func Execute(regs []uint32, memory []uint8) uint8 {
	// Fill in values
	regs[REG_IP] = 1024
	
	// CPUID
	memory[20] = uint8(0x2) << 2 | 0x000
	memory[21] = 0x00
	memory[22] = 0x10
	memory[23] = 0x00
	
	//MSI
	memsz := uint32(len(memory))
	memory[24] = uint8((memsz >> 24) & 0xFF)
	memory[25] = uint8((memsz >> 16) & 0xFF)
	memory[26] = uint8((memsz >> 8) & 0xFF)
	memory[27] = uint8((memsz >> 0) & 0xFF)

	for {
		if regs[REG_IP] >= uint32(len(memory)) {
			return 5
		}

		i := memory[regs[REG_IP]]
		regs[REG_IP]++

		if i <= OP_LDCC {
			// fetch uint32
			if regs[REG_IP]+3 >= uint32(len(memory)) {
				return 4
			}
			c0 := memory[regs[REG_IP]]
			regs[REG_IP]++

			c1 := memory[regs[REG_IP]]
			regs[REG_IP]++

			c2 := memory[regs[REG_IP]]
			regs[REG_IP]++

			c3 := memory[regs[REG_IP]]
			regs[REG_IP]++

			value := uint32(c0)<<24 | uint32(c1)<<16 | uint32(c2)<<8 | uint32(c3)

			switch i {
			case OP_LDCA:
				regs[REG_A] = value
			case OP_LDCB:
				regs[REG_B] = value
			case OP_LDCC:
				regs[REG_C] = value
			default:
				return 3
			}
		} else {
			if regs[REG_IP] >= uint32(len(memory)) {
				return 2
			}

			b1 := memory[regs[REG_IP]]
			regs[REG_IP]++

			src := b1 & 0x0F
			dst := (b1 >> 4) & 0x0F

			switch i {
			case OP_NOP:
				/* nop */
			case OP_XCHG:
				t := regs[dst]
				regs[dst] = regs[src]
				regs[src] = t
			case OP_ADD:
				regs[dst] += regs[src]
			case OP_SUB:
				regs[dst] -= regs[src]
			case OP_MUL:
				regs[dst] *= regs[src]
			case OP_DIV:
				regs[dst] /= regs[src]
			case OP_INC:
				regs[dst]++
			case OP_DEC:
				regs[dst]--
			case OP_OUT:
				if regs[dst] == 0x12345678 {
					fmt.Printf("%c", regs[src]&0xFF)
				}
			case OP_JMP:
				regs[REG_IP] = regs[dst]
			case OP_JNZ:
				if regs[src] != 0 {
					regs[REG_IP] = regs[dst]
				}
			case OP_JIZ:
				if regs[src] == 0 {
					regs[REG_IP] = regs[dst]
				}
			case OP_LDB:
				regs[dst] = uint32(memory[regs[src]])
			case OP_STB:
				memory[regs[dst]] = uint8(regs[src] & 0xFF)
			case OP_MOV:
				regs[dst] = regs[src]
			case OP_HLT:
				return 0
			case OP_FAIL:
				return 255
			case OP_CALL:
				b0 := uint8((regs[REG_IP] >> 24) & 0xFF)
				b1 := uint8((regs[REG_IP] >> 16) & 0xFF)
				b2 := uint8((regs[REG_IP] >>  8) & 0xFF)
				b3 := uint8((regs[REG_IP] >>  0) & 0xFF)
				memory[regs[src]+0] = b0
				memory[regs[src]+1] = b1
				memory[regs[src]+2] = b2
				memory[regs[src]+3] = b3
				regs[src]+=4
				regs[REG_IP] = regs[dst]
			case OP_RET:
				regs[dst] -= 4
				b0 := memory[regs[dst]+0]
				b1 := memory[regs[dst]+1]
				b2 := memory[regs[dst]+2]
				b3 := memory[regs[dst]+3]
				regs[REG_IP] = uint32(b0) << 24 | uint32(b1) << 16 | uint32(b2) << 8 | uint32(b3)
			case OP_PUSH:
				b0 := uint8((regs[src] >> 24) & 0xFF)
				b1 := uint8((regs[src] >> 16) & 0xFF)
				b2 := uint8((regs[src] >>  8) & 0xFF)
				b3 := uint8((regs[src] >>  0) & 0xFF)
				memory[regs[dst]+0] = b0
				memory[regs[dst]+1] = b1
				memory[regs[dst]+2] = b2
				memory[regs[dst]+3] = b3
				regs[dst] += 4
			case OP_POP:
				regs[src] -= 4
				b0 := memory[regs[src]+0]
				b1 := memory[regs[src]+1]
				b2 := memory[regs[src]+2]
				b3 := memory[regs[src]+3]
				regs[dst] = uint32(b0) << 24 | uint32(b1) << 16 | uint32(b2) << 8 | uint32(b3)
			case OP_SHL:
				regs[dst] <<= regs[src]
			case OP_SHR:
				regs[dst] >>= regs[src]
			default:
				return 1
			}
		}
	}
}
