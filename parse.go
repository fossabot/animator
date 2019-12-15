package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type rule struct {
	// Replacement Rule
	re     string              // What to search for
	expand func([]byte) []byte // What to do on it (replaceAllFunc)
}

func expandRange(inText []byte) []byte {
	// Expand range e.g. 4-7, into comma separated e.g 4,5,6,7
	s := ""
	text := string(inText)
	start, serr := strconv.Atoi(strings.Split(text, "-")[0]) // 4-7 -> 4
	end, eerr := strconv.Atoi(strings.Split(text, "-")[1])   // 4-7 -> 7
	if serr != nil || eerr != nil {
		fmt.Printf("err: %q\n", text)
	}
	// RANGE
	counter := 0
	if end-start >= 0 {
		counter = 1
	} else {
		counter = -1
	}
	for i := start; i != end+counter; i += counter {
		// TODO factor
		if i == end {
			s += strconv.Itoa(i) // no comma at end
		} else {
			s += strconv.Itoa(i) + ","
		}
	}
	return []byte(s)
}

func expandRepeat(inText []byte) []byte {
	// Expand repeat e.g 4(3), into comma separated e.g 3,3,3,3
	s := ""
	text := string(inText)
	timesString := regexp.MustCompile("(\\d+)\\(.+\\)").FindStringSubmatch(text)[1] // 6(7-90) -> 6
	times, _ := strconv.Atoi(timesString)
	// do not do 0 times (repeat till end)
	// TODO expandRepeatTillEnd
	if times == 0 {
		return inText
	}

	inside := regexp.MustCompile("\\d+\\((.+)\\)").FindStringSubmatch(text)[1] // 6(7-90) -> 7-90
	// append inside to returning string `times` times
	for i := 0; i != times; i += 1 {
		s += inside
		if i != times-1 {
			s += ","
		}
	}
	return []byte(s)
}

func ReplaceAll(s string) string {
	// Replace (expand?) all according to a list of `rule`s
	rules := []rule{rule{"[0-9]+\\-[0-9]+", expandRange}, rule{"\\d+\\([^\\)]+\\)", expandRepeat}}
	for r := range rules {
		s = string(regexp.MustCompile(rules[r].re).ReplaceAllFunc([]byte(s), rules[r].expand))
	}
	return s
}
