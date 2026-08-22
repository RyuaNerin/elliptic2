#!/usr/bin/env sage

from __future__ import annotations

import json
import sys
import traceback
from collections.abc import Callable
from typing import TYPE_CHECKING, Final, cast

from type_sage import (
    EdwardsCurveParams,
    GeneratedTestCases,
    MontgomeryCurveParams,
    TestPoint,
    TwistedEdwardsCurveParams,
    UnifiedCurve,
    WeierstrassBinaryCurveParams,
    WeierstrassPrimeCurveParams,
)
from type_std import CurveTypes

if TYPE_CHECKING:
    from sage.all import (  # pyright: ignore[reportMissingModuleSource]
        EllipticCurveObject as TypingEllipticCurve,
    )
    from sage.all import (  # pyright: ignore[reportMissingModuleSource]
        FiniteFieldElement as TypingPointElement,
    )
    from sage.all import (  # pyright: ignore[reportMissingModuleSource]
        Point as TypingSagePoint,
    )

    TypingPoint = (
        TypingSagePoint | tuple[TypingPointElement | int, TypingPointElement | int]
    )

from sage.all import (  # pyright: ignore[reportMissingModuleSource]
    GF,
    EllipticCurve,
    set_random_seed,
)
from sage.misc.prandom import randint  # pyright: ignore[reportMissingModuleSource]

set_random_seed(20251129)

NUM_OF_TESTCASE: Final = 10

INPUT_CURVE: Final = "/input/curve.json"
OUTPUT_RESULT: Final = "/output/result.json"
OUTPUT_ERROR: Final = "/output/error.txt"


def point_to_xy(P: TypingPoint) -> tuple[int, int]:
    return (elem_to_int(P[0]), elem_to_int(P[1]))


def point_to_dict(P: TypingPoint) -> TestPoint:
    return {
        "x": format_int(elem_to_int(P[0])),
        "y": format_int(elem_to_int(P[1])),
    }


def elem_to_int(elem: TypingPointElement | int) -> int:
    try:
        return elem.to_integer()  # pyright: ignore[reportAttributeAccessIssue]
    except Exception:  # noqa: BLE001
        return int(elem)


def format_int(value: int) -> str:
    if value == 0:
        return "0x0"
    return hex(value)


def format_point(p: TypingPoint | None) -> str:
    if p is None:
        return "None"

    pd = point_to_dict(p)
    return f"(x: {pd['x']}, y: {pd['y']})" if pd else "Zero"


def to_int(x: TypingPointElement | str | int | None) -> int:
    if x is None:
        raise ValueError("Cannot convert None to int")

    if isinstance(x, int):
        return x

    if isinstance(x, str):
        isNegative = x.startswith("-")
        if isNegative:
            x = x[1:]

        if x.startswith(("0x", "0X")):
            return int(x, 16) * (-1 if isNegative else 1)
        else:
            return int(x) * (-1 if isNegative else 1)

    if isinstance(x, dict) and "raw" in x:
        return to_int(x["raw"])

    try:
        x = elem_to_int(x)
        return int(x)
    except Exception as ex:  # noqa: BLE001
        raise ValueError(f"Cannot convert {x} to int", ex)


def err_print(*values: object):
    print(*values, file=sys.stderr, flush=True)


def create_curve(
    input: UnifiedCurve,
) -> tuple[
    TypingEllipticCurve,
    TypingSagePoint | None,
    Callable[[TypingPoint], TypingPoint],
    Callable[[TypingPoint], TypingPoint],
]:
    err_print("Creating curve...")
    err_print(json.dumps(input, indent=2))

    match CurveTypes(input["type"]):
        case CurveTypes.WEIERSTRASS_BINARY_POLY:
            return create_weierstrass_binary_poly_curve(input)

        case CurveTypes.WEIERSTRASS_PRIME:
            return create_weierstrass_prime_curve(input)

        case CurveTypes.MONTGOMERY_PRIME:
            return create_montgomery_curve(input)

        case CurveTypes.TWISTEDEDWARDS_PRIME:
            return create_twistededwards_curve(input)

        case CurveTypes.EDWARDS_PRIME:
            return create_edwards_curve(input)

    raise ValueError(f"Unsupported curve type: {input['type']}")


