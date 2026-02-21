package solver

import (
	"errors"
	"lem-in/internal/graph"
)

func BFS(g *graph.Graph) ([]string, error) {
	queue := [][]string{{g.Start}}
	visited := make(map[string]bool)

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]

		room := path[len(path)-1]

		if room == g.End {
			return path, nil
		}

		if visited[room] {
			continue
		}
		visited[room] = true

		for _, neighbor := range g.Edges[room] {
			newPath := append([]string{}, path...)
			newPath = append(newPath, neighbor)
			queue = append(queue, newPath)
		}
	}

	return nil, errors.New("no path found")
}
