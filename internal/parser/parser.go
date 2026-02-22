// Package parser implements functions to read and validate the input file.
// Пакет parser реализует функции для чтения и валидации входного файла.
package parser

import (
	"bufio"
	"fmt"
	"lem-in/internal/models"
	"os"
	"strconv"
	"strings"
)

// Parse reads a file and converts it into a Farm structure.
// Parse читает файл и преобразует его в структуру Farm.
func Parse(filename string) (*models.Farm, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format")
	}
	defer file.Close()

	farm := &models.Farm{
		Rooms:    make(map[string]*models.Room),
		RawLines: make([]string, 0),
		Links:    make([]string, 0),
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0
	var isStart, isEnd bool

	for scanner.Scan() {
		line := scanner.Text() // Убираем TrimSpace здесь, чтобы сохранить формат для RawLines
		trimmed := strings.TrimSpace(line)

		farm.RawLines = append(farm.RawLines, line)

		// Skip comments (except commands) and empty lines
		if trimmed == "" || (strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, "##")) {
			continue
		}

		lineNum++

		// First non-comment line must be the number of ants
		if farm.Ants == 0 && !isStart && !isEnd {
			ants, err := strconv.Atoi(trimmed)
			if err != nil || ants <= 0 {
				return nil, fmt.Errorf("ERROR: invalid data format, invalid number of ants")
			}
			farm.Ants = ants
			continue
		}

		// Handle commands
		if trimmed == "##start" {
			isStart = true
			continue
		} else if trimmed == "##end" {
			isEnd = true
			continue
		}

		// Parse Links or Rooms
		if strings.Contains(trimmed, "-") {
			if err := parseLink(farm, trimmed); err != nil {
				return nil, err
			}
		} else {
			if err := parseRoom(farm, trimmed, &isStart, &isEnd); err != nil {
				return nil, err
			}
		}
	}

	if err := validateFarm(farm); err != nil {
		return nil, err
	}

	return farm, nil
}

// parseRoom handles the extraction of room names and coordinates.
// parseRoom обрабатывает извлечение имен комнат и их координат.
func parseRoom(farm *models.Farm, line string, isStart, isEnd *bool) error {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return fmt.Errorf("ERROR: invalid data format")
	}

	name := parts[0]
	if strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") {
		return fmt.Errorf("ERROR: invalid data format")
	}

	x, err1 := strconv.Atoi(parts[1])
	y, err2 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil {
		return fmt.Errorf("ERROR: invalid data format")
	}

	if _, exists := farm.Rooms[name]; exists {
		return fmt.Errorf("ERROR: invalid data format, duplicate room")
	}

	// Coordinate uniqueness check
	for _, r := range farm.Rooms {
		if r.X == x && r.Y == y {
			return fmt.Errorf("ERROR: invalid data format, duplicate coordinates")
		}
	}

	farm.Rooms[name] = &models.Room{Name: name, X: x, Y: y}

	if *isStart {
		if farm.Start != "" {
			return fmt.Errorf("ERROR: multiple start rooms")
		}
		farm.Start = name
		*isStart = false
	}
	if *isEnd {
		if farm.End != "" {
			return fmt.Errorf("ERROR: multiple end rooms")
		}
		farm.End = name
		*isEnd = false
	}
	return nil
}

// parseLink handles connection strings like "A-B".
// parseLink обрабатывает строки связей вида "A-B".
func parseLink(farm *models.Farm, line string) error {
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		return fmt.Errorf("ERROR: invalid data format")
	}
	// Basic validation: do rooms exist?
	// Note: Detailed validation can be done after parsing all rooms.
	farm.Links = append(farm.Links, line)
	return nil
}

// validateFarm ensures the minimum requirements for a valid colony.
// validateFarm проверяет минимальные требования для валидной колонии.
func validateFarm(farm *models.Farm) error {
	if farm.Ants <= 0 || farm.Start == "" || farm.End == "" || len(farm.Links) == 0 {
		return fmt.Errorf("ERROR: invalid data format")
	}
	return nil
}
