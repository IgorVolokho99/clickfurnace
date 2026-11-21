package worker

import "testing"

func TestWorker(t *testing.T) {
	var _ Worker = (*QueryWorker)(nil)
}
