// Package models defines core data structures for the ant farm project.
// Пакет models определяет основные структуры данных для проекта муравьиной фермы.
package models

// Room represents a node in the colony with coordinates.
// Room представляет собой узел колонии с координатами.
type Room struct {
	Name string
	X, Y int
}

// Path represents a sequence of rooms from start to end.
// Path представляет последовательность комнат от старта до финиша.
type Path struct {
	Rooms []string
	Len   int
}

// Ant represents an individual ant in the simulation.
// Ant представляет отдельного муравья в симуляции.
type Ant struct {
	ID        int
	PathIndex int
	Position  int
	Path      []string
	EndRoom   string
	Finished  bool
}

// Farm represents the entire colony configuration.
// Farm представляет полную конфигурацию колонии.
type Farm struct {
	Ants     int
	Rooms    map[string]*Room
	Start    string
	End      string
	Links    []string // Raw links like "A-B"
	RawLines []string // Original file content for output
}
