package defs

const REG_IP = 0x00
const REG_A = 0x01
const REG_B = 0x02
const REG_C = 0x03
const REG_D = 0x04
const REG_E = 0x05
const REG_SP = 0x0F

const OP_LDCA = 0x00
const OP_LDCB = 0x01
const OP_LDCC = 0x02
const OP_XCHG = 0x03
const OP_ADD = 0x04
const OP_SUB = 0x05
const OP_MUL = 0x06
const OP_DIV = 0x07
const OP_INC = 0x08
const OP_DEC = 0x09
const OP_MOV = 0x0A

const OP_JMP = 0x30
const OP_JNZ = 0x31

const OP_LDB = 0x40

const OP_FAIL = 0xFC
const OP_OUT = 0xFD
const OP_HLT = 0xFE
const OP_NOP = 0xFF