def create_weierstrass_binary_poly_curve(
    curve_data: UnifiedCurve,
) -> tuple[
    TypingEllipticCurve,
    TypingSagePoint | None,
    Callable[[TypingPoint], TypingPoint],
    Callable[[TypingPoint], TypingPoint],
]:
    params = cast(WeierstrassBinaryCurveParams, curve_data["params"])

    m = params["bits"]
    fx = params["fx"]
    a = params["a2"]
    b = params["a6"]
    n = params["n"]
    h = params["h"]

    err_print(f"m = {hex(m)}")
    err_print(f"fx = {[hex(power) for power in fx]}")
    err_print(f"a = {hex(a)}")
    err_print(f"b = {hex(b)}")
    err_print(f"n = {hex(n)}")
    err_print(f"h = {hex(h)}")

    """
    K-163:
        Parameters:
            m       163
            f(x)    x^163 + x^7 + x^6 + x^3 + 1
            a       0x000000000000000000000000000000000000000001
            b       0x000000000000000000000000000000000000000001
            G       (0x02fe13c0537bbc11acaa07d793de4e6d5e5c94eee8, 0x0289070fb05d38ff58321f2e800536d538ccdaa3d9)
            n       0x04000000000000000000020108a2e0cc0d99f8a5ef
            h       0x2
        SAGE:
            F.<x> = GF(2)[]
            K = GF(2^163, name="x", modulus= x^163 +  x^7 +  x^6 +  x^3 + 1)
            E = EllipticCurve(K, (1, K.from_integer(0x000000000000000000000000000000000000000001), 0, 0, K.from_integer(0x000000000000000000000000000000000000000001)))
            E.set_order(0x04000000000000000000020108a2e0cc0d99f8a5ef * 0x2)
            G = E(K.from_integer(0x02fe13c0537bbc11acaa07d793de4e6d5e5c94eee8), K.from_integer(0x0289070fb05d38ff58321f2e800536d538ccdaa3d9))

    K-233:
        Parameters:
            m       233
            f(x)    x^233 + x^74 + 1
            a       0x000000000000000000000000000000000000000000000000000000000000
            b       0x000000000000000000000000000000000000000000000000000000000001
            G       (0x017232ba853a7e731af129f22ff4149563a419c26bf50a4c9d6eefad6126, 0x01db537dece819b7f70f555a67c427a8cd9bf18aeb9b56e0c11056fae6a3)
            n       0x8000000000000000000000000000069d5bb915bcd46efb1ad5f173abdf
            h       0x4
        SAGE:
            F.<x> = GF(2)[]
            K = GF(2^233, name="x", modulus= x^233 +  x^74 + 1)
            E = EllipticCurve(K, (1, K.from_integer(0x000000000000000000000000000000000000000000000000000000000000), 0, 0, K.from_integer(0x000000000000000000000000000000000000000000000000000000000001)))
            E.set_order(0x8000000000000000000000000000069d5bb915bcd46efb1ad5f173abdf * 0x4)
            G = E(K.from_integer(0x017232ba853a7e731af129f22ff4149563a419c26bf50a4c9d6eefad6126), K.from_integer(0x01db537dece819b7f70f555a67c427a8cd9bf18aeb9b56e0c11056fae6a3))
    """

    err_print("Defining field K...")
    F = GF(2)["x"]
    (x,) = F._first_ngens(1)

    err_print("Calculating modulus...")
    modulus = sum([x**power for power in fx])
    err_print(f"modulus = {modulus}")

    err_print("Defining field K...")
    K = GF(2**m, name="x", modulus=modulus)

    err_print("Defining curve E...")
    E = EllipticCurve(K, [K.from_integer(x) for x in [1, a, 0, 0, b]])

    err_print("Setting curve order...")
    E.set_order(n * h)

    G: TypingSagePoint | None = None
    Gx = params.get("gx")
    Gy = params.get("gy")

    if Gx and Gy:
        err_print("Defining generator point G...")
        err_print(f"Gx = {hex(Gx)}")
        err_print(f"Gy = {hex(Gy)}")

        G = E(K.from_integer(Gx), K.from_integer(Gy))

    return E, G, lambda x: x, lambda x: x


