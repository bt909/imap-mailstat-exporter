# 0.6.6 / 2025-07-06

* [CHORE] update module github.com/prometheus/common from 0.63.0 to v0.65.0 update module github.com/prometheus/common from 0.61.0 to v0.62.0
* [FIX] fix CVE-2025-22874, CVE-2025-0913 and CVE-2025-4673 by creating new release using latest golang version

# 0.6.5 / 2025-05-02

* [CHORE] update module github.com/burntsushi/toml from v1.4.0 to v1.5.0 https://github.com/bt909/imap-mailstat-exporter/pull/120
* [CHORE] update module github.com/prometheus/client_golang to from v1.21.1 to v1.22.0 https://github.com/bt909/imap-mailstat-exporter/pull/121
* [CHORE] update all dependecies to fix CVE-2025-22871 and CVE-2025-22872 https://github.com/bt909/imap-mailstat-exporter/pull/122

# 0.6.4 / 2025-03-17

*  [FIX] fix CVE-2025-22870 by updating dependencies https://github.com/bt909/imap-mailstat-exporter/pull/119

# 0.6.3 / 2025-02-09

* [CHORE] update module github.com/prometheus/exporter-toolkit from v0.13.2 to v0.14.0 https://github.com/bt909/imap-mailstat-exporter/pull/115
* [CHORE] update module github.com/prometheus/common from 0.61.0 to v0.62.0 https://github.com/bt909/imap-mailstat-exporter/pull/114
* [CHORE] add labels to dockerfile https://github.com/bt909/imap-mailstat-exporter/pull/113
* [CHORE] migrate renovate configuration https://github.com/bt909/imap-mailstat-exporter/pull/112

* [FIX] fix several CVEs by release with latest Go version (1.23.6),fixes: CVE-2025-22866, CVE-2024-45341, CVE-2024-45336
 
# 0.6.2 / 2024-12-22

* [CHORE] update several transitive dependencies https://github.com/bt909/imap-mailstat-exporter/pull/111

* [FIX] fix CVE-2024-45338 by updating golang.org/x/net to v0.33.0 https://github.com/bt909/imap-mailstat-exporter/pull/111

# 0.6.1 / 2024-12-13

* [CHORE] update module github.com/prometheus/exporter-toolkit from v0.13.1 to v0.13.2 https://github.com/bt909/imap-mailstat-exporter/pull/110
* [CHORE] update module github.com/prometheus/client_golang to from v1.20.4 to v1.20.5 https://github.com/bt909/imap-mailstat-exporter/pull/106
* [CHORE] update module github.com/prometheus/exporter-toolkit from v0.13.0 to v0.13.1 https://github.com/bt909/imap-mailstat-exporter/pull/108

* [FIX] fix CVE-2024-45337 by updating golang.org/x/crypto to v0.31.0 https://github.com/bt909/imap-mailstat-exporter/pull/110

# 0.6.0 / 2024-10-10

* [CHANGE] adapt to slog logging (changes log output) which is introduced by exporter-toolkit https://github.com/bt909/imap-mailstat-exporter/pull/101
* [CHANGE] adapt duration log entries to strings for json output like versions worked before switch to slog https://github.com/bt909/imap-mailstat-exporter/pull/104

* [CHORE] update module github.com/prometheus/exporter-toolkit from v0.12.0 to v0.13.0 https://github.com/bt909/imap-mailstat-exporter/pull/101
* [CHORE] update module github.com/prometheus/client_golang to from v1.20.3 to v1.20.4 https://github.com/bt909/imap-mailstat-exporter/pull/102
* [CHORE] update module github.com/prometheus/common from 0.59.1 to v0.60.0 https://github.com/bt909/imap-mailstat-exporter/pull/103

* [FIX] fix flags handling for log.level and log.format, which was broken with refactoring in release 0.5.0 https://github.com/bt909/imap-mailstat-exporter/pull/101

# 0.5.0 / 2024-09-09

* [FEATURE] enable pprof profiling from Go runtime https://github.com/bt909/imap-mailstat-exporter/pull/94

