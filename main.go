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
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "File does not exist: %v\n", file)
			os.Exit(1)
		}

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

// Matches post-rename date: 2017-04-01
var dateRX = regexp.MustCompile(`^[1-9][0-9]{3}-[01][0-9]-[0123][0-9]$`)

// Matches: 3, 12, 3rd, 8th, 29th, etc.
var dayRX = regexp.MustCompile(`^[0-9]{1,2}[A-Za-z]{0,2}$`)

var months = map[string]string{
	"jan":       "01",
	"january":   "01",
	"feb":       "02",
	"february":  "02",
	"mar":       "03",
	"march":     "03",
	"apr":       "04",
	"april":     "04",
	"may":       "05",
	"jun":       "06",
	"june":      "06",
	"jul":       "07",
	"july":      "07",
	"aug":       "08",
	"august":    "08",
	"sep":       "09",
	"september": "09",
	"oct":       "10",
	"october":   "10",
	"nov":       "11",
	"november":  "11",
	"dec":       "12",
	"december":  "12",
}

// Splits on a "." or "-" so we handle the weird hyphenation between words and
// date.
var partsRX = regexp.MustCompile(`[()[\] .-]+`)

var yearRX = regexp.MustCompile(`^[0-9]{4}$`)

func extractDate(original string) ([]string, string) {
	original = strings.ToLower(original)

	// Note that parts is all the components of the filename *including* the
	// extension.
	parts := partsRX.Split(original, -1)
	l := len(parts)

	fmt.Printf("parts: %+v (length %d)\n", parts, len(parts))

	if l > 11 && isDay(parts[l-5]) && isMonth(parts[l-4]) && isYear(parts[l-3]) {
		// The Mercantilist (UK) - Vol. 444 No. 9314 [24 Sep 2022] (TruePDF).pdf
		return parts[0 : l-9],
			parts[l-3] + "-" + extractMonth(parts[l-4]) + "-" + extractDay(parts[l-5])
	} else if l > 7 && isDay(parts[l-7]) && isMonth(parts[l-6]) && isYear(parts[l-2]) {
		// The.Mercantilist.11TH.July.17TH.TruePDF-July.2015.pdf
		return parts[0 : l-7],
			parts[l-2] + "-" + extractMonth(parts[l-6]) + "-" + extractDay(parts[l-7])
	} else if l > 7 && isMonth(parts[l-7]) && isDay(parts[l-6]) && isYear(parts[l-2]) {
		// The.Mercantilist.Europe.July.29.TruePDF-4.August.2017.pdf
		return parts[0 : l-7],
			parts[l-2] + "-" + extractMonth(parts[l-7]) + "-" + extractDay(parts[l-6])
	} else if l > 6 && isMonth(parts[l-6]) && isDay(parts[l-5]) && isYear(parts[l-2]) {
		// The.Mercantilist.Europe.April.1.7.TruePDF-2017.pdf
		return parts[0 : l-6],
			parts[l-2] + "-" + extractMonth(parts[l-6]) + "-" + extractDay(parts[l-5])
	} else if l > 5 && isDay(parts[l-5]) && isMonth(parts[l-3]) && isYear(parts[l-2]) {
		// The.Old.Yorker.TruePDF-22.29.December.2014.pdf
		return parts[0 : l-5],
			parts[l-2] + "-" + extractMonth(parts[l-3]) + "-" + extractDay(parts[l-5])
	} else if l > 4 && isDay(parts[l-4]) && isMonth(parts[l-3]) && isYear(parts[l-2]) {
		// The.Mercantilist-09.August.2014.pdf
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
		if isDate(part) {
			return true
		}
	}

	return false
}

func isDay(part string) bool {
	return dayRX.MatchString(part)
}

func isDate(part string) bool {
	return dateRX.MatchString(part)
}

func isMonth(part string) bool {
	_, ok := months[part]
	return ok
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
	// "The.Mercantilist-09.August.2014.pdf".
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
		correctedParts = append(correctedParts[1:], "the")
	}

	ext := path.Ext(original)
	ext = ext[1:] // strip leading dot

	correctedParts = append(correctedParts, dateString, ext)
	return strings.Join(correctedParts, ".")
}
