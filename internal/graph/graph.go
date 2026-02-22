package graph

//преобразует муравейник в граф - математическую структуру

import (
	"lem-in/internal/parser"
	"strings"
)

type Graph struct {
	Rooms         map[string]bool
	Links         map[string]map[string]bool
	Start         string
	End           string
	AdjacencyList map[string][]string
}

func Build(farm *parser.Farm) *Graph {
	g := &Graph{
		Rooms:         make(map[string]bool),
		Links:         make(map[string]map[string]bool),
		Start:         farm.Start,
		End:           farm.End,
		AdjacencyList: make(map[string][]string),
	}

	// Добавляем комнаты
	for name := range farm.Rooms {
		g.Rooms[name] = true
		g.AdjacencyList[name] = make([]string, 0)
	}

	// Добавляем связи
	for _, link := range farm.Links {
		parts := strings.Split(link, "-")
		room1, room2 := parts[0], parts[1]

		// Проверяем существование комнат
		if !g.Rooms[room1] || !g.Rooms[room2] {
			continue
		}

		// Инициализируем карту связей
		if g.Links[room1] == nil {
			g.Links[room1] = make(map[string]bool)
		}
		if g.Links[room2] == nil {
			g.Links[room2] = make(map[string]bool)
		}

		// Добавляем двунаправленные связи
		if !g.Links[room1][room2] {
			g.Links[room1][room2] = true
			g.Links[room2][room1] = true
			g.AdjacencyList[room1] = append(g.AdjacencyList[room1], room2)
			g.AdjacencyList[room2] = append(g.AdjacencyList[room2], room1)
		}
	}

	return g
}
