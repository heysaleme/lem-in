package formatter

import (
	"fmt"
	"lem-in/internal/simulation"
)

func Print(raw []string, moves [][]simulation.Move) {
	for _, line := range raw {
		fmt.Println(line)
	}
	fmt.Println()

	for _, turn := range moves {
		fmt.Println(simulation.FormatTurn(turn))
	}
}
