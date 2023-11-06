# imap-mailstat-exporter

[![publish](https://github.com/bt909/imap-mailstat-exporter/actions/workflows/publish_latest_oci_image.yaml/badge.svg)](https://github.com/bt909/imap-mailstat-exporter/actions/workflows/publish_latest_oci_image.yaml)
 [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

This is a prometheus exporter which gives you metrics for how many emails you have in your INBOX and in additional configured folders.  

Connections to IMAP are only TLS enrypted supported, either via TLS or STARTTLS.

> [!NOTE]
> This exporter is in early development and at the moment highly adjusted for my personal usecase. As it not reached 1.0.0 yet, there may are breaking changes at any time. Keep an eye on the [CHANGELOG](https://github.com/bt909/imap-mailstat-exporter/blob/main/CHANGELOG.md) for information.

As this exporter is using [exporter-toolkit](https://github.com/prometheus/exporter-toolkit) since 0.1.0 (not yet released), you can also configure basic auth, or TLS secured connection to the exporter using http/2, for more information visit the [configuration page of exporter-toolkit](https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md).

The exporter provides 14 metrics, three main metrics are provided for all accounts, one metric can be enabled using a feature flag `--oldestunseen.feature`, eigth metrics are quota related and only provided if the server supports imap quota and two informatial metrics.

If your account supports quota you can see in loglevel INFO (default) with the following log entry:

```output
ts=2023-10-30T22:12:27.376Z caller=valuecollector.go:266 level=info msg="Fetching quota related metrics" address=jane.doe@example.com

```

The exposed metrics were the following in version 0.0.1 and can be enabled by using commandline flag `--migration.mode`:

`imap_mailstat_mails_all_quantity`  
`imap_mailstat_mails_unseen_quantity`  
`imap_mailstat_mails_levelquotaavail_quantity` (only imap with quota support)  
`imap_mailstat_mails_levelquotaused_quantity` (only imap with quota support)  
`imap_mailstat_mails_mailboxquotaavail_quantity` (only imap with quota support)  
`imap_mailstat_mails_mailboxquotaused_quantity` (only imap with quota support)  
`imap_mailstat_mails_messagequotaavail_quantity` (only imap with quota support)  
`imap_mailstat_mails_messagequotaused_quantity` (only imap with quota support)  
`imap_mailstat_mails_storagequotaavail_kilobytes` (only imap with quota support)  
`imap_mailstat_mails_storagequotaused_kilobytes` (only imap with quota support)  
`imap_mailstat_mails_oldestunseen_timestamp` (only with enabled feature flag `--oldestunseen.feature`)

> [!IMPORTANT]
> In version 0.1.0 (not yet released) the metric names were changed. First because they were hard to read and now I hope I follow more best practices in naming metrics. As 0.1.0 comes with more than one breaking change my decision was to rename the metrics at this point as well. The exporter allows you for migration to get the old metrics as well using commandline flag `--migration.mode` or the also available environment variable `MAILSTAT_EXPORTER_MIGRATIONMODE=true`. This flag/variable is deprecated and the old metrics will be removed in version 0.2.0.

The exposed metrics since 0.1.0 (not yet released) are the following:

metric | type | description | remarks
-------|------|-------------|---------
`mailstat_info` | gauge | Info metric for imap-mailstat-exporter | 
`mailstat_up` | gauge | Was talking to all accounts imap successfully | if value is 0: any account has a problem, check logs
`mailstat_fetch_duration_seconds` | gauge | Duration for fetching the metrics for the given account |
`mailstat_mails_all` | gauge | The total number of mails in folder |
`mailstat_mails_unseen` | gauge | The total number of unseen mails in folder |
`mailstat_level_quota_avail` | gauge | How many levels are available according your quota | only imap with quota support
`mailstat_level_quota_used` | gauge | How many levels are used | only imap with quota support
`mailstat_mailbox_quota_avail` | gauge | How many mailboxes are available according your quota | only imap with quota support
`mailstat_mailbox_quota_used` | gauge |  How many mailboxes are used | only imap with quota support
`mailstat_message_quota_avail` | gauge | How many messages available according your quota | only imap with quota support
`mailstat_message_quota_used` | gauge | How many messages are used | only imap with quota support
`mailstat_storage_quota_avail_bytes` | gauge | How many storage is available according your quota | only imap with quota support
`mailstat_storage_quota_used_bytes` | gauge | How many storage is used | only imap with quota support
`mailstat_mails_oldest_unseen_timestamp` | gauge | Timestamp in unix format of oldest unseen mail | only with enabled feature flag `--oldestunseen.feature`  

Example output:

```output
# HELP mailstat_fetch_duration_seconds Duration for fetching the metrics for the given account
# TYPE mailstat_fetch_duration_seconds gauge
mailstat_fetch_duration_seconds{mailboxname="Jane_Doe_Mailbox"} 1.303695723
mailstat_fetch_duration_seconds{mailboxname="Jane_Mailbox"} 0.612008505
# HELP mailstat_info Info metric for imap-mailstat-exporter.
# TYPE mailstat_info gauge
mailstat_info{version="0.1.0"} 1
# HELP mailstat_level_quota_avail How many levels are available according your quota
# TYPE mailstat_level_quota_avail gauge
mailstat_level_quota_avail{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 3000
# HELP mailstat_level_quota_used How many levels are used
# TYPE mailstat_level_quota_used gauge
mailstat_level_quota_used{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 1000
# HELP mailstat_mailbox_quota_avail How many mailboxes are available according your quota
# TYPE mailstat_mailbox_quota_avail gauge
mailstat_mailbox_quota_avail{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 31000
# HELP mailstat_mailbox_quota_used How many mailboxes are used
# TYPE mailstat_mailbox_quota_used gauge
mailstat_mailbox_quota_used{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 0
# HELP mailstat_mails_all The total number of mails in folder
# TYPE mailstat_mails_all gauge
mailstat_mails_all{mailboxfoldername="INBOX",mailboxname="Jane_Doe_Mailbox"} 404
mailstat_mails_all{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 2
mailstat_mails_all{mailboxfoldername="INBOX_Spam",mailboxname="Jane_Doe_Mailbox"} 5
mailstat_mails_all{mailboxfoldername="INBOX_Spam",mailboxname="Jane_Mailbox"} 0
mailstat_mails_all{mailboxfoldername="INBOX_Trash",mailboxname="Jane_Doe_Mailbox"} 32
mailstat_mails_all{mailboxfoldername="INBOX_Trash",mailboxname="Jane_Mailbox"} 0
# HELP mailstat_mails_oldest_unseen_timestamp Timestamp in unix format of oldest unseen mail
# TYPE mailstat_mails_oldest_unseen_timestamp gauge
mailstat_mails_oldest_unseen_timestamp{mailboxfoldername="INBOX",mailboxname="Jane_Doe_Mailbox"} 1.69878845e+09
mailstat_mails_oldest_unseen_timestamp{mailboxfoldername="INBOX_Spam",mailboxname="Jane_Doe_Mailbox"} 1.698318945e+09
# HELP mailstat_mails_unseen The total number of unseen mails in folder
# TYPE mailstat_mails_unseen gauge
mailstat_mails_unseen{mailboxfoldername="INBOX",mailboxname="Jane_Doe_Mailbox"} 1
mailstat_mails_unseen{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 0
mailstat_mails_unseen{mailboxfoldername="INBOX_Spam",mailboxname="Jane_Doe_Mailbox"} 5
mailstat_mails_unseen{mailboxfoldername="INBOX_Spam",mailboxname="Jane_Mailbox"} 0
mailstat_mails_unseen{mailboxfoldername="INBOX_Trash",mailboxname="Jane_Doe_Mailbox"} 0
mailstat_mails_unseen{mailboxfoldername="INBOX_Trash",mailboxname="Jane_Mailbox"} 0
# HELP mailstat_message_quota_avail How many messages available according your quota
# TYPE mailstat_message_quota_avail gauge
mailstat_message_quota_avail{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 6.2e+07
# HELP mailstat_message_quota_used How many messages are used
# TYPE mailstat_message_quota_used gauge
mailstat_message_quota_used{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 2000
# HELP mailstat_storage_quota_avail_bytes How many storage is available according your quota
# TYPE mailstat_storage_quota_avail_bytes gauge
mailstat_storage_quota_avail_bytes{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 1.048576e+09
# HELP mailstat_storage_quota_used_bytes How many storage is used
# TYPE mailstat_storage_quota_used_bytes gauge
mailstat_storage_quota_used_bytes{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 35000
# HELP mailstat_up Was talking to all accounts imap successfully.
# TYPE mailstat_up gauge
mailstat_up 1

```

Metrics are available via http (or https if configured) on port 8081/tcp on path `/metrics` as default, but you can configure this of you want to change.

## Commandline Options

### version 0.1.0 (not yet released, about to come)

You have several commandline options. Three of them can also be set via environment variables, if you like.  

```shell
usage: imap-mailstat-exporter [<flags>]

a prometheus-exporter to expose metrics about your mailboxes


Flags:
  -h, --[no-]help                Show context-sensitive help (also try --help-long and --help-man).
  -c, --config.file="./config/config.toml"  
                                 Provide the configfile ($MAILSTAT_EXPORTER_CONFIGFILE)
      --[no-]oldestunseen.feature  
                                 Enable metric with timestamp of oldest unseen mail, default false ($MAILSTAT_EXPORTER_OLDESTUNSEEN)
      --[no-]migration.mode      Enable old metric format, default false, WILL BE REMOVED IN version 0.2.0 ($MAILSTAT_EXPORTER_MIGRATIONMODE)
      --[no-]web.systemd-socket  Use systemd socket activation listeners instead of port listeners (Linux only).
      --web.listen-address=:8081 ...  
                                 Addresses on which to expose metrics and web interface. Repeatable for multiple addresses.
      --web.config.file=""       [EXPERIMENTAL] Path to configuration file that can enable TLS or authentication. See:
                                 https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md
      --web.telemetry-path="/metrics"  
                                 Path under which to expose the IMAP mailstat Prometheus metrics ($MAILSTAT_EXPORTER_WEB_TELEMETRY_PATH)
  -v, --[no-]version             Show application version.
      --log.level=info           Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --log.format=logfmt        Output format of log messages. One of: [logfmt, json]
```

### Version 0.0.1

You have three available commandline options.

```shell
Usage of imap-mailstat-exporter:
  -config string
        provide the configfile (default "./config/config.toml")
  -loglevel string
        provide the desired loglevel, INFO and ERROR are supported (default "INFO")
  -oldestunseendate
        enable metric with timestamp of oldest unseen mail
```

## Configuration

You can configure your accounts in a configfile in [toml](https://toml.io) format. You can find the example file in the folder `examples`. You can use
commandline flag `-config=<path/configfile>` (version 0.0.1) or `--config.file="<path/configfile>"` (version 0.1.0, not yet released) to specify where your configfile is located.

> [!IMPORTANT]
> If you are using the container image, the default configfile used where you need to mount your config is `/home/nonroot/config/config.toml`.

Example configuration, for one account, use only one account definition.

```config
# This is a example configfile.  
# You need only one account configured, but all keys need to be defined except username which can be empty and mailaddress is used as username value instead.
# place this file named as config.toml in a folder named config along your imap-mailstat-exporter binary or mount this file as config.toml in folder /config/ in the container.
# If you put you config elsewhere you can use the commandline flag --config.file="<path/configfile>" to specify where your config is.

[[Accounts]]
name = "Jane Mailbox" # mailbox, you can set as you like, will be used as metric label (whitespace are replaced by underscore)
mailaddress = "jane@example.com" # your email address (at the moment used as login name)
username = "your_user_name" # if empty string mailaddress value is used
password = "your_password" # beware of escaping characters like \ or "
serveraddress = "mail.example.com" # mailserver name or address
serverport = 993 # imap port number (at the moment only tls connection is supported)
starttls = false # set to true for using standard port 143 and STARTTLS to start a TLS connection
additionalfolders = ["Trash", "Spam"] # additional mailfolders you want to have metrics for

[[Accounts]] # you can configure more accounts if you like
name = "Jane Doe Mailbox"
mailaddress = "jane.doe@example.com"
username = ""
password = ""
serveraddress = "mail.example.com"
serverport = 143
starttls = true
additionalfolders = ["Trash", "Spam"]
```

## Loglevel

At the moment INFO (default), WARN and ERROR are used. DEBUG is available, but I don't output anything on this level yet. INFO tells you when metrics are fetched and give you additional information how long the connection setup, the login process and the whole metric fetch takes.
If INFO is too noisy you can switch to WARN or ERROR level and only get information about warnings or errors by using e.g. commandline flag `-loglevel WARN` (version 0.0.1), or `--log.level="WARN"` (version 0.1.0, not yet released).

## OCI Container Image

Image is available on: `ghcr.io/bt909/imap-mailstat-exporter`. Images are build for linux/amd64 and linux/arm64 for every release. Release versions are v*.*.* and the Container Images are without the `v` in front of the version, so use:

```shell
docker pull ghcr.io/bt909/imap-mailstat-exporter:*.*.*
```

The tag `latest` is following main branch and not related to the releases. This behavior will stay until release 1.0.0.

## License

This project is licensed using MIT license, see [LICENSE](https://github.com/bt909/imap-mailstat-exporter/blob/main/LICENSE)

## Trivia

This exporter is used personally with e-mail accounts provided by my webhosting provider [1984.is](https://1984.is/) (IMAP without quota) and provided by my ISP [Deutsche Telekom AG](https://www.telekom.de) (T-Online Freemail, IMAP with quota support).
