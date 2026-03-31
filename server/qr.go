package main

import (
	"fmt"
	"strings"

	qrcode "github.com/yeqown/go-qrcode/v2"
)

type stdoutWriter struct {
	bitmap [][]bool
}

func (w *stdoutWriter) Write(mat qrcode.Matrix) error {
	w.bitmap = mat.Bitmap()
	return nil
}

func (w *stdoutWriter) Close() error {
	if w.bitmap == nil {
		return nil
	}

	height := len(w.bitmap)
	width := len(w.bitmap[0])
	quiet := 2

	paddedHeight := height + quiet*2
	paddedWidth := width + quiet*2

	isBlack := func(row, col int) bool {
		r := row - quiet
		c := col - quiet
		if r < 0 || r >= height || c < 0 || c >= width {
			return false
		}
		return w.bitmap[r][c]
	}

	var sb strings.Builder
	sb.WriteByte('\n')

	for row := 0; row < paddedHeight; row += 2 {
		for col := 0; col < paddedWidth; col++ {
			top := isBlack(row, col)
			bottom := row+1 < paddedHeight && isBlack(row+1, col)

			switch {
			case top && bottom:
				sb.WriteString("█")
			case top:
				sb.WriteString("▀")
			case bottom:
				sb.WriteString("▄")
			default:
				sb.WriteString(" ")
			}
		}
		sb.WriteByte('\n')
	}

	fmt.Print(sb.String())
	return nil
}
