package worker

import (
	"clickfurnace/config"
)

var _ Worker = (*QueryWorker)(nil)

type QueryWorker struct {
	cfg config.Config
}

func NewQueryWorker(cfg config.Config) *QueryWorker {
	return &QueryWorker{
		cfg: cfg,
	}
}

func (q *QueryWorker) GenerateQuery() string {
	// Function needs to generate query using cfg.
	return ""
}
