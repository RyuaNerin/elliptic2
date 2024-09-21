package curvetesting

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/RyuaNerin/elliptic2"
	"github.com/RyuaNerin/elliptic2/internal"
	"github.com/RyuaNerin/elliptic2/internal/curve"
)

type Point struct {
	X, Y *big.Int
}

// ParsePoint
func P(xHex, yHex string) Point {
	return Point{
		X: internal.HI(xHex),
		Y: internal.HI(yHex),
	}
}

func (p *Point) Set(other Point) {
	if p.X == nil {
		p.X = new(big.Int)
	}
	if p.Y == nil {
		p.Y = new(big.Int)
	}
	p.X.Set(other.X)
	p.Y.Set(other.Y)
}

type CurveTestCases struct {
	Curve      elliptic2.Curve
	K          [][]byte
	P          []Point // P[0] = G * k[0], P[1] = P[0] * k[1], ...
	Add        []Point // Adds[0] = P[0] + P[0], Adds[1] = Adds[0] + P[1], ...
	Double     []Point // Double[0] = 2P[0], Double[1] = 2P[1], ...
	ScalarMult []Point // ScalarMult[0] = P[0] * k[0], ScalarMult[1] = ScalarMult[0] * k[1], ...
	InvalidP   []Point
}

type ScalarBaseMultTestCase struct {
	K []byte
	P Point
}

func TC(t *testing.T, curves []CurveTestCases, f func(*testing.T, elliptic2.Curve, CurveTestCases)) {
	for _, curveTests := range curves {
		curveTests := curveTests

		t.Run(
			fmt.Sprintf("%s/%s", GetCurveType(curveTests.Curve), GetName(curveTests.Curve)),
			func(t *testing.T) {
				f(t, curveTests.Curve, curveTests)
			},
		)
	}
}

func TcOp[
	TCurveParams any,
	TCurveArithmetic curve.CurveArithmeticBase,
	TOperator any,
](
	t *testing.T,
	fnBuild func(params *TCurveParams, fnNewOp func(c TCurveArithmetic) TOperator) TCurveArithmetic,
	fnNewOp func(c TCurveArithmetic) TOperator,
	fnTest func(t *testing.T, curve elliptic2.Curve, testCases CurveTestCases),
	tcs ...CurveTestCases,
) {
	for _, tc := range tcs {
		base := curve.GetBase(tc.Curve)
		c := curve.NewCurve(fnBuild(base.RawParams().(*TCurveParams), fnNewOp))

		t.Run(GetName(c), func(t *testing.T) { fnTest(t, c, tc) })
	}
}

// Test IsOnCurve
func TcIsOnCurve(t *testing.T, curveTestCases []CurveTestCases) {
	for _, curveTests := range curveTestCases {
		curveTests := curveTests

		t.Run(
			fmt.Sprintf("%s/%s", GetCurveType(curveTests.Curve), GetName(curveTests.Curve)),
			func(t *testing.T) {
				testIsOnCurve(t, curveTests.Curve, curveTests)
			},
		)
	}
}

// Test ComputeY
func TcComputeY(t *testing.T, curveTestCases []CurveTestCases) {
	for _, curveTests := range curveTestCases {
		curveTests := curveTests

		t.Run(
			fmt.Sprintf("%s/%s", GetCurveType(curveTests.Curve), GetName(curveTests.Curve)),
			func(t *testing.T) {
				testComputeY(t, curveTests.Curve, curveTests)
			},
		)
	}
}

func testIsOnCurve(t *testing.T, curve elliptic2.Curve, testCase CurveTestCases) {
	var pSave Point
	G := GetGenerator(curve)

	if G != nil {
		pSave.Set(*G)
		if !curve.IsOnCurve(G.X, G.Y) {
			t.Errorf("G is not on curve\nX: %s\nY: %s", G.X.Text(16), G.Y.Text(16))
			return
		}
		if pSave.X.Cmp(G.X) != 0 {
			t.Errorf("G.X modified")
			return
		}
		if pSave.Y.Cmp(G.Y) != 0 {
			t.Errorf("G.Y modified")
			return
		}
	}

	test := func(pointName string, p Point) bool {
		if p.X == nil || p.Y == nil || (p.X.Sign() == 0 && p.Y.Sign() == 0) {
			return true
		}

		pSave.Set(p)

		if !curve.IsOnCurve(p.X, p.Y) {
			t.Errorf("%s is not on curve\nX: %s\nY: %s", pointName, p.X.Text(16), p.Y.Text(16))
			return false
		}

		if pSave.X.Cmp(p.X) != 0 {
			t.Errorf("x modified")
			return false
		}
		if pSave.Y.Cmp(p.Y) != 0 {
			t.Errorf("y modified")
			return false
		}
		return true
	}

	// test points
	for idx, tc := range testCase.P {
		if !test(fmt.Sprintf("P[%d]", idx), tc) {
			return
		}
	}
	// Adds
	for idx, tc := range testCase.Add {
		if !test(fmt.Sprintf("Adds[%d]", idx), tc) {
			return
		}
	}
	// Doubles
	for idx, tc := range testCase.Double {
		if !test(fmt.Sprintf("Doubles[%d]", idx), tc) {
			return
		}
	}
	// ScalarMult
	for idx, tc := range testCase.ScalarMult {
		if !test(fmt.Sprintf("ScalarMult[%d]", idx), tc) {
			return
		}
	}

	for idx, tc := range testCase.InvalidP {
		pSave.Set(tc)

		if curve.IsOnCurve(tc.X, tc.Y) {
			t.Errorf("InvalidP[%d] is on curve\nX: %s\nY: %s", idx, tc.X.Text(16), tc.Y.Text(16))
			return
		}

		if pSave.X.Cmp(tc.X) != 0 {
			t.Errorf("x modified")
			return
		}
		if pSave.Y.Cmp(tc.Y) != 0 {
			t.Errorf("y modified")
			return
		}
	}
}

