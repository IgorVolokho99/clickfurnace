package main

import (
	"fmt"

	"clickfurnace/internal/generator"
)

func main() {
	var g generator.Generator

	g = generator.NewRandomIntegerGenerator(
		1,
		10,
	)

	for i := 0; i < 10; i++ {
		value := g.Generate()
		fmt.Printf("%T: %v\n", value, value)
	}
}