def create_weierstrass_prime_curve(
    curve_data: UnifiedCurve,
) -> tuple[
    TypingEllipticCurve,
    TypingSagePoint | None,
    Callable[[TypingPoint], TypingPoint],
    Callable[[TypingPoint], TypingPoint],
]:
    params = cast(WeierstrassPrimeCurveParams, curve_data["params"])

    p = params["p"]
    a = params["a"]
    b = params["b"]
    n = params["n"]
    h = params["h"]

    err_print(f"p = {hex(p)}")
    err_print(f"a = {hex(a)}")
    err_print(f"b = {hex(b)}")
    err_print(f"n = {hex(n)}")
    err_print(f"h = {hex(h)}")

    # https://doc.sagemath.org/html/en/reference/arithmetic_curves/sage/schemes/elliptic_curves/ell_finite_field.html#sage.schemes.elliptic_curves.ell_finite_field.EllipticCurve_finite_field.set_order
    """
    P-192
        Parameters:
            p   0xfffffffffffffffffffffffffffffffeffffffffffffffff
            a   0xfffffffffffffffffffffffffffffffefffffffffffffffc
            b   0x64210519e59c80e70fa7e9ab72243049feb8deecc146b9b1
            G   (0x188da80eb03090f67cbf20eb43a18800f4ff0afd82ff1012, 0x07192b95ffc8da78631011ed6b24cdd573f977a11e794811)
            n   0xffffffffffffffffffffffff99def836146bc9b1b4d22831
            h   0x1
        SAGE:
            p = 0xfffffffffffffffffffffffffffffffeffffffffffffffff
            K = GF(p)
            a = K(0xfffffffffffffffffffffffffffffffefffffffffffffffc)
            b = K(0x64210519e59c80e70fa7e9ab72243049feb8deecc146b9b1)
            E = EllipticCurve(K, (a, b))
            E.set_order(0xffffffffffffffffffffffff99def836146bc9b1b4d22831 * 0x1)

            G = E(0x188da80eb03090f67cbf20eb43a18800f4ff0afd82ff1012, 0x07192b95ffc8da78631011ed6b24cdd573f977a11e794811)
    """

    # p = p
    err_print("Defining field K...")
    K = GF(p)

    err_print("Defining curve E...")
    a = K(a)
    b = K(b)
    E = EllipticCurve(K, (a, b))

    err_print("Setting curve order...")
    E.set_order(n * h)

    G: TypingSagePoint | None = None
    Gx = params.get("gx")
    Gy = params.get("gy")

    if Gx and Gy:
        err_print("Defining generator point G...")
        err_print(f"Gx = {hex(Gx)}")
        err_print(f"Gy = {hex(Gy)}")

        G = E(Gx, Gy)

    return E, G, lambda x: x, lambda x: x


def create_montgomery_curve(
    curve_data: UnifiedCurve,
) -> tuple[
    TypingEllipticCurve,
    TypingSagePoint | None,
    Callable[[TypingPoint], TypingPoint],
    Callable[[TypingPoint], TypingPoint],
]:
    params = cast(MontgomeryCurveParams, curve_data["params"])

    p = params["p"]
    a = params["a"]
    b = params["b"]
    n = params["n"]
    h = params["h"]

    err_print(f"p = {hex(p)}")
    err_print(f"a = {hex(a)}")
    err_print(f"b = {hex(b)}")
    err_print(f"n = {hex(n)}")
    err_print(f"h = {hex(h)}")

    """
    M-221
        Parameters:
            p	0x1FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFD
            a	0x01c93a
            b	0x01
            G	(0x04, 0x0f7acdd2a4939571d1cef14eca37c228e61dbff10707dc6c08c5056d)
            n	0x040000000000000000000000000015A08ED730E8A2F77F005042605B
            h	0x8
        SAGE:
            p = 0x1FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFD
            K = GF(p)
            A = K(0x01c93a)
            B = K(0x01)
            E = EllipticCurve(K, ((3 - A^2)/(3 * B^2), (2 * A^3 - 9 * A)/(27 * B^3)))
            def to_weierstrass(A, B, x, y):
                return (x/B + A/(3*B), y/B)
            def to_montgomery(A, B, u, v):
                return (B * (u - A/(3*B)), B*v)
            G = E(*to_weierstrass(A, B, K(0x04), K(0x0f7acdd2a4939571d1cef14eca37c228e61dbff10707dc6c08c5056d)))
            E.set_order(0x040000000000000000000000000015A08ED730E8A2F77F005042605B * 0x8)
            # This curve is a Weierstrass curve (SAGE does not support Montgomery curves) birationally equivalent to the intended curve.
            # You can use the to_weierstrass and to_montgomery functions to convert the points.
    """

    err_print("Defining field K...")
    K = GF(p)

    err_print("Defining curve E...")
    A = K(a)
    B = K(b)
    E = EllipticCurve(K, ((3 - A**2) / (3 * B**2), (2 * A**3 - 9 * A) / (27 * B**3)))

    def to_weierstrass(P: TypingPoint) -> TypingPoint:
        x, y = P[0], P[1]
        u = x / B + A / (3 * B)
        v = y / B
        return (u, v)  # type: ignore[return-value]

    def to_montgomery(P: TypingPoint) -> TypingPoint:
        u, v = P[0], P[1]
        x = B * (u - A / (3 * B))
        y = B * v
        return (x, y)  # type: ignore[return-value]

    err_print("Setting curve order...")
    E.set_order(n * h)

    G: TypingSagePoint | None = None
    Gx = params.get("gx")
    Gy = params.get("gy")

    if Gx and Gy:
        err_print("Defining generator point G...")
        err_print(f"Gx = {hex(Gx)}")
        err_print(f"Gy = {hex(Gy)}")

        G = E(*to_weierstrass((K(Gx), K(Gy))))  # type: ignore[return-value]

    return E, G, to_weierstrass, to_montgomery


