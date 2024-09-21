//go:build gc && riscv64 && !purego && go1.18

#include "textflag.h"

TEXT ·_clmul(SB), NOSPLIT, $0-32
    MOV x+0(FP), A0
    MOV y+8(FP), A1
    
    // clmul A2, A0, A1 (low 64 bits)
    // funct7=0000101, rs2=A1, rs1=A0, funct3=001, rd=A2, opcode=0110011
    WORD $0x0AB51633
    
    // clmulh A3, A0, A1 (high 64 bits)
    // funct7=0000101, funct3=011
    WORD $0x0AB536B3
    
    MOV A2, lo+16(FP)
    MOV A3, hi+24(FP)
    
    RET
