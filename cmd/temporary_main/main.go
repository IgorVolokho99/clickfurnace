package main

import (
	"fmt"

	"clickfurnace/internal/generator"
)

func main() {
	var g generator.Generator

	g = generator.NewTimestampGenerator()

	value := g.Generate()

	fmt.Printf("value: %v\n", value)
	fmt.Printf("type:  %T\n", value)
}