def create_twistededwards_curve(
    curve_data: UnifiedCurve,
) -> tuple[
    TypingEllipticCurve,
    TypingSagePoint | None,
    Callable[[TypingPoint], TypingPoint],
    Callable[[TypingPoint], TypingPoint],
]:
    params = cast(TwistedEdwardsCurveParams, curve_data["params"])

    p = params["p"]
    a = params["a"]
    d = params["d"]
    n = params["n"]
    h = params["h"]

    err_print(f"p = {hex(p)}")
    err_print(f"a = {hex(a)}")
    err_print(f"d = {hex(d)}")
    err_print(f"n = {hex(n)}")
    err_print(f"h = {hex(h)}")

    """
    numsp256t1
        Parameters:
            p	0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff43
            a	0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff42
            b	0x3bee
            G	(0x0d, 0x7d0ab41e2a1276dba3d330b39fa046bfbe2a6d63824d303f707f6fb5331cadba)
            n	0x3fffffffffffffffffffffffffffffffbe6aa55ad0a6bc64e5b84e6f1122b4ad
            h	0x04
        SAGE:
            p = 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff43
            K = GF(p)
            a = K(0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff42)
            d = K(0x3bee)
            E = EllipticCurve(K, (K(-1/48) * (a^2 + 14*a*d + d^2),K(1/864) * (a + d) * (-a^2 + 34*a*d - d^2)))
            def to_weierstrass(a, d, x, y):
                return ((5*a + a*y - 5*d*y - d)/(12 - 12*y), (a + a*y - d*y -d)/(4*x - 4*x*y))
            def to_twistededwards(a, d, u, v):
                y = (5*a - 12*u - d)/(-12*u - a + 5*d)
                x = (a + a*y - d*y -d)/(4*v - 4*v*y)
                return (x, y)
            G = E(*to_weierstrass(a, d, K(0x0d), K(0x7d0ab41e2a1276dba3d330b39fa046bfbe2a6d63824d303f707f6fb5331cadba)))
            E.set_order(0x3fffffffffffffffffffffffffffffffbe6aa55ad0a6bc64e5b84e6f1122b4ad * 0x04)
            # This curve is a Weierstrass curve (SAGE does not support TwistedEdwards curves) birationally equivalent to the intended curve.
            # You can use the to_weierstrass and to_twistededwards functions to convert the points.
    """

    err_print("Defining field K...")
    K = GF(p)

    err_print("Defining curve E...")
    a = K(a)
    d = K(d)
    E = EllipticCurve(
        K,
        (
            (K(-1) / K(48)) * (a**2 + 14 * a * d + d**2),
            (K(1) / K(864)) * (a + d) * (-(a**2) + 34 * a * d - d**2),
        ),
    )

    def to_weierstrass(P: TypingPoint) -> TypingPoint:
        x, y = P[0], P[1]
        u = (5 * a + a * y - 5 * d * y - d) / (12 - 12 * y)
        v = (a + a * y - d * y - d) / (4 * x - 4 * x * y)
        return (u, v)  # type: ignore[return-value]

    def to_twistededwards(P: TypingPoint) -> TypingPoint:
        u, v = P[0], P[1]
        y = (5 * a - 12 * u - d) / (-12 * u - a + 5 * d)
        x = (a + a * y - d * y - d) / (4 * v - 4 * v * y)
        return (x, y)  # type: ignore[return-value]

    err_print("Setting curve order...")
    E.set_order(n * h)

    G: TypingSagePoint | None = None
    Gx = params.get("gx")
    Gy = params.get("gy")

    if Gx and Gy:
        err_print("Defining generator point G...")
        err_print(f"Gx = {hex(Gx)}")
        err_print(f"Gy = {hex(Gy)}")

        G = E(*to_weierstrass((K(Gx), K(Gy))))  # type: ignore[return-value]

    return E, G, to_weierstrass, to_twistededwards


