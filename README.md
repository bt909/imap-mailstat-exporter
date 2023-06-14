# imap-mailstat-exporter

[![publish](https://github.com/bt909/imap-mailstat-exporter/actions/workflows/publish.yaml/badge.svg)](https://github.com/bt909/imap-mailstat-exporter/actions/workflows/publish.yaml)
 [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

This is a prometheus exporter which gives you metrics for how many emails you have in your INBOX and in additional configured folders.

> **Note**
> This exporter is in early development and at the moment highly adjusted for my personal usecase.

The exporter provides two metrics:

`imap_entire_mails` is a gauge metric for how many emails are in this folder labeled with the foldername of your mailbox and your configured mailboxname  
`imap_unseen_mails` is a gauge metric for how many emails are in this folder labeled with the foldername of your mailbox and your configured mailboxname 

Metrics are available via http on port 8081/tcp on path `/metrics`.

## Configuration

You can configure your accounts in a configfile in [toml](https://toml.io) format. You can find an example file in the folder `examples`. You can use
commandline flag `-config=<path/configfile`> to specify where your configfile is located.

This exporter is in early development and at the moment highly adjusted for my personal usecase.

## OCI Container Image

Image is available on: `ghcr.io/bt909/imap-mailstat-exporter`.
