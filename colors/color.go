package colors

import (
	"fmt"
	"strconv"
	"strings"
)

type Color struct {
	r, g, b uint8
}

func FromHex(hex string) *Color {
	hex, _ = strings.CutPrefix(hex, "#")
	r, err := strconv.ParseUint(hex[:2], 16, 8)
	if err != nil {
		return nil
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return nil
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return nil
	}

	return &Color{
		r: uint8(r),
		g: uint8(g),
		b: uint8(b),
	}
}

func TextBackground(text string, c Color) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm%s\033[0m", c.r, c.g, c.b, text)
}