def create_edwards_curve(
    curve_data: UnifiedCurve,
) -> tuple[
    TypingEllipticCurve,
    TypingSagePoint | None,
    Callable[[TypingPoint], TypingPoint],
    Callable[[TypingPoint], TypingPoint],
]:
    params = cast(EdwardsCurveParams, curve_data["params"])

    p = params["p"]
    c = params["c"]
    d = params["d"]
    n = params["n"]
    h = params["h"]

    err_print(f"p = {hex(p)}")
    err_print(f"c = {hex(c)}")
    err_print(f"d = {hex(d)}")
    err_print(f"n = {hex(n)}")
    err_print(f"h = {hex(h)}")

    """
    E-222
        Parameters:
            p	0x3fffffffffffffffffffffffffffffffffffffffffffffffffffff8b
            c	0x01
            d	0x27166
            G	(0x19b12bb156a389e55c9768c303316d07c23adab3736eb2bc3eb54e51, 0x1c)
            n	0xffffffffffffffffffffffffffff70cbc95e932f802f31423598cbf
            h	0x04
        SAGE:
            p = 0x3fffffffffffffffffffffffffffffffffffffffffffffffffffff8b
            K = GF(p)
            d = K(0x27166)
            E = EllipticCurve(K, (0, K(2 * (1 + d)/(1 - d)^2), 0, K(1/(1 - d)^2), 0))
            E.set_order(0xffffffffffffffffffffffffffff70cbc95e932f802f31423598cbf * 0x04)
            # This curve is a Weierstrass curve (SAGE does not support Edwards curves) birationally equivalent to the intended curve.
            # You can use the to_weierstrass and to_edwards functions to convert the points.
    """

    err_print("Defining field K...")
    K = GF(p)

    err_print("Defining curve E...")
    c_K = K(c)
    d_K = K(d)

    d_normalized = d_K * c_K**4

    A = K(2 * (1 + d_normalized) / (1 - d_normalized))
    B = K(4 / (1 - d_normalized))

    # Weierstrass short form: y² = x³ + ax + b
    a_w = (3 - A**2) / (3 * B**2)
    b_w = (2 * A**3 - 9 * A) / (27 * B**3)

    E = EllipticCurve(K, [a_w, b_w])

    def to_weierstrass(P: TypingPoint) -> TypingPoint:
        """Edwards (x, y) → Weierstrass (x_w, y_w)"""
        x, y = K(P[0]) / c_K, K(P[1]) / c_K
        # Edwards → Montgomery
        u = (1 + y) / (1 - y)
        v = u / x
        # Montgomery → Weierstrass
        x_w = u / B + A / (3 * B)
        y_w = v / B
        return (x_w, y_w)  # pyright: ignore[reportReturnType]

    def to_edwards(P: TypingPoint) -> TypingPoint:
        """Weierstrass (x_w, y_w) → Edwards (x, y)"""
        x_w, y_w = K(P[0]), K(P[1])
        # Weierstrass → Montgomery
        u = B * (x_w - A / (3 * B))
        v = B * y_w
        # Montgomery → Edwards
        y = c_K * (u - 1) / (u + 1)
        x = c_K * u / v
        return (x, y)  # pyright: ignore[reportReturnType]

    err_print("Setting curve order...")
    E.set_order(n * h)

    G: TypingSagePoint | None = None
    Gx = params.get("gx")
    Gy = params.get("gy")

    if Gx and Gy:
        err_print("Defining generator point G...")
        err_print(f"Gx = {hex(Gx)}")
        err_print(f"Gy = {hex(Gy)}")

        G = E(*to_weierstrass((Gx, Gy)))  # type: ignore[return-value]

    return E, G, to_weierstrass, to_edwards


