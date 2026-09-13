from gogen import GoGenerator

############################################################


def add2007bl2(c, d, X1, Y1, Z1, X2, Y2, Z2):
    """
    https://hyperelliptic.org/EFD/g1p/auto-edwards-projective.html#addition-add-2007-bl-2

    The "add-2007-bl-2" addition formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Cost: 10M + 1S + 1*c + 1*d + 7add.
    - Cost: 10M + 1S + 1*c + 1*d + 6add dependent upon the first point.
    - Source: 2007 Bernstein-Lange.
    - Strongly unified.
    """

    R1 = X1
    R2 = Y1
    R3 = Z1
    R4 = X2
    R5 = Y2
    R6 = Z2
    R3 = R3 * R6
    R7 = R1 + R2
    R8 = R4 + R5
    R1 = R1 * R4
    R2 = R2 * R5
    R7 = R7 * R8
    R7 = R7 - R1
    R7 = R7 - R2
    R7 = R7 * R3
    R8 = R1 * R2
    R8 = d * R8
    R2 = R2 - R1
    R2 = R2 * R3
    R3 = R3 ^ 2
    R1 = R3 - R8
    R3 = R3 + R8
    R2 = R2 * R3
    R3 = R3 * R1
    R1 = R1 * R7
    R3 = c * R3
    X3 = R1
    Y3 = R2
    Z3 = R3

    return X3, Y3, Z3


def dbl2007bl2(c, X1, Y1, Z1):
    """
    https://hyperelliptic.org/EFD/g1p/auto-edwards-projective.html#doubling-dbl-2007-bl-2

    The "dbl-2007-bl-2" doubling formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Cost: 3M + 4S + 3*c + 5add + 1*2.
    - Source: 2007 Bernstein-Lange; source comments that these formulas use two temporary registers.
    """

    R1 = X1
    R2 = Y1
    R3 = Z1
    R4 = R1 + R2
    R3 = c * R3
    R1 = R1 ^ 2
    R2 = R2 ^ 2
    R3 = R3 ^ 2
    R4 = R4 ^ 2
    R3 = 2 * R3
    R5 = R1 + R2
    R2 = R1 - R2
    R4 = R4 - R5
    R3 = R5 - R3
    R1 = R3 * R4
    R3 = R3 * R5
    R2 = R2 * R5
    R1 = c * R1
    R2 = c * R2
    X3 = R1
    Y3 = R2
    Z3 = R3

    return X3, Y3, Z3


def madd2007bl2(c, d, X1, Y1, Z1, X2, Y2):
    """
    https://hyperelliptic.org/EFD/g1p/auto-edwards-projective.html#addition-madd-2007-bl-2

    The "madd-2007-bl-2" addition formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Assumptions: Z2=1.
    - Cost: 9M + 1S + 1*c + 1*d + 7add.
    - Cost: 9M + 1S + 1*c + 1*d + 6add dependent upon the first point.
    - Source: 2007 Bernstein-Lange.
    - Strongly unified.
    """

    R1 = X1
    R2 = Y1
    R3 = Z1
    R4 = X2
    R5 = Y2
    R7 = R1 + R2
    R6 = R4 + R5
    R1 = R1 * R4
    R2 = R2 * R5
    R7 = R7 * R6
    R7 = R7 - R1
    R7 = R7 - R2
    R7 = R7 * R3
    R6 = R1 * R2
    R6 = d * R6
    R2 = R2 - R1
    R2 = R2 * R3
    R3 = R3 ^ 2
    R1 = R3 - R6
    R3 = R3 + R6
    R2 = R2 * R3
    R3 = R3 * R1
    R1 = R1 * R7
    R3 = c * R3
    X3 = R1
    Y3 = R2
    Z3 = R3

    return X3, Y3, Z3


############################################################


def is_on_curve(c, d, X1, Y1):
    # xx + yy = cc * (1 + d * xx * yy)
    XX = X1**2
    YY = Y1**2
    lhs = XX + YY
    rhs = c * c * (1 + d * XX * YY)

    return lhs, rhs


def compute_y(c, d, X1):
    # y = sqrt( (cc - xx) / (1 - ccdxx) )
    CC = c * c
    XX = X1 * X1

    y1 = ((CC - XX) / (1 - CC * d * XX)) ** 0.5
    return y1


############################################################


def main():
    g = GoGenerator(package_name="edwards", gf="GFp")

    g.gen_projection_fn(add2007bl2, projection="projective")
    g.gen_projection_fn(dbl2007bl2, projection="projective")
    g.gen_projection_fn(madd2007bl2, projection="projective")

    # g.gen_projection_fn(is_on_curve)
    # g.gen_projection_fn(compute_y)


if __name__ == "__main__":
    main()
