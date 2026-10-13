package generator

import (
	"clickfurnace/config"
	"testing"
)

func TestNewRandomGenerator(t *testing.T) {
	tests := []struct {
		name    string
		min     float64
		max     float64
		wantErr bool
	}{
		{
			name:    "valid positive range",
			min:     1,
			max:     100,
			wantErr: false,
		},
		{
			name:    "valid negative range",
			min:     -100,
			max:     -1,
			wantErr: false,
		},
		{
			name:    "valid mixed range",
			min:     -100,
			max:     100,
			wantErr: false,
		},
		{
			name:    "equal boundaries",
			min:     10,
			max:     10,
			wantErr: true,
		},
		{
			name:    "inverted boundaries",
			min:     100,
			max:     10,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRandomGenerator(config.ColumnTypeUInt64, tt.min, tt.max)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"NewRandomGenerator(%v, %v) error = %v, wantErr %v",
					tt.min,
					tt.max,
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestRandomGenerator_Generate(t *testing.T) {
	const (
		min        = "10"
		max        = "20"
		iterations = 10000
	)

	g, err := NewRandomGenerator(config.ColumnTypeUInt64, min, max)
	if err != nil {
		t.Fatalf("failed to create generator: %v", err)
	}

	for i := 0; i < iterations; i++ {
		value, ok := g.Generate().(float64)
		if !ok {
			t.Fatal("Generate() returned unexpected type")
		}

		if value < min || value >= max {
			t.Fatalf(
				"generated value %v is outside [%v, %v)",
				value,
				min,
				max,
			)
		}
	}

}
