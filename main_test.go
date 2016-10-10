package main

import (
	"testing"
)

func TestRename(t *testing.T) {
	for _, tuple := range renameTestCases {
		original := tuple[0]
		expected := tuple[1]
		actual := rename(original)

		if actual != expected {
			t.Errorf(`Expected "%v" to be renamed to "%v", but was renamed to "%v"`,
				original, expected, actual)
		}
	}
}

var renameTestCases = [][]string{
	// test for empty
	{"", ""},

	// files that are totally unrelated
	{"hello", "hello"},
	{"hello.pdf", "hello.pdf"},
	{"2016", "2016"},
	{"2016.pdf", "2016.pdf"},
	{"hello.2016.pdf", "hello.2016.pdf"},

	// files that have already been renamed
	{"economist.the.2014-08-09.pdf", "economist.the.2014-08-09.pdf"},
	{"new.yorker.the.2015-10-12.pdf", "new.yorker.the.2015-10-12.pdf"},

	// files that should be renamed (in various formats)
	{"The.Economist-09.August.2014.pdf", "economist.the.2014-08-09.pdf"},
	{"The.Economist-25.July.2015.pdf", "economist.the.2015-07-25.pdf"},
	{"The.Economist-30.January.2016.pdf", "economist.the.2016-01-30.pdf"},

	{"The.Economist.TruePDF-16.January.2016.pdf", "economist.the.2016-01-16.pdf"},

	{"The.Economist.Europe.TruePDF-30.July.2016.pdf", "economist.europe.the.2016-07-30.pdf"},

	{"The.Economist.USA-3.September.2016.pdf", "economist.usa.the.2016-09-03.pdf"},
	{"The.Economist.USA.TruePDF-10.September.2016.pdf", "economist.usa.the.2016-09-10.pdf"},

	{"The.Economist.11TH.July.17TH.TruePDF-July.2015.pdf", "economist.the.2015-07-11.pdf"},
	{"The.Economist.1ST.November.7TH.TruePDF-November.2014.pdf", "economist.the.2014-11-01.pdf"},
	{"The.Economist.28TH.November.4TH.TruePDF-December.2015.pdf", "economist.the.2015-11-28.pdf"},

	{"The.New.Yorker-01.September.2014.pdf", "new.yorker.the.2014-09-01.pdf"},
	{"The.New.Yorker-18.January.2016.pdf", "new.yorker.the.2016-01-18.pdf"},

	{"The.New.Yorker.TruePDF-12.October.2015.pdf", "new.yorker.the.2015-10-12.pdf"},
	{"The.New.Yorker.TruePDF-22.29.December.2014.pdf", "new.yorker.the.2014-12-22.pdf"},
}
