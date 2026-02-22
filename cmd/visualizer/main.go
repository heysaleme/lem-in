package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Point struct {
	X, Y int
}

type model struct {
	rooms                  map[string]Point
	links                  [][2]string
	steps                  [][]string
	currStep               int
	minX, minY, maxX, maxY int
}

// Теперь Init ничего не запускает (таймер не нужен)
func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		// Вперед: Стрелка вправо или Пробел
		case "right", " ":
			if m.currStep < len(m.steps)-1 {
				m.currStep++
			}

		// Назад: Стрелка влево
		case "left":
			if m.currStep > 0 {
				m.currStep--
			}

		// Сброс: Клавиша "r" (reset)
		case "r":
			m.currStep = 0
		}
	}
	return m, nil
}

func (m model) View() string {
	width, height := 100, 25
	canvas := make([][]string, height)
	for i := range canvas {
		canvas[i] = make([]string, width)
		for j := range canvas[i] {
			canvas[i][j] = " "
		}
	}

	// Вычисление масштаба
	rangeX, rangeY := m.maxX-m.minX, m.maxY-m.minY
	if rangeX == 0 {
		rangeX = 1
	}
	if rangeY == 0 {
		rangeY = 1
	}
	scaleX, scaleY := (width-20)/rangeX, (height-8)/rangeY
	if scaleX < 5 {
		scaleX = 5
	}
	if scaleY < 2 {
		scaleY = 2
	}

	// Связи
	for _, link := range m.links {
		p1, ok1 := m.rooms[link[0]]
		p2, ok2 := m.rooms[link[1]]
		if ok1 && ok2 {
			drawConnection(canvas, p1, p2, m.minX, m.minY, scaleX, scaleY)
		}
	}

	// Позиции муравьев
	antsInRooms := make(map[string]string)
	movesInfo := "Начало (муравьи в старте)"
	if m.currStep < len(m.steps) && len(m.steps[m.currStep]) > 0 {
		movesInfo = strings.Join(m.steps[m.currStep], " ")
		for _, move := range m.steps[m.currStep] {
			parts := strings.Split(move, "-")
			if len(parts) == 2 {
				antsInRooms[parts[1]] = parts[0]
			}
		}
	}

	// Комнаты
	for name, pos := range m.rooms {
		x := (pos.X-m.minX)*scaleX + 2
		y := (pos.Y-m.minY)*scaleY + 2
		display := fmt.Sprintf("[%s]", name)
		if antID, ok := antsInRooms[name]; ok {
			display = fmt.Sprintf("[%s 🐜 (%s)]", name, antID)
		}
		if y < height && x < width {
			for i, char := range display {
				if x+i < width {
					canvas[y][x+i] = string(char)
				}
			}
		}
	}

	var out strings.Builder
	out.WriteString("┌── LEM-IN ИНТЕРАКТИВНЫЙ ВИЗУАЛИЗАТОР ──────────────────┐\n")
	out.WriteString(fmt.Sprintf("│ Шаг: %d/%d | [→/Space] Вперед | [←] Назад | [r] Сброс │\n", m.currStep+1, len(m.steps)))
	out.WriteString("└────────────────────────────────────────────────────────┘\n")

	for _, row := range canvas {
		out.WriteString(strings.TrimRight(strings.Join(row, ""), " ") + "\n")
	}

	out.WriteString("\n🎬 Перемещения на текущем шаге:\n")
	out.WriteString("   " + movesInfo + "\n")

	if m.currStep == len(m.steps)-1 {
		out.WriteString("\n🏁 ФИНИШ! Все муравьи дома.")
	}

	return out.String()
}

func drawConnection(canvas [][]string, p1, p2 Point, minX, minY, scaleX, scaleY int) {
	x1, y1 := (p1.X-minX)*scaleX+3, (p1.Y-minY)*scaleY+2
	x2, y2 := (p2.X-minX)*scaleX+3, (p2.Y-minY)*scaleY+2
	steps := 8
	for i := 1; i < steps; i++ {
		cx, cy := x1+(x2-x1)*i/steps, y1+(y2-y1)*i/steps
		if cy >= 0 && cy < len(canvas) && cx >= 0 && cx < len(canvas[0]) {
			if canvas[cy][cx] == " " {
				canvas[cy][cx] = "·"
			}
		}
	}
}

func main() {
	m := model{
		rooms: make(map[string]Point),
		minX:  100000, minY: 100000, maxX: -100000, maxY: -100000,
	}
	scanner := bufio.NewScanner(os.Stdin)
	parsingMoves := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || (strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##")) {
			continue
		}
		if strings.HasPrefix(line, "L") {
			parsingMoves = true
			m.steps = append(m.steps, strings.Fields(line))
			continue
		}
		parts := strings.Fields(line)
		if !parsingMoves && len(parts) == 3 {
			var x, y int
			fmt.Sscanf(parts[1], "%d", &x)
			fmt.Sscanf(parts[2], "%d", &y)
			m.rooms[parts[0]] = Point{X: x, Y: y}
			if x < m.minX {
				m.minX = x
			}
			if x > m.maxX {
				m.maxX = x
			}
			if y < m.minY {
				m.minY = y
			}
			if y > m.maxY {
				m.maxY = y
			}
		} else if !parsingMoves && strings.Contains(line, "-") {
			l := strings.Split(line, "-")
			if len(l) == 2 {
				m.links = append(m.links, [2]string{l[0], l[1]})
			}
		}
	}
	if len(m.steps) == 0 {
		m.steps = [][]string{{}}
	}
	tea.NewProgram(m).Run()
}