* [CHORE] update module github.com/prometheus/common to from 0.56.0 to v0.57.0 https://github.com/bt909/imap-mailstat-exporter/pull/91
* [CHORE] update module github.com/prometheus/common to from 0.57.0 to v0.58.0 https://github.com/bt909/imap-mailstat-exporter/pull/92
* [CHORE] update module github.com/prometheus/common to from 0.58.0 to v0.59.0 https://github.com/bt909/imap-mailstat-exporter/pull/96
* [CHORE] update module github.com/prometheus/common to from 0.59.0 to v0.59.1 https://github.com/bt909/imap-mailstat-exporter/pull/97
* [CHORE] update module github.com/prometheus/exporter-toolkit from v0.11.0 to v0.12.0 https://github.com/bt909/imap-mailstat-exporter/pull/93
* [CHORE] update module github.com/prometheus/client_golang to from v1.20.2 to v1.20.3 https://github.com/bt909/imap-mailstat-exporter/pull/95
* [CHORE] add go mod tidy to renovate https://github.com/bt909/imap-mailstat-exporter/pull/98

* [REFACTOR] centralize the promlogconfig assignment https://github.com/bt909/imap-mailstat-exporter/pull/99

* [FIX] fix CVE-2024-34156, CVE-2024-34155 and CVE-2024-34158 by rebuilding image with latest golang version by building a new release. As this release contains a new feature, version will be increased in minor version https://github.com/bt909/imap-mailstat-exporter/pull/100

# 0.4.1 / 2024-08-28

* [CHORE] update module github.com/prometheus/common to from 0.55.0 to v0.56.0 https://github.com/bt909/imap-mailstat-exporter/pull/89
* [CHORE] update module github.com/prometheus/client_golang to from v1.19.1 to v1.20.1 https://github.com/bt909/imap-mailstat-exporter/pull/87
* [CHORE] update module github.com/prometheus/client_golang to from v1.20.1 to v1.20.2 https://github.com/bt909/imap-mailstat-exporter/pull/88
* [CHORE] update golang from 1.22 to 1.23 https://github.com/bt909/imap-mailstat-exporter/pull/86

* [FIX] fix CVE-2024-24791 by updating golang from 1.22 to 1.23 https://github.com/bt909/imap-mailstat-exporter/pull/86

# 0.4.0 / 2024-06-26

* [FEATURE] allow web.telemetry-path to be set to / https://github.com/bt909/imap-mailstat-exporter/pull/82

* [CHORE] update module github.com/burntsushi/toml from v1.3.2 to v1.4.0 https://github.com/bt909/imap-mailstat-exporter/pull/76
* [CHORE] update module github.com/prometheus/common from v0.53.0 to v0.54.0 https://github.com/bt909/imap-mailstat-exporter/pull/78
* [CHORE] update golang toolchain from 1.22.3 to 1.22.4 https://github.com/bt909/imap-mailstat-exporter/pull/79
* [CHORE] update golang from 1.21 to 1.22 https://github.com/bt909/imap-mailstat-exporter/pull/82
* [CHORE] update module github.com/prometheus/common from v0.54.0 to v0.55.0 https://github.com/bt909/imap-mailstat-exporter/pull/84

* [CI] update goreleaser/goreleaser-action action from v5 to v6 https://github.com/bt909/imap-mailstat-exporter/pull/80
* [CI] update docker/build-push-action action from v5 to v6 https://github.com/bt909/imap-mailstat-exporter/pull/83

# 0.3.1 / 2024-05-12

* [CHORE] update module github.com/prometheus/common from v0.52.2 to v0.52.3 https://github.com/bt909/imap-mailstat-exporter/pull/70 
* [CHORE] update module github.com/prometheus/common from v0.52.3 to v0.53.0 https://github.com/bt909/imap-mailstat-exporter/pull/71
* [CHORE] update golang toolchain from 1.22.0 to 1.22.2 https://github.com/bt909/imap-mailstat-exporter/pull/72
* [CHORE] update golang toolchain from 1.22.2 to 1.22.3 https://github.com/bt909/imap-mailstat-exporter/pull/73
* [CHORE] update module github.com/prometheus/client_golang from v1.19.0 to v1.19.1 https://github.com/bt909/imap-mailstat-exporter/pull/74
* [CHORE] update several transitive dependencies https://github.com/bt909/imap-mailstat-exporter/pull/75

* [FIX] fix CVE-2023-45288 by updating golang.org/x/net from 0.22 to 0.25 https://github.com/bt909/imap-mailstat-exporter/pull/75

