"""curve_generator
python 으로 기술한 타원곡선 공식을 Go 코드로 변환한다.

  1) AST -> AC(결합/교환) 체인 평탄화
  2) 전역 pair 재결합으로 중복 연산 제거 (해시컨싱 CSE + 2-pass 사전정보)
  3) CP-SAT 로 레지스터 압력 최소 스케줄 (출력 dst 를 초기 스크래치로 사용)
  4) 구간그래프 선형스캔 채색 -> op.t[N] 최소화
"""

import ast
import inspect
import re
import shutil
import subprocess
import textwrap
from itertools import combinations
from pathlib import Path
from typing import TypedDict

from ortools.sat.python import cp_model

POW2 = {1: 0, 2: 1, 4: 2, 8: 3, 16: 4, 32: 5, 64: 6}
COMM = {"add", "mul"}


# 확장 가능한 단항 연산 레지스트리
#   go      : Go 메서드명
#   tag     : 비용 표기
#   text    : 주석/테이블 표기 포맷
#   boolean : 실패할 수 있는 연산 (Go 에서 nil 반환) -> 함수도 포인터를 반환
#   heavy   : 피연산자가 완전 축약되어야 하고 결과 magnitude = 1
class UnaryOp(TypedDict):
    go: str
    tag: str
    text: str
    boolean: bool
    heavy: bool


UNARY: dict[str, UnaryOp] = {
    "inv": UnaryOp(go="Inv", tag="I", text="1 / {0}", boolean=False, heavy=True),
    "sqrt": UnaryOp(go="Sqrt", tag="sqrt", text="{0}**(1/2)", boolean=True, heavy=True),
    "half": UnaryOp(go="Half", tag="half", text="{0} / 2", boolean=False, heavy=False),
    "neg": UnaryOp(go="Neg", tag="add", text="-{0}", boolean=False, heavy=False),
}


def register_op(name, go, tag, text, boolean=False, heavy=True):
    UNARY[name] = UnaryOp(go=go, tag=tag, text=text, boolean=boolean, heavy=heavy)


class DAG:
    def __init__(self):
        self.nodes, self._c = {}, {}

    @staticmethod
    def key(op, children, extra=None):
        if op in COMM:
            children = tuple(sorted(children))
        return (op, tuple(children), extra)

    def mk(self, op, children=(), extra=None):
        k = self.key(op, children, extra)
        if k in self._c:
            return self._c[k], False
        i = len(self.nodes)
        self.nodes[i] = k
        self._c[k] = i
        return i, True

    def has(self, op, children, extra=None):
        return self.key(op, children, extra) in self._c

    def op(self, i):
        return self.nodes[i][0]

    def kids(self, i):
        return list(self.nodes[i][1])

    def extra(self, i):
        return self.nodes[i][2]


class Row:
    __slots__ = ("idx", "node", "label", "text", "expr")

    def __init__(self, idx, node, label):
        self.idx, self.node, self.label = idx, node, label
        self.text, self.expr = "", None


class _NormalizeEfdSyntax(ast.NodeTransformer):
    """EFD 의 ``x ^ 2`` 제곱 표기를 Python 거듭제곱 AST 로 변환한다."""

    def visit_BinOp(self, node):
        node = self.generic_visit(node)
        if not isinstance(node.op, ast.BitXor):
            return node
        if not isinstance(node.right, ast.Constant) or node.right.value != 2:
            raise NotImplementedError("'^' 연산자는 EFD 제곱 표기인 '^ 2'만 지원합니다")
        node.op = ast.Pow()
        return node


