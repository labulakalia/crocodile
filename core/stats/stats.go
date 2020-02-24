package stats

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	// "github.com/prometheus/client_golang/prometheus"

	"log"
	"net/http"
)

// 监控数据项
// goroutinue total
// 占用内存
// 任务总数
//
var (

// GuageVecApiDuration = prometheus.NewGaugeFunc(opts prometheus.GaugeOpts, function func() float64)
)

// Stats start listen port 9100,prometheus will pull data from this url
// http://ip:9100/metrics
func Stats() {

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe("0.0.0.0:9100", nil))
	// promhttp.Handler
}
