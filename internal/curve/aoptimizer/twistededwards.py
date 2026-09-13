from gogen import GoGenerator

############################################################


def add2008hwcd(a, X1, Y1, Z1, T1, X2, Y2, Z2, T2, d):
    """
    https://hyperelliptic.org/EFD/g1p/auto-twisted-extended.html#addition-add-2008-hwcd

    The "add-2008-hwcd" addition formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Cost: 9M + 1*a + 1*d + 7add.
    - Cost: 9M + 1*a + 6add dependent upon the first point.
    - Source: 2008 Hisil-Wong-Carter-Dawson, http://eprint.iacr.org/2008/522, Section 3.1.
    - Strongly unified.
    """

    A = X1 * X2
    B = Y1 * Y2
    t0 = d * T2
    C = T1 * t0
    D = Z1 * Z2
    t1 = X1 + Y1
    t2 = X2 + Y2
    t3 = t1 * t2
    t4 = t3 - A
    E = t4 - B
    F = D - C
    G = D + C
    t5 = a * A
    H = B - t5
    X3 = E * F
    Y3 = G * H
    T3 = E * H
    Z3 = F * G

    return X3, Y3, Z3, T3


def dbl2008hwcd(a, X1, Y1, Z1):
    """
    https://hyperelliptic.org/EFD/g1p/auto-twisted-extended.html#doubling-dbl-2008-hwcd

    The "dbl-2008-hwcd" doubling formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Cost: 4M + 4S + 1*a + 6add + 1*2.
    - Source: 2008 Hisil-Wong-Carter-Dawson, http://eprint.iacr.org/2008/522, Section 3.3.
    """

    A = X1 ^ 2
    B = Y1 ^ 2
    t0 = Z1 ^ 2
    C = 2 * t0
    D = a * A
    t1 = X1 + Y1
    t2 = t1 ^ 2
    t3 = t2 - A
    E = t3 - B
    G = D + B
    F = G - C
    H = D - B
    X3 = E * F
    Y3 = G * H
    T3 = E * H
    Z3 = F * G

    return X3, Y3, Z3, T3


# Z=-1
def madd2008hwcd(a, X1, Y1, Z1, T1, X2, Y2, T2, d):
    """
    https://hyperelliptic.org/EFD/g1p/auto-twisted-extended.html#addition-madd-2008-hwcd

    The "madd-2008-hwcd" addition formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Assumptions: Z2=1.
    - Cost: 8M + 1*a + 1*d + 7add.
    - Cost: 8M + 1*a + 6add dependent upon the first point.
    - Source: 2008 Hisil-Wong-Carter-Dawson, http://eprint.iacr.org/2008/522, Section 3.1.
    - Strongly unified.
    """
    A = X1 * X2
    B = Y1 * Y2
    t0 = d * T2
    C = T1 * t0
    D = Z1
    t1 = X1 + Y1
    t2 = X2 + Y2
    t3 = t1 * t2
    t4 = t3 - A
    E = t4 - B
    F = D - C
    G = D + C
    t5 = a * A
    H = B - t5
    X3 = E * F
    Y3 = G * H
    T3 = E * H
    Z3 = F * G

    return X3, Y3, Z3, T3


############################################################


# a=-1, k=2d
def add2008hwcd3(X1, Y1, Z1, T1, X2, Y2, Z2, T2, k):
    """
    https://hyperelliptic.org/EFD/g1p/auto-twisted-extended-1.html#addition-add-2008-hwcd-3

    The "add-2008-hwcd-3" addition formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Assumptions: k=2*d.
    - Cost: 8M + 1*k + 8add + 1*2.
    - Cost: 8M + 6add dependent upon the first point.
    - Source: 2008 Hisil-Wong-Carter-Dawson, http://eprint.iacr.org/2008/522, Section 3.1.
    - Strongly unified.
    """

    t0 = Y1 - X1
    t1 = Y2 - X2
    A = t0 * t1
    t2 = Y1 + X1
    t3 = Y2 + X2
    B = t2 * t3
    t4 = k * T2
    C = T1 * t4
    t5 = 2 * Z2
    D = Z1 * t5
    E = B - A
    F = D - C
    G = D + C
    H = B + A
    X3 = E * F
    Y3 = G * H
    T3 = E * H
    Z3 = F * G

    return X3, Y3, Z3, T3


# a=-1, k=2d, A2=1
def madd2008hwcd3(X1, Y1, Z1, T1, X2, Y2, T2, k):
    """
    https://hyperelliptic.org/EFD/g1p/auto-twisted-extended-1.html#addition-madd-2008-hwcd-3

    The "madd-2008-hwcd-3" addition formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Assumptions: Z2=1 and k=2*d.
    - Cost: 7M + 1*k + 8add + 1*2.
    - Cost: 7M + 6add + 1*2 dependent upon the first point.
    - Source: 2008 Hisil-Wong-Carter-Dawson, http://eprint.iacr.org/2008/522, Section 3.1.
    - Strongly unified.
    """
    t0 = Y1 - X1
    t1 = Y2 - X2
    A = t0 * t1
    t2 = Y1 + X1
    t3 = Y2 + X2
    B = t2 * t3
    t4 = k * T2
    C = T1 * t4
    D = 2 * Z1
    E = B - A
    F = D - C
    G = D + C
    H = B + A
    X3 = E * F
    Y3 = G * H
    T3 = E * H
    Z3 = F * G

    return X3, Y3, Z3, T3


############################################################


def is_on_curve(a, b, X, Y):
    XX = X**2
    YY = Y**2
    lhs = a * XX + YY
    rhs = 1 + b * XX * YY

    return lhs, rhs


def compute_y(a, b, X):
    """
    axx + yy = 1 + dxxyy

    yy - dxxyy = 1 - axx
    yy(1 - dxx) = 1 - axx
    yy = (1 - axx) / (1 - dxx)
    y = sqrt( (1 - axx) / (1 - dxx) )
    """

    XX = X**2
    u = 1 - a * XX
    v = 1 - b * XX
    yy = u / v
    yy = yy**0.5

    return yy


############################################################


def main():
    g = GoGenerator(package_name="twistededwards", gf="GFp")

    g.gen_projection_fn(add2008hwcd, projection="extended")
    g.gen_projection_fn(dbl2008hwcd, projection="extended")
    g.gen_projection_fn(madd2008hwcd, projection="extended")

    g.gen_projection_fn(add2008hwcd3, projection="extended")
    g.gen_projection_fn(madd2008hwcd3, projection="extended")

    # g.gen_projection_fn(is_on_curve)
    # g.gen_projection_fn(compute_y)


if __name__ == "__main__":
    main()
