package solver

import (
	"errors"
	"lem-in/internal/graph"
	"sort"
)

type Path struct {
	Rooms []string
	Len   int
}

func FindAllPaths(g *graph.Graph, antCount int) ([]Path, [][]int, error) {
	allPaths := findAllPaths(g)
	if len(allPaths) == 0 {
		return nil, nil, errors.New("no path found")
	}

	// Сортируем по длине
	sort.SliceStable(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	// Выбираем непересекающиеся оптимальные пути
	optimalPaths := selectOptimalPathsByTime(allPaths, antCount)

	// Распределяем муравьев (используем новую логику для порядка ID)
	distribution := distributeAntsCorrectOrder(optimalPaths, antCount)

	return optimalPaths, distribution, nil
}

func findAllPaths(g *graph.Graph) [][]string {
	var paths [][]string
	var dfs func(current string, path []string, visited map[string]bool)

	dfs = func(current string, path []string, visited map[string]bool) {
		if len(path) > 100 {
			return
		} // Ограничитель глубины

		if current == g.End {
			newPath := make([]string, len(path))
			copy(newPath, path)
			paths = append(paths, newPath)
			return
		}

		for _, neighbor := range g.AdjacencyList[current] {
			if neighbor == g.Start {
				continue
			}
			if !visited[neighbor] {
				visited[neighbor] = true
				dfs(neighbor, append(path, neighbor), visited)
				delete(visited, neighbor)
			}
		}
	}

	visited := make(map[string]bool)
	visited[g.Start] = true
	dfs(g.Start, []string{g.Start}, visited)
	return paths
}

func selectOptimalPathsByTime(paths [][]string, antCount int) []Path {
	if len(paths) == 0 {
		return nil
	}

	var bestCombination []Path
	bestTime := int(^uint(0) >> 1)

	// Пробуем комбинации из разного количества путей
	for numPaths := 1; numPaths <= len(paths); numPaths++ {
		combinations := generateCombinations(paths, numPaths)
		for _, combo := range combinations {
			if !pathsCompatible(combo) {
				continue
			}

			// Временный расчет времени для этой комбинации
			maxTime := calculateMaxTime(combo, antCount)

			if maxTime < bestTime {
				bestTime = maxTime
				bestCombination = combo
			}
		}
	}
	return bestCombination
}

func calculateMaxTime(paths []Path, antCount int) int {
	counts := make([]int, len(paths))
	for ant := 0; ant < antCount; ant++ {
		bestIdx := 0
		minT := paths[0].Len + counts[0]
		for i := 1; i < len(paths); i++ {
			if paths[i].Len+counts[i] < minT {
				minT = paths[i].Len + counts[i]
				bestIdx = i
			}
		}
		counts[bestIdx]++
	}

	maxT := 0
	for i := 0; i < len(paths); i++ {
		if paths[i].Len+counts[i]-1 > maxT {
			maxT = paths[i].Len + counts[i] - 1
		}
	}
	return maxT
}

func distributeAntsCorrectOrder(paths []Path, antCount int) [][]int {
	counts := make([]int, len(paths))
	for ant := 0; ant < antCount; ant++ {
		bestIdx := 0
		minT := paths[0].Len + counts[0]
		for i := 1; i < len(paths); i++ {
			if paths[i].Len+counts[i] < minT {
				minT = paths[i].Len + counts[i]
				bestIdx = i
			}
		}
		counts[bestIdx]++
	}

	distribution := make([][]int, len(paths))
	currentAnt := 1
	for {
		added := false
		for i := 0; i < len(paths); i++ {
			if counts[i] > 0 {
				distribution[i] = append(distribution[i], currentAnt)
				counts[i]--
				currentAnt++
				added = true
			}
		}
		if !added {
			break
		}
	}
	return distribution
}

func pathsCompatible(paths []Path) bool {
	used := make(map[string]bool)
	for _, path := range paths {
		for i := 1; i < len(path.Rooms)-1; i++ {
			if used[path.Rooms[i]] {
				return false
			}
			used[path.Rooms[i]] = true
		}
	}
	return true
}

func generateCombinations(paths [][]string, k int) [][]Path {
	var result [][]Path
	var generate func(start int, current []Path)
	generate = func(start int, current []Path) {
		if len(current) == k {
			res := make([]Path, k)
			copy(res, current)
			result = append(result, res)
			return
		}
		for i := start; i < len(paths); i++ {
			generate(i+1, append(current, Path{Rooms: paths[i], Len: len(paths[i]) - 1}))
		}
	}
	generate(0, []Path{})
	return result
}
