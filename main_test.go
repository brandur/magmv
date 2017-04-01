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

	// other file types
	{"The.Mercantilist-09.August.2014.nzb", "mercantilist.the.2014-08-09.nzb"},

	// files that have already been renamed
	{"mercantilist.the.2014-08-09.pdf", "mercantilist.the.2014-08-09.pdf"},
	{"new.yorker.the.2015-10-12.pdf", "new.yorker.the.2015-10-12.pdf"},

	// files that should be renamed (in various formats)
	{"The.Mercantilist-09.August.2014.pdf", "mercantilist.the.2014-08-09.pdf"},
	{"The.Mercantilist-25.July.2015.pdf", "mercantilist.the.2015-07-25.pdf"},
	{"The.Mercantilist-30.January.2016.pdf", "mercantilist.the.2016-01-30.pdf"},

	{"The.Mercantilist.TruePDF-16.January.2016.pdf", "mercantilist.the.2016-01-16.pdf"},

	{"The.Mercantilist.Europe.TruePDF-30.July.2016.pdf", "mercantilist.europe.the.2016-07-30.pdf"},

	{"The.Mercantilist.USA-3.September.2016.pdf", "mercantilist.usa.the.2016-09-03.pdf"},
	{"The.Mercantilist.USA.TruePDF-10.September.2016.pdf", "mercantilist.usa.the.2016-09-10.pdf"},

	{"The.Mercantilist.11TH.July.17TH.TruePDF-July.2015.pdf", "mercantilist.the.2015-07-11.pdf"},
	{"The.Mercantilist.1ST.November.7TH.TruePDF-November.2014.pdf", "mercantilist.the.2014-11-01.pdf"},
	{"The.Mercantilist.28TH.November.4TH.TruePDF-December.2015.pdf", "mercantilist.the.2015-11-28.pdf"},
	{"The.Mercantilist.Europe.April.1.7.TruePDF-2017.pdf", "mercantilist.europe.the.2017-04-01.pdf"},

	{"The.Old.Yorker-01.September.2014.pdf", "old.yorker.the.2014-09-01.pdf"},
	{"The.Old.Yorker-18.January.2016.pdf", "old.yorker.the.2016-01-18.pdf"},

	{"The.Old.Yorker.TruePDF-12.October.2015.pdf", "old.yorker.the.2015-10-12.pdf"},
	{"The.Old.Yorker.TruePDF-22.29.December.2014.pdf", "old.yorker.the.2014-12-22.pdf"},
}
