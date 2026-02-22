package parser

//преобразует текстовый файл в структуру Farm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Farm struct {
	Ants     int
	Rooms    map[string]*Room
	Start    string
	End      string
	Links    []string
	RawLines []string
}

type Room struct {
	Name string
	X, Y int
}

func Parse(filename string) (*Farm, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("invalid data format")
	}
	defer file.Close()

	farm := &Farm{
		Rooms:    make(map[string]*Room),
		RawLines: make([]string, 0),
		Links:    make([]string, 0),
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0
	var isStart, isEnd bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNum++

		// Сохраняем оригинальную строку
		farm.RawLines = append(farm.RawLines, line)

		if line == "" || strings.HasPrefix(line, "#") && line != "##start" && line != "##end" {
			continue
		}

		// Парсим количество муравьев
		if lineNum == 1 {
			ants, err := strconv.Atoi(line)
			if err != nil || ants <= 0 {
				return nil, fmt.Errorf("invalid data format")
			}
			farm.Ants = ants
			continue
		}

		switch line {
		case "##start":
			isStart = true
			continue
		case "##end":
			isEnd = true
			continue
		}

		// Парсим комнаты и связи
		if strings.Contains(line, "-") {
			// Это связь
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid data format")
			}
			farm.Links = append(farm.Links, line)
		} else {
			// Это комната
			parts := strings.Fields(line)
			if len(parts) != 3 {
				return nil, fmt.Errorf("invalid data format")
			}

			name := parts[0]
			if strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") {
				return nil, fmt.Errorf("invalid data format")
			}

			x, err1 := strconv.Atoi(parts[1])
			y, err2 := strconv.Atoi(parts[2])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid data format")
			}

			// Проверяем уникальность комнаты
			if _, exists := farm.Rooms[name]; exists {
				return nil, fmt.Errorf("invalid data format")
			}

			// Проверяем уникальность координат
			for _, room := range farm.Rooms {
				if room.X == x && room.Y == y {
					return nil, fmt.Errorf("invalid data format")
				}
			}

			farm.Rooms[name] = &Room{Name: name, X: x, Y: y}

			if isStart {
				farm.Start = name
				isStart = false
			}
			if isEnd {
				farm.End = name
				isEnd = false
			}
		}
	}

	// Валидация
	if farm.Ants == 0 {
		return nil, fmt.Errorf("invalid data format")
	}
	if farm.Start == "" {
		return nil, fmt.Errorf("invalid data format")
	}
	if farm.End == "" {
		return nil, fmt.Errorf("invalid data format")
	}
	if len(farm.Links) == 0 {
		return nil, fmt.Errorf("invalid data format")
	}

	return farm, nil
}
