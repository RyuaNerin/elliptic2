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

		RequireIsOnCurve(t, curve, G.X, G.Y, "G")
		RequireXYUnmodified(t, &pSave, G, "G")
	}

	test := func(pointName string, p Point) bool {
		if p.X == nil || p.Y == nil || (p.X.Sign() == 0 && p.Y.Sign() == 0) {
			return true
		}

		pSave.Set(p)

		RequireIsOnCurve(t, curve, p.X, p.Y, pointName)
		RequireXYUnmodified(t, &pSave, &p, pointName)

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

		RequireNotIsOnCurve(t, curve, tc.X, tc.Y, fmt.Sprintf("InvalidP[%d]", idx))
		RequireXYUnmodified(t, &pSave, &tc, fmt.Sprintf("InvalidP[%d]", idx))
	}
}

func testComputeY(t *testing.T, curve elliptic2.Curve, testCase CurveTestCases) {
	var pSaveX big.Int

	c, ok := curve.(elliptic2.CurveExtended)
	if !ok {
		t.Skipf("Curve %s does not implement CurveExtended", GetName(curve))
		return
	}

	test := func(pointName string, p Point) {
		if p.X == nil || (p.X.Sign() == 0) {
			return
		}
		pSaveX.Set(p.X)

		ySmall := c.ComputeY(p.X, false)
		require.NotNil(t, ySmall, "%s: ComputeY returned nil\n  X: %s", pointName, p.X.String())

		yLarge := c.ComputeY(p.X, true)
		require.NotNil(t, yLarge, "%s: ComputeY returned nil\n  X: %s", pointName, p.X.String())

		RequireIsOnCurve(t, curve, p.X, ySmall, fmt.Sprintf("%s.Y (small)", pointName))
		RequireIsOnCurve(t, curve, p.X, yLarge, fmt.Sprintf("%s.Y (large)", pointName))

		require.LessOrEqual(t,
			ySmall.Cmp(yLarge), 0,
			"%s: small Y is bigger than large Y\n  X:  %s\n  Y1: %s\n  Y2: %s",
			pointName,
			p.X.String(),
			ySmall.String(),
			yLarge.String(),
		)
		require.True(t,
			ySmall.Cmp(p.Y) == 0 || yLarge.Cmp(p.Y) == 0,
			"%s: Computed Y mismatch\n  X:      %s\n  got Y1: %s\n  got Y2: %s\n  want Y: %s",
			pointName,
			p.X.String(),
			ySmall.String(),
			yLarge.String(),
			p.Y.String(),
		)

		RequireUnmodified(t, &pSaveX, p.X, "X")
	}

	for idx, tc := range testCase.P {
		test(fmt.Sprintf("ScalarBaseMult[%d].P", idx), tc)
	}
	for idx, tc := range testCase.Add {
		test(fmt.Sprintf("Adds[%d]", idx), tc)
	}
	for idx, tc := range testCase.Double {
		test(fmt.Sprintf("Doubles[%d]", idx), tc)
	}
	for idx, tc := range testCase.ScalarMult {
		test(fmt.Sprintf("ScalarMult[%d]", idx), tc)
	}
}

// TAdd
func TAdd(t *testing.T, curve elliptic2.Curve, testCase CurveTestCases) {
	var p1, p1Save, p2Save Point

	p1.Set(testCase.P[0])
	for idx, pWant := range testCase.Add {
		/**
		Adds[0] = Points[0] + Points[0]
		Adds[1] = Adds[0]   + Points[1]
		Adds[2] = Adds[1]   + Points[2]
		...
		*/
		p2 := testCase.P[idx]

		p1Save.Set(p1)
		p2Save.Set(p2)

		p3x, p3y := curve.Add(p1.X, p1.Y, p2.X, p2.Y)

		RequireXYEquals(t, &pWant, &Point{X: p3x, Y: p3y}, fmt.Sprintf("Add[%d]", idx))

		RequireXYUnmodified(t, &p1Save, &p1, fmt.Sprintf("Add[%d].P1", idx))
		RequireXYUnmodified(t, &p2Save, &p2, fmt.Sprintf("Add[%d].P2", idx))

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

		RequireXYEquals(t, &pWant, &Point{X: p3x, Y: p3y}, fmt.Sprintf("Double[%d]", idx))
		RequireXYUnmodified(t, &p1Save, &p1, fmt.Sprintf("Double[%d].P1", idx))
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

		RequireXYEquals(t, &pWant, &Point{X: p3x, Y: p3y}, fmt.Sprintf("ScalarMult[%d]", idx))
		RequireXYUnmodified(t, &p1Save, &p1, fmt.Sprintf("ScalarMult[%d].P1", idx))
		require.Equal(t, kSave, k, "k modified")

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

		RequireXYEquals(t, &want, &Point{X: x, Y: y}, fmt.Sprintf("ScalarBaseMult[%d]", idx))
		RequireXYUnmodified(t, &GSave, G, "G")
		require.Equal(t, kSave, k, "k modified")
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
					k := GetRandomK(t, c)
					wantX, wantY := c.ScalarBaseMult(k)

					op := newOp(cb)
					op.ToCoordinate(&coords, wantX, wantY)
					op.ToAffinePoint(&gotX, &gotY, &coords)

					RequireXYEquals(t, &Point{X: wantX, Y: wantY}, &Point{X: &gotX, Y: &gotY}, "ConvertPoint")
				}
			},
		)
	}
}

