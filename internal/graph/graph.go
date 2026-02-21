package graph

import "lem-in/internal/parser"

type Room struct {
	Name string
	X    int
	Y    int
}

type Graph struct {
	Rooms map[string]*Room
	Edges map[string][]string
	Start string
	End   string
}

func Build(farm *parser.Farm) *Graph {
	g := &Graph{
		Rooms: make(map[string]*Room),
		Edges: make(map[string][]string),
		Start: farm.Start,
		End:   farm.End,
	}

	for _, r := range farm.Rooms {
		g.Rooms[r.Name] = &Room{
			Name: r.Name,
			X:    r.X,
			Y:    r.Y,
		}
	}

	for _, l := range farm.Links {
		g.Edges[l.From] = append(g.Edges[l.From], l.To)
		g.Edges[l.To] = append(g.Edges[l.To], l.From)
	}

	return g
}
