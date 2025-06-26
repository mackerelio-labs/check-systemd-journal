# check-systemd-journal

[日本語](#説明)

**check-systemd-journal** checks journals whether new logs are available, then reports them. It can filter logs with any of systemd unit, priority, syslog facility and/or regexp.

behaves as like **grep(1)**. Logs will be printed to standard output, reporting via exit status.

## synopsis

```
Usage of check-systemd-journal:
  -check int
    	threshold[=NUM] (default 1)
  -e PATTERN
    	PATTERN(s) to search for
  -facility FACILITY
    	FACILITY(s) name
  -icase
    	Run a case insensitive match
  -priority PRIORITY
    	PRIORITY name
  -quiet
    	quiet
  -state-file file
    	state file path
  -unit unit
    	unit name
  -user
    	user scope
  -v PATTERN
    	NOT matched PATTERN(s) to search for

Constants available in PRIORITY
	emerg alert crit err warning notice info debug
Constants available in FACILITY
	kern user mail daemon auth syslog lpr news uucp cron authpriv ftp
	local0 local1 local2 local3 local4 local5 local6 local7
```


*-state-file* option is passed, **check-systemd-journal** saves a last cursor position to *FILE*. Subsequent execution after first **check-systemd-journal** execution, they will use the cursor to skip until new available logs.

*-unit* option selects only logs belongs to *UNIT*.

*-priority* option selects logs by *PRIORITY* or higher. The available options are see [PRIORITY](#Priorities)

*-facility* option selects logs by *FACILITY*. If one or more *-facility* options, all *FACILITY*s combines with **OR** operator.　The available options are see [FACILITY](#Facilities)

*-e* option selects logs matched by *PATTERN*. If one or more *-e* options, all *PATTERN*s combines with **AND** operator.

*-icase* option indicates *PATTERN*s are case-insensitive.

*-v* option selects logs matched **NOT** by *PATTERN*. If one or more *-v* options, all *PATTERN*s combines with **OR** operator.

*-quiet* option suppress outputs of selected logs.

*-check* option reached NUM times, default by 1, **check-systemd-journal** reports a critical alert.

### example

```
check-systemd-journal -e pam_unix -priority info -facility authpriv -state-file statefile --user
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


# License

© 2025 Hatena Co., Ltd.

Apache License (see LICENSE file)

---

# 説明

**check-systemd-journal** は、systemd ジャーナルに新しいログがあるかどうかを確認し、報告します。
systemd の ユニット、優先度、syslog ファシリティ、正規表現のいずれかでログをフィルタリングできます。

**grep(1)** と同様に動作します。ログは標準出力に出力され、終了ステータスで報告されます。

## 概要

```
Usage of check-systemd-journal:
  -check int
    	threshold[=NUM] (default 1)
  -e PATTERN
    	PATTERN(s) to search for
  -facility FACILITY
    	FACILITY(s) name
  -icase
    	Run a case insensitive match
  -priority PRIORITY
    	PRIORITY name
  -quiet
    	quiet
  -state-file file
    	state file path
  -unit unit
    	unit name
  -user
    	user scope
  -v PATTERN
    	NOT matched PATTERN(s) to search for

Constants available in PRIORITY
	emerg alert crit err warning notice info debug
Constants available in FACILITY
	kern user mail daemon auth syslog lpr news uucp cron authpriv ftp
	local0 local1 local2 local3 local4 local5 local6 local7
```
*-state-file* オプションが指定されると、**check-systemd-journal** は最後のカーソル位置を *FILE* に保存します。
最初の **check-systemd-journal** 実行後の以降の実行では、カーソル位置を使用して新しい利用可能なログまでスキップします。

*-unit* オプションは、*UNIT* に属するログのみを選択します。

*-priority* オプションは、*PRIORITY* 以上の優先順位でログを選択します。 使用できる選択肢については [PRIORITY](#Priorities) を参照ください。

*-facility* オプションは、*FACILITY* でログを選択します。1 つ以上の *-facility* オプションが指定された場合、すべての *FACILITY* は **OR** 演算子で結合されます。 使用できる選択肢については [FACILITY](#Facilities) を参照ください。

*-e* オプションは、*PATTERN* に一致するログを選択します。1 つ以上の *-e* オプションが指定された場合、すべての *PATTERN* は **AND** 演算子で結合されます。

*-icase* オプションは、*PATTERN* が大文字と小文字を区別しないことを示します。

*-v* オプションは、*PATTERN* に **NOT** 一致するログを選択します。*-v* オプションを1つ以上指定した場合、すべての *PATTERN* は **OR** 演算子で結合されます。

*-quiet* オプションは、選択したログの出力を抑制します。

*-check* オプションが NUM 回 (デフォルトは 1 回) に達すると、**check-systemd-journal** は重大なアラート(Critical)を報告します。

### 使用例

```
check-systemd-journal -e pam_unix -priority info -facility authpriv -state-file statefile --user
```

# ライセンス

© 2025 Hatena Co., Ltd.

Apache License (詳細は LICENSE ファイルを参照ください)
