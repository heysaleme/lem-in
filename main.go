package main

import "os"

func main() {
	input := parser.Parse(os.Args[1])
	graph := graph.Build(input)
	paths := solver.FindOptimalPaths(graph)
	moves := simulation.Run(paths, input.Ants)
	formatter.Print(input, moves)
}
