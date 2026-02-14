//go:build gc && s390x && !purego

#include "textflag.h"

TEXT ·_clmul(SB), NOSPLIT, $0-32
    MOVD x+0(FP), R1
    MOVD y+8(FP), R2
    MOVD $0, R3
    
    VLVGG $0, R1, V0
    VLVGG $1, R3, V0
    VLVGG $0, R2, V1
    VLVGG $1, R3, V1
    
    VGFMG V0, V1, V2
    
    VLGVG $1, V2, R3
    MOVD R3, lo+16(FP)
    
    VLGVG $0, V2, R3
    MOVD R3, hi+24(FP)
    
    RET
