package generator

import "time"

type TimestampGenerator struct{}

var _ Generator = (*TimestampGenerator)(nil)

func NewTimestampGenerator() *TimestampGenerator {
	return &TimestampGenerator{}
}

func (g *TimestampGenerator) Generate() any {
	return time.Now()
}
