//go:build gc && (ppc64 || ppc64le) && !purego

#include "textflag.h"

TEXT ·_clmul(SB), NOSPLIT, $0-32
    MOVD x+0(FP), R3
    MOVD y+8(FP), R4
    
    MOVD R3, R5
    MTVSRD V0, R5
    
    MOVD R4, R5
    MTVSRD V1, R5
    
    VPMSUMD V2, V0, V1
    WORD $0x10400C48
    
    MFVSRD R5, V2
    MOVD R5, lo+16(FP)
    
    VSLDOI V3, V2, V2, $8
    MFVSRD R5, V3
    MOVD R5, hi+24(FP)
    
    RET
