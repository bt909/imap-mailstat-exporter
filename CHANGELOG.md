# Unreleased

* [BREAKING CHANGE]: change logging behavior by switching to [exporter-toolkit](https://github.com/prometheus/exporter-toolkit) and removing go.uber.org/zap as logging framework https://github.com/bt909/imap-mailstat-exporter/pull/33
* [BREAKING CHANGE]: change command line handling by switching command line parsing to [kingpin v2](https://github.com/alecthomas/kingpin) https://github.com/bt909/imap-mailstat-exporter/pull/32

* [FEATURE]: add basic auth and http/2 (TLS secured) connection by using exporter-toolkit https://github.com/bt909/imap-mailstat-exporter/pull/33
* [FEATURE]: add possibility to configure metrics path and listen address and port by using exporter-toolkit https://github.com/bt909/imap-mailstat-exporter/pull/33
* [CHORE]: bump to golang 1.21 and rename some internal things https://github.com/bt909/imap-mailstat-exporter/pull/31
* [CHORE]: update module github.com/prometheus/client_golang from v0.16.0 to v1.17.0 https://github.com/bt909/imap-mailstat-exporter/pull/30
* [CHORE]: update module go.uber.org/zap from v1.25.0 to v1.26.0 https://github.com/bt909/imap-mailstat-exporter/pull/29

# 0.0.1

* [FEAT]: first release