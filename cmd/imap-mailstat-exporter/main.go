// Package imap-mailstat-exporter provides metrics for imap mailboxes
package main

import (
	"flag"
	"fmt"
	"imap-mailstat-exporter/internal/valuecollect"
	"imap-mailstat-exporter/utils"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// main function just for the main prometheus exporter functions
func main() {

	flag.StringVar(&valuecollect.Configfile, "config", "./config/config.toml", "provide the configfile")
	flag.StringVar(&valuecollect.Loglevel, "loglevel", "INFO", "provide the desired loglevel, INFO and ERROR are supported")
	flag.Parse()

	utils.InitializeLogger(valuecollect.Loglevel)
	utils.Logger.Info("imap-mailstat-exporter started")

	reg := prometheus.NewRegistry()
	d := valuecollect.NewImapStatsCollector()
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