# 0.3.0 / 2024-04-04

* [DOCS] clarify the naming 'unseen' by setting link to IMAP flags description in RFC (thanks to [kekscode](https://github.com/kekscode)) https://github.com/bt909/imap-mailstat-exporter/pull/57

* [FEATURE] allow setting mailbox password via commandline or environment variable for one account setup https://github.com/bt909/imap-mailstat-exporter/pull/67

* [CHORE] update module github.com/prometheus/common from v0.47.0 to v0.48.0 https://github.com/bt909/imap-mailstat-exporter/pull/56
* [CHORE] update module github.com/prometheus/client_golang from v1.18.0 to v1.19.0 https://github.com/bt909/imap-mailstat-exporter/pull/60
* [CHORE] update module github.com/prometheus/common from v0.48.0 to v0.49.0 https://github.com/bt909/imap-mailstat-exporter/pull/61
* [CHORE] for module update of github.com/prometheus/common to v0.49.0 a switch from golang 1.20 to 1.21 was performed https://github.com/bt909/imap-mailstat-exporter/pull/61
* [CHORE] update module github.com/prometheus/common from v0.49.0 to v0.50.0 https://github.com/bt909/imap-mailstat-exporter/pull/62
* [CHORE] update devcontainer mcr.microsoft.com/devcontainers/go from 1-1.21-bullseye to 1-1.22-bullseye https://github.com/bt909/imap-mailstat-exporter/pull/63
* [CHORE] update module github.com/prometheus/common from v0.50.0 to v0.51.0 https://github.com/bt909/imap-mailstat-exporter/pull/65
* [CHORE] update module github.com/prometheus/common from v0.51.0 to v0.51.1 https://github.com/bt909/imap-mailstat-exporter/pull/66
* [CHORE] update module github.com/prometheus/common from v0.51.1 to v0.52.2 https://github.com/bt909/imap-mailstat-exporter/pull/68

Thank you for contribution [kekscode](https://github.com/kekscode).

# 0.2.0 / 2024-02-22

As warnings, log entries and the readme indicated the old metrics and the related migration mode is removed in this release!

* [BREAKING CHANGE]: remove old metrics and migration mode https://github.com/bt909/imap-mailstat-exporter/pull/44

* [FIX]: refactored configfile handling and fixed a bug regarding setting configured empty username to mailaddress https://github.com/bt909/imap-mailstat-exporter/pull/52

* [DOCS]: add some word to the example dashboard https://github.com/bt909/imap-mailstat-exporter/pull/44
* [DOCS]: clarify that only release images are build as Docker multi-platform images https://github.com/bt909/imap-mailstat-exporter/pull/55
* [DOCS]: added some info for 0.2.0 release https://github.com/bt909/imap-mailstat-exporter/pull/55

* [CHORE] update module github.com/alecthomas/kingpin/v2 from v2.3.2 to v2.4.0 https://github.com/bt909/imap-mailstat-exporter/pull/43
* [CHORE] update module github.com/prometheus/exporter-toolkit to v0.11.0 https://github.com/bt909/imap-mailstat-exporter/pull/47
* [CHORE] cleanup go mods and adjust changelog https://github.com/bt909/imap-mailstat-exporter/pull/49
* [CHORE] update module github.com/prometheus/client_golang from v1.17.0 to v.1.18.0 https://github.com/bt909/imap-mailstat-exporter/pull/50
* [CHORE] update module github.com/prometheus/common from v.0.45.0 to v.0.46.0 https://github.com/bt909/imap-mailstat-exporter/pull/51
* [CHORE] update build container version from golang:1.21 to golang:1.22 https://github.com/bt909/imap-mailstat-exporter/pull/53
* [CHORE] update module github.com/prometheus/common from v.0.46.0 to v.0.47.0 https://github.com/bt909/imap-mailstat-exporter/pull/54

* [CI] update actions/setup-go action to v5 https://github.com/bt909/imap-mailstat-exporter/pull/46
* [CI] adjust setup go to use latest stable go https://github.com/bt909/imap-mailstat-exporter/pull/48
* [CI] disable docker build provenance (fix unknown/unknown arch) https://github.com/bt909/imap-mailstat-exporter/pull/55

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
