package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [-live] <file> [<file> ...]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	live := flag.Bool("live", false, "Perform operations (as opposed to dry run)")
	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
	}

	for _, file := range flag.Args() {
		name := path.Base(file)
		newName := rename(name)
		newFile := path.Join(path.Dir(file), newName)

		change := ""
		if name == newName {
			change = " [unchanged]"
		}
		fmt.Printf("%v -> %v%v\n", file, newFile, change)

		if *live {
			err := os.Rename(file, newFile)
			if err != nil {
				panic(err)
			}
		}
	}

	if !(*live) {
		fmt.Printf("Dry run. Use -live to move files.\n")
	}
}

//
// Helpers
//

var bannedParts = map[string]bool{
	"truepdf": true,
}

// Matches: 3, 12, 3rd, 8th, 29th, etc.
var dayRX = regexp.MustCompile(`^[0-9]{1,2}[A-Za-z]{0,2}$`)

var months = map[string]string{
	"january":   "01",
	"february":  "02",
	"march":     "03",
	"april":     "04",
	"may":       "05",
	"june":      "06",
	"july":      "07",
	"august":    "08",
	"september": "09",
	"october":   "10",
	"november":  "11",
	"december":  "12",
}

// Splits on a "." or "-" so we handle the weird hyphenation between words and
// date.
var partsRX = regexp.MustCompile(`[.-]`)

var yearRX = regexp.MustCompile(`^[0-9]{4}$`)

func extractDate(original string) ([]string, string) {
	original = strings.ToLower(original)
	parts := partsRX.Split(original, -1)
	l := len(parts)

	if l > 7 && isDay(parts[l-7]) && isMonth(parts[l-6]) && isYear(parts[l-2]) {
		// The.Economist.11TH.July.17TH.TruePDF-July.2015.pdf
		return parts[0 : l-7],
			parts[l-2] + "-" + extractMonth(parts[l-6]) + "-" + extractDay(parts[l-7])
	} else if l > 5 && isDay(parts[l-5]) && isMonth(parts[l-3]) && isYear(parts[l-2]) {
		// The.New.Yorker.TruePDF-22.29.December.2014.pdf
		return parts[0 : l-5],
			parts[l-2] + "-" + extractMonth(parts[l-3]) + "-" + extractDay(parts[l-5])
	} else if l > 4 && isDay(parts[l-4]) && isMonth(parts[l-3]) && isYear(parts[l-2]) {
		// The.Economist-09.August.2014.pdf
		return parts[0 : l-4],
			parts[l-2] + "-" + extractMonth(parts[l-3]) + "-" + extractDay(parts[l-4])
	}

	return []string{}, ""
}

func extractDay(part string) string {
	if len(part) == 4 {
		// 22nd
		return part[0:2]
	} else if len(part) == 3 {
		// 3rd
		return "0" + part[0:1]
	} else if len(part) == 2 {
		// 19
		return part
	}

	// 4
	return "0" + part
}

func extractMonth(part string) string {
	return months[part]
}

func isAlreadyRenamed(original string) bool {
	dotParts := strings.Split(original, ".")
	for _, part := range dotParts {
		if isYear(part) {
			return false
		}
	}

	return true
}

func isDay(part string) bool {
	return dayRX.MatchString(part)
}

func isMonth(part string) bool {
	for month, _ := range months {
		if part == month {
			return true
		}
	}
	return false
}

func isYear(part string) bool {
	return yearRX.MatchString(part)
}

func rename(original string) string {
	// We should only be passed a filename.
	if strings.Contains(original, "/") {
		panic("Expected file but got path")
	}

	// Determine whether the file's already been renamed. To do so we look for
	// a cluster of digits that looks like a year. For example "2014" in
	// "The.Economist-09.August.2014.pdf".
	if isAlreadyRenamed(original) {
		return original
	}

	parts, dateString := extractDate(original)

	// Not a name that we can handle.
	if dateString == "" {
		return original
	}

	var correctedParts []string
	for _, part := range parts {
		if _, ok := bannedParts[part]; !ok {
			correctedParts = append(correctedParts, part)
		}
	}

	// If there's a "the" on the front, moved it to the back.
	if correctedParts[0] == "the" {
		correctedParts = append(correctedParts[1:len(correctedParts)], "the")
	}

	correctedParts = append(correctedParts, dateString, "pdf")

	return strings.Join(correctedParts, ".")
}
