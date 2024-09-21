//go:build gc && arm64 && !purego && go1.18

#include "textflag.h"

TEXT ·_clmul(SB), NOSPLIT, $0-32
    MOVD x+0(FP), R0
    MOVD y+8(FP), R1
    
    // FMOV D0, R0
    WORD $0x9E670000
    // FMOV D1, R1
    WORD $0x9E670021
    
    // PMULL V2.1Q, V0.1D, V1.1D
    WORD $0x0EE0E002
    
    // FMOV R2, D2
    WORD $0x9E660042
    MOVD R2, lo+16(FP)
    
    // MOV V3.D[0], V2.D[1]
    WORD $0x6E180443
    // FMOV R3, D3
    WORD $0x9E660063
    MOVD R3, hi+24(FP)
    
    RET
