package main

import (
	"fmt"
	valuecollector "imap-mailstat-exporter/pkg/valuecollect"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	reg := prometheus.NewRegistry()
	d := valuecollector.NewImapStatsCollector()
	reg.MustRegister(d)

	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHandler)

	port := 8081
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatalf("cannot start exporter: %s", err)
	}
}
