package prometheusutil

import (
	"minik8s/pkg/apiobj"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cpuPercent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "node_cpu_percent",
		Help: "CPU usage percent of the node.",
	}, []string{"condition"})

	memPercent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "node_memory_percent",
		Help: "Memory usage percent of the node.",
	}, []string{"condition"})

	podNum = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "node_pod_num",
		Help: "Number of pods running on the node.",
	}, []string{"condition"})
)

func ExposeNodeStatusToPrometheus(nodeStatus *apiobj.NodeStatus) {
	labels := prometheus.Labels{"condition": string(nodeStatus.Condition)}
	cpuPercent.With(labels).Set(nodeStatus.CpuPercent)
	memPercent.With(labels).Set(nodeStatus.MemPercent)
	podNum.With(labels).Set(float64(nodeStatus.PodNum))
}

func StartPrometheusMetricsServer(port string) {
	addr := ":" + port
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			panic(err)
		}
	}()
}
