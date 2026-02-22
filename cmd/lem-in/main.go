package main

import (
	"fmt"
	"os"

	"lem-in/internal/formatter"
	"lem-in/internal/graph"
	"lem-in/internal/parser"
	"lem-in/internal/simulation"
	"lem-in/internal/solver"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <filename>")
		return
	}

	// 1. Parsing / Парсинг
	farm, err := parser.Parse(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// 2. Graph Building / Построение графа
	g := graph.Build(farm)

	// 3. Solving / Поиск путей и распределение
	paths, distribution, err := solver.Solve(g, farm.Ants)
	if err != nil {
		fmt.Println("ERROR: invalid data format, no paths found")
		return
	}

	// 4. Simulation / Симуляция движений
	moves := simulation.Run(paths, distribution)

	// 5. Output / Форматированный вывод
	formatter.Print(farm.RawLines, moves)
}
