package solver

import (
	"errors"
	"lem-in/internal/graph"
)

func FindAllPaths(g *graph.Graph) ([][]string, error) {
	var paths [][]string

	used := make(map[string]bool)

	for {
		path, found := bfsWithBlock(g, used)
		if !found {
			break
		}

		paths = append(paths, path)

		// блокируем внутренние комнаты
		for i := 1; i < len(path)-1; i++ {
			used[path[i]] = true
		}
	}

	if len(paths) == 0 {
		return nil, errors.New("no path found")
	}

	return paths, nil
}

func bfsWithBlock(g *graph.Graph, blocked map[string]bool) ([]string, bool) {
	queue := [][]string{{g.Start}}
	visited := make(map[string]bool)

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]

		room := path[len(path)-1]

		if room == g.End {
			return path, true
		}

		if visited[room] {
			continue
		}
		visited[room] = true

		for _, next := range g.Edges[room] {
			if blocked[next] && next != g.End {
				continue
			}

			newPath := append([]string{}, path...)
			newPath = append(newPath, next)
			queue = append(queue, newPath)
		}
	}

	return nil, false
}
