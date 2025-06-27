# check-systemd-journal

[日本語](#説明)

**check-systemd-journal** is a tool that tracks systemd journal logs and checks for new logs based on filters. It is intended to be run as a check plugin for [mackerel-agent](https://mackerel.io/docs/entry/howto/install-agent), but it can also be used as a standalone command.

You can filter logs based on systemd units, priorities, syslog facilities, or regular expressions. Regular expressions follow the Go language regexp library's [regular expression notation](https://pkg.go.dev/regexp/syntax).

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

- *-check int*: If valid results after filtering are obtained consecutively the specified number of times (default is 1), **check-systemd-journal** returns status code 2. When used as a mackerel-agent plugin, this is treated as a Critical alert.
- *-e PATTERN*: Performs a regular expression match on each line of the journal and selects the matching lines. This option can be specified multiple times and is evaluated as an AND (all must match) condition. To evaluate as an OR (any match is selected) condition, use `|` within the regular expression.
- *-facility FACILITY*: Selects lines from the journal that match the specified facility level. This option can be specified multiple times and is evaluated as OR (selects any matching items). For a list of strings that can be specified for facility, refer to [FACILITY](#Facilities).
- *-icase*: When this option is specified, the *-e* option and *-v* option do not distinguish between uppercase and lowercase letters in regular expressions.
- *-priority PRIORITY*: Selects lines in the journal with a priority level equal to or higher than the specified priority. For details on the strings that can be specified for priority, refer to [PRIORITY](#Priorities).
- *-quiet*: Suppresses the output of the selected results to standard output. Specify this option when you only want to use the exit status code.
- *-state-file file*: Saves the last cursor position in the journal to a state file. On the next execution, the collection will skip to the newly available log starting from the saved cursor position. If this option is not specified, the collection will always start from the beginning of the saved journal. Therefore, this option must be specified when using this as a mackerel-agent plugin (e.g., `/var/tmp/mackerel-agent/ssh.state`). When using the `--user` option to target the journal of a specific user's service, the file path must be writable by that user.
- *-unit unit*: Selects only the lines in the journal that belong to the specified systemd unit.
- *-user*: Specify when the unit is linked to a user rather than the system. Since the units of the user who called *check-systemd-journal* are targeted, when using it as a check plugin for mackerel-agent, you must specify the user parameter in the mackerel-agent configuration file (see the [Configuration items](https://mackerel.io/docs/entry/custom-checks#items) in mackerel-agent.conf). The state file path for *-state-file* must also be writable by the user.
- *-v PATTERN*: Performs regular expression matching on each line of the journal and selects those that do NOT match. This option can be specified multiple times and is evaluated as OR (select if any match is found).

### example

Here is an example of using check-systemd-journal to monitor journal logs with mackerel-agent.

To generate an alert when the string `authentication failure` is found in the journal log of the ssh unit, add the following to mackerel-agent.conf and restart mackerel-agent.

```
[plugin.checks.ssh_authentication_failure]
command = ["check-systemd-journal", "-unit", "ssh", "-e", "authentication failure", "--state-file", "/var/tmp/mackerel-agent/ssh.state"]
```

The monitoring settings to search for all journal logs of units with a priority of `err` or higher that contain the strings `failed` or `error` (case-insensitive) and exclude lines containing the string `debug` are as follows.

```
[plugin.checks.failed_or_error]
command = ["check-systemd-journal", "--priority", "err", "-e", "failed|error", "-icase", "-v", "debug", "--state-file", "/var/tmp/mackerel-agent/failed_or_error.state"]
```

In addition, examples of standalone execution are also shown.

```
check-systemd-journal -e pam_unix -priority info -facility authpriv -state-file statefile --user
```

In this example, logs associated with the user that have a priority of `info` or higher, a facility of `authpriv`, and contain the string `pam_unix` are selected and displayed. The cursor position is written to the state file `statefile`.

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

Apache License (see [LICENSE](./LICENSE) file)

---

# 説明

**check-systemd-journal** は、systemd ジャーナルのログを追跡し、フィルタに基づいて新規のログがあるかどうかを検査できるツールです。[mackerel-agent](https://mackerel.io/ja/docs/entry/howto/install-agent) のチェックプラグインとして実行することを想定していますが、単体のコマンドとしても利用可能です。

systemd のユニット、プライオリティ、syslog ファシリティ、正規表現のいずれかでログをフィルタリングできます。正規表現は Go 言語の regexp ライブラリの[正規表現記法](https://pkg.go.dev/regexp/syntax)に従います。

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

- *-check 回数*：フィルタリングされた有効な結果が指定の回数（デフォルトは 1 回）連続で繰り返されると、**check-systemd-journal** はステータスコード 2 を返します。mackerel-agent のプラグインとして利用している場合、これは重大なアラート（Critical）として扱われます。
- *-e 正規表現*：ジャーナルの各行に対して正規表現マッチを行い、マッチしたものを選出します。このオプションは複数指定することができ、AND（すべて一致したものが選出）で評価されます。OR（いずれか一致したものが選出）評価をするには正規表現内で `|` を利用してください。
- *-facility ファシリティ*：ジャーナルの各行のうち指定のファシリティレベルのものを選出します。このオプションは複数指定することができ、OR（いずれか一致したものが選出）評価されます。ファシリティに指定できる文字列については [FACILITY](#Facilities) を参照してください。
- *-icase*：このオプションを指定した場合、*-e* オプションおよび *-v* オプションの正規表現における大文字・小文字を区別しないようにします。
- *-priority プライオリティ*：ジャーナルの各行のうち指定のプライオリティレベル以上のものを選出します。 プライオリティに指定できる文字列については [PRIORITY](#Priorities) を参照してください。
- *-quiet*：選出した結果の標準出力への出力を抑制します。終了ステータスコードのみを利用したいときに指定してください。
- *-state-file 状態ファイルパス*：ジャーナルの最後のカーソル位置を状態ファイルに保存します。次回の実行時には、保存されたカーソル位置を使用して新しく利用可能なログまでスキップされます。これを指定しない場合、保存されているジャーナルの先頭から常に収集されることになるので、mackerel-agent のプラグインとして利用する場合は必ず指定してください（たとえば `/var/tmp/mackerel-agent/ssh.state` など）。`--user` オプションを利用して特定のユーザーのサービスのジャーナルを対象とする場合、そのユーザーが書き込めるファイルパスである必要があります。
- *-unit ユニット*：ジャーナルの各行のうち指定の systemd ユニットに属するもののみを選出します。
- *-user*：ユニットがシステムではなくユーザーに紐づく場合に指定します。*check-systemd-journal* を呼び出したユーザーのユニットが対象となるため、mackerel-agent のチェックプラグインとして使う場合は mackerel-agent の設定ファイルで user パラメータを指定する必要があります（mackerel-agent.conf の[設定項目](https://mackerel.io/ja/docs/entry/custom-checks#items)を参照してください。 *-state-file* の状態ファイルパスをそのユーザーが書き込める場所にする必要もあります）。
- *-v 正規表現*：ジャーナルの各行に対して正規表現マッチを行い、マッチ**しない**ものを選出します。このオプションは複数指定することができ、OR（いずれか一致したら選出しない）で評価されます。

### 使用例

check-systemd-journal を使ってジャーナルログ監視を mackerel-agent で行う例を示します。

ssh ユニットのジャーナルログで `authentication failure` という文字列を発見したときにアラートを発生させるには、以下のように mackerel-agent.conf に記述し、mackerel-agent を再起動します。

```
[plugin.checks.ssh_authentication_failure]
command = ["check-systemd-journal", "-unit", "ssh", "-e", "authentication failure", "--state-file", "/var/tmp/mackerel-agent/ssh.state"]
```

プライオリティが err 以上のすべてのユニットのジャーナルログを対象に `failed` か `error` という文字列を大文字小文字問わず含んでいるものを探し、その中から `debug` という文字列を含む行は除くという監視の設定は、次のようになります。

```
[plugin.checks.failed_or_error]
command = ["check-systemd-journal", "--priority", "err", "-e", "failed|error", "-icase", "-v", "debug", "--state-file", "/var/tmp/mackerel-agent/failed_or_error.state"]
```

このほか、単体で実行する例も示しておきます。

```
check-systemd-journal -e pam_unix -priority info -facility authpriv -state-file statefile --user
```

ユーザーに紐づくジャーナルのうち、プライオリティが info 以上、ファシリティが authpriv で、かつ `pam_unix` という文字列を含むログを抽出し、表示します。カーソル位置を状態ファイル `statefile` に書き出します。

# ライセンス

© 2025 Hatena Co., Ltd.

Apache License（詳細は [LICENSE](./LICENSE) ファイルを参照してください）
