// Package imap-mailstat-exporter provides metrics for imap mailboxes
package main

import (
	"fmt"
	"imap-mailstat-exporter/internal/valuecollect"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// main function just for the main prometheus exporter functions
func main() {

	reg := prometheus.NewRegistry()
	d := valuecollect.NewImapStatsCollector()
	reg.MustRegister(d)

	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHandler)

	port := 8081
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatalf("cannot start exporter: %s", err)
	}
}
