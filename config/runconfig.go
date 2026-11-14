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
	Min *yaml.Node `yaml:"min,omitempty"`
	Max *yaml.Node `yaml:"max,omitempty"`

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

	if err := c.ClickHouse.validate(); err != nil {
		return err
	}

	if err := c.Table.validate(); err != nil {
		return err
	}

	if err := c.Workload.validate(); err != nil {
		return err
	}

	return nil
}

func (c ClickHouseConfiguration) validate() error {
	if c.Host == "" {
		return fmt.Errorf("clickhouse.host must not be empty")
	}

	if c.Port == 0 {
		return fmt.Errorf("clickhouse.port must be greater than 0")
	}

	if c.Database == "" {
		return fmt.Errorf("clickhouse.database must not be empty")
	}

	return nil
}

func (w WorkloadConfiguration) validate() error {
	if w.Rate <= 0 {
		return fmt.Errorf("workload.rate must be greater than 0")
	}

	if w.BatchSize <= 0 {
		return fmt.Errorf("workload.batch_size must be greater than 0")
	}

	if w.Workers <= 0 {
		return fmt.Errorf("workload.workers must be greater than 0")
	}

	if w.Duration <= 0 {
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
			return fmt.Errorf(
				"table.columns[%d] (%s): %w",
				i,
				column.Name,
				err,
			)
		}
	}
	return nil
}

func (c ColumnConfiguration) validate() error {
	if c.Name == "" {
		return fmt.Errorf("name must not be empty")
	}

	if err := c.ColumnType.validate(); err != nil {
		return err
	}

	if err := c.validateGeneratorCompatibility(); err != nil {
		return err
	}

	return nil
}

func (c ColumnConfiguration) validateGeneratorCompatibility() error {
	switch c.Generator.Type {

	case GeneratorTypeRandom:
		if c.Generator.Min == nil {
			return fmt.Errorf("random generator requires min")
		}

		if c.Generator.Max == nil {
			return fmt.Errorf("random generator requires max")
		}
		switch c.ColumnType {
		case ColumnTypeUInt8:
			var minValue, maxValue uint8
			if err := c.Generator.Min.Decode(&minValue); err != nil {
				return fmt.Errorf("invalid min for UInt8: %w", err)
			}

			if err := c.Generator.Max.Decode(&maxValue); err != nil {
				return fmt.Errorf("invalid max for UInt8: %w", err)
			}
			if minValue >= maxValue {
				return fmt.Errorf(
					"random generator requires min < max, got min=%d max=%d",
					minValue,
					maxValue,
				)
			}
			return nil
		case ColumnTypeInt8:
			var minValue, maxValue int8
			if err := c.Generator.Min.Decode(&minValue); err != nil {
				return fmt.Errorf("invalid min for Int8: %w", err)
			}

			if err := c.Generator.Max.Decode(&maxValue); err != nil {
				return fmt.Errorf("invalid max for Int8: %w", err)
			}
			if minValue >= maxValue {
				return fmt.Errorf(
					"random generator requires min < max, got min=%d max=%d",
					minValue,
					maxValue,
				)
			}
			return nil
		}
	}
	return fmt.Errorf(
		"unsupported generator type: %w",
		c.Generator.Type,
	)
}
