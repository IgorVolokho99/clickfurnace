package generator

import (
	"math/rand/v2"
	"strconv"
)

type RandomIntegerGenerator struct {
	minValue int64
	maxValue int64
}

var _ Generator = (*RandomIntegerGenerator)(nil)

func NewRandomIntegerGenerator(
	minValue, maxValue int64,
) *RandomIntegerGenerator {
	return &RandomIntegerGenerator{
		minValue: minValue,
		maxValue: maxValue,
	}
}

func (g *RandomIntegerGenerator) Generate() string {
	num := int(g.minValue) + rand.IntN(int(g.maxValue)-int(g.minValue)+1)
	return strconv.Itoa(num)
}