def generate_tests(input: UnifiedCurve) -> GeneratedTestCases:
    err_print(f"Generating test cases for curve: {input['name']} - {input['type']}")

    curve, G, to_weierstrass, to_curvepoint = create_curve(input)
    n = input["params"]["n"]

    # Start in Sage's Weierstrass model, convert to the target model, and back.
    P_w = curve.random_point()
    P_c = to_curvepoint(P_w)
    P_w_roundtrip = to_weierstrass(P_c)
    if point_to_xy(P_w) != point_to_xy(P_w_roundtrip):
        raise ValueError(
            "to_weierstrass and to_curvepoint functions are not consistent"
        )

    k_list: list[int] = []
    p_list: list[TypingSagePoint] = []
    add_list: list[TypingSagePoint] = []
    double_list: list[TypingSagePoint] = []
    scalar_mult_list: list[TypingSagePoint] = []
    invalid_p_list: list[TypingPoint] = []

    err_print("Generating K")
    for idx in range(NUM_OF_TESTCASE):
        k = randint(1, n - 1)
        k_list.append(k)
        err_print(f"k[{idx + 1}]: {format_int(k)}")

    if G:
        err_print("Generating points")
        for idx in range(NUM_OF_TESTCASE):
            p1 = k_list[idx] * G
            p_list.append(p1)
            err_print(f"P  [{idx + 1}]: {format_point(p1)}")
    else:
        err_print("Generating points")
        for idx in range(NUM_OF_TESTCASE):
            p1 = curve.random_point()
            p_list.append(p1)
            err_print(f"P  [{idx + 1}]: {format_point(p1)}")

    err_print("Generating addition")
    p1 = p_list[0]
    for idx in range(NUM_OF_TESTCASE):
        p1 = p1 + p_list[idx]
        add_list.append(p1)
        err_print(f"Add {idx + 1}: {format_point(p1)}")

    err_print("Generating doubling")
    for idx in range(NUM_OF_TESTCASE):
        p1 = p_list[idx] + p_list[idx]
        double_list.append(p1)
        err_print(f"Double {idx + 1}: {format_point(p1)}")

    err_print("Generating scalar multiplication")
    p1 = p_list[0]
    for idx in range(NUM_OF_TESTCASE):
        p1 = p1 * k_list[idx]
        err_print(f"ScalarMult {idx + 1}: {format_point(p1)}")
        scalar_mult_list.append(p1)

    err_print("Generating invalid points")
    while len(invalid_p_list) < NUM_OF_TESTCASE:
        p1 = curve.random_point()
        x, y = p1[0] + 1, p1[1]

        if curve.is_on_curve(x, y):
            err_print(f"IsNotOnCurve {len(invalid_p_list) + 1} failed, retrying...")
            continue

        try:
            invalid_point = to_curvepoint((x, y))
            roundtrip = to_weierstrass(invalid_point)
        except ZeroDivisionError:
            err_print(
                f"IsNotOnCurve {len(invalid_p_list) + 1} failed (ZeroDivisionError), retrying..."
            )
            continue
        except ArithmeticError:
            err_print(
                f"IsNotOnCurve {len(invalid_p_list) + 1} failed (ArithmeticError), retrying..."
            )
            continue

        if point_to_xy(roundtrip) != point_to_xy((x, y)):
            continue

        err_print(f"IsNotOnCurve {len(invalid_p_list) + 1} succeeded")
        invalid_p_list.append(invalid_point)

    return {
        "k": [format_int(k) for k in k_list],
        "scalar_base_mult": [point_to_dict(to_curvepoint(p)) for p in p_list],
        "add": [point_to_dict(to_curvepoint(p)) for p in add_list],
        "double": [point_to_dict(to_curvepoint(p)) for p in double_list],
        "scalar_mult": [point_to_dict(to_curvepoint(p)) for p in scalar_mult_list],
        "invalid_p": [point_to_dict(p) for p in invalid_p_list],
    }


# 메인 실행
def main():
    try:
        with open(INPUT_CURVE, "r") as f:
            curve_data: UnifiedCurve = json.load(f)

        tc = generate_tests(curve_data)

        with open(OUTPUT_RESULT, "w") as f:
            json.dump(tc, f, indent=2)

    except Exception as e:  # noqa: BLE001
        s = str(e) + "\n\n" + traceback.format_exc()
        err_print(s)

        with open(OUTPUT_ERROR, "w") as f:
            f.write(s)


if __name__ == "__main__":
    main()
