package parser

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

type RoomInput struct {
	Name string
	X    int
	Y    int
}

type LinkInput struct {
	From string
	To   string
}

type Farm struct {
	Ants     int
	Rooms    []RoomInput
	Links    []LinkInput
	Start    string
	End      string
	RawLines []string
}

func Parse(filename string) (*Farm, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	farm := &Farm{}

	stage := "ants"
	nextIsStart := false
	nextIsEnd := false

	for scanner.Scan() {
		line := scanner.Text()
		farm.RawLines = append(farm.RawLines, line)

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			if line == "##start" {
				nextIsStart = true
			} else if line == "##end" {
				nextIsEnd = true
			}
			continue
		}

		if stage == "ants" {
			ants, err := strconv.Atoi(line)
			if err != nil || ants <= 0 {
				return nil, errors.New("invalid number of ants")
			}
			farm.Ants = ants
			stage = "rooms"
			continue
		}

		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, errors.New("invalid link")
			}
			farm.Links = append(farm.Links, LinkInput{
				From: parts[0],
				To:   parts[1],
			})
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			return nil, errors.New("invalid room format")
		}

		x, err1 := strconv.Atoi(parts[1])
		y, err2 := strconv.Atoi(parts[2])
		if err1 != nil || err2 != nil {
			return nil, errors.New("invalid room coordinates")
		}

		room := RoomInput{
			Name: parts[0],
			X:    x,
			Y:    y,
		}

		farm.Rooms = append(farm.Rooms, room)

		if nextIsStart {
			farm.Start = room.Name
			nextIsStart = false
		}
		if nextIsEnd {
			farm.End = room.Name
			nextIsEnd = false
		}
	}

	if farm.Start == "" || farm.End == "" {
		return nil, errors.New("no start or end room found")
	}

	return farm, nil
}
