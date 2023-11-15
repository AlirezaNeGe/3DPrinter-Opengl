package interpreter

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

type GCode struct {
	G float32
	X float32
	Y float32
	Z float32
}

func ExtractValues(scanner *bufio.Scanner) []GCode {
	var values []GCode

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		var value GCode

		for _, part := range parts {
			switch {
			case strings.HasPrefix(part, "G"):
				fmt.Sscanf(part, "G%f", &value.G)
			case strings.HasPrefix(part, "X"):
				fmt.Sscanf(part, "X%f", &value.X)
			case strings.HasPrefix(part, "Y"):
				fmt.Sscanf(part, "Y%f", &value.Y)
			case strings.HasPrefix(part, "Z"):
				fmt.Sscanf(part, "Z%f", &value.Z)
			}
		}
		values = append(values, value)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return values
}
