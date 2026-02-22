package solver

import (
	"errors"
	"fmt"
	"lem-in/internal/graph"
	"sort"
)

type Path struct {
	Rooms []string
	Len   int
}

func FindAllPaths(g *graph.Graph, antCount int) ([]Path, [][]int, error) {
	// Находим ВСЕ возможные пути
	allPaths := findAllPaths(g)
	if len(allPaths) == 0 {
		return nil, nil, errors.New("no path found")
	}

	// Выводим все найденные пути для отладки
	fmt.Println("\n=== ВСЕ НАЙДЕННЫЕ ПУТИ (В ПОРЯДКЕ ОБНАРУЖЕНИЯ) ===")
	for i, path := range allPaths {
		fmt.Printf("Путь %d: %v (длина %d)\n", i+1, path, len(path)-1)
	}

	// Сортируем по длине, НО СОХРАНЯЕМ ПОРЯДОК для одинаковых длин
	sort.SliceStable(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	// Выбираем непересекающиеся пути, СОХРАНЯЯ порядок
	optimalPaths := selectOptimalPathsByTime(allPaths, antCount)

	fmt.Println("\n=== ВЫБРАННЫЕ ПУТИ (В ПОРЯДКЕ ОБНАРУЖЕНИЯ) ===")
	for i, path := range optimalPaths {
		fmt.Printf("Путь %d: %v (длина %d)\n", i+1, path.Rooms, path.Len)
	}

	// Распределяем муравьев
	distribution := distributeAntsRoundRobin(optimalPaths, antCount)

	return optimalPaths, distribution, nil
}

func findAllPaths(g *graph.Graph) [][]string {
	var paths [][]string
	var dfs func(current string, path []string, visited map[string]bool)

	dfs = func(current string, path []string, visited map[string]bool) {
		if len(path) > 100 {
			return
		}

		if current == g.End {
			newPath := make([]string, len(path))
			copy(newPath, path)
			paths = append(paths, newPath)
			return
		}

		// Перебираем соседей в том порядке, в котором они в графе
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

	// Пробуем разные комбинации
	for numPaths := 1; numPaths <= 4 && numPaths <= len(paths); numPaths++ {
		// Генерируем комбинации из numPaths путей
		combinations := generateCombinations(paths, numPaths)

		for _, combo := range combinations {
			// Проверяем, можно ли использовать эти пути вместе
			if !pathsCompatible(combo) {
				continue
			}

			// Распределяем муравьев
			distribution := distributeAntsRoundRobin(combo, antCount)

			// Считаем время
			maxTime := 0
			for i, ants := range distribution {
				time := combo[i].Len + len(ants) - 1
				if time > maxTime {
					maxTime = time
				}
			}

			if maxTime < bestTime {
				bestTime = maxTime
				bestCombination = combo
			}
		}
	}

	return bestCombination
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
	// Генерирует все combinations из k путей
	var result [][]Path
	var generate func(start int, current []Path)

	generate = func(start int, current []Path) {
		if len(current) == k {
			result = append(result, current)
			return
		}
		for i := start; i < len(paths); i++ {
			newCurrent := make([]Path, len(current))
			copy(newCurrent, current)
			newCurrent = append(newCurrent, Path{Rooms: paths[i], Len: len(paths[i]) - 1})
			generate(i+1, newCurrent)
		}
	}

	generate(0, []Path{})
	return result
}

func distributeAntsRoundRobin(paths []Path, antCount int) [][]int {
	// 1. Сначала просто считаем количество (емкость) для каждого пути
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

	// 2. Теперь распределяем ID муравьев (1, 2, 3...) по этим путям
	distribution := make([][]int, len(paths))
	currentAnt := 1

	// Распределяем "слоями", чтобы ID шли по порядку в каждом ходу
	for {
		movedInThisLayer := false
		for i := 0; i < len(paths); i++ {
			if counts[i] > 0 {
				distribution[i] = append(distribution[i], currentAnt)
				counts[i]--
				currentAnt++
				movedInThisLayer = true
			}
		}
		if !movedInThisLayer {
			break
		}
	}
	return distribution
}
