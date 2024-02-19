# Unreleased

As warnings, log entries and the readme indicated the old metrics and the related migration mode will removed in the next release 0.2.0.

* [BREAKING CHANGE]: remove old metrics and migration mode https://github.com/bt909/imap-mailstat-exporter/pull/44

* [FIX]: refactored configfile handling and fixed a bug regarding setting configured empty username to mailaddress https://github.com/bt909/imap-mailstat-exporter/pull/52

* [DOCS]: add some word to the example dashboard https://github.com/bt909/imap-mailstat-exporter/pull/44

* [CHORE] update module github.com/alecthomas/kingpin/v2 from v2.3.2 to v2.4.0 https://github.com/bt909/imap-mailstat-exporter/pull/43
* [CHORE] update module github.com/prometheus/exporter-toolkit to v0.11.0 https://github.com/bt909/imap-mailstat-exporter/pull/47
* [CHORE] cleanup go mods and adjust changelog https://github.com/bt909/imap-mailstat-exporter/pull/49
* [CHORE] update module github.com/prometheus/client_golang from v1.17.0 to v.1.18.0 https://github.com/bt909/imap-mailstat-exporter/pull/50
* [CHORE] update module github.com/prometheus/common from v.0.45.0 to v.0.46.0 https://github.com/bt909/imap-mailstat-exporter/pull/51
* [CHORE] update build container version from golang:1.21 to golang:1.22 https://github.com/bt909/imap-mailstat-exporter/pull/53

* [CI] update actions/setup-go action to v5 https://github.com/bt909/imap-mailstat-exporter/pull/46
* [CI] adjust setup go to use latest stable go https://github.com/bt909/imap-mailstat-exporter/pull/48

# 0.1.0 / 2023-11-07

This version has some heavy changes, but these were things a thought about for a while now and think this shouldn't be deferred for a long time. So now we have a big bang but the next ideas I have in mind are not that destructive.

* [BREAKING CHANGE]: add up and info metric (maybe not breaking change, but commited as) https://github.com/bt909/imap-mailstat-exporter/pull/37
* [BREAKING CHANGE]: rename metrics for readability and follow some best practices https://github.com/bt909/imap-mailstat-exporter/pull/36
* [BREAKING CHANGE]: change logging behavior by switching to [exporter-toolkit](https://github.com/prometheus/exporter-toolkit) and removing go.uber.org/zap as logging framework https://github.com/bt909/imap-mailstat-exporter/pull/33
* [BREAKING CHANGE]: change command line handling by switching command line parsing to [kingpin v2](https://github.com/alecthomas/kingpin) https://github.com/bt909/imap-mailstat-exporter/pull/32

* [CHANGE]: add deprecation warning in logs for migration mode https://github.com/bt909/imap-mailstat-exporter/pull/39
* [CHANGE]: switch from CMD to ENTRYPOINT in Dockerfile to allow commandline arguments easily to be passed to the container https://github.com/bt909/imap-mailstat-exporter/pull/38

* [DOCS]: add example dashboard https://github.com/bt909/imap-mailstat-exporter/pull/41

* [FEATURE]: add up and info metric https://github.com/bt909/imap-mailstat-exporter/pull/37
* [FEATURE]: add a new metric named mailstat_fetch_duration_seconds https://github.com/bt909/imap-mailstat-exporter/pull/36
* [FEATURE]: add basic auth and http/2 (TLS secured) connection by using exporter-toolkit https://github.com/bt909/imap-mailstat-exporter/pull/33
* [FEATURE]: add possibility to configure metrics path and listen address and port by using exporter-toolkit https://github.com/bt909/imap-mailstat-exporter/pull/33

* [FIX]: bump to golang 1.20 because of goreleaser problems ins pipeline https://github.com/bt909/imap-mailstat-exporter/pull/42

* [CHORE]: update module github.com/prometheus/common to from v0.44.0 to v0.45.0 https://github.com/bt909/imap-mailstat-exporter/pull/34
* [CHORE]: bump to golang 1.21 and rename some internal things https://github.com/bt909/imap-mailstat-exporter/pull/31
* [CHORE]: update module github.com/prometheus/client_golang from v0.16.0 to v1.17.0 https://github.com/bt909/imap-mailstat-exporter/pull/30
* [CHORE]: update module go.uber.org/zap from v1.25.0 to v1.26.0 https://github.com/bt909/imap-mailstat-exporter/pull/29

# 0.0.1 / 2023-09-14

* [FEAT]: first release