var kSamples = [][]byte{
	{},
	{0x01},
	{0x7F},
	{0x80},
	{0xFF},
	{0x55, 0xAA},
	{0x7F, 0xFF},
	{0x80, 0x00},
	{0x80, 0x01},
	{0x80, 0x80, 0x80},
	{0xAA, 0x55},
	{0xFF, 0xFF, 0xFF},
	{0xFF, 0xFF},
	{0x00, 0x00, 0x00, 0x00},
	{0x00, 0x00, 0x00, 0x01},
	{0x01, 0x02, 0x03, 0x04},
	{0x01, 0x00, 0x00, 0x00, 0x00},
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

				// hide deprecated warning
				cSimple := cSimple.(elliptic2.Curve)
				cMadd := cMadd.(elliptic2.Curve)

				for idx, k := range kSamples {
					xWant, yWant := cSimple.ScalarBaseMult(k)
					xGot, yGot := cMadd.ScalarBaseMult(k)

					RequireXYEquals(t, &Point{X: xWant, Y: yWant}, &Point{X: xGot, Y: yGot}, fmt.Sprintf("kSamples[%d]", idx))
				}

				{
					var one big.Int
					one.SetInt64(1)

					var kn big.Int
					kn.Set(GetN(t, cSimple))
					kn.Sub(&kn, &one)

					k := kn.Bytes()
					xWant, yWant := cSimple.ScalarBaseMult(k)
					xGot, yGot := cMadd.ScalarBaseMult(k)

					RequireXYEquals(t, &Point{X: xWant, Y: yWant}, &Point{X: xGot, Y: yGot}, "k_sub_1")
				}

				for range 100 {
					k := GetRandomK(t, cSimple)

					var xWant, yWant *big.Int
					var xGot, yGot *big.Int

					xWant, yWant = cSimple.ScalarBaseMult(k)
					xGot, yGot = cMadd.ScalarBaseMult(k)

					RequireXYEquals(t, &Point{X: xWant, Y: yWant}, &Point{X: xGot, Y: yGot}, "ScalarBaseMult")
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

				k := GetRandomK(t, c)
				x1, y1 := c.ScalarBaseMult(k)
				op.ToCoordinate(&p1, x1, y1)

				var gotX, gotY, wantX, wantY big.Int

				for idx, k := range kSamples {
					x2, y2 := c.ScalarBaseMult(k)
					op.ToCoordinate(&p2, x2, y2)

					op.Add(&dstAdd, &p1, &p2)
					op.Madd(&dstMadd, &p1, &p2)

					op.ToAffinePoint(&gotX, &gotY, &dstAdd)
					op.ToAffinePoint(&wantX, &wantY, &dstMadd)

					RequireXYEquals(t, &Point{X: &wantX, Y: &wantY}, &Point{X: &gotX, Y: &gotY}, fmt.Sprintf("kSamples[%d]", idx))
				}

				{
					var one big.Int
					one.SetInt64(1)

					var kn big.Int
					kn.Set(GetN(t, c))
					kn.Sub(&kn, &one)

					k := kn.Bytes()
					x2, y2 := c.ScalarBaseMult(k)
					op.ToCoordinate(&p2, x2, y2)

					op.Add(&dstAdd, &p1, &p2)
					op.Madd(&dstMadd, &p1, &p2)

					op.ToAffinePoint(&gotX, &gotY, &dstAdd)
					op.ToAffinePoint(&wantX, &wantY, &dstMadd)

					RequireXYEquals(t, &Point{X: &wantX, Y: &wantY}, &Point{X: &gotX, Y: &gotY}, "k_sub_1")
				}

				for range 100 {
					k = GetRandomK(t, c)

					x2, y2 := c.ScalarBaseMult(k)
					op.ToCoordinate(&p2, x2, y2)

					op.Add(&dstAdd, &p1, &p2)
					op.Madd(&dstMadd, &p1, &p2)

					op.ToAffinePoint(&gotX, &gotY, &dstAdd)
					op.ToAffinePoint(&wantX, &wantY, &dstMadd)

					RequireXYEquals(t, &Point{X: &wantX, Y: &wantY}, &Point{X: &gotX, Y: &gotY}, "Madd")
				}
			},
		)
	}
}
