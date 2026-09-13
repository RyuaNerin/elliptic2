from gogen import GoGenerator

############################################################


def add2007bl(X1, Y1, Z1, X2, Y2, Z2):
    """
    https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian.html#addition-add-2007-bl

    The "add-2007-bl" addition formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Cost: 11M + 5S + 9add + 4*2.
    - Cost: 10M + 4S + 9add + 4*2 dependent upon the first point.
    - Source: 2007 Bernstein-Lange; note that the improvement from 12M+4S to 11M+5S was already mentioned in 2001 Bernstein http://cr.yp.to/talks.html#2001.10.29.
    """

    Z1Z1 = Z1 ^ 2
    Z2Z2 = Z2 ^ 2
    U1 = X1 * Z2Z2
    U2 = X2 * Z1Z1
    t0 = Z2 * Z2Z2
    S1 = Y1 * t0
    t1 = Z1 * Z1Z1
    S2 = Y2 * t1
    H = U2 - U1
    t2 = 2 * H
    I = t2 ^ 2
    J = H * I
    t3 = S2 - S1
    r = 2 * t3
    V = U1 * I
    t4 = r ^ 2
    t5 = 2 * V
    t6 = t4 - J
    X3 = t6 - t5
    t7 = V - X3
    t8 = S1 * J
    t9 = 2 * t8
    t10 = r * t7
    Y3 = t10 - t9
    t11 = Z1 + Z2
    t12 = t11 ^ 2
    t13 = t12 - Z1Z1
    t14 = t13 - Z2Z2
    Z3 = t14 * H

    return X3, Y3, Z3


def dbl2007bl(a, X1, Y1, Z1):
    """
    https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian.html#doubling-dbl-2007-bl

    The "dbl-2007-bl" doubling formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Cost: 1M + 8S + 1*a + 10add + 2*2 + 1*3 + 1*8.
    - Source: 2007 Bernstein-Lange.
    """

    XX = X1 ^ 2
    YY = Y1 ^ 2
    YYYY = YY ^ 2
    ZZ = Z1 ^ 2
    t0 = X1 + YY
    t1 = t0 ^ 2
    t2 = t1 - XX
    t3 = t2 - YYYY
    S = 2 * t3
    t4 = ZZ ^ 2
    t5 = a * t4
    t6 = 3 * XX
    M = t6 + t5
    t7 = M ^ 2
    t8 = 2 * S
    T = t7 - t8
    X3 = T
    t9 = S - T
    t10 = 8 * YYYY
    t11 = M * t9
    Y3 = t11 - t10
    t12 = Y1 + Z1
    t13 = t12 ^ 2
    t14 = t13 - YY
    Z3 = t14 - ZZ

    return X3, Y3, Z3


def madd2007bl(X1, Y1, Z1, X2, Y2):
    """
    https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian.html#addition-madd-2007-bl

    The "madd-2007-bl" addition formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Assumptions: Z2=1.
    - Cost: 7M + 4S + 9add + 3*2 + 1*4.
    - Source: 2007 Bernstein-Lange.
    """

    Z1Z1 = Z1 ^ 2
    U2 = X2 * Z1Z1
    t0 = Z1 * Z1Z1
    S2 = Y2 * t0
    H = U2 - X1
    HH = H ^ 2
    I = 4 * HH
    J = H * I
    t1 = S2 - Y1
    r = 2 * t1
    V = X1 * I
    t2 = r ^ 2
    t3 = 2 * V
    t4 = t2 - J
    X3 = t4 - t3
    t5 = V - X3
    t6 = Y1 * J
    t7 = 2 * t6
    t8 = r * t5
    Y3 = t8 - t7
    t9 = Z1 + H
    t10 = t9 ^ 2
    t11 = t10 - Z1Z1
    Z3 = t11 - HH

    return X3, Y3, Z3


# a4 = -3
def dbl2001b(X1, Y1, Z1):
    """
    https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#doubling-dbl-2001-b

    The "dbl-2001-b" doubling formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Cost: 3M + 5S + 8add + 1*3 + 1*4 + 2*8.
    - Source: 2001 Bernstein "A software implementation of NIST P-224".
    """

    delta = Z1 ^ 2
    gamma = Y1 ^ 2
    beta = X1 * gamma
    t0 = X1 - delta
    t1 = X1 + delta
    t2 = t0 * t1
    alpha = 3 * t2
    t3 = alpha ^ 2
    t4 = 8 * beta
    X3 = t3 - t4
    t5 = Y1 + Z1
    t6 = t5 ^ 2
    t7 = t6 - gamma
    Z3 = t7 - delta
    t8 = 4 * beta
    t9 = t8 - X3
    t10 = gamma ^ 2
    t11 = 8 * t10
    t12 = alpha * t9
    Y3 = t12 - t11

    return X3, Y3, Z3


# a4 = 0
def dbl2009l(X1, Y1, Z1):
    """
    https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l

    The "dbl-2009-l" doubling formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Cost: 2M + 5S + 6add + 3*2 + 1*3 + 1*8.
    - Source: 2009.04.01 Lange.
    """

    A = X1 ^ 2
    B = Y1 ^ 2
    C = B ^ 2
    t0 = X1 + B
    t1 = t0 ^ 2
    t2 = t1 - A
    t3 = t2 - C
    D = 2 * t3
    E = 3 * A
    F = E ^ 2
    t4 = 2 * D
    X3 = F - t4
    t5 = D - X3
    t6 = 8 * C
    t7 = E * t5
    Y3 = t7 - t6
    t8 = Y1 * Z1
    Z3 = 2 * t8

    return X3, Y3, Z3


# a4 = 0
def dbl2009lv2(X1, Y1, Z1):
    """
    https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l

    The "dbl-2009-l" doubling formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Cost: 2M + 5S + 6add + 3*2 + 1*3 + 1*8.
    - Source: 2009.04.01 Lange.
    """

    A = X1**2
    B = Y1**2
    C = B**2
    D = 2 * ((X1 + B) ** 2 - A - C)
    E = 3 * A
    F = E**2
    X3 = F - 2 * D
    Y3 = E * (D - X3) - 8 * C
    Z3 = 2 * Y1 * Z1

    return X3, Y3, Z3


############################################################


def is_on_curve(a, b, X, Y):
    lhs = Y**2
    rhs = X**3 + a * X + b

    return lhs, rhs


def compute_y(a, b, X):
    """
    yy = xxx + ax + b
    y = sqrt(xxx + ax + b)
    """

    y1 = (X**3 + a * X + b) ** 0.5
    return y1


############################################################


def main():
    g = GoGenerator(package_name="weierstrassprime", gf="GFp")

    g.gen_projection_fn(add2007bl, projection="jacobian")
    g.gen_projection_fn(dbl2007bl, projection="jacobian")
    g.gen_projection_fn(madd2007bl, projection="jacobian")
    g.gen_projection_fn(dbl2001b, projection="jacobian")
    g.gen_projection_fn(dbl2009l, projection="jacobian")

    g.gen_projection_fn(dbl2009lv2, projection="jacobian")

    # g.gen_projection_fn(is_on_curve)
    # g.gen_projection_fn(compute_y)


if __name__ == "__main__":
    main()
