# imap-mailstat-exporter

[![publish](https://github.com/bt909/imap-mailstat-exporter/actions/workflows/publish.yaml/badge.svg)](https://github.com/bt909/imap-mailstat-exporter/actions/workflows/publish.yaml)
 [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

This is a prometheus exporter which gives you metrics for how many emails you have in your INBOX and in additional configured folders.

> **Note**  
> This exporter is in early development and at the moment highly adjusted for my personal usecase.

The exporter provides two metrics:

```output
# HELP imap_entire_mails The total number of mails in folder
# TYPE imap_entire_mails gauge
imap_entire_mails{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 537
imap_entire_mails{mailboxfoldername="INBOX",mailboxname="Jane_Doe_Mailbox"} 1308
imap_entire_mails{mailboxfoldername="INBOX_Spam",mailboxname="Jane_Mailbox"} 0
imap_entire_mails{mailboxfoldername="INBOX_Spam",mailboxname="Jane_Doe_Mailbox"} 1
imap_entire_mails{mailboxfoldername="INBOX_Trash",mailboxname="Jane_Mailbox"} 10
imap_entire_mails{mailboxfoldername="INBOX_Trash",mailboxname="Jane_Doe_Mailbox"} 9
# HELP imap_unseen_mails The total number of unseen mails in folder
# TYPE imap_unseen_mails gauge
imap_unseen_mails{mailboxfoldername="INBOX",mailboxname="Jane_Mailbox"} 2
imap_unseen_mails{mailboxfoldername="INBOX",mailboxname="Jane_Doe_Mailbox"} 0
imap_unseen_mails{mailboxfoldername="INBOX_Spam",mailboxname="Jane_Mailbox"} 0
imap_unseen_mails{mailboxfoldername="INBOX_Spam",mailboxname="Jane_Doe_Mailbox"} 0
imap_unseen_mails{mailboxfoldername="INBOX_Trash",mailboxname="Jane_Mailbox"} 0
imap_unseen_mails{mailboxfoldername="INBOX_Trash",mailboxname="Jane_Doe_Mailbox"} 0
```

Metrics are available via http on port 8081/tcp on path `/metrics`.

## Configuration

You can configure your accounts in a configfile in [toml](https://toml.io) format. You can find the example file in the folder `examples`. You can use
commandline flag `-config=<path/configfile`> to specify where your configfile is located.

Example configuration, for one account, use only one account definition.

```config
# This is a example configfile. You need only one account configured, but all keys need to be defined (or empty if not yet implemented).
# place this file named as config.toml in a folder named config along your imap-mailstat-exporter binary or mount this file as config.toml in folder /config/ in the container.
# If you put you config elsewhere you can use the commandline flag -config=<path/configfile> to specify where your config is.

[[Accounts]]
name = "Jane Mailbox" # mailbox, you can set as you like, will be used as metric label (whitespace are replaced by underscore)
mailaddress = "jane@example.com" # your email address (at the moment used as login name)
username = "your_user_name" # if empty string mailaddress value is used
password = "your_password" # beware of escaping characters like \ or "
serveraddress = "mail.example.com" # mailserver name or address
serverport = 993 # imap port number (at the moment only tls connection is supported)
starttls = false # not yet implemented, will be available if you use STARTTLS
additionalfolders = ["Trash", "Spam"] # additional mailfolders you want to have metrics for

[[Accounts]] # you can configure more accounts if you like
name = "Jane Doe Mailbox"
mailaddress = "jane.doe@example.com"
username = ""
password = ""
serveraddress = "mail.example.com"
serverport = 993
starttls = false
additionalfolders = ["Trash", "Spam"]
```

## OCI Container Image

Image is available on: `ghcr.io/bt909/imap-mailstat-exporter`. At the moment there are no releases, just latest or you can use the digest.
Releases will be available soon.