class Formula:
    def __init__(self, func, cse=True):
        self.func, self.name = func, func.__name__
        self.doc = inspect.getdoc(func) or ""
        tree = ast.parse(textwrap.dedent(inspect.getsource(func))).body[0]
        self.tree = _NormalizeEfdSyntax().visit(tree)
        self.params = [a.arg for a in self.tree.args.args]
        self.cse = cse
        known = set()
        for _ in range(2):
            self._run(known)
            known = set(self.d._c)
        self._count()

    def _run(self, known):
        self.d = DAG()
        self.known = known
        self.env = {p: self.d.mk("in", (), p)[0] for p in self.params}
        self.label = {self.env[p]: p for p in self.params}
        self.rows, self.milestones = [], []
        self._anon, self._new = 0, []
        for st in self.tree.body:
            if isinstance(st, ast.Expr):
                continue
            if isinstance(st, ast.Assign):
                self._new = []
                node = self._visit(st.value)
                tgt = st.targets[0].id
                self.env[tgt] = node
                if self._new and self._new[-1] == node:
                    self.rows[-1].label = tgt
                    self.label[node] = tgt
                    self._retext(self.rows[-1])
                self.milestones.append(f"{tgt} = {ast.unparse(st.value)}")
                if self.rows:
                    self.rows[-1].expr = (len(self.milestones), self.milestones[-1])
            elif isinstance(st, ast.Return):
                v = st.value
                elts = v.elts if isinstance(v, ast.Tuple) else [v]
                self.roots = [self.env[e.id] for e in elts]
                self.outs = [e.id for e in elts]

    def _emit(self, op, children=(), extra=None):
        i, fresh = self.d.mk(op, children, extra)
        if not self.cse:
            fresh = True
        if fresh or i not in self.label:
            self.label[i] = f"_{self._anon}"
            self._anon += 1
            r = Row(len(self.rows) + 1, i, self.label[i])
            self._retext(r)
            self.rows.append(r)
            self._new.append(i)
        return i

    def _ref(self, i):
        return str(self.d.extra(i)) if self.d.op(i) == "const" else self.label[i]

    def _retext(self, row):
        i = row.node
        k = [self._ref(c) for c in self.d.kids(i)]
        o, e = self.d.op(i), self.d.extra(i)
        if o in UNARY:
            t = UNARY[o]["text"].format(*k)
        else:
            t = {
                "mul": lambda: f"{k[0]} * {k[1]}",
                "sqr": lambda: f"{k[0]}**2",
                "add": lambda: f"{k[0]} + {k[1]}",
                "sub": lambda: f"{k[0]} - {k[1]}",
                "kmul": lambda: f"{e} * {k[0]}",
            }[o]()
        row.text = f"{row.label} = {t}"

    def _flat_mul(self, n, fac, con):
        if isinstance(n, ast.BinOp) and isinstance(n.op, ast.Mult):
            self._flat_mul(n.left, fac, con)
            self._flat_mul(n.right, fac, con)
        elif isinstance(n, ast.BinOp) and isinstance(n.op, ast.Div):
            # 분모는 통째로 1회 역원 (역원은 가장 비싼 연산이므로 분배하지 않는다)
            self._flat_mul(n.left, fac, con)
            fac.append(self._emit("inv", (self._visit(n.right),)))
        elif isinstance(n, ast.BinOp) and isinstance(n.op, ast.Pow):
            e = n.right.value if isinstance(n.right, ast.Constant) else None
            if isinstance(e, int) and e >= 1:
                for _ in range(e):
                    self._flat_mul(n.left, fac, con)
            else:
                fac.append(self._pow_special(self._visit(n.left), e))
        elif isinstance(n, ast.Constant):
            con.append(n.value)
        elif isinstance(n, ast.UnaryOp) and isinstance(n.op, ast.USub):
            con.append(-1)
            self._flat_mul(n.operand, fac, con)
        else:
            fac.append(self._visit(n))

    def _flat_add(self, n, terms, con, sign=1):
        if isinstance(n, ast.BinOp) and isinstance(n.op, ast.Add):
            self._flat_add(n.left, terms, con, sign)
            self._flat_add(n.right, terms, con, sign)
        elif isinstance(n, ast.BinOp) and isinstance(n.op, ast.Sub):
            self._flat_add(n.left, terms, con, sign)
            self._flat_add(n.right, terms, con, -sign)
        elif isinstance(n, ast.UnaryOp) and isinstance(n.op, ast.USub):
            self._flat_add(n.operand, terms, con, -sign)
        elif isinstance(n, ast.Constant):
            con.append(sign * n.value)
        else:
            terms.append((sign, self._visit(n)))

    def _pow_special(self, base, e):
        """정수 이외 / 음수 지수 처리"""
        if e == 0:  # x**0
            return self.d.mk("const", (), 1)[0]
        if isinstance(e, float):
            if abs(e - 0.5) < 1e-12:  # x**0.5 -> Sqrt
                return self._emit("sqrt", (base,))
            if abs(e - 0.25) < 1e-12:
                return self._emit("sqrt", (self._emit("sqrt", (base,)),))
            raise NotImplementedError(f"지원하지 않는 지수: {e}")
        if isinstance(e, int) and e < 0:  # x**-n -> Inv(x**n)
            f, c = [base] * (-e), []
            return self._emit("inv", (self._finish_mul(f, c),))
        raise NotImplementedError(f"지원하지 않는 지수: {e}")

    def _score(self, op, ch):
        if self.d.has(op, ch):
            return 2
        if self.d.key(op, ch) in self.known:
            return 1
        return 0

    def _pair_mul(self, items):
        items = list(items)
        while len(items) > 1:
            best = None
            for a, b in combinations(range(len(items)), 2):
                x, y = items[a], items[b]
                op, ch = ("sqr", (x,)) if x == y else ("mul", (x, y))
                s = self._score(op, ch) * 2 + (1 if x == y else 0)
                if best is None or s > best[0]:
                    best = (s, a, b, op, ch)
            _, a, b, op, ch = best
            n = self._emit(op, ch)
            items = [items[i] for i in range(len(items)) if i not in (a, b)] + [n]
        return items[0]

    def _pair_add(self, terms):
        terms = list(terms)
        while len(terms) > 1:
            best = None
            for a, b in combinations(range(len(terms)), 2):
                (sx, x), (sy, y) = terms[a], terms[b]
                if sx == sy:
                    op, ch, sg = "add", (x, y), sx
                elif sx > 0:
                    op, ch, sg = "sub", (x, y), 1
                else:
                    op, ch, sg = "sub", (y, x), 1
                s = self._score(op, ch) * 2 + (1 if sx != sy else 0)
                if best is None or s > best[0]:
                    best = (s, a, b, op, ch, sg)
            _, a, b, op, ch, sg = best
            n = self._emit(op, ch)
            terms = [terms[i] for i in range(len(terms)) if i not in (a, b)] + [(sg, n)]
        s, n = terms[0]
        return n if s > 0 else self._emit("neg", (n,))

    def _visit(self, n):
        if isinstance(n, ast.Name):
            return self.env[n.id]
        if isinstance(n, ast.Call) and isinstance(n.func, ast.Name):
            fn = n.func.id
            if fn == "pow" and len(n.args) == 2:
                b = self._visit(n.args[0])
                e = n.args[1].value
                return (
                    self._finish_mul([b] * e, [])
                    if isinstance(e, int) and e >= 1
                    else self._pow_special(b, e)
                )
            if fn in UNARY:
                return self._emit(fn, tuple(self._visit(a) for a in n.args))
            raise NotImplementedError(f"알 수 없는 함수: {fn}")
        if isinstance(n, ast.Constant):
            return self.d.mk("const", (), n.value)[0]
        if isinstance(n, ast.UnaryOp) and isinstance(n.op, ast.USub):
            t, c = [], []
            self._flat_add(n, t, c)
            return self._finish_add(t, c)
        if isinstance(n, ast.BinOp):
            o = type(n.op).__name__
            if o in ("Add", "Sub"):
                t, c = [], []
                self._flat_add(n, t, c)
                return self._finish_add(t, c)
            if o in ("Mult", "Div", "Pow"):
                f, c = [], []
                self._flat_mul(n, f, c)
                return self._finish_mul(f, c)
        raise NotImplementedError(ast.dump(n))

    def _finish_mul(self, fac, con):
        k = 1
        for v in con:
            k *= v
        if not fac:
            return self.d.mk("const", (), k)[0]
        n = self._pair_mul(fac)
        if k == 1:
            return n
        if k == -1:
            return self._emit("neg", (n,))
        return self._emit("kmul", (n,), k)

    def _finish_add(self, terms, con):
        c = sum(con)
        if not terms:
            return self.d.mk("const", (), c)[0]
        if c:
            terms.append((1 if c > 0 else -1, self.d.mk("const", (), abs(c))[0]))
        return self._pair_add(terms)

    def _count(self):
        seen, st = set(), list(self.roots)
        while st:
            i = st.pop()
            if i in seen:
                continue
            seen.add(i)
            st += self.d.kids(i)
        c = {}
        for i in seen:
            o = self.d.op(i)
            if o in ("in", "const"):
                continue
            if o == "kmul":
                o = f"*{self.d.extra(i)}"
            c[o] = c.get(o, 0) + 1
        self.has_bool = any(
            self.d.op(i) in UNARY and UNARY[self.d.op(i)]["boolean"]
            for i in seen
            if self.d.op(i) not in ("in", "const")
        )
        self.const_nodes = [i for i in seen if self.d.op(i) == "const"]
        self.counts = c


