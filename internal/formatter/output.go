package formatter

import (
	"fmt"
)

func Print(rawLines []string, moves []string) {
	// Выводим исходные данные
	for _, line := range rawLines {
		fmt.Println(line)
	}

	// Пустая строка между картой и движением
	fmt.Println()

	// Выводим движения муравьев
	for _, move := range moves {
		fmt.Println(move)
	}
}
