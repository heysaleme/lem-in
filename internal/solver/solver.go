package solver

import (
	"errors"
	"lem-in/internal/graph"
	"lem-in/internal/models"
	"sort"
)

func Solve(g *graph.Graph, antCount int) ([]models.Path, [][]int, error) {
	// 1. Находим ВООБЩЕ все возможные пути от Start до End
	allPaths := findAllPathsDFS(g)
	if len(allPaths) == 0 {
		return nil, nil, errors.New("no path found")
	}

	// 2. Генерируем комбинации непересекающихся путей и выбираем лучшую
	bestCombination := findBestPathCombo(allPaths, antCount)

	// 3. Распределяем муравьев
	distribution := distributeAnts(bestCombination, antCount)

	return bestCombination, distribution, nil
}

// findAllPathsDFS находит все пути без циклов
func findAllPathsDFS(g *graph.Graph) [][]string {
	var paths [][]string
	var dfs func(curr string, visited map[string]bool, path []string)

	dfs = func(curr string, visited map[string]bool, path []string) {
		if curr == g.End {
			temp := make([]string, len(path))
			copy(temp, path)
			paths = append(paths, temp)
			return
		}
		for _, next := range g.AdjacencyList[curr] {
			if !visited[next] {
				visited[next] = true
				dfs(next, visited, append(path, next))
				visited[next] = false
			}
		}
	}

	visited := map[string]bool{g.Start: true}
	dfs(g.Start, visited, []string{g.Start})
	return paths
}

// findBestPathCombo перебирает комбинации путей, которые не пересекаются по комнатам
func findBestPathCombo(allPaths [][]string, antCount int) []models.Path {
	var bestCombo []models.Path
	minSteps := int(^uint(0) >> 1)

	// Превращаем в структуру Path и сортируем для стабильности
	var paths []models.Path
	for _, p := range allPaths {
		paths = append(paths, models.Path{Rooms: p, Len: len(p) - 1})
	}

	// Рекурсивно ищем наборы непересекающихся путей
	var backtrack func(index int, currentCombo []models.Path)
	backtrack = func(index int, currentCombo []models.Path) {
		if len(currentCombo) > 0 {
			steps := calculateSteps(currentCombo, antCount)
			if steps < minSteps {
				minSteps = steps
				bestCombo = make([]models.Path, len(currentCombo))
				copy(bestCombo, currentCombo)
			}
		}

		for i := index; i < len(paths); i++ {
			if isCompatible(currentCombo, paths[i]) {
				backtrack(i+1, append(currentCombo, paths[i]))
			}
		}
	}

	backtrack(0, []models.Path{})

	// Сортируем пути в комбинации по длине (важно для распределения)
	sort.Slice(bestCombo, func(i, j int) bool {
		return bestCombo[i].Len < bestCombo[j].Len
	})

	return bestCombo
}

func isCompatible(combo []models.Path, newPath models.Path) bool {
	for _, p := range combo {
		for _, r1 := range p.Rooms[1 : len(p.Rooms)-1] {
			for _, r2 := range newPath.Rooms[1 : len(newPath.Rooms)-1] {
				if r1 == r2 {
					return false
				}
			}
		}
	}
	return true
}

// Математический расчет количества строк
func calculateSteps(paths []models.Path, antCount int) int {
	if len(paths) == 0 {
		return 1000000
	}
	sumLens := 0
	for _, p := range paths {
		sumLens += p.Len
	}
	return (antCount + sumLens - 1) / len(paths)
}

// Логика распределения ID (оставляем ту же, она работает верно)
func distributeAnts(paths []models.Path, antCount int) [][]int {
	counts := make([]int, len(paths))
	for ant := 0; ant < antCount; ant++ {
		bestIdx := 0
		minTime := paths[0].Len + counts[0]
		for i := 1; i < len(paths); i++ {
			if paths[i].Len+counts[i] < minTime {
				minTime = paths[i].Len + counts[i]
				bestIdx = i
			}
		}
		counts[bestIdx]++
	}

	distribution := make([][]int, len(paths))
	currentID := 1
	for {
		added := false
		for i := 0; i < len(paths); i++ {
			if counts[i] > 0 {
				distribution[i] = append(distribution[i], currentID)
				counts[i]--
				currentID++
				added = true
			}
		}
		if !added {
			break
		}
	}
	return distribution
}
