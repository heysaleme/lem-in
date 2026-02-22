package simulation

import (
	"fmt"
	"lem-in/internal/solver"
	"sort"
	"strings"
)

type Ant struct {
	ID        int
	PathIndex int
	Position  int
	Path      []string
}

func Run(paths []solver.Path, distribution [][]int) []string {
	var moves []string

	// Создаем муравьев
	ants := make([]*Ant, 0)
	for pathIdx, antIDs := range distribution {
		for _, antID := range antIDs {
			ants = append(ants, &Ant{
				ID:        antID,
				PathIndex: pathIdx,
				Position:  0,
				Path:      paths[pathIdx].Rooms,
			})
		}
	}

	// Сортируем по ID
	sort.Slice(ants, func(i, j int) bool {
		return ants[i].ID < ants[j].ID
	})

	// Для отладки
	fmt.Println("\n=== НАЧАЛО СИМУЛЯЦИИ ===")
	for i, path := range paths {
		fmt.Printf("Путь %d: %v\n", i+1, path.Rooms)
	}
	fmt.Println("Распределение муравьев по путям:")
	for i, antsOnPath := range distribution {
		fmt.Printf("  Путь %d: муравьи %v\n", i+1, antsOnPath)
	}
	fmt.Println()

	finished := make(map[int]bool)
	totalFinished := 0
	turn := 1

	for totalFinished < len(ants) {
		turnMoves := make([]string, 0)
		occupied := make(map[string]bool)

		// ВАЖНО: Сначала двигаем муравьев в порядке ИХ ID
		// НЕ сортируем муравьев каждый ход - они уже отсортированы

		for _, ant := range ants {
			if finished[ant.ID] {
				continue
			}

			nextPos := ant.Position + 1
			if nextPos >= len(ant.Path) {
				finished[ant.ID] = true
				totalFinished++
				continue
			}

			nextRoom := ant.Path[nextPos]

			// Проверяем, свободна ли комната
			// Для end всегда свободно
			if !occupied[nextRoom] || nextRoom == "end" {
				ant.Position = nextPos
				if nextRoom == "end" {
					finished[ant.ID] = true
					totalFinished++
					turnMoves = append(turnMoves, fmt.Sprintf("L%d-end", ant.ID))
				} else {
					occupied[nextRoom] = true
					turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", ant.ID, nextRoom))
				}
			}
			// Если комната занята, муравей просто ждет (не двигается в этом ходу)
		}

		if len(turnMoves) > 0 {
			// Сортируем ходы по ID муравья
			sort.Slice(turnMoves, func(i, j int) bool {
				id1 := extractID(turnMoves[i])
				id2 := extractID(turnMoves[j])
				return id1 < id2
			})

			moves = append(moves, strings.Join(turnMoves, " "))
			turn++
		}
	}

	return moves
}

func extractID(move string) int {
	var id int
	fmt.Sscanf(move, "L%d", &id)
	return id
}
