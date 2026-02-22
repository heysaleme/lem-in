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
	// Проверка аргументов
	if len(os.Args) != 2 {
		fmt.Println("ERROR: invalid data format")
		return
	}

	// Парсим файл
	farm, err := parser.Parse(os.Args[1])
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	// Строим граф
	g := graph.Build(farm)

	// Находим пути + распределение муравьёв
	// Больше не передаем farm, только граф и количество муравьев
	paths, assign, err := solver.FindAllPaths(g, farm.Ants)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	// Симуляция
	moves := simulation.Run(paths, assign)

	// Вывод
	formatter.Print(farm.RawLines, moves)
}
