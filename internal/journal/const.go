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
