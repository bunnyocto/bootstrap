package defs

const REG_IP = 0x00
const REG_A = 0x01
const REG_B = 0x02
const REG_C = 0x03
const REG_D = 0x04
const REG_E = 0x05
const REG_F = 0x06
const REG_G = 0x07
const REG_H = 0x08
const REG_I = 0x09
const REG_J = 0x0A
const REG_K = 0x0B
const REG_L = 0x0C
const REG_M = 0x0D
const REG_N = 0x0E
const REG_O = 0x0F

const OP_LDCA = 0x00
const OP_LDCB = 0x01
const OP_LDCC = 0x02

const OP_ADD = 0x04
const OP_SUB = 0x05
const OP_MUL = 0x06
const OP_DIV = 0x07
const OP_INC = 0x08
const OP_DEC = 0x09
const OP_MOV = 0x0A
const OP_MOD = 0x0B
const OP_SADD = 0x0C
const OP_SSUB = 0x0D
const OP_SMUL = 0x0E
const OP_SDIV = 0x0F
const OP_XCHG = 0x10
const OP_SHR = 0x11
const OP_SHL = 0x12
const OP_BYT = 0x13

const OP_JMP = 0x30
const OP_JNZ = 0x31
const OP_JIZ = 0x32
const OP_CALL = 0x33
const OP_RET = 0x34
const OP_JNE = 0x35
const OP_JEQ = 0x36
const OP_JLT = 0x37
const OP_JGT = 0x38
const OP_JBL = 0x39
const OP_JAB = 0x3A

const OP_LDB = 0x40
const OP_STB = 0x41
const OP_LDH = 0x42
const OP_STH = 0x43
const OP_LDW = 0x44
const OP_STW = 0x45
const OP_PUSH = 0x46
const OP_POP = 0x47

const OP_LDP = 0xE0
const OP_LDN = 0xE1

const OP_FAIL = 0xFC
const OP_OUT = 0xFD
const OP_HLT = 0xFE
const OP_NOP = 0xFF
