from gogen import GoGenerator

############################################################


def add(a, b, x1, y1, x2, y2):
    """
    https://hyperelliptic.org/EFD/g1p/auto-montgom.html

    Affine addition formulas: (x1,y1)+(x2,y2)=(x3,y3) where
    x3 = b*(y2-y1)**2/(x2-x1)**2-a-x1-x2
    y3 = (2*x1+x2+a)*(y2-y1)/(x2-x1)-b*(y2-y1)**3/(x2-x1)**3-y1
    """

    dxy = (y2 - y1) / (x2 - x1)
    bdxy2 = b * dxy**2

    x3 = bdxy2 - a - x1 - x2
    y3 = (2 * x1 + x2 + a - bdxy2) * dxy - y1

    return x3, y3


def dbl(a, b, x1, y1):
    """
    https://hyperelliptic.org/EFD/g1p/auto-montgom.html

    Affine doubling formulas: 2(x1,y1)=(x3,y3) where
    x3 = b*(3*x1**2+2*a*x1+1)**2/(2*b*y1)**2-a-x1-x1
    y3 = (2*x1+x1+a)*(3*x1**2+2*a*x1+1)/(2*b*y1)-b*(3*x1**2+2*a*x1+1)**3/(2*b*y1)**3-y1
    """

    D = (3 * x1**2 + 2 * a * x1 + 1) / (2 * b * y1)
    bD2 = b * D**2

    x3 = bD2 - a - x1 - x1
    y3 = (3 * x1 + a - bD2) * D - y1

    return x3, y3


############################################################


def is_on_curve(a, b, X, Y):
    # byy = xxx + axx + x
    lhs = b * Y**2
    rhs = X**3 + a * X**2 + X

    return lhs, rhs


def compute_y(a, b, X):
    # y = sqrt((xxx + axx + x) / b)
    y1 = ((X**3 + a * X**2 + X) / b) ** 0.5
    return y1


############################################################


def main():
    g = GoGenerator(package_name="montgomery", gf="GFp")

    g.gen_projection_fn(add, projection="affine")
    g.gen_projection_fn(dbl, projection="affine")

    # g.gen_projection_fn(is_on_curve)
    # g.gen_projection_fn(compute_y)


if __name__ == "__main__":
    main()
