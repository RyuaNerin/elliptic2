//go:build gc && arm64 && !purego

#include "textflag.h"

TEXT ·_clmul(SB), NOSPLIT, $0-32
    MOVD x+0(FP), R0
    MOVD y+8(FP), R1
    
    FMOV D0, X0
    FMOV D1, X1
    
    PMULL V2.1Q, V0.1D, V1.1D
    
    FMOV X2, D2
    MOVD R2, lo+16(FP)
    
    DUP D3, V2.D[1]
    FMOV X3, D3
    MOVD R3, hi+24(FP)
    
    RET
