// Package graph converts the parsed colony data into a mathematical graph structure.
// Пакет graph преобразует данные распарсенной колонии в математическую структуру графа.
package graph

import (
	"lem-in/internal/models"
	"strings"
)

// Graph represents the ant colony as an adjacency list for efficient pathfinding.
// Graph представляет муравьиную колонию в виде списка смежности для эффективного поиска путей.
type Graph struct {
	Rooms         map[string]bool
	Start         string
	End           string
	AdjacencyList map[string][]string
}

// Build creates a graph structure from the Farm data provided by the parser.
// Build создает структуру графа на основе данных Farm, предоставленных парсером.
func Build(farm *models.Farm) *Graph {
	g := &Graph{
		Rooms:         make(map[string]bool),
		Start:         farm.Start,
		End:           farm.End,
		AdjacencyList: make(map[string][]string),
	}

	// Initialize rooms in the adjacency list
	// Инициализируем комнаты в списке смежности
	for name := range farm.Rooms {
		g.Rooms[name] = true
		g.AdjacencyList[name] = make([]string, 0)
	}

	// Add links between rooms
	// Добавляем связи между комнатами
	addedLinks := make(map[string]map[string]bool)

	for _, link := range farm.Links {
		parts := strings.Split(link, "-")
		if len(parts) != 2 {
			continue
		}
		u, v := parts[0], parts[1]

		// Ensure both rooms exist and aren't linking to themselves
		// Проверяем, что обе комнаты существуют и связь не ведет к самой себе
		if !g.Rooms[u] || !g.Rooms[v] || u == v {
			continue
		}

		// Avoid duplicate links
		// Избегаем дублирования связей
		if addedLinks[u] == nil {
			addedLinks[u] = make(map[string]bool)
		}
		if addedLinks[v] == nil {
			addedLinks[v] = make(map[string]bool)
		}

		if !addedLinks[u][v] {
			g.AdjacencyList[u] = append(g.AdjacencyList[u], v)
			g.AdjacencyList[v] = append(g.AdjacencyList[v], u)
			addedLinks[u][v] = true
			addedLinks[v][u] = true
		}
	}

	return g
}