func testComputeY(t *testing.T, curve elliptic2.Curve, testCase CurveTestCases) {
	var pSaveX big.Int

	c, ok := curve.(elliptic2.CurveExtended)
	if !ok {
		t.Skipf("Curve %s does not implement CurveExtended", GetName(curve))
		return
	}

	test := func(pointName string, p Point) bool {
		if p.X == nil || (p.X.Sign() == 0) {
			return true
		}
		pSaveX.Set(p.X)

		ySmall := c.ComputeY(p.X, false)
		if ySmall == nil {
			t.Errorf("%s: ComputeY returned nil\nX: %s", pointName, p.X.Text(16))
			return false
		}
		yLarge := c.ComputeY(p.X, true)
		if yLarge == nil {
			t.Errorf("%s: ComputeY returned nil\nX: %s", pointName, p.X.Text(16))
			return false
		}

		if !c.IsOnCurve(p.X, ySmall) {
			t.Errorf("%s: Computed Y (small) is not on curve\nX: %s\nY: %s", pointName, p.X.Text(16), ySmall.Text(16))
			return false
		}
		if !c.IsOnCurve(p.X, yLarge) {
			t.Errorf("%s: Computed Y (large) is not on curve\nX: %s\nY: %s", pointName, p.X.Text(16), yLarge.Text(16))
			return false
		}
		if ySmall.Cmp(yLarge) != 0 && ySmall.Cmp(yLarge) != -1 {
			t.Errorf("%s: small Y is bigger than large Y\nX: %s\nY1: %s\nY2: %s",
				pointName,
				p.X.Text(16),
				ySmall.Text(16),
				yLarge.Text(16),
			)
			return false
		}

		if ySmall.Cmp(p.Y) != 0 && yLarge.Cmp(p.Y) != 0 {
			t.Errorf("%s: Computed Y mismatch\nX:      %s\ngot Y1: %s\ngot Y2: %s\nwant Y: %s",
				pointName,
				p.X.Text(16),
				ySmall.Text(16),
				yLarge.Text(16),
				p.Y.Text(16),
			)
			return false
		}

		if pSaveX.Cmp(p.X) != 0 {
			t.Errorf("x modified")
			return false
		}

		return true
	}

	for idx, tc := range testCase.P {
		if !test(fmt.Sprintf("ScalarBaseMult[%d].P", idx), tc) {
			return
		}
	}
	for idx, tc := range testCase.Add {
		if !test(fmt.Sprintf("Adds[%d]", idx), tc) {
			return
		}
	}
	for idx, tc := range testCase.Double {
		if !test(fmt.Sprintf("Doubles[%d]", idx), tc) {
			return
		}
	}
	for idx, tc := range testCase.ScalarMult {
		if !test(fmt.Sprintf("ScalarMult[%d]", idx), tc) {
			return
		}
	}
}

// TAdd
func TAdd(t *testing.T, curve elliptic2.Curve, testCase CurveTestCases) {
	var p1, p1s, p2s Point

	p1.Set(testCase.P[0])
	for idx, pWant := range testCase.Add {
		/**
		Adds[0] = Points[0] + Points[0]
		Adds[1] = Adds[0]   + Points[1]
		Adds[2] = Adds[1]   + Points[2]
		...
		*/
		p2 := testCase.P[idx]

		p1s.Set(p1)
		p2s.Set(p2)

		p3x, p3y := curve.Add(p1.X, p1.Y, p2.X, p2.Y)
		if p3x.Cmp(pWant.X) != 0 {
			t.Errorf("Add[%d]: Add X mismatch:\ngot:  %s\nwant: %s", idx, p3x.Text(16), pWant.X.Text(16))
			return
		}
		if p3y.Cmp(pWant.Y) != 0 {
			t.Errorf("Add[%d]: Add Y mismatch:\ngot:  %s\nwant: %s", idx, p3y.Text(16), pWant.Y.Text(16))
			return
		}

		if p1.X.Cmp(p1s.X) != 0 {
			t.Errorf("x1 modified")
			return
		}
		if p1.Y.Cmp(p1s.Y) != 0 {
			t.Errorf("y1 modified")
			return
		}
		if p2.X.Cmp(p2s.X) != 0 {
			t.Errorf("x2 modified")
			return
		}
		if p2.Y.Cmp(p2s.Y) != 0 {
			t.Errorf("y2 modified")
			return
		}

		p1.Set(pWant)
	}
}

