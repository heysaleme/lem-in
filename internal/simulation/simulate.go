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
	Finished  bool
}

func Run(paths []solver.Path, distribution [][]int) []string {
	var moves []string

	// 1. Формируем очередь муравьев правильно.
	// Вместо того чтобы брать всех муравьев пути 1, потом всех пути 2,
	// мы берем первого муравья из каждого пути, потом второго и т.д.
	ants := make([]*Ant, 0)
	maxAntsInPath := 0
	for _, d := range distribution {
		if len(d) > maxAntsInPath {
			maxAntsInPath = len(d)
		}
	}

	for i := 0; i < maxAntsInPath; i++ {
		for pathIdx := 0; pathIdx < len(distribution); pathIdx++ {
			if i < len(distribution[pathIdx]) {
				antID := distribution[pathIdx][i]
				ants = append(ants, &Ant{
					ID:        antID,
					PathIndex: pathIdx,
					Position:  0,
					Path:      paths[pathIdx].Rooms,
					EndRoom:   paths[pathIdx].Rooms[len(paths[pathIdx].Rooms)-1],
					Finished:  false,
				})
			}
		}
	}

	for {
		turnMoves := make([]string, 0)
		occupied := make(map[string]bool)
		usedTunnels := make(map[string]bool)
		anyMoved := false

		// 2. Сортируем: те, кто уже в пути и дальше всех, ходят первыми.
		// Если позиция одинаковая (например, оба в старте), сохраняем порядок очереди.
		sort.SliceStable(ants, func(i, j int) bool {
			return ants[i].Position > ants[j].Position
		})

		for _, ant := range ants {
			if ant.Finished {
				continue
			}

			currentRoom := ant.Path[ant.Position]
			nextRoom := ant.Path[ant.Position+1]

			tunnelKey := currentRoom + "-" + nextRoom
			if currentRoom > nextRoom {
				tunnelKey = nextRoom + "-" + currentRoom
			}

			// Условие хода: туннель свободен И (комната свободна ИЛИ это финиш)
			if !usedTunnels[tunnelKey] && (nextRoom == ant.EndRoom || !occupied[nextRoom]) {
				ant.Position++
				turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", ant.ID, nextRoom))
				anyMoved = true
				usedTunnels[tunnelKey] = true

				if nextRoom == ant.EndRoom {
					ant.Finished = true
				} else {
					occupied[nextRoom] = true
				}
			}
		}

		if !anyMoved {
			break
		}

		// Сортируем только для красивого вывода в консоль
		sort.Slice(turnMoves, func(i, j int) bool {
			return extractID(turnMoves[i]) < extractID(turnMoves[j])
		})

		moves = append(moves, strings.Join(turnMoves, " "))
	}

	return moves
}

func extractID(move string) int {
	var id int
	fmt.Sscanf(move, "L%d", &id)
	return id
}
