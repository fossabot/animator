package main

import (
	"fmt"
	"strconv"
)

func Render(mainExpr string, otherFrames map[string][]float64, out string) string {
	// Take a file pattern and a list of filenames and in-out times and return an `ffmpeg` command to turn them into an output file
	otherExpr := ""

	for k := range otherFrames {
		otherExpr += string(fmt.Sprintf("-i %v ", k)) // Input all otherFrames
	}

	otherFilter := ""

	prev := "0"
	count := 0
	for _, v := range otherFrames {
		// Add text to `filter_complex` overlay
		count++
		otherFilter += string(fmt.Sprintf("[%v][%v] overlay=0:0:enable='between(t,%v,%v)'", prev, count, v[0], v[1]))
		prev = "v" + strconv.Itoa(count)
		otherFilter += fmt.Sprintf("[%v]", prev)
		if count != len(otherFrames) {
			otherFilter += "," // No comma at end
		}
	}
	// format `ffmpeg` command
	return fmt.Sprintf("ffmpeg -y -f image2 -pattern_type sequence -i %q %v -filter_complex %q -map \"[%v]\" %v", mainExpr, otherExpr, otherFilter, prev, out)
}
