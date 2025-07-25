package caca

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

/*
   TODO: Parse metadata from file

   Metadata:
     # Start date: <date>
     # Repository: <path-to-repository>
     # Author: <name> <email>
*/

func ReadCanvas(path string) (*Canvas, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metadata CanvasMetadata
	// TODO: Factor out number of columns and rows
	canvas2D := make([][]uint8, 0, 7)
	scanner := bufio.NewScanner(file)
	lineNo, count := 0, 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		stripped := strings.TrimSpace(line)

		if len(stripped) == 0 {
			continue
		}

		if strings.HasPrefix(stripped, "#") {
			// TODO: Parse metadata
			continue
		}

		weekdayData := make([]uint8, 0, 53)

		for col, c := range line {
			if unicode.IsSpace(c) {
				continue
			}
			if '0' <= c && c <= '4' {
				weekdayData = append(weekdayData, uint8(c-'0'))
				count++
			} else {
				return nil, fmt.Errorf("unexpected character '%c' at %d:%d", c, lineNo, col)
			}
		}
		canvas2D = append(canvas2D, weekdayData)
	}

	return &Canvas{canvas2D, count, &metadata}, nil
}
