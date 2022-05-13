package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"

	"github.com/cqroot/s3-active-exporter/collector"
	"github.com/cqroot/s3-active-exporter/internal"
	"github.com/cqroot/s3-active-exporter/logger"
)

func main() {
	defer logger.Sync()
	internal.InitConfig()

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector.NewCollector())

	http.Handle(viper.GetString("web.telemetry-path"), promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
<html>
<head><title>Swift Exporter v` + "0.0.1" + `</title></head>
<body>
<h1>Swift Exporter ` + internal.BuildVersion + `</h1>
<p><a href='` + viper.GetString("web.telemetry-path") + `'>Metrics</a></p>
</body>
</html>
        `))
	})

	logger.Info(fmt.Sprintf("Providing metrics at %s%s", viper.GetString("web.listen-address"), viper.GetString("web.telemetry-path")))
	logger.Fatal(http.ListenAndServe(viper.GetString("web.listen-address"), nil).Error())
}
