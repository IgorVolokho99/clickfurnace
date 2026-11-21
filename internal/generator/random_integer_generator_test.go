package generator

import "testing"

func TestGenerator(t *testing.T) {
	var _ Generator = (*RandomIntegerGenerator)(nil)
}
