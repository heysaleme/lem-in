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
	Finished  bool
	Path      []string
}

func Run(paths []solver.Path, distribution [][]int) []string {
	var moves []string
	ants := make([]*Ant, 0)

	// Создаем муравьев
	for pathIdx, antIDs := range distribution {
		for _, antID := range antIDs {
			ants = append(ants, &Ant{
				ID:        antID,
				PathIndex: pathIdx,
				Position:  0,
				Finished:  false,
				Path:      paths[pathIdx].Rooms,
			})
		}
	}

	// Сортируем по ID
	sort.Slice(ants, func(i, j int) bool {
		return ants[i].ID < ants[j].ID
	})

	allFinished := false
	for !allFinished {
		allFinished = true
		turnMoves := make([]string, 0)

		// Карта занятых комнат в этом ходу
		occupied := make(map[string]bool)

		// Сначала определяем, кто куда пойдет
		type move struct {
			ant    *Ant
			room   string
			finish bool
		}
		possibleMoves := make([]move, 0)

		for _, ant := range ants {
			if ant.Finished {
				continue
			}

			nextPos := ant.Position + 1
			if nextPos >= len(ant.Path) {
				ant.Finished = true
				continue
			}

			nextRoom := ant.Path[nextPos]
			possibleMoves = append(possibleMoves, move{
				ant:    ant,
				room:   nextRoom,
				finish: nextRoom == "end",
			})
		}

		// Сортируем возможные ходы (сначала муравьи с меньшими ID)
		sort.Slice(possibleMoves, func(i, j int) bool {
			return possibleMoves[i].ant.ID < possibleMoves[j].ant.ID
		})

		// Выполняем ходы
		for _, m := range possibleMoves {
			// Проверяем, свободна ли комната
			if !occupied[m.room] || m.finish {
				m.ant.Position++
				if m.finish {
					m.ant.Finished = true
					turnMoves = append(turnMoves, fmt.Sprintf("L%d-end", m.ant.ID))
				} else {
					occupied[m.room] = true
					turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", m.ant.ID, m.room))
				}
			}
		}

		// Проверяем, все ли закончили
		for _, ant := range ants {
			if !ant.Finished {
				allFinished = false
				break
			}
		}

		if len(turnMoves) > 0 {
			// Сортируем ходы по ID муравья
			sort.Slice(turnMoves, func(i, j int) bool {
				id1 := extractID(turnMoves[i])
				id2 := extractID(turnMoves[j])
				return id1 < id2
			})
			moves = append(moves, strings.Join(turnMoves, " "))
		}
	}

	return moves
}

func extractID(move string) int {
	var id int
	fmt.Sscanf(move, "L%d", &id)
	return id
}
