# check-systemd-journal

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

