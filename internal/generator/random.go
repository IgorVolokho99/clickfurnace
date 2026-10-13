package generator

import (
	"fmt"
	"math"
	"strconv"

	"clickfurnace/config"
)

type RandomGenerator struct {
	generate func() any
}

var _ Generator = (*RandomGenerator)(nil)

func NewRandomGenerator(
	columnType config.ColumnType,
	minRaw, maxRaw string,
) (*RandomGenerator, error) {
	switch columnType {
	case config.ColumnTypeUInt64:
		min, err := strconv.ParseUint(minRaw, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid min: %w", err)
		}

		max, err := strconv.ParseUint(maxRaw, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid max: %w", err)
		}

		if min > max {
			return nil, fmt.Errorf("min must be <= max")
		}

		return &RandomGenerator{
			generate: func() any {
				return randomUint64(min, max)
			},
		}, nil
	case config.ColumnTypeFloat64:
		min, err := strconv.ParseFloat(minRaw, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid min: %w", err)
		}

		max, err := strconv.ParseFloat(maxRaw, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid max: %w", err)
		}

		if math.IsNaN(min) || math.IsNaN(max) ||
			math.IsInf(min, 0) || math.IsInf(max, 0) {
			return nil, fmt.Errorf("random bounds must be finite")
		}

		if min >= max {
			return nil, fmt.Errorf("min must be < max")
		}

		if math.IsInf(max-min, 0) {
			return nil, fmt.Errorf("random range is too large")
		}

		return &RandomGenerator{
			generate: func() any {
				return randomFloat64(min, max)
			},
		}, nil
	default:
		return nil, fmt.Errorf(
			"random generator does not support column type %q",
			columnType,
		)
	}

}

func (g *RandomGenerator) Generate() any {
	return g.generate()
}
