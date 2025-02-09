// Package imap-mailstat-exporter provides metrics for imap mailboxes
package main

import (
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"

	"github.com/bt909/imap-mailstat-exporter/internal/configread"
	"github.com/bt909/imap-mailstat-exporter/internal/valuecollect"
)

var (
	name                = "imap-mailstat-exporter"
	Version             = "0.6.3"
	configfile          *string
	oldestunseenfeature *bool
	mailboxpassword     *string
	promslogConfig      = &promslog.Config{}
)

// main function just for the main prometheus exporter functions
func main() {

	var app = kingpin.New(name, "a prometheus-exporter to expose metrics about your mailboxes")
	configfile = app.Flag("config.file", "Provide the configfile").Envar("MAILSTAT_EXPORTER_CONFIGFILE").Default("./config/config.toml").Short('c').String()
	oldestunseenfeature = app.Flag("oldestunseen.feature", "Enable metric with timestamp of oldest unseen mail, default false").Envar("MAILSTAT_EXPORTER_OLDESTUNSEEN").Default("false").Bool()
	mailboxpassword = app.Flag("mailboxpassword", "Password for mailbox, available if only one mailbox is configured").Envar("MAILSTAT_EXPORTER_MAILBOX_PASSWORD").Default("\x00").Short('p').String()
	toolkitFlags := kingpinflag.AddFlags(app, ":8081")
	metricsPath := app.Flag("web.telemetry-path", "Path under which to expose the IMAP mailstat Prometheus metrics").Envar("MAILSTAT_EXPORTER_WEB_TELEMETRY_PATH").Default("/metrics").String()
	app.Version(Version)
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')

	flag.AddFlags(app, promslogConfig)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	// initialize logger
	logger := promslog.New(promslogConfig)

	logger.Info("Starting imap-mailstat-exporter", "Version", Version)
	reg := prometheus.NewRegistry()
	logger.Info("Reading config file", "File", *configfile)
	config, err := configread.GetConfig(*configfile)
	if err != nil {
		logger.Error("Error in reading config", "Error", err)
		os.Exit(1)
	}
	switch {
	case len(config.Accounts) == 1 && *mailboxpassword != "\x00":
		config.Accounts[0].Password = *mailboxpassword
	case len(config.Accounts) > 1 && *mailboxpassword != "\x00":
		logger.Error("Configfile has set more than one mailbox and password via is set commandline or environment variable, but will be ignored!")
	case len(config.Accounts) == 1 && *mailboxpassword == "\x00" && config.Accounts[0].Password == "":
		logger.Warn("Configfile has empty password and you don't have set a password via commandline or environment variable")
	}
	d := valuecollect.NewImapStatsCollector(config, logger, *oldestunseenfeature, Version)
	reg.MustRegister(d)

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle(*metricsPath, promHandler)
	if *metricsPath != "/" {
		landingConfig := web.LandingConfig{
			Name:        "IMAP Mailstat Exporter",
			HeaderColor: "purple",
			Description: "Prometheus Exporter for IMAP Mailboxes",
			Version:     Version,
			Links: []web.LandingLinks{
				{
					Address: *metricsPath,
					Text:    "Metrics",
				},
				{
					Address: "/healthz",
					Text:    "Health Endpoint",
				},
				{
					Address: "https://github.com/bt909/imap-mailstat-exporter",
					Text:    "Code",
				},
			},
		}
		landingPage, err := web.NewLandingPage(landingConfig)
		if err != nil {
			logger.Error("Error on web page handling", "Error", err.Error())
			os.Exit(1)
		}
		mux.Handle("/", landingPage)
	}
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
		logger.Error("Cannot start exporter", "Error", err)
	}
}
