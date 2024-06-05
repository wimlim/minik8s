package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage_percent",
		Help: "Current CPU usage percentage.",
	})
	memoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_usage_percent",
		Help: "Current memory usage percentage.",
	})
)

func init() {
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(memoryUsage)
}

func recordMetrics() {
	go func() {
		cpu := 0.0
		mem := 100.0
		for {
			cpuUsage.Set(cpu)
			memoryUsage.Set(mem)
			cpu += 1.0
			mem -= 1.0
			if cpu > 100 {
				cpu = 0
			}
			if mem < 0 {
				mem = 100
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

func main() {
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Beginning to serve on port :12345")
	log.Fatal(http.ListenAndServe(":12345", nil))
}
