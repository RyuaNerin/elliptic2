package curve

import (
	"github.com/RyuaNerin/elliptic2/internal/field"
)

type (
	Coordinate[Self any] interface {
		*Self
		Cmp(other *Self) int
		Set(other *Self)
		SetModulus(poly *field.Modulus)
	}

	GFpCoordinate  [4]field.GFp
	GF2mCoordinate [4]field.GF2m
)

func (p *GFpCoordinate) Cmp(other *GFpCoordinate) int {
	for idx := range p {
		if c := p[idx].Cmp(&other[idx]); c != 0 {
			return c
		}
	}
	return 0
}

func (p *GFpCoordinate) Set(other *GFpCoordinate) {
	for idx := range p {
		p[idx].Set(&other[idx])
	}
}

func (p *GFpCoordinate) SetModulus(poly *field.Modulus) {
	for idx := range p {
		p[idx].SetModulus(poly)
	}
}

func (p *GF2mCoordinate) Cmp(other *GF2mCoordinate) int {
	for idx := range p {
		if c := p[idx].Cmp(&other[idx]); c != 0 {
			return c
		}
	}
	return 0
}

func (p *GF2mCoordinate) Set(other *GF2mCoordinate) {
	for idx := range p {
		p[idx].Set(&other[idx])
	}
}

func (p *GF2mCoordinate) SetModulus(poly *field.Modulus) {
	for idx := range p {
		p[idx].SetModulus(poly)
	}
}
