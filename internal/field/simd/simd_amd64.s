//go:build gc && amd64 && !purego

#include "textflag.h"

TEXT ·_clmulWords(SB), NOSPLIT, $0-72
    MOVQ z_base+0(FP), DI
    MOVQ z_len+8(FP), R8
    MOVQ x_base+24(FP), SI
    MOVQ x_len+32(FP), R9
    MOVQ y_base+48(FP), DX
    MOVQ y_len+56(FP), R10
    
    TESTQ R9, R9
    JZ done
    TESTQ R10, R10
    JZ done
    
    XORQ R11, R11

outer_loop:
    CMPQ R11, R10
    JGE done
    
    MOVQ (DX)(R11*8), X1
    MOVQ X1, AX
    TESTQ AX, AX
    JZ next_y
    
    XORQ R12, R12

inner_loop:
    CMPQ R12, R9
    JGE next_y
    
    MOVQ (SI)(R12*8), X0
    
    PCLMULQDQ $0x00, X1, X0
    
    LEAQ (R11)(R12*1), R13
    
    MOVQ X0, AX
    XORQ AX, (DI)(R13*8)
    
    PSRLDQ $8, X0
    MOVQ X0, AX
    INCQ R13
    XORQ AX, (DI)(R13*8)
    
    INCQ R12
    JMP inner_loop

next_y:
    INCQ R11
    JMP outer_loop

done:
    RET

TEXT ·_clmul(SB), NOSPLIT, $0-32
    MOVQ a+0(FP), AX
    MOVQ b+8(FP), BX
    
    MOVQ AX, X0
    MOVQ BX, X1
    
    PCLMULQDQ $0x00, X1, X0
    
    MOVQ X0, lo+16(FP)
    PSRLDQ $8, X0
    MOVQ X0, hi+24(FP)
    
    RET

TEXT ·_expandBitsBMI2(SB), NOSPLIT, $0-24
    MOVQ x+0(FP), AX
    MOVQ $0x5555555555555555, CX
    
    MOVL AX, DX
    // PDEP BX, DX, CX
    BYTE $0xC4; BYTE $0xE2; BYTE $0xEB; BYTE $0xF5; BYTE $0xD9
    MOVQ BX, lo+8(FP)
    
    SHRQ $32, AX
    MOVL AX, DX
    // PDEP BX, DX, CX
    BYTE $0xC4; BYTE $0xE2; BYTE $0xEB; BYTE $0xF5; BYTE $0xD9
    MOVQ BX, hi+16(FP)
    
    RET

TEXT ·_compressBitsBMI2(SB), NOSPLIT, $0-24
    MOVQ x+0(FP), AX
    MOVQ $0x5555555555555555, CX
    
    // PEXTQ BX, AX, CX
    BYTE $0xC4; BYTE $0xE2; BYTE $0xFA; BYTE $0xF5; BYTE $0xD9
    MOVQ BX, even+8(FP)
    
    SHRQ $1, AX    
    // PEXTQ BX, AX, CX
    BYTE $0xC4; BYTE $0xE2; BYTE $0xFA; BYTE $0xF5; BYTE $0xD9
    MOVQ BX, odd+16(FP)

    RET
