# check-systemd-journal

**check-systemd-journal** checks journals whether new logs are available, then reports them. It can filter logs with any of systemd unit, priority, syslog facility and/or regexp.

behaves as like **grep(1)**. Logs will be printed to standard output, reporting via exit status.

## synopsis

```
Usage of ./check-systemd-journal:
  -check int
    	threshold[=NUM] (default 1)
  -e PATTERN
    	PATTERN(s) to search for
  -facility value
    	facility(s) name
  -icase
    	Run a case insensitive match
  -priority string
    	priority name
  -quiet
    	quiet
  -state-file string
    	state file path
  -unit string
    	unit
  -user
    	user scope?
  -v PATTERN
    	NOT matched PATTERN(s) to search for
```


*-state-file* option is passed, **check-systemd-journal** saves a last cursor position to *FILE*. Subsequent execution after first **check-systemd-journal** execution, they will use the cursor to skip until new available logs.

*-unit* option selects only logs belongs to *UNIT*.

*-priority* option selects logs by *PRIORITY* or higher.

*-facility* option selects logs by *FACILITY*. If one or more *-facility* options, all *FACILITY*s combines with **OR** operator.

*-e* option selects logs matched by *PATTERN*. If one or more *-e* options, all *PATTERN*s combines with **AND** operator.

*-icase* option indicates *PATTERN*s are case-insensitive.

*-v* option selects logs matched **NOT** by *PATTERN*. If one or more *-v* options, all *PATTERN*s combines with **AND** operator.

*-quiet* option suppress outputs of selected logs.

*-check* option indicates to behave as Sensu plugin mode. If selected logs by above options reached NUM times, default by 1, **check-systemd-journal** reports a critical alert.

### example

``` 
./check-systemd-journal -e pam_unix -priority info -facility authpriv -state-file statefile --user
```

### Priorities

- emerg
- alert
- crit
- err
- warning
- notice
- info
- debug

### Facilities

- kern
- user
- mail
- daemon
- auth
- syslog
- lpr
- news
- uucp
- cron
- authpriv
- ftp
- local0
- local1
- local2
- local3
- local4
- local5
- local6
- local7

