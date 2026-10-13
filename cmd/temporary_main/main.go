package main

import (
	"fmt"

	"clickfurnace/config"
	"clickfurnace/internal/generator"
)

func main() {
	var g generator.Generator

	g, err := generator.NewRandomGenerator(config.ColumnTypeUInt64, "1", "20")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	value := g.Generate()

	fmt.Printf("value: %v\n", value)
	fmt.Printf("type:  %T\n", value)
}
