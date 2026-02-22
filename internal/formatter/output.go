// Package formatter handles the final output display of the colony and moves.
// Пакет formatter управляет финальным отображением колонии и ходов.
package formatter

import (
	"fmt"
)

// Print displays the original file content followed by the ant movement steps.
// Print выводит исходное содержание файла, а затем шаги передвижения муравьев.
func Print(rawLines []string, moves []string) {
	// Output original farm data
	// Выводим оригинальные данные фермы
	for _, line := range rawLines {
		fmt.Println(line)
	}

	// Print a newline between the farm data and the simulation results
	// Печатаем пустую строку между данными фермы и результатами симуляции
	fmt.Println()

	// Output ant moves step by step
	// Выводим ходы муравьев шаг за шагом
	for _, move := range moves {
		fmt.Println(move)
	}
}
