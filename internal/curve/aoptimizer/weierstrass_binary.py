from gogen import GoGenerator

############################################################


def add2005dl(a2, X1, Y1, Z1, X2, Y2, Z2):
    """
    https://hyperelliptic.org/EFD/g12o/auto-shortw-jacobian.html#addition-add-2005-dl

    The "add-2005-dl" addition formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Cost: 14M + 5S + 1*a2 + 7add.
    - Cost: 13M + 4S + 1*a2 + 7add dependent upon the first point.
    - Source: 2005 Doche-Lange.
    """

    O1 = Z1 ^ 2
    O2 = Z2 ^ 2
    A = X1 * O2
    B = X2 * O1
    t0 = O2 * Z2
    C = Y1 * t0
    t1 = O1 * Z1
    D = Y2 * t1
    E = A + B
    F = C + D
    G = E * Z1
    t2 = G * Y2
    t3 = F * X2
    H = t3 + t2
    Z3 = G * Z2
    I = F + Z3
    t4 = Z3 ^ 2
    t5 = E ^ 2
    t6 = E * t5
    t7 = F * I
    t8 = a2 * t4
    t9 = t8 + t7
    X3 = t9 + t6
    t10 = G ^ 2
    t11 = t10 * H
    t12 = I * X3
    Y3 = t12 + t11

    return X3, Y3, Z3


def dbl2005dl(a6, X1, Y1, Z1):
    """
    https://hyperelliptic.org/EFD/g12o/auto-shortw-jacobian.html#doubling-dbl-2005-dl

    The "dbl-2005-dl" doubling formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Cost: 4M + 5S + 1*a6 + 4add.
    - Source: 2005 Doche-Lange.
    - Explicit formulas:
    """

    A = X1 ^ 2
    B = A ^ 2
    C = Z1 ^ 2
    D = C ^ 2
    t0 = D ^ 2
    t1 = a6 * t0
    X3 = B + t1
    Z3 = X1 * C
    t2 = Y1 * Z1
    t3 = A + t2
    t4 = t3 + Z3
    t5 = t4 * X3
    t6 = B * Z3
    Y3 = t6 + t5

    return X3, Y3, Z3


def madd2008bl(a2, X1, Y1, Z1, X2, Y2):
    """
    https://hyperelliptic.org/EFD/g12o/auto-shortw-jacobian.html#addition-madd-2008-bl

    The "madd-2008-bl" addition formulas [database entry; Sage verification script; Sage output; three-operand code]:
    - Assumptions: Z2=1.
    - Cost: 10M + 3S + 1*a2 + 7add.
    - Source: 2005 Doche-Lange.
    """

    O1 = Z1**2
    B = X2 * O1
    D = Y2 * O1 * Z1
    E = X1 + B
    F = Y1 + D
    Z3 = E * Z1
    H = F * X2 + Z3 * Y2
    I = F + Z3
    G = Z3**2
    X3 = a2 * G + F * I + E * E**2
    Y3 = I * X3 + G * H

    return X3, Y3, Z3


############################################################


def is_on_curve(a2, a6, X, Y):
    """
    lhs = yy+xy
    rhs = x*x*x + a_2 * x * x + a_6
    """

    lhs = Y**2 + X * Y
    rhs = X**3 + a2 * X**2 + a6

    return lhs, rhs


############################################################


def main():
    g = GoGenerator(package_name="weierstrassbinary", gf="GF2m")

    g.gen_projection_fn(add2005dl, projection="jacobian")
    g.gen_projection_fn(dbl2005dl, projection="jacobian")
    g.gen_projection_fn(madd2008bl, projection="jacobian")

    # g.gen_projection_fn(is_on_curve)


if __name__ == "__main__":
    main()
