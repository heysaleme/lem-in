// Package simulation implements the ant movement logic, ensuring no collisions.
// Пакет simulation реализует логику передвижения муравьев, предотвращая столкновения.
package simulation

import (
	"fmt"
	"lem-in/internal/models"
	"sort"
	"strings"
)

// Run executes the movement simulation step-by-step until all ants reach the end.
// Run пошагово выполняет симуляцию движения, пока все муравьи не достигнут финиша.
func Run(paths []models.Path, distribution [][]int) []string {
	var moves []string

	// Initialize ant objects based on the distribution layers
	// Инициализируем объекты муравьев на основе слоев распределения
	ants := initializeAnts(paths, distribution)

	for {
		turnMoves := make([]string, 0)
		occupied := make(map[string]bool)
		usedTunnels := make(map[string]bool)
		anyMoved := false

		// Sort ants: those closer to the end move first to free up rooms
		// Сортировка: те, кто ближе к концу, ходят первыми, освобождая комнаты
		sort.SliceStable(ants, func(i, j int) bool {
			return ants[i].Position > ants[j].Position
		})

		for _, ant := range ants {
			if ant.Finished {
				continue
			}

			// Try to move to the next room in the ant's assigned path
			// Попытка перейти в следующую комнату согласно назначенному пути
			if moveAnt(ant, &turnMoves, occupied, usedTunnels) {
				anyMoved = true
			}
		}

		if !anyMoved {
			break
		}

		// Sort output moves by Ant ID for consistent formatting
		// Сортируем ходы по ID муравья для единообразия вывода
		sortMoves(turnMoves)
		moves = append(moves, strings.Join(turnMoves, " "))
	}

	return moves
}

// initializeAnts creates an ordered slice of ants to ensure fair start line exit.
func initializeAnts(paths []models.Path, distribution [][]int) []*models.Ant {
	ants := make([]*models.Ant, 0)
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
				ants = append(ants, &models.Ant{
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
	return ants
}

// moveAnt attempts to advance a single ant to its next room.
func moveAnt(ant *models.Ant, turnMoves *[]string, occupied, usedTunnels map[string]bool) bool {
	currentRoom := ant.Path[ant.Position]
	nextRoom := ant.Path[ant.Position+1]

	tunnelKey := generateTunnelKey(currentRoom, nextRoom)

	if !usedTunnels[tunnelKey] && (nextRoom == ant.EndRoom || !occupied[nextRoom]) {
		ant.Position++
		*turnMoves = append(*turnMoves, fmt.Sprintf("L%d-%s", ant.ID, nextRoom))
		usedTunnels[tunnelKey] = true

		if nextRoom == ant.EndRoom {
			ant.Finished = true
		} else {
			occupied[nextRoom] = true
		}
		return true
	}
	return false
}

func generateTunnelKey(r1, r2 string) string {
	if r1 > r2 {
		return r2 + "-" + r1
	}
	return r1 + "-" + r2
}

func sortMoves(moves []string) {
	sort.Slice(moves, func(i, j int) bool {
		var id1, id2 int
		fmt.Sscanf(moves[i], "L%d", &id1)
		fmt.Sscanf(moves[j], "L%d", &id2)
		return id1 < id2
	})
}
