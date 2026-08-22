from typing import NotRequired, TypedDict

from type_std import GitCurve


class Params(TypedDict):
    n: int


class CurveParams(TypedDict):
    bits: int  # = bits = m
    n: int
    h: int
    gx: NotRequired[int | None]
    gy: NotRequired[int | None]
    oid: NotRequired[str | None]


class PrimeCurveParams(CurveParams):
    p: int


class BinaryCurveParams(CurveParams):
    fx: list[int]


class ExtensionCurvePrams(CurveParams):
    fx: list[int]


class WeierstrassBinaryCurveParams(BinaryCurveParams):
    a2: int  # A
    a6: int  # B


class WeierstrassPrimeCurveParams(PrimeCurveParams):
    a: int
    b: int


class TwistedEdwardsCurveParams(PrimeCurveParams):
    a: int
    d: int


class EdwardsCurveParams(PrimeCurveParams):
    c: int
    d: int


class MontgomeryCurveParams(PrimeCurveParams):
    a: int
    b: int


class TestPoint(TypedDict):
    x: str
    y: str


CURVE_PARAMS = (
    Params
    | WeierstrassPrimeCurveParams
    | WeierstrassBinaryCurveParams
    | TwistedEdwardsCurveParams
    | EdwardsCurveParams
    | MontgomeryCurveParams
)


class UnifiedCurve(TypedDict):
    name: str
    desc: str
    type: str
    aliases: list[str]
    params: CURVE_PARAMS
    raw: GitCurve


class IgnoredCurve(TypedDict):
    name: str
    desc: str
    type: str
    aliases: list[str]
    reason: str
    raw: GitCurve


class UnifiedNamespace(TypedDict):
    name: str
    namespace: str
    desc: str
    curve_list: list[UnifiedCurve]
    ignored_curve_list: list[IgnoredCurve]


class GeneratedTestCases(TypedDict):
    k: list[str]
    scalar_base_mult: list[TestPoint]
    add: list[TestPoint]
    double: list[TestPoint]
    scalar_mult: list[TestPoint]
    invalid_p: list[TestPoint]