// Double
func TDouble(t *testing.T, curve elliptic2.Curve, testCase CurveTestCases) {
	var p1Save Point

	for idx, pWant := range testCase.Double {
		/**
		Doubles[0] = Points[0] doubled
		Doubles[1] = Points[1] doubled
		...
		*/
		p1 := testCase.P[idx]
		p1Save.Set(p1)

		p3x, p3y := curve.Double(p1.X, p1.Y)
		if p3x.Cmp(pWant.X) != 0 {
			t.Errorf("Double[%d]: X mismatch:\ngot:  %s\nwant: %s", idx, p3x.Text(16), pWant.X.Text(16))
			return
		}
		if p3y.Cmp(pWant.Y) != 0 {
			t.Errorf("Double[%d]: Y mismatch:\ngot:  %s\nwant: %s", idx, p3y.Text(16), pWant.Y.Text(16))
			return
		}

		if p1.X.Cmp(p1Save.X) != 0 {
			t.Errorf("x1 modified")
			return
		}
		if p1.Y.Cmp(p1Save.Y) != 0 {
			t.Errorf("y1 modified")
			return
		}
	}
}

// ScalarMult
func TMult(t *testing.T, curve elliptic2.Curve, testCase CurveTestCases) {
	var p1, p1Save Point
	var kSave []byte

	p1.Set(testCase.P[0])
	for idx, pWant := range testCase.ScalarMult {
		/**
		p3			  = p1            * k
		ScalarMult[0] = G             * k[0]
		ScalarMult[1] = ScalarMult[0] * k[1]
		ScalarMult[2] = ScalarMult[1] * k[2]
		...
		*/
		k := testCase.K[idx]

		if cap(kSave) < len(k) {
			kSave = make([]byte, len(k))
		}
		kSave = kSave[:len(k)]
		copy(kSave, k)
		p1Save.Set(p1)

		p3x, p3y := curve.ScalarMult(p1.X, p1.Y, k)
		if p3x.Cmp(pWant.X) != 0 {
			t.Errorf("ScalarMult[%d]: Add X mismatch:\ngot:  %s\nwant: %s", idx, p3x.Text(16), pWant.X.Text(16))
			return
		}
		if p3y.Cmp(pWant.Y) != 0 {
			t.Errorf("ScalarMult[%d]: Add Y mismatch:\ngot:  %s\nwant: %s", idx, p3y.Text(16), pWant.Y.Text(16))
			return
		}

		if p1.X.Cmp(p1Save.X) != 0 {
			t.Errorf("x1 modified")
			return
		}
		if p1.Y.Cmp(p1Save.Y) != 0 {
			t.Errorf("y1 modified")
			return
		}
		if !internal.Equals(k, k) {
			t.Errorf("k modified")
			return
		}

		p1.X.Set(pWant.X)
		p1.Y.Set(pWant.Y)
	}
}

// ScalarBaseMult
func TBaseMult(t *testing.T, curve elliptic2.Curve, testCase CurveTestCases) {
	RequireGenerator(t, curve)

	G := GetGenerator(curve)

	var GSave Point
	GSave.Set(Point{X: G.X, Y: G.Y})

	var kSave []byte

	for idx, want := range testCase.P {
		k := testCase.K[idx]

		if cap(kSave) < len(k) {
			kSave = make([]byte, len(k))
		}
		kSave = kSave[:len(k)]
		copy(kSave, k)

		x, y := curve.ScalarBaseMult(k)
		if x.Cmp(want.X) != 0 {
			t.Errorf("ScalarBaseMult[%d]: X mismatch:\ngot:  %s\nwant: %s", idx, x.Text(16), want.X.Text(16))
			return
		}
		if y.Cmp(want.Y) != 0 {
			t.Errorf("ScalarBaseMult[%d]: Y mismatch:\ngot:  %s\nwant: %s", idx, y.Text(16), want.Y.Text(16))
			return
		}

		if GSave.X.Cmp(G.X) != 0 {
			t.Errorf("G.X modified")
			return
		}
		if GSave.Y.Cmp(G.Y) != 0 {
			t.Errorf("G.Y modified")
			return
		}
		if !internal.Equals(k, kSave) {
			t.Errorf("k modified")
			return
		}
	}
}

