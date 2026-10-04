package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version    string                  `yaml:"version"`
	ClickHouse ClickHouseConfiguration `yaml:"clickhouse"`
	Table      TableConfiguration      `yaml:"table"`
	Workload   WorkloadConfiguration   `yaml:"workload"`
}

type ClickHouseConfiguration struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
}

type TableConfiguration struct {
	Name    string                `yaml:"name"`
	Columns []ColumnConfiguration `yaml:"columns"`
}

type ColumnConfiguration struct {
	Name       string                 `yaml:"name"`
	ColumnType ColumnType             `yaml:"type"`
	Generator  GeneratorConfiguration `yaml:"generator"`
}

type GeneratorConfiguration struct {
	Type GeneratorType `yaml:"type"`

	// random
	Min *float64 `yaml:"min,omitempty"`
	Max *float64 `yaml:"max,omitempty"`

	// normal
	Mean   *float64 `yaml:"mean,omitempty"`
	Stddev *float64 `yaml:"stddev,omitempty"`

	// choice
	Values []string `yaml:"values,omitempty"`
}

type WorkloadConfiguration struct {
	Rate      int           `yaml:"rate"`
	BatchSize int           `yaml:"batch_size"`
	Workers   int           `yaml:"workers"`
	Duration  time.Duration `yaml:"-"`
}

type rawWorkloadConfiguration struct {
	Rate      int    `yaml:"rate"`
	BatchSize int    `yaml:"batch_size"`
	Workers   int    `yaml:"workers"`
	Duration  string `yaml:"duration"`
}

type rawConfig struct {
	Version    string                   `yaml:"version"`
	ClickHouse ClickHouseConfiguration  `yaml:"clickhouse"`
	Table      TableConfiguration       `yaml:"table"`
	Workload   rawWorkloadConfiguration `yaml:"workload"`
}

func New(configName string) (Config, error) {
	data, err := os.ReadFile(configName)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var raw rawConfig

	if err := yaml.Unmarshal(data, &raw); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	duration, err := time.ParseDuration(raw.Workload.Duration)
	if err != nil {
		return Config{}, fmt.Errorf(
			"invalid workload.duration %q: %w",
			raw.Workload.Duration,
			err,
		)
	}

	cfg := Config{
		Version:    raw.Version,
		ClickHouse: raw.ClickHouse,
		Table:      raw.Table,
		Workload: WorkloadConfiguration{
			Rate:      raw.Workload.Rate,
			BatchSize: raw.Workload.BatchSize,
			Workers:   raw.Workload.Workers,
			Duration:  duration,
		},
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) validate() error {
	if c.Version != "1" {
		return fmt.Errorf("unsupported config version: %q", c.Version)
	}

	if c.ClickHouse.Host == "" {
		return fmt.Errorf("clickhouse.host must not be empty")
	}

	if c.ClickHouse.Port == 0 { // todo how to check port correctly
		return fmt.Errorf("clickhouse.port must be greater than 0")
	}

	if c.ClickHouse.Database == "" {
		return fmt.Errorf("clickhouse.database must not be empty")
	}

	if c.Table.Name == "" {
		return fmt.Errorf("table.name must not be empty")
	}

	if len(c.Table.Columns) == 0 {
		return fmt.Errorf("table.columns must not be empty")
	}

	if c.Workload.Rate <= 0 {
		return fmt.Errorf("workload.rate must be greater than 0")
	}

	if c.Workload.BatchSize <= 0 {
		return fmt.Errorf("workload.batch_size must be greater than 0")
	}

	if c.Workload.Workers <= 0 {
		return fmt.Errorf("workload.workers must be greater than 0")
	}

	if c.Workload.Duration <= 0 {
		return fmt.Errorf("workload.duration must be greater than 0")
	}

	return nil
}

func (t TableConfiguration) validate() error {
	if t.Name == "" {
		return fmt.Errorf("table.name must not be empty")
	}

	if len(t.Columns) == 0 {
		return fmt.Errorf("table.columns must not be empty")
	}

	for i, column := range t.Columns {
		if err := column.validate(); err != nil {
			return fmt.Errorf("table.columns[%d]: %w", i, err)
		}
	}
	return nil
}

func (c ColumnConfiguration) validate() error {
	if c.Name == "" {
		return fmt.Errorf("name must not be empty")
	}

	if c.ColumnType == "" {
		return fmt.Errorf("type must not be empty")
	}

	if err := c.Generator.validate(); err != nil {
		return fmt.Errorf("generator: %w", err)
	}

	if err := c.validateGeneratorCompatibility(); err != nil {
		return nil
	}

	return nil
}

func (c ColumnConfiguration) validateGeneratorCompatibility() error {
	switch c.Generator.Type {
	case GeneratorTypeTimestamp:
		switch c.ColumnType {
		case ColumnTypeDateTime, ColumnTypeDateTime64:
			return nil
		}

	case GeneratorTypeRandom:
		switch c.ColumnType {
		case ColumnTypeUInt8,
			ColumnTypeUInt16,
			ColumnTypeUInt32,
			ColumnTypeUInt64,
			ColumnTypeInt8,
			ColumnTypeInt16,
			ColumnTypeInt32,
			ColumnTypeInt64,
			ColumnTypeFloat32,
			ColumnTypeFloat64:
			return nil
		}

	case GeneratorTypeNormal:
		switch c.ColumnType {
		case ColumnTypeFloat32, ColumnTypeFloat64:
			return nil
		}

	case GeneratorTypeChoice:
		if c.ColumnType == ColumnTypeString {
			return nil
		}
	}

	return fmt.Errorf(
		"generator %q is not compatible with column type %q",
		c.Generator.Type,
		c.ColumnType,
	)
}

func (g GeneratorConfiguration) validate() error {
	switch g.Type {
	case "timestamp":
		return nil
	case "random":
		return g.validateRandom()
	case "normal":
		return g.validateNormal()
	case "choice":
		return g.validateChoice()
	default:
		return fmt.Errorf("unsupported generator time: %q", g.Type)
	}
}

func (g GeneratorConfiguration) validateRandom() error {
	if g.Min == nil {
		return fmt.Errorf("random generator requires min")
	}

	if g.Max == nil {
		return fmt.Errorf("random generator requires max")
	}

	if *g.Min >= *g.Max {
		return fmt.Errorf(
			"random generator requires min < max, got min=%v max=%v",
			*g.Min,
			*g.Max,
		)
	}

	return nil
}

func (g GeneratorConfiguration) validateNormal() error {
	if g.Mean == nil {
		return fmt.Errorf("normal generator requires mean")
	}

	if g.Stddev == nil {
		return fmt.Errorf("normal generator requires stddev")
	}

	if *g.Stddev <= 0 {
		return fmt.Errorf(
			"normal generator requires stddev > 0, got %v",
			*g.Stddev,
		)
	}

	return nil
}

func (g GeneratorConfiguration) validateChoice() error {
	if len(g.Values) == 0 {
		return fmt.Errorf("choice generator requires at least one value")
	}

	return nil
}