def _reachable(d, roots):
    seen, st = set(), list(roots)
    while st:
        i = st.pop()
        if i in seen:
            continue
        seen.add(i)
        st += d.kids(i)
    return sorted(i for i in seen if d.op(i) not in ("in", "const"))


def _distinct(d, i):
    return d.op(i) == "kmul" and d.extra(i) not in POW2


def color(d, order, roots, dstn):
    use = {i: 0 for i in d.nodes}
    for i in order:
        for c in d.kids(i):
            use[c] += 1
    for r in roots:
        use[r] += 1
    free, reg, nc = [], {}, 0

    def new():
        nonlocal nc
        if free:
            return free.pop(0)
        nc += 1
        return nc - 1

    for i in order:
        ks = d.kids(i)

        def rel():
            for c in ks:
                use[c] -= 1
                if use[c] == 0 and c in reg:
                    free.append(reg[c])

        if not _distinct(d, i):
            rel()
        reg[i] = new()
        if _distinct(d, i):
            rel()

    name, t = {}, 0
    for r, nm in zip(roots, dstn):
        name[reg[r]] = nm
    for i in order:
        if reg[i] not in name:
            name[reg[i]] = f"_{t}"
            t += 1
    return {i: name[reg[i]] for i in order}, t


def schedule_cpsat(d, roots, time_limit=30.0, workers=8):
    N = _reachable(d, roots)
    n = len(N)
    idx, rs = set(N), set(roots)
    users = {i: [] for i in N}
    for i in N:
        for c in d.kids(i):
            if c in idx:
                users[c].append(i)
    m = cp_model.CpModel()
    s = {i: m.new_int_var(0, n - 1, f"s{i}") for i in N}
    m.add_all_different(list(s.values()))
    for i in N:
        for c in d.kids(i):
            if c in idx:
                m.add(s[c] < s[i])
    R = m.new_int_var(1, n, "R")
    iv = []
    for i in N:
        ends = []
        for u in users[i]:
            if _distinct(d, u):
                e = m.new_int_var(0, n, f"e{i}_{u}")
                m.add(e == s[u] + 1)
                ends.append(e)
            else:
                ends.append(s[u])
        if i in rs:
            ends.append(m.new_constant(n))
        if not ends:
            continue
        e = m.new_int_var(0, n, f"e{i}")
        m.add_max_equality(e, ends)
        du = m.new_int_var(1, n, f"d{i}")
        m.add(du == e - s[i])
        iv.append(m.new_interval_var(s[i], du, e, f"iv{i}"))
    m.add_cumulative(iv, [1] * len(iv), R)
    m.minimize(R)
    so = cp_model.CpSolver()
    so.parameters.max_time_in_seconds = time_limit
    so.parameters.num_search_workers = workers
    st = so.Solve(m)
    if st not in (cp_model.OPTIMAL, cp_model.FEASIBLE):
        raise RuntimeError("no schedule")
    return sorted(N, key=lambda i: so.Value(s[i])), st == cp_model.OPTIMAL


