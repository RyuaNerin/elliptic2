//go:build gc && 386 && !purego && go1.18

#include "textflag.h"

TEXT ·_clmul32(SB), NOSPLIT, $0-16
    // Load a into low 32 bits of X0
    MOVL a+0(FP), AX
    MOVD AX, X0
    
    // Load b into low 32 bits of X1
    MOVL b+4(FP), AX
    MOVD AX, X1
    
    // PCLMULQDQ $0x00, X1, X0
    // Performs carry-less multiplication of low 64 bits
    // Since we only loaded 32 bits, this gives us 32x32->64
    BYTE $0x66; BYTE $0x0F; BYTE $0x3A; BYTE $0x44; BYTE $0xC1; BYTE $0x00
    
    // Extract low 32 bits
    MOVD X0, AX
    MOVL AX, lo+8(FP)
    
    // Extract high 32 bits
    PSRLDQ $4, X0
    MOVD X0, AX
    MOVL AX, hi+12(FP)
    
    RET