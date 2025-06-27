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
  -facility value
    	facility(s) name
  -icase
    	Run a case insensitive match
  -priority string
    	priority name
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
```


*-state-file* option is passed, **check-systemd-journal** saves a last cursor position to *FILE*. Subsequent execution after first **check-systemd-journal** execution, they will use the cursor to skip until new available logs.

*-unit* option selects only logs belongs to *UNIT*.

*-priority* option selects logs by *PRIORITY* or higher. The available options are see [PRIORITY](#Priorities)

*-facility* option selects logs by *FACILITY*. If one or more *-facility* options, all *FACILITY*s combines with **OR** operator.　The available options are see [FACILITY](#Facilities)

*-e* option selects logs matched by *PATTERN*. If one or more *-e* options, all *PATTERN*s combines with **AND** operator.

*-icase* option indicates *PATTERN*s are case-insensitive.

*-v* option selects logs matched **NOT** by *PATTERN*. If one or more *-v* options, all *PATTERN*s combines with **AND** operator.

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

**check-systemd-journal** は、systemd ジャーナルのログを追跡し、フィルタに基づいて新規のログがあるかどうかを検査できるツールです。[mackerel-agent](https://mackerel.io/ja/docs/entry/howto/install-agent) のチェックプラグインとして実行することを想定していますが、単体のコマンドとしても利用可能です。

systemd のユニット、プライオリティ、syslog ファシリティ、正規表現のいずれかでログをフィルタリングできます。正規表現は Go 言語の regexp ライブラリの[正規表現記法](https://pkg.go.dev/regexp/syntax)に従います。

## 概要

```
Usage of check-systemd-journal:
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
  -state-file file
    	state file path
  -unit unit
    	unit name
  -user
    	user scope
  -v PATTERN
    	NOT matched PATTERN(s) to search for
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

Apache License (詳細は [LICENSE](./LICENSE) ファイルを参照してください)