def optimize(d, roots, dstn, time_limit=30.0):
    order, opt = schedule_cpsat(d, roots, time_limit)
    reg, nt = color(d, order, roots, dstn)
    return order, reg, nt, opt


class GoGenerator:
    def __init__(
        self,
        package_name,
        gf,
        max_magnitude=4,
        cse=True,
    ):
        self.pkg = package_name
        self.gf = gf
        self.maxmag, self.cse = max_magnitude, cse
        self.coord = f"curve.{gf}Coordinate"
        self.binary = gf == "GF2m"

    @staticmethod
    def _kname(v):
        names = {0: "zero", 1: "one", 2: "two", 3: "three"}
        return names.get(v, f"constant{'N' if v < 0 else ''}{abs(v)}")

    def _kexpr(self, v):
        predefined = {0: "Zero", 1: "One", 2: "Two", 3: "Three"}
        if not self.binary and v in predefined:
            return f"field.GFp{predefined[v]}"
        return f"field.New{self.gf}Int64({v})"

    def gen_projection_fn(self, func, projection=""):
        out = Path("..") / self.pkg / f"{projection}_fn_{func.__name__}.go"
        code = self.generate_code(func, opname=f"{projection}Op")
        imports = ['"github.com/RyuaNerin/elliptic2/internal/curve"']
        if "field." in code:
            imports.append('"github.com/RyuaNerin/elliptic2/internal/field"')

        with out.open("w", encoding="utf-8") as f:
            f.write(f"""// Code generated by {Path(__file__).name}. DO NOT EDIT.
package {self.pkg}

import (
    {chr(10).join(imports)}
)

""")
            f.write(code)

        run_gofmt(str(out))

    def generate_code(self, func, opname: str):
        f = Formula(func, cse=self.cse)
        d = f.d
        IDX = {"x": 0, "y": 1, "z": 2, "t": 3}
        groups, consts = {}, []
        for p in f.params:
            m = re.fullmatch(r"([A-Za-z]+)(\d+)", p)
            if m and m.group(1).lower() in IDX:
                groups.setdefault(m.group(2), []).append((IDX[m.group(1).lower()], p))
            else:
                consts.append(p)
        gk = sorted(groups)
        for k in gk:
            groups[k].sort()
        # 출력 변수: 좌표 문자로 dst 인덱스 결정, 입력과 이름 충돌 시 리네임
        outidx, outname = [], []
        for o in f.outs:
            m = re.fullmatch(r"([A-Za-z]+)(\d+)", o)
            outidx.append(
                IDX[m.group(1).lower()]
                if m and m.group(1).lower() in IDX
                else len(outidx)
            )
            nm = o
            while any(nm.lower() == q.lower() for q in f.params) or nm in outname:
                nm = "d" + nm
            outname.append(nm)
        f.goname = dict(zip(f.outs, outname))
        order, reg, ntemp, opt = optimize(d, f.roots, outname)
        row_of = {}
        for r in f.rows:
            row_of.setdefault(r.node, r)
        L = []
        A = L.append
        A(
            f"func (op *{opname}) {f.name}("
            f"{', '.join(['dst'] + [f'p{k}' for k in gk])} *{self.coord})"
            f"{' *' + self.coord if f.has_bool else ''} {{"
        )

        if True:
            A("/**")
            for line in f.doc.splitlines():
                A(line.strip())
            A("")
            A(
                f"\tRegisters: {len(f.outs)} dst + {ntemp} temp"
                + (" (proven minimal)" if opt else "")
            )
            A("")
            A("| Operation")
            w = max(22, max(len(r.text) for r in f.rows))
            for r in f.rows:
                A(f"| {r.idx:>2}: {r.text:<{w}}")
            A("*/")

        A("")
        if consts:
            A(f"{', '.join(consts)} := {', '.join('op.' + c.upper() for c in consts)}")
            A("")

        for k in gk:
            v = groups[k]
            A(
                f"{', '.join(n for _, n in v)} := "
                f"{', '.join(f'&p{k}[{i}]' for i, _ in v)}"
            )
        A(f"{', '.join(outname)} := {', '.join(f'&dst[{i}]' for i in outidx)}")
        kdecl = sorted({d.extra(i) for i in f.const_nodes})
        if kdecl:
            for v in kdecl:
                A(f"{self._kname(v)} := {self._kexpr(v)} // {v}")
        if ntemp:
            A("")
            A(
                f"{', '.join(f'_{i}' for i in range(ntemp))} := "
                f"{', '.join(f'&op.t[{i}]' for i in range(ntemp))}"
            )
        A("")
        L += [b for b in self._body(f, d, order, reg, row_of)]
        if f.has_bool:
            A("")
            A("\treturn dst")
        A("}")
        return "\n".join(L)

    def _body(self, f, d, order, reg, row_of):
        mag, rootset, pend = {}, set(f.roots), []

        def ref(i):
            o = d.op(i)
            if o == "in":
                return f.goname.get(f.label[i], f.label[i])
            if o == "const":
                return self._kname(d.extra(i))
            return reg[i]

        def push(c, n=""):
            pend.append((c, n))

        for i in order:
            o, ks = d.op(i), d.kids(i)
            heavy = o in ("mul", "sqr") or (o in UNARY and UNARY[o]["heavy"])
            if heavy and not self.binary:
                for c in ks:
                    if mag.get(c, 1) > 1:
                        push(f"{ref(c)}.Reduce()")
                        mag[c] = 1
            s = [ref(c) for c in ks]
            t, k = reg[i], d.extra(i)
            r = row_of.get(i)
            note = f"{r.idx}: {r.text}" if r else ""
            add = "Add"
            sub = "Add" if self.binary else "Sub"
            if o == "mul":
                push(f"{t}.Mul({s[0]}, {s[1]})", note)
                mag[i] = 1
            elif o == "sqr":
                push(f"{t}.Sqr({s[0]})", note)
                mag[i] = 1
            elif o in UNARY:
                u = UNARY[o]
                call = f"{t}.{u['go']}({', '.join(s)})"
                push(
                    f"if {call} == nil {{ return nil }}" if u["boolean"] else call, note
                )
                mag[i] = 1 if u["heavy"] else mag.get(ks[0], 1)
            elif o in ("add", "sub"):
                push(f"{t}.{add if o == 'add' else sub}({s[0]}, {s[1]})", note)
                mag[i] = mag.get(ks[0], 1) + mag.get(ks[1], 1)
            elif o == "kmul":
                if self.binary:
                    raise ValueError(f"GF2m: constant multiply {k} unsupported")
                if k in POW2:
                    push(f"{t}.Lsh({s[0]}, {POW2[k]})", note)
                    mag[i] = mag.get(ks[0], 1) * k
                elif k == 3:
                    push(f"{t}.Lsh({s[0]}, 1)", note)
                    push(f"{t}.{add}({t}, {s[0]})")
                    mag[i] = mag.get(ks[0], 1) * 3
                else:
                    push(f"{t}.Mul({s[0]}, {k})", note)
                    mag[i] = 1
            if not self.binary:
                if i in rootset:
                    push(f"{t}.Reduce()")
                    mag[i] = 1
                elif heavy:
                    push(f"{t}.ReduceSoft()")
                    mag[i] = 1
                elif mag[i] > self.maxmag:
                    push(f"{t}.Reduce()")
                    mag[i] = 1
        w = max((len(c) for c, _ in pend), default=0)
        return [f"{c:<{w}} // {n}".rstrip() for c, n in pend]


def run_gofmt(path: str) -> None:
    gofmt = shutil.which("gofmt")
    if gofmt is None:
        return
    subprocess.run([gofmt, "-w", path], check=True)
