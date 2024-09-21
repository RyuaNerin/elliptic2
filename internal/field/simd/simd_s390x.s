//go:build gc && s390x && !purego && go1.18

#include "textflag.h"

TEXT ·_clmul(SB), NOSPLIT, $0-32
    MOVD x+0(FP), R1
    MOVD y+8(FP), R2
    
    VLVGG $0, R1, V0
    VLVGG $0, R2, V1
    VZERO V2
    
    // VGFMAG V3, V0, V1, V2
    WORD $0xE700100000BC
    
    VLGVG $1, V3, R3
    MOVD R3, lo+16(FP)
    
    VLGVG $0, V3, R3
    MOVD R3, hi+24(FP)
    
    RET
