package config

type GeneratorType string

const (
	GeneratorTypeTimestamp GeneratorType = "timestamp"
	GeneratorTypeRandom    GeneratorType = "random"
	GeneratorTypeNormal    GeneratorType = "normal"
	GeneratorTypeChoice    GeneratorType = "choice"
)
