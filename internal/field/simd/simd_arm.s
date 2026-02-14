//go:build gc && arm && !purego

#include "textflag.h"

TEXT ·_clmul32(SB), NOSPLIT, $0-16
    MOVW a+0(FP), R0        // R0 = a
    MOVW b+4(FP), R1        // R1 = b
    
    MOVW $0, R2             // lo = 0
    MOVW $0, R3             // hi = 0
    MOVW $0, R4             // bit index i = 0

loop:
    CMP $32, R4
    BGE done
    
    // if (b >> i) & 1 != 0
    MOVW R1, R5
    SRL R4, R5              // R5 = b >> i
    AND $1, R5
    CMP $0, R5
    BEQ next
    
    // lo ^= a << i
    MOVW R0, R5
    SLL R4, R5              // R5 = a << i
    EOR R5, R2              // lo ^= (a << i)
    
    // if i > 0: hi ^= a >> (32 - i)
    CMP $0, R4
    BEQ next
    RSB $32, R4, R6         // R6 = 32 - i
    MOVW R0, R5
    SRL R6, R5              // R5 = a >> (32 - i)
    EOR R5, R3              // hi ^= (a >> (32 - i))

next:
    ADD $1, R4
    B loop

done:
    MOVW R2, ret+8(FP)
    MOVW R3, ret_hi+12(FP)
    
    RET
