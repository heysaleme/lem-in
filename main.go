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
		fmt.Println("ERROR: invalid data format")
		return
	}

	farm, err := parser.Parse(os.Args[1])
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	g := graph.Build(farm)

	path, err := solver.BFS(g)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	moves := simulation.Run(path, farm.Ants)

	formatter.Print(farm.RawLines, moves)
}
