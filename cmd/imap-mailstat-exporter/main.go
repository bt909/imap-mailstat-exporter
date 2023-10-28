// Package imap-mailstat-exporter provides metrics for imap mailboxes
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/bt909/imap-mailstat-exporter/internal/valuecollect"
	"github.com/bt909/imap-mailstat-exporter/utils"
)

var (
	name                = "imap-mailstat-exporter"
	Version             = "0.1.0-alpha"
	configfile          *string
	loglevel            *string
	oldestunseenfeature *bool
)

// main function just for the main prometheus exporter functions
func main() {

	var app = kingpin.New(name, "a prometheus-exporter to expose metrics about your mailboxes")
	configfile = app.Flag("config.file", "provide the configfile").Envar("MAILSTAT_EXPORTER_CONFIGFILE").Default("./config/config.toml").Short('c').String()
	loglevel = app.Flag("log.level", "provide the desired loglevel, INFO and ERROR are supported").Envar("MAILSTAT_EXPORTER_LOGLEVEL").Default("INFO").String()
	oldestunseenfeature = app.Flag("oldestunseen.feature", "enable metric with timestamp of oldest unseen mail, default false").Envar("MAILSTAT_EXPORTER_OLDESTUNSEEN").Default("false").Bool()
	app.Version(Version)
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	utils.InitializeLogger(*loglevel)
	utils.Logger.Info("imap-mailstat-exporter started")

	reg := prometheus.NewRegistry()
	d := valuecollect.NewImapStatsCollector(*configfile, *loglevel, *oldestunseenfeature)
	reg.MustRegister(d)

	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, ` <!DOCTYPE html>
		<html>
		<body>
		Hello, this is imap-mailstat-exporter, for metrics check
		<a href="/metrics">metrics endpoint</a>
        for code check <a href="https://github.com/bt909/imap-mailstat-exporter" target="_blank"> the github repository</a>
		and for healthchecks you can use <a href="/healthz"> a healthz endpoint</>.
		</body>
		</html>`)
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, ` <!DOCTYPE html>
		<html>
		<body>
		imap-mailstat-exporter is healthy
		</body>
		</html>`)
	})

	port := 8081
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		utils.Logger.Fatal("cannot start exporter: %s", zap.Error(err))
	}
}
