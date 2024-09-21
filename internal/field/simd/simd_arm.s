//go:build gc && arm && !purego && go1.18

#include "textflag.h"

TEXT ·_clmul32(SB), NOSPLIT, $0-16
    MOVW a+0(FP), R0
    MOVW b+4(FP), R1
    
    // R0 → D0, R1 → D1
    WORD $0xee000b10      // VMOV D0, R0, R0
    WORD $0xee001b10      // VMOV D1, R1, R1
    
    // VMULL.P8 Q0, D0, D1
    WORD $0xf2800e01
    
    // Extract result
    WORD $0xec510b10      // VMOV R0, R1, D0
    MOVW R0, ret+8(FP)
    MOVW R1, ret_hi+12(FP)
    
    RET
