package main

import (
	"fmt"
	"net"
	"time"

	"net/http"

	"github.com/go-ping/ping"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	fqdn = []string{"ipv4.google.com", "ipv6.google.com"}
)

func main() {
	pingSuccess := make(map[string]prometheus.Gauge)
	pingDuration := make(map[string]prometheus.Gauge)

	for _, fqdn := range fqdn {
		pingSuccess[fqdn] = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name:        "ping_success",
				Help:        "Indicates if the ping was successful or not",
				ConstLabels: prometheus.Labels{"fqdn": fqdn},
			},
		)
		pingDuration[fqdn] = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name:        "ping_duration_seconds",
				Help:        "Duration of the ping in seconds",
				ConstLabels: prometheus.Labels{"fqdn": fqdn},
			},
		)

		prometheus.MustRegister(pingSuccess[fqdn])
		prometheus.MustRegister(pingDuration[fqdn])
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			panic(err)
		}
	}()

	for {
		for _, fqdn := range fqdn {
			ipv4Addr, err := net.ResolveIPAddr("ip4", fqdn)
			if err == nil {
				pinger, err := ping.NewPinger(ipv4Addr.String())
				if err != nil {
					pingSuccess[fqdn].Set(0)
					pingDuration[fqdn].Set(0)
					fmt.Printf("Error creating IPv4 pinger for %s: %v\n", fqdn, err)
					continue
				}

				pinger.SetPrivileged(true)
				pinger.OnRecv = func(pkt *ping.Packet) {
					pingSuccess[fqdn].Set(1)
					pingDuration[fqdn].Set(pkt.Rtt.Seconds())
				}
				pinger.OnFinish = func(stats *ping.Statistics) {
					if stats.PacketsRecv == 0 {
						pingSuccess[fqdn].Set(0)
						pingDuration[fqdn].Set(0)
						fmt.Printf("Error pinging IPv4 %s: Request timed out\n", fqdn)
					}
				}
				pinger.Count = 1
				pinger.Timeout = time.Second * 5
				go pinger.Run()
			}

			ipv6Addr, err := net.ResolveIPAddr("ip6", fqdn)
			if err == nil {
				pinger, err := ping.NewPinger(ipv6Addr.String())
				if err != nil {
					pingSuccess[fqdn].Set(0)
					pingDuration[fqdn].Set(0)
					fmt.Printf("Error creating IPv6 pinger for %s: %v\n", fqdn, err)
					continue
				}

				pinger.SetPrivileged(true)
				pinger.OnRecv = func(pkt *ping.Packet) {
					pingSuccess[fqdn].Set(1)
					pingDuration[fqdn].Set(pkt.Rtt.Seconds())
				}
				pinger.OnFinish = func(stats *ping.Statistics) {
					if stats.PacketsRecv == 0 {
						pingSuccess[fqdn].Set(0)
						pingDuration[fqdn].Set(0)
						fmt.Printf("Error pinging IPv6 %s: Request timed out\n", fqdn)
					}
				}
				pinger.Count = 1
				pinger.Timeout = time.Second * 5
				go pinger.Run()
			}

		}
		time.Sleep(time.Second * 30)
	}

}
