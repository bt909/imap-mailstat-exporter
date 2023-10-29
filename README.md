# imap-mailstat-exporter

[![publish](https://github.com/bt909/imap-mailstat-exporter/actions/workflows/publish_latest_oci_image.yaml/badge.svg)](https://github.com/bt909/imap-mailstat-exporter/actions/workflows/publish_latest_oci_image.yaml)
 [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

This is a prometheus exporter which gives you metrics for how many emails you have in your INBOX and in additional configured folders.  

Connections to IMAP are only TLS enrypted supported, either via TLS or STARTTLS.

> [!NOTE]
> This exporter is in early development and at the moment highly adjusted for my personal usecase. As it not reached 1.0.0 yet, there may are breaking changes at any time. Keep an eye on the [CHANGELOG](https://github.com/bt909/imap-mailstat-exporter/blob/main/CHANGELOG.md) for information.

The exporter provides nine metrics, two main metrics are provided for all accounts, one metric can be enabled using a feature flag `-oldestunseendate` and six metrics are quota related and only provided if the server supports imap quota.

If your account supports quota you can see in loglevel INFO (default) with the following log entry:

```output
{"level":"info","ts":"2023-08-17T14:45:21.399Z","caller":"valuecollect/valuecollector.go:257","msg":"Fetching quota related metrics","address":"jane.doe@example.com"}

```

The exposed metrics are the following:

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

Example output:

```output
# HELP imap_mailstat_mails_all_quantity The total number of mails in folder
# TYPE imap_mailstat_mails_all_quantity gauge
imap_mailstat_mails_all_quantity{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 537
imap_mailstat_mails_all_quantity{mailboxfoldername="INBOX",mailboxname="Jane_Doe_Mailbox"} 1308
imap_mailstat_mails_all_quantity{mailboxfoldername="INBOX_Spam",mailboxname="Jane_Mailbox"} 0
imap_mailstat_mails_all_quantity{mailboxfoldername="INBOX_Spam",mailboxname="Jane_Doe_Mailbox"} 1
imap_mailstat_mails_all_quantity{mailboxfoldername="INBOX_Trash",mailboxname="Jane_Mailbox"} 10
imap_mailstat_mails_all_quantity{mailboxfoldername="INBOX_Trash",mailboxname="Jane_Doe_Mailbox"} 9
# HELP imap_mailstat_mails_levelquotaavail_quantity How many levels are available according your quota
# TYPE imap_mailstat_mails_levelquotaavail_quantity gauge
imap_mailstat_mails_levelquotaavail_quantity{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 3
imap_mailstat_mails_levelquotaavail_quantity{mailboxfoldername="INBOX",mailboxname="Jane_Doe_Mailbox"} 3
# HELP imap_mailstat_mails_levelquotaused_quantity How many levels are used
# TYPE imap_mailstat_mails_levelquotaused_quantity gauge
imap_mailstat_mails_levelquotaused_quantity{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 1
imap_mailstat_mails_levelquotaused_quantity{mailboxfoldername="INBOX",mailboxname="Jane_Doe_Mailbox"} 1
# HELP imap_mailstat_mails_mailboxquotaavail_quantity How many mailboxes are available according your quota
# TYPE imap_mailstat_mails_mailboxquotaavail_quantity gauge
imap_mailstat_mails_mailboxquotaavail_quantity{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 31
imap_mailstat_mails_mailboxquotaavail_quantity{mailboxfoldername="INBOX",mailboxname="Jane_Doe_Mailbox"} 31
# HELP imap_mailstat_mails_mailboxquotaused_quantity How many mailboxes are used
# TYPE imap_mailstat_mails_mailboxquotaused_quantity gauge
imap_mailstat_mails_mailboxquotaused_quantity{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 0
imap_mailstat_mails_mailboxquotaused_quantity{mailboxfoldername="INBOX",mailboxname="Jane_Doe_Mailbox"} 0
# HELP imap_mailstat_mails_messagequotaavail_quantity How many messages available according your quota
# TYPE imap_mailstat_mails_messagequotaavail_quantity gauge
imap_mailstat_mails_messagequotaavail_quantity{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 62000
imap_mailstat_mails_messagequotaavail_quantity{mailboxfoldername="INBOX",mailboxname="Jane_Doe_Mailbox"} 62000
# HELP imap_mailstat_mails_messagequotaused_quantity How many messages are used
# TYPE imap_mailstat_mails_messagequotaused_quantity gauge
imap_mailstat_mails_messagequotaused_quantity{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 2
imap_mailstat_mails_messagequotaused_quantity{mailboxfoldername="INBOX",mailboxname="Jane_Doe_Mailbox"} 4
# HELP imap_mailstat_mails_oldestunseen_timestamp Timestamp in unix format of oldest unseen mail
# TYPE imap_mailstat_mails_oldestunseen_timestamp gauge
imap_mailstat_mails_oldestunseen_timestamp{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 1.693660714e+09
imap_mailstat_mails_oldestunseen_timestamp{mailboxfoldername="INBOX",mailboxname="Jane_Doe_Mailbox"} 1.694538222e+09
imap_mailstat_mails_oldestunseen_timestamp{mailboxfoldername="INBOX_Spam",mailboxname="Jane_Mailbox"} 1.69398128e+09
# HELP imap_mailstat_mails_storagequotaavail_kilobytes How many storage is available according your quota
# TYPE imap_mailstat_mails_storagequotaavail_kilobytes gauge
imap_mailstat_mails_storagequotaavail_kilobytes{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 1.048576e+06
imap_mailstat_mails_storagequotaavail_kilobytes{mailboxfoldername="INBOX",mailboxname="Jane_Doe_Mailbox"} 1.048576e+06
# HELP imap_mailstat_mails_storagequotaused_kilobytes How many storage is used
# TYPE imap_mailstat_mails_storagequotaused_kilobytes gauge
imap_mailstat_mails_storagequotaused_kilobytes{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 35
imap_mailstat_mails_storagequotaused_kilobytes{mailboxfoldername="INBOX",mailboxname="Jane_Doe_Mailbox"} 71
# HELP imap_mailstat_mails_unseen_quantity The total number of unseen mails in folder
# TYPE imap_mailstat_mails_unseen_quantity gauge
imap_mailstat_mails_unseen_quantity{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 2
imap_mailstat_mails_unseen_quantity{mailboxfoldername="INBOX",mailboxname="Jane_Doe_Mailbox"} 0
imap_mailstat_mails_unseen_quantity{mailboxfoldername="INBOX_Spam",mailboxname="Jane_Mailbox"} 0
imap_mailstat_mails_unseen_quantity{mailboxfoldername="INBOX_Spam",mailboxname="Jane_Doe_Mailbox"} 0
imap_mailstat_mails_unseen_quantity{mailboxfoldername="INBOX_Trash",mailboxname="Jane_Mailbox"} 0
imap_mailstat_mails_unseen_quantity{mailboxfoldername="INBOX_Trash",mailboxname="Jane_Doe_Mailbox"} 0
```

Metrics are available via http on port 8081/tcp on path `/metrics`.

## Commandline Options

### version 0.1.0 (not yet released, about to come)

You have three important commandline options. The information of the commandline flags can also be provided as environment variables.  

```shell
usage: imap-mailstat-exporter [<flags>]

a prometheus-exporter to expose metrics about your mailboxes


Flags:
  -h, --[no-]help         Show context-sensitive help (also try --help-long and --help-man).
  -c, --config.file="./config/config.toml"  
                          provide the configfile ($MAILSTAT_EXPORTER_CONFIGFILE)
      --log.level="INFO"  provide the desired loglevel, INFO and ERROR are supported ($MAILSTAT_EXPORTER_LOGLEVEL)
      --[no-]oldestunseen.feature  
                          enable metric with timestamp of oldest unseen mail, default false ($MAILSTAT_EXPORTER_OLDESTUNSEEN)
  -v, --[no-]version      Show application version.
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

At the moment INFO (default) and ERROR are available. INFO tells you when metrics are fetched and give you additional information how long the connection setup, the login process and the whole metric fetch takes.
If INFO is too noisy you can switch to ERROR level and only get information about errors by using commandline flag `-loglevel ERROR` (version 0.0.1), or `--log.level="ERROR"` (version 0.1.0, not yet released).

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
