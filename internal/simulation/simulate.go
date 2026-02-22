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
	EndRoom   string
}

func Run(paths []solver.Path, distribution [][]int) []string {
	var moves []string

	// Создаем муравьев
	ants := make([]*Ant, 0)
	for pathIdx, antIDs := range distribution {
		endRoom := paths[pathIdx].Rooms[len(paths[pathIdx].Rooms)-1]
		for _, antID := range antIDs {
			ants = append(ants, &Ant{
				ID:        antID,
				PathIndex: pathIdx,
				Position:  0,
				Path:      paths[pathIdx].Rooms,
				EndRoom:   endRoom,
			})
		}
	}

	// Сортируем по ID
	sort.Slice(ants, func(i, j int) bool {
		return ants[i].ID < ants[j].ID
	})

	finished := make(map[int]bool)

	for {
		allFinished := true
		for _, ant := range ants {
			if !finished[ant.ID] {
				allFinished = false
				break
			}
		}
		if allFinished {
			break
		}

		turnMoves := make([]string, 0)
		occupied := make(map[string]bool)
		usedTunnels := make(map[string]bool)

		for _, ant := range ants {
			if finished[ant.ID] {
				continue
			}

			if ant.Position == len(ant.Path)-1 {
				finished[ant.ID] = true
				continue
			}

			currentRoom := ant.Path[ant.Position]
			nextRoom := ant.Path[ant.Position+1]

			// Ключ туннеля
			tunnelKey := currentRoom + "-" + nextRoom
			if currentRoom > nextRoom {
				tunnelKey = nextRoom + "-" + currentRoom
			}

			// Туннель уже использован?
			if usedTunnels[tunnelKey] {
				continue
			}

			// Выход в END (всегда можно, туннель только проверяем)
			if nextRoom == ant.EndRoom {
				ant.Position++
				finished[ant.ID] = true
				usedTunnels[tunnelKey] = true
				turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", ant.ID, nextRoom))
				continue
			}

			// Выход из старта
			if ant.Position == 0 {
				// Для выхода из старта проверяем только occupied (не willBeFree)
				if !occupied[nextRoom] {
					ant.Position++
					occupied[nextRoom] = true
					usedTunnels[tunnelKey] = true
					turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", ant.ID, nextRoom))
				}
				continue
			}

			// Обычное движение
			if !occupied[nextRoom] {
				ant.Position++
				occupied[nextRoom] = true
				usedTunnels[tunnelKey] = true
				turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", ant.ID, nextRoom))
			}
		}

		if len(turnMoves) > 0 {
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
