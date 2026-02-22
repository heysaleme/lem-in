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

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "right", " ":
			if m.currStep < len(m.steps)-1 {
				m.currStep++
			}
		case "left":
			if m.currStep > 0 {
				m.currStep--
			}
		case "r":
			m.currStep = 0
		}
	}
	return m, nil
}

func (m model) View() string {
	// Максимальная ширина для предотвращения обрезания боковых комнат
	width, height := 160, 30
	canvas := make([][]string, height)
	for i := range canvas {
		canvas[i] = make([]string, width)
		for j := range canvas[i] {
			canvas[i][j] = " "
		}
	}

	rangeX, rangeY := m.maxX-m.minX, m.maxY-m.minY
	if rangeX == 0 {
		rangeX = 1
	}
	if rangeY == 0 {
		rangeY = 1
	}

	// Оптимальный масштаб для предотвращения наложения комнат
	scaleX := (width - 60) / rangeX
	scaleY := (height - 10) / rangeY
	if scaleX < 15 {
		scaleX = 15
	}
	if scaleY < 3 {
		scaleY = 3
	}

	// 1. СНАЧАЛА РИСУЕМ СВЯЗИ (фоновый слой)
	for _, link := range m.links {
		p1, ok1 := m.rooms[link[0]]
		p2, ok2 := m.rooms[link[1]]
		if ok1 && ok2 {
			drawConnection(canvas, p1, p2, m.minX, m.minY, scaleX, scaleY)
		}
	}

	// Сбор данных о текущих позициях муравьев
	antsInRooms := make(map[string]string)
	movesInfo := "Start / Начало"
	if m.currStep < len(m.steps) && len(m.steps[m.currStep]) > 0 {
		movesInfo = strings.Join(m.steps[m.currStep], " ")
		for _, move := range m.steps[m.currStep] {
			parts := strings.Split(move, "-")
			if len(parts) == 2 {
				antsInRooms[parts[1]] = parts[0]
			}
		}
	}

	// 2. ЗАТЕМ РИСУЕМ КОМНАТЫ (затирая точки фона)
	for name, pos := range m.rooms {
		// Смещение x+5 и y+2 для центрирования графа
		x := (pos.X-m.minX)*scaleX + 5
		y := (pos.Y-m.minY)*scaleY + 2

		var display string
		if antID, ok := antsInRooms[name]; ok {
			// Ультра-компактный формат: [Имя🐜ID]
			display = fmt.Sprintf("[%s🐜%s]", name, antID)
		} else {
			display = fmt.Sprintf("[%s]", name)
		}

		if y < height && x < width {
			for i, char := range display {
				if x+i < width {
					// Принудительная запись (удаляет точки внутри комнаты)
					canvas[y][x+i] = string(char)
				}
			}
		}
	}

	var out strings.Builder
	// ПОЛНЫЙ ЗАГОЛОВОК (ИНТЕРФЕЙС)
	out.WriteString("┌─────── LEM-IN INTERACTIVE VISUALIZER / ИНТЕРАКТИВНЫЙ ВИЗУАЛИЗАТОР ─────────┐\n")
	out.WriteString(fmt.Sprintf("│  Шаг/Step: %d/%d | [→/Space] Next/Вперед | [←] Back/Назад | [r] Reset/Сброс  │\n", m.currStep+1, len(m.steps)))
	out.WriteString("└────────────────────────────────────────────────────────────────────────────┘\n")

	// Рендеринг игрового поля
	for _, row := range canvas {
		line := strings.TrimRight(strings.Join(row, ""), " ")
		if line != "" {
			out.WriteString(line + "\n")
		}
	}

	// НИЖНЯЯ ПАНЕЛЬ С ИНФОРМАЦИЕЙ
	out.WriteString("\n🎬 Moves on this step / Перемещения на шаге:\n")
	out.WriteString("   " + movesInfo + "\n")

	if m.currStep == len(m.steps)-1 && len(m.steps) > 1 {
		out.WriteString("\n🏁 FINISH! All ants are home / ФИНИШ! Все муравьи дома.")
	}

	return out.String()
}

func drawConnection(canvas [][]string, p1, p2 Point, minX, minY, scaleX, scaleY int) {
	// Точки входа туннелей в комнаты (с учетом смещения)
	x1, y1 := (p1.X-minX)*scaleX+6, (p1.Y-minY)*scaleY+2
	x2, y2 := (p2.X-minX)*scaleX+6, (p2.Y-minY)*scaleY+2

	steps := 12
	for i := 1; i < steps; i++ {
		cx, cy := x1+(x2-x1)*i/steps, y1+(y2-y1)*i/steps
		if cy >= 0 && cy < len(canvas) && cx >= 0 && cx < len(canvas[0]) {
			// Рисуем точки только там, где еще нет текста
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
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
