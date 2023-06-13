# imap-mailstat-exporter

[![publish](https://github.com/bt909/imap-mailstat-exporter/actions/workflows/publish.yaml/badge.svg)](https://github.com/bt909/imap-mailstat-exporter/actions/workflows/publish.yaml)

This is a prometheus exporter which gives you metrics for how many emails you have in your INBOX and in additional configured folders.
You will get 2 metrics exposed.

`imap_entire_mails` is a gauge metric for how many emails are in this folder labeled with the foldername of your mailbox and your configured mailboxname  
`imap_unseen_mails` is a gauge metric for how many emails are in this folder labeled with the foldername of your mailbox and your configured mailboxname 

Metrics are available on port 8081/tcp on path `/metrics`.

This exporter is in early development and at the moment highly adjusted for my personal usecase.

Docker image is available on: 
` ghcr.io/bt909/imap-mailstat-exporter`
