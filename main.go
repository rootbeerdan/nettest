package main

import (
	"time"

	"net/http"

	"github.com/go-ping/ping"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	fqdn = []string{"test-v4.maymobility.com", "test-v6.maymobility.com", "zero-trust-client.cloudflareclient.com", "ipv6.google.com", "ipv4.google.com"}
)

func main() {
	for {
		for _, fqdn := range fqdn {
			pingSuccess := prometheus.NewGauge(
				prometheus.GaugeOpts{
					Name:        "ping_success",
					Help:        "Indicates if the ping was successful or not",
					ConstLabels: prometheus.Labels{"fqdn": fqdn},
				},
			)
			pingDuration := prometheus.NewGauge(
				prometheus.GaugeOpts{
					Name:        "ping_duration_seconds",
					Help:        "Duration of the ping in seconds",
					ConstLabels: prometheus.Labels{"fqdn": fqdn},
				},
			)

			prometheus.MustRegister(pingSuccess)
			prometheus.MustRegister(pingDuration)

			pinger, err := ping.NewPinger(fqdn)
			if err != nil {
				pingSuccess.Set(0)
				pingDuration.Set(0)
				continue
			}
			pinger.OnRecv = func(pkt *ping.Packet) {
				pingSuccess.Set(1)
				pingDuration.Set(pkt.Rtt.Seconds())
			}
			pinger.OnFinish = func(stats *ping.Statistics) {
				pingSuccess.Set(0)
				pingDuration.Set(0)
			}
			pinger.Count = 1
			pinger.Timeout = time.Second * 5
			pinger.Run()
		}
		time.Sleep(time.Second * 30)
	}

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}
}
