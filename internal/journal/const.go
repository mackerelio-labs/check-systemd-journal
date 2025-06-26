package journal

//go:generate go tool stringer -output const_string.go -type=Facility,Priority

type Priority int16

const (
	emerg Priority = iota
	alert
	crit
	err
	warning
	notice
	info
	debug
)

type Facility int16

const (
	kern Facility = iota
	user
	mail
	daemon
	auth
	syslog
	lpr
	news
	uucp
	cron
	authpriv
	ftp
)

const (
	local0 Facility = iota + 16
	local1
	local2
	local3
	local4
	local5
	local6
	local7
)

func FacilityNames() (f1 []string, f2 []string) {
	for i := 0; i <= 11; i++ {
		f1 = append(f1, Facility(i).String())
	}
	for i := 16; i <= 23; i++ {
		f2 = append(f2, Facility(i).String())
	}
	return
}

func PriorityNames() (p []string) {
	for i := 0; i <= 7; i++ {
		p = append(p, Priority(i).String())
	}
	return
}
