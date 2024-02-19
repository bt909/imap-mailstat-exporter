// Package imap-mailstat-exporter provides metrics for imap mailboxes
package main

import (
	"net/http"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"

	"github.com/bt909/imap-mailstat-exporter/internal/configread"
	"github.com/bt909/imap-mailstat-exporter/internal/valuecollect"
)

var (
	name                = "imap-mailstat-exporter"
	Version             = "0.2.0-beta"
	configfile          *string
	oldestunseenfeature *bool
	logger              = promlog.New(&promlog.Config{})
)

// main function just for the main prometheus exporter functions
func main() {

	var app = kingpin.New(name, "a prometheus-exporter to expose metrics about your mailboxes")
	configfile = app.Flag("config.file", "Provide the configfile").Envar("MAILSTAT_EXPORTER_CONFIGFILE").Default("./config/config.toml").Short('c').String()
	oldestunseenfeature = app.Flag("oldestunseen.feature", "Enable metric with timestamp of oldest unseen mail, default false").Envar("MAILSTAT_EXPORTER_OLDESTUNSEEN").Default("false").Bool()
	toolkitFlags := kingpinflag.AddFlags(app, ":8081")
	metricsPath := app.Flag("web.telemetry-path", "Path under which to expose the IMAP mailstat Prometheus metrics").Envar("MAILSTAT_EXPORTER_WEB_TELEMETRY_PATH").Default("/metrics").String()
	app.Version(Version)
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')

	promlogConfig := &promlog.Config{}
	flag.AddFlags(app, promlogConfig)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	logger = promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting imap-mailstat-exporter", "Version", Version)
	reg := prometheus.NewRegistry()
	level.Info(logger).Log("msg", "Reading config file", "File", *configfile)
	config, err := configread.GetConfig(*configfile)
	if err != nil {
		level.Error(logger).Log("msg", "Error in reading config", "Error", err)
		os.Exit(1)
	}
	d := valuecollect.NewImapStatsCollector(config, logger, *oldestunseenfeature, Version)
	reg.MustRegister(d)

	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle(*metricsPath, promHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(` <!DOCTYPE html>
		<html>
		<head><title>imap-mailstat-exporter</title></head>
		<body>
		<h1>imap-mailstat-exporter</h1>
		<p><a href='` + *metricsPath + `'>Metrics</a></p>
		<p><a href="/healthz">Health</a></p>
		<p><a href="https://github.com/bt909/imap-mailstat-exporter" target="_blank">Code</a></p>
		</body>
		</html>`))
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(` <!DOCTYPE html>
		<html>
		<body>
		imap-mailstat-exporter is healthy
		</body>
		</html>`))
	})

	server := &http.Server{Handler: mux}
	if err := web.ListenAndServe(server, toolkitFlags, logger); err != nil {
		level.Error(logger).Log("msg", "Cannot start exporter")
	}
}
