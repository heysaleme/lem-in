
type Room struct {
	Name string
	X, Y int
}

type Graph struct {
	Rooms map[string]*Room
	Edges map[string][]string
	Start *Room
	End   *Room
}