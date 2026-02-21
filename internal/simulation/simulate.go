package simulation

import (
	"fmt"
	"sort"
)

type Ant struct {
	ID   int
	Path []string
	Pos  int
}

type Move struct {
	Ant  int
	Room string
}

func Run(paths [][]string, antsCount int) [][]Move {
	var result [][]Move
	var ants []*Ant

	nextAntID := 1

	for {
		var turn []Move
		occupied := make(map[string]bool)

		// двигаем существующих муравьёв
		for i := len(ants) - 1; i >= 0; i-- {
			ant := ants[i]

			if ant.Pos < len(ant.Path)-1 {
				nextRoom := ant.Path[ant.Pos+1]

				if nextRoom == ant.Path[len(ant.Path)-1] || !occupied[nextRoom] {
					ant.Pos++
					turn = append(turn, Move{
						Ant:  ant.ID,
						Room: nextRoom,
					})

					if nextRoom != ant.Path[len(ant.Path)-1] {
						occupied[nextRoom] = true
					}
				}
			}
		}

		// запускаем новых муравьёв
		for _, path := range paths {
			if nextAntID > antsCount {
				break
			}

			firstRoom := path[1]
			if !occupied[firstRoom] {
				ant := &Ant{
					ID:   nextAntID,
					Path: path,
					Pos:  1,
				}

				ants = append(ants, ant)

				turn = append(turn, Move{
					Ant:  ant.ID,
					Room: firstRoom,
				})

				occupied[firstRoom] = true
				nextAntID++
			}
		}

		if len(turn) == 0 {
			break
		}

		result = append(result, turn)
	}

	return result
}

func FormatTurn(turn []Move) string {
	sort.Slice(turn, func(i, j int) bool {
		return turn[i].Ant < turn[j].Ant
	})

	s := ""
	for _, m := range turn {
		s += fmt.Sprintf("L%d-%s ", m.Ant, m.Room)
	}
	return s[:len(s)-1]
}
