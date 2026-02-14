//go:build gc && 386 && !purego

#include "textflag.h"

TEXT ·_clmul32(SB), NOSPLIT, $0-16
    // Load a into low 32 bits of X0
    MOVL a+0(FP), AX
    MOVL AX, X0
    
    // Load b into low 32 bits of X1
    MOVL b+4(FP), AX
    MOVL AX, X1
    
    PCLMULQDQ $0x00, X1, X0
    
    // Extract low 32 bits
    MOVL X0, AX
    MOVL AX, lo+8(FP)
    
    // Extract high 32 bits
    PSRLDQ $4, X0
    MOVL X0, AX
    MOVL AX, hi+12(FP)
    
    RET