func TestCoordinates[
	OpType curve.Operator[C, CP],
	C any,
	CP curve.Coordinate[C],
](
	t *testing.T,
	newOp func(c curve.CurveArithmeticBase) OpType,
	curves ...elliptic2.Curve,
) {
	for _, c := range curves {
		c := c
		t.Run(
			fmt.Sprintf("%s/%s", GetCurveType(c), GetName(c)),
			func(t *testing.T) {
				cb := curve.GetBase(c)

				modulus := cb.Modulus()

				var coords C
				CP(&coords).SetModulus(modulus)

				var gotX, gotY big.Int

				for i := 0; i < 100; i++ {
					k := GetRandomK(c)
					wantX, wantY := c.ScalarBaseMult(k)

					op := newOp(cb)
					op.ToCoordinate(&coords, wantX, wantY)
					op.ToAffinePoint(&gotX, &gotY, &coords)

					if wantX.Cmp(&gotX) != 0 || wantY.Cmp(&gotY) != 0 {
						t.Fatalf("ConvertPoint failed:\ngot:  (%s, %s)\nwant: (%s, %s)",
							gotX.String(), gotY.String(),
							wantX.String(), wantY.String(),
						)
						return
					}
				}
			},
		)
	}
}

func TestCurveMadd(t *testing.T, curves ...curve.CurveArithmeticBase) {
	for _, cb := range curves {
		cSimple := curve.NewCurveSimple(cb)
		cMadd := curve.NewCurveMadd(cb)

		t.Run(
			fmt.Sprintf("%s/%s", GetCurveType(cSimple), GetName(cSimple)),
			func(t *testing.T) {
				if cMadd == nil {
					t.Skipf("Curve %s does not support Madd", GetName(cSimple))
					return
				}

				// for hide deprecated warning
				cSimple := cSimple.(elliptic2.Curve)
				cMadd := cMadd.(elliptic2.Curve)

				for i := 0; i < 100; i++ {
					k := GetRandomK(cSimple)

					var wantx1, wanty1 *big.Int
					var gotx1, goty1 *big.Int

					wantx1, wanty1 = cSimple.ScalarBaseMult(k)
					gotx1, goty1 = cMadd.ScalarBaseMult(k)

					if wantx1.Cmp(gotx1) != 0 || wanty1.Cmp(goty1) != 0 {
						t.Fatalf("ScalarBaseMult mismatch:\ngot:  (%s, %s)\nwant: (%s, %s)",
							gotx1.Text(16), goty1.Text(16),
							wantx1.Text(16), wanty1.Text(16),
						)
					}
				}
			},
		)
	}
}

func TestMadd[
	MaddOpType curve.MaddOperator[C, CP],
	OpType curve.Operator[C, CP],
	C any,
	CP curve.Coordinate[C],
](
	t *testing.T,
	newOp func(c curve.CurveArithmeticBase) OpType,
	curves ...elliptic2.Curve,
) {
	for _, c := range curves {
		c := c
		t.Run(
			fmt.Sprintf("%s/%s", GetCurveType(c), GetName(c)),
			func(t *testing.T) {
				cb := curve.GetBase(c)
				op, ok := any(newOp(cb)).(MaddOpType)
				if !ok {
					t.Fatalf("Curve %s does not support GFpExtendedOperator", GetName(c))
					return
				}

				modulus := cb.Modulus()

				var p1, p2, dstAdd, dstMadd C
				CP(&p1).SetModulus(modulus)
				CP(&p2).SetModulus(modulus)
				CP(&dstAdd).SetModulus(modulus)
				CP(&dstMadd).SetModulus(modulus)

				k := GetRandomK(c)
				x1, y1 := c.ScalarBaseMult(k)
				op.ToCoordinate(&p1, x1, y1)

				var gotX, gotY, wantX, wantY big.Int

				for i := 0; i < 100; i++ {
					k[1] ^= k[0]
					k[0] += 1

					x2, y2 := c.ScalarBaseMult(k)
					op.ToCoordinate(&p2, x2, y2)

					op.Add(&dstAdd, &p1, &p2)
					op.Madd(&dstMadd, &p1, &p2)

					op.ToAffinePoint(&gotX, &gotY, &dstAdd)
					op.ToAffinePoint(&wantX, &wantY, &dstMadd)

					if gotX.Cmp(&wantX) != 0 || gotY.Cmp(&wantY) != 0 {
						t.Fatalf("Madd failed:\ngot:  (%s, %s)\nwant: (%s, %s)",
							gotX.String(), gotY.String(),
							wantX.String(), wantY.String(),
						)
						return
					}
				}
			},
		)
	}
}
