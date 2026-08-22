from enum import Enum
from typing import NotRequired, TypedDict


class GitTorsionDegree(TypedDict):
    r: int
    least: int
    full: int


class GitPoly(TypedDict):
    power: int
    coeff: str


class GitFieldBase(TypedDict):
    type: str


class GitFieldBinaryPoly(GitFieldBase):
    bits: int
    degree: int
    poly: list[GitPoly]
    basis: str


class GitFieldPrime(GitFieldBase):
    bits: str | int
    p: str | int


class GitRaw(TypedDict):
    raw: NotRequired[str | None]


class GitGenerator(TypedDict):
    x: NotRequired[GitRaw | None]
    y: NotRequired[GitRaw | None]


class GitParams(TypedDict):
    a: NotRequired[GitRaw | None]
    b: NotRequired[GitRaw | None]
    c: NotRequired[GitRaw | None]
    d: NotRequired[GitRaw | None]


class GitCurve(TypedDict):
    name: str
    category: str
    desc: str
    field: GitFieldBase | GitFieldBinaryPoly | GitFieldPrime
    form: str
    params: GitParams
    order: str  # N
    cofactor: str  # H
    generator: NotRequired[GitGenerator | None]
    oid: NotRequired[str | None]
    aliases: NotRequired[list[str] | None]


class GitDocument(TypedDict):
    name: str
    desc: str
    curves: list[GitCurve]


class UnsupportedCurveTypeError(Exception):
    @staticmethod
    def _get_message(form: str, field_type: str, basis: str | None = None) -> str:
        return f"Unsupported curve type: form={form}, field_type={field_type}" + (
            f", basis={basis}" if basis else ""
        )

    def __init__(self, form: str, field_type: str, basis: str | None = None):
        super().__init__(
            UnsupportedCurveTypeError._get_message(form, field_type, basis)
        )

        self.form = form
        self.field_type = field_type
        self.basis = basis

    @property
    def message(self) -> str:
        return UnsupportedCurveTypeError._get_message(
            self.form, self.field_type, self.basis
        )


class CurveTypes(Enum):
    WEIERSTRASS_BINARY_POLY = "WEIERSTRASS_BINARY_POLY"
    WEIERSTRASS_BINARY_NORMAL = "WEIERSTRASS_BINARY_NORMAL"
    WEIERSTRASS_PRIME = "WEIERSTRASS_PRIME"
    TWISTEDEDWARDS_PRIME = "TWISTEDEDWARDS_PRIME"
    MONTGOMERY_PRIME = "MONTGOMERY_PRIME"
    EDWARDS_PRIME = "EDWARDS_PRIME"

    @staticmethod
    def get(curve_data: GitCurve) -> "CurveTypes":
        CURVE_TYPES_MAP = {
            ("weierstrass", "binary", "poly"): CurveTypes.WEIERSTRASS_BINARY_POLY,
            # ("weierstrass", "binary", "normal"): CurveTypes.WEIERSTRASS_BINARY_NORMAL,
            ("weierstrass", "prime", ""): CurveTypes.WEIERSTRASS_PRIME,
            ("montgomery", "prime", ""): CurveTypes.MONTGOMERY_PRIME,
            ("twistededwards", "prime", ""): CurveTypes.TWISTEDEDWARDS_PRIME,
            ("edwards", "prime", ""): CurveTypes.EDWARDS_PRIME,
        }

        form = curve_data["form"].lower()
        field_type = curve_data["field"]["type"].lower()
        basis = curve_data["field"].get("basis", "").lower()

        t = CURVE_TYPES_MAP.get((form, field_type, basis), None)
        if t is None:
            raise UnsupportedCurveTypeError(form, field_type, basis)

        return t
