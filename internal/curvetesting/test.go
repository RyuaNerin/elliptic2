package curvetesting

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/RyuaNerin/elliptic2"
	"github.com/RyuaNerin/elliptic2/internal"
	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/stretchr/testify/require"
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

		require.True(t, curve.IsOnCurve(G.X, G.Y), "Generator is not on curve\nX: %s\nY: %s", G.X.Text(16), G.Y.Text(16))

		require.True(t, pSave.X.Cmp(G.X) == 0, "G.X modified")
		require.True(t, pSave.Y.Cmp(G.Y) == 0, "G.Y modified")
	}

	test := func(pointName string, p Point) bool {
		if p.X == nil || p.Y == nil || (p.X.Sign() == 0 && p.Y.Sign() == 0) {
			return true
		}

		pSave.Set(p)

		require.True(t, curve.IsOnCurve(p.X, p.Y), "%s is not on curve\nX: %s\nY: %s", pointName, p.X.Text(16), p.Y.Text(16))

		require.True(t, pSave.X.Cmp(p.X) == 0, "x modified")
		require.True(t, pSave.Y.Cmp(p.Y) == 0, "y modified")

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

		require.False(t, curve.IsOnCurve(tc.X, tc.Y), "InvalidP[%d] is on curve\nX: %s\nY: %s", idx, tc.X.Text(16), tc.Y.Text(16))

		require.True(t, pSave.X.Cmp(tc.X) == 0, "x modified")
		require.True(t, pSave.Y.Cmp(tc.Y) == 0, "y modified")
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
		require.NotNil(t, ySmall, "%s: ComputeY returned nil\nX: %s", pointName, p.X.Text(16))

		yLarge := c.ComputeY(p.X, true)
		require.NotNil(t, yLarge, "%s: ComputeY returned nil\nX: %s", pointName, p.X.Text(16))

		require.True(t, c.IsOnCurve(p.X, ySmall), "%s: Computed Y (small) is not on curve\nX: %s\nY: %s", pointName, p.X.Text(16), ySmall.Text(16))
		require.True(t, c.IsOnCurve(p.X, yLarge), "%s: Computed Y (large) is not on curve\nX: %s\nY: %s", pointName, p.X.Text(16), yLarge.Text(16))

		require.True(t,
			ySmall.Cmp(yLarge) <= 0,
			"%s: small Y is bigger than large Y\nX: %s\nY1: %s\nY2: %s",
			pointName,
			p.X.Text(16),
			ySmall.Text(16),
			yLarge.Text(16),
		)

		require.True(t,
			ySmall.Cmp(p.Y) == 0 || yLarge.Cmp(p.Y) == 0,
			"%s: Computed Y mismatch\nX:      %s\ngot Y1: %s\ngot Y2: %s\nwant Y: %s",
			pointName,
			p.X.Text(16),
			ySmall.Text(16),
			yLarge.Text(16),
			p.Y.Text(16),
		)

		require.True(t, pSaveX.Cmp(p.X) == 0, "x modified")

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

		require.True(t, p3x.Cmp(pWant.X) == 0, "Add[%d]: Add X mismatch:\ngot:  %s\nwant: %s", idx, p3x.Text(16), pWant.X.Text(16))
		require.True(t, p3y.Cmp(pWant.Y) == 0, "Add[%d]: Add Y mismatch:\ngot:  %s\nwant: %s", idx, p3y.Text(16), pWant.Y.Text(16))

		require.True(t, p1.X.Cmp(p1s.X) == 0, "x1 modified")
		require.True(t, p1.Y.Cmp(p1s.Y) == 0, "y1 modified")
		require.True(t, p2.X.Cmp(p2s.X) == 0, "x2 modified")
		require.True(t, p2.Y.Cmp(p2s.Y) == 0, "y2 modified")

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

		require.True(t, p3x.Cmp(pWant.X) == 0, "Double[%d]: X mismatch:\ngot:  %s\nwant: %s", idx, p3x.Text(16), pWant.X.Text(16))
		require.True(t, p3y.Cmp(pWant.Y) == 0, "Double[%d]: Y mismatch:\ngot:  %s\nwant: %s", idx, p3y.Text(16), pWant.Y.Text(16))

		require.True(t, p1.X.Cmp(p1Save.X) == 0, "x1 modified")
		require.True(t, p1.Y.Cmp(p1Save.Y) == 0, "y1 modified")
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

		require.True(t, p3x.Cmp(pWant.X) == 0, "ScalarMult[%d]: Add X mismatch:\ngot:  %s\nwant: %s", idx, p3x.Text(16), pWant.X.Text(16))
		require.True(t, p3y.Cmp(pWant.Y) == 0, "ScalarMult[%d]: Add Y mismatch:\ngot:  %s\nwant: %s", idx, p3y.Text(16), pWant.Y.Text(16))

		require.True(t, p1.X.Cmp(p1Save.X) == 0, "x1 modified")
		require.True(t, p1.Y.Cmp(p1Save.Y) == 0, "y1 modified")
		require.Equal(t, k, kSave, "k modified")

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

		require.True(t, x.Cmp(want.X) == 0, "ScalarBaseMult[%d]: X mismatch:\ngot:  %s\nwant: %s", idx, x.Text(16), want.X.Text(16))
		require.True(t, y.Cmp(want.Y) == 0, "ScalarBaseMult[%d]: Y mismatch:\ngot:  %s\nwant: %s", idx, y.Text(16), want.Y.Text(16))

		require.True(t, GSave.X.Cmp(G.X) == 0, "G.X modified")
		require.True(t, GSave.Y.Cmp(G.Y) == 0, "G.Y modified")
		require.Equal(t, k, kSave, "k modified")
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
		t.Run(
			fmt.Sprintf("%s/%s", GetCurveType(c), GetName(c)),
			func(t *testing.T) {
				cb := curve.GetBase(c)

				modulus := cb.Modulus()

				var coords C
				CP(&coords).SetModulus(modulus)

				var gotX, gotY big.Int

				for range 100 {
					k := GetRandomK(c)
					wantX, wantY := c.ScalarBaseMult(k)

					op := newOp(cb)
					op.ToCoordinate(&coords, wantX, wantY)
					op.ToAffinePoint(&gotX, &gotY, &coords)

					require.True(t, wantX.Cmp(&gotX) == 0, "ConvertPoint X mismatch:\ngot:  %s\nwant: %s", gotX.String(), wantX.String())
					require.True(t, wantY.Cmp(&gotY) == 0, "ConvertPoint Y mismatch:\ngot:  %s\nwant: %s", gotY.String(), wantY.String())
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

				for range 100 {
					k := GetRandomK(cSimple)

					var wantx1, wanty1 *big.Int
					var gotx1, goty1 *big.Int

					wantx1, wanty1 = cSimple.ScalarBaseMult(k)
					gotx1, goty1 = cMadd.ScalarBaseMult(k)

					require.True(t, wantx1.Cmp(gotx1) == 0, "ScalarBaseMult X mismatch:\ngot:  %s\nwant: %s", gotx1.Text(16), wantx1.Text(16))
					require.True(t, wanty1.Cmp(goty1) == 0, "ScalarBaseMult Y mismatch:\ngot:  %s\nwant: %s", goty1.Text(16), wanty1.Text(16))
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
				require.True(t, ok, "Curve %s does not support MaddOperator", GetName(c))

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

				for range 100 {
					k[1] ^= k[0]
					k[0] += 1

					x2, y2 := c.ScalarBaseMult(k)
					op.ToCoordinate(&p2, x2, y2)

					op.Add(&dstAdd, &p1, &p2)
					op.Madd(&dstMadd, &p1, &p2)

					op.ToAffinePoint(&gotX, &gotY, &dstAdd)
					op.ToAffinePoint(&wantX, &wantY, &dstMadd)

					require.True(t, gotX.Cmp(&wantX) == 0 && gotY.Cmp(&wantY) == 0, "Madd failed:\ngot:  (%s, %s)\nwant: (%s, %s)", gotX.String(), gotY.String(), wantX.String(), wantY.String())
				}
			},
		)
	}
}
