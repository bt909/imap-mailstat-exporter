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

	port := 8081
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		utils.Logger.Fatal("cannot start exporter: %s", zap.Error(err))
	}
}
