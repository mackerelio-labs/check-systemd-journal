package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/mackerelio-labs/check-systemd-journal/internal/journal"
)

type sliceString []string

func (i *sliceString) String() string {
	return fmt.Sprint(*i)
}
func (i *sliceString) Set(v string) error {
	*i = append(*i, v)
	return nil
}

var (
	patterns        sliceString
	ignorePatterns  sliceString
	stateFile       string
	unitName        string
	priorityName    string
	facilitiesName  sliceString
	quiet           bool
	caseInsensitive bool
	journalUser     bool
	threshold       int
)

func init() {
	flag.Var(&patterns, "e", "`PATTERN`(s) to search for")
	flag.Var(&ignorePatterns, "v", "NOT matched `PATTERN`(s) to search for")
	flag.StringVar(&stateFile, "state-file", "", "state `file` path")
	flag.StringVar(&unitName, "unit", "", "`unit` name")
	flag.StringVar(&priorityName, "priority", "", "priority name")
	flag.Var(&facilitiesName, "facility", "facility(s) name")
	flag.BoolVar(&quiet, "quiet", false, "quiet")
	flag.BoolVar(&caseInsensitive, "icase", false, "Run a case insensitive match")
	flag.BoolVar(&journalUser, "user", false, "user scope")
	flag.IntVar(&threshold, "check", 1, "threshold[=NUM]")

	flag.Parse()
}

func regCompileWithCase(ptn string, caseInsensitive bool) (*regexp.Regexp, error) {
	if caseInsensitive {
		ptn = "(?i)" + ptn
	}
	return regexp.Compile(ptn)
}

func facilityNameCheck(v string) (journal.Facility, error) {
	for i := 0; i <= 23; i++ {
		if journal.Facility(i).String() == v {
			return journal.Facility(i), nil
		}
	}

	return -1, fmt.Errorf("invalid facility name : %q", v)
}

func priorityNameCheck(v string) (journal.Priority, error) {
	for i := 0; i <= 7; i++ {
		if journal.Priority(i).String() == v {
			return journal.Priority(i), nil
		}
	}

	return -1, fmt.Errorf("invalid priority name : %q", v)
}

func isCritical(n int) bool {
	return (threshold > 0 && (n) >= threshold)
}

func isWarning(n int) bool {
	return (threshold > 0 && (n) > 0 && (n) < threshold)
}

func main() {
	var regs, ignoreRegs []*regexp.Regexp
	for _, pattern := range patterns {
		reg, err := regCompileWithCase(pattern, caseInsensitive)
		if err != nil {
			log.Fatalf("pattern is invalid : %q", pattern)
		}
		regs = append(regs, reg)
	}

	for _, pattern := range ignorePatterns {
		reg, err := regCompileWithCase(pattern, caseInsensitive)
		if err != nil {
			log.Fatalf("pattern is invalid : %q", pattern)
		}
		ignoreRegs = append(ignoreRegs, reg)
	}

	var cursor string
	if stateFile != "" {
		if _, err := os.Stat(stateFile); !os.IsNotExist(err) {
			b, err := os.ReadFile(stateFile)
			if err != nil {
				log.Fatal(err)
			}
			cursor = string(b)
		}
	}

	scope := journal.FlagSystem
	if journalUser {
		scope = journal.FlagUser
	}

	var facilities []journal.Facility
	for i := range facilitiesName {
		facility, err := facilityNameCheck(facilitiesName[i])
		if err != nil {
			log.Fatal(err)
		}
		facilities = append(facilities, facility)
	}

	var priority *journal.Priority
	if priorityName != "" {
		var err error
		p, err := priorityNameCheck(priorityName)
		if err != nil {
			log.Fatal(err)
		}
		priority = &p
	}

	matched, cursor, err := readJournal(arg{
		Scope:      scope,
		Priority:   priority,
		Facilities: facilities,
		Unit:       unitName,
		regs:       regs,
		ignoreRegs: ignoreRegs,
		quiet:      quiet,
		Cursor:     cursor,
	})

	if err != nil {
		log.Fatal(err)
	}

	if stateFile != "" && cursor != "" {
		err := os.WriteFile(stateFile, []byte(cursor), 0666)
		if err != nil {
			log.Fatal(err)
		}
	}

	if isCritical(matched) {
		os.Exit(2)
	} else if isWarning(matched) {
		os.Exit(1)
	}
	os.Exit(0)
}

type arg struct {
	Scope      journal.Flag
	Unit       string
	Facilities []journal.Facility
	Priority   *journal.Priority

	regs       []*regexp.Regexp
	ignoreRegs []*regexp.Regexp

	quiet bool

	Cursor string
}

func readJournal(arg arg) (nmatched int, cursor string, err error) {
	j, err := journal.New(arg.Scope)
	if err != nil {
		return
	}
	defer j.Close()

	if arg.Unit != "" {
		// generates `(_SYSTEMD_UNIT=<u> OR UNIT=<u>) AND ...` matches

		buf := fmt.Sprintf("_SYSTEMD_UNIT=%s", arg.Unit)
		j.AddMatch(buf)
		j.AddDisjunction()

		buf2 := fmt.Sprintf("UNIT=%s", arg.Unit)
		j.AddMatch(buf2)
		j.AddConjunction()
	}

	if arg.Priority != nil {
		for i := 0; i <= int(*arg.Priority); i++ {
			buf := fmt.Sprintf("PRIORITY=%d", i)
			j.AddMatch(buf)
		}
	}
	for p := range arg.Facilities {
		buf := fmt.Sprintf("SYSLOG_FACILITY=%d", arg.Facilities[p])
		j.AddMatch(buf)
	}

	if arg.Cursor != "" {
		if err = j.SeekCursor(arg.Cursor); err != nil {
			return
		}

		var n int
		// A position pointed with cursor has been read in previous operation.
		n, err = j.SeekNext()
		if err != nil {
			return
		}
		if n == 0 { // no more data
			return
		}

		if err = j.TestCursor(arg.Cursor); err != nil {
			return
		}
	}

	var n int

	var foundData bool
	n, err = j.SeekNext()
	if err != nil {
		return
	}

	for n > 0 {
		foundData = true

		var s, u string
		s, err = j.Getdata("MESSAGE")
		if err != nil {
			return
		}

		if matchLine(s, arg.regs, arg.ignoreRegs) {
			if !quiet {
				u, err = j.Getdata("UNIT")
				if err != nil {
					return
				}
				if u == "" {
					u, err = j.Getdata("_SYSTEMD_UNIT")
					if err != nil {
						return
					}
				}

				if u != "" {
					fmt.Printf("%s:", u)
				} else {
					fmt.Print("(null):")
				}
				fmt.Printf("%s\n", s)
			}
			nmatched++
		}

		n, err = j.SeekNext()
		if err != nil {
			return
		}
	}

	if !foundData {
		return
	}

	cursor, err = j.GetCursor()

	return
}

func matchLine(s string, regs, ignoreRegs []*regexp.Regexp) bool {
	for _, reg := range regs {
		if len(reg.FindStringSubmatch(s)) == 0 {
			return false
		}
	}

	if len(ignoreRegs) == 0 {
		return true
	}

	exclude := true
	for _, reg := range ignoreRegs {
		if len(reg.FindStringSubmatch(s)) > 0 {
			exclude = false
		}
	}
	return exclude
}
