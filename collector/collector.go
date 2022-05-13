package collector

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"

	"github.com/cqroot/s3-active-exporter/internal"
)

var namespace = "s3active"

type Collector struct {
	metrics map[string]*prometheus.Desc
	mutex   sync.Mutex
}

func NewCollector() *Collector {
	return &Collector{
		metrics: map[string]*prometheus.Desc{
			"put_metric": prometheus.NewDesc(namespace+"_"+"put", "s3 active put metric", []string{"endpoint", "filename"}, nil),
			"get_metric": prometheus.NewDesc(namespace+"_"+"get", "s3 active get metric", []string{"endpoint", "filename"}, nil),
			"del_metric": prometheus.NewDesc(namespace+"_"+"del", "s3 active del metric", []string{"endpoint", "filename"}, nil),
		},
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	endpoints := viper.GetStringSlice("s3.endpoints")
	wg := sync.WaitGroup{}
	wg.Add(len(endpoints))

	for _, endpoint := range endpoints {
		go func(wg *sync.WaitGroup, ch chan<- prometheus.Metric, endpoint string) {
			am := internal.S3ActiveMonitor{}
			putResult, getResult, delResult, putFilename, getFilename, delFilename := am.Run(endpoint)
			ch <- prometheus.MustNewConstMetric(c.metrics["put_metric"], prometheus.GaugeValue, float64(putResult), endpoint, putFilename)
			ch <- prometheus.MustNewConstMetric(c.metrics["get_metric"], prometheus.GaugeValue, float64(getResult), endpoint, getFilename)
			ch <- prometheus.MustNewConstMetric(c.metrics["del_metric"], prometheus.GaugeValue, float64(delResult), endpoint, delFilename)
			wg.Done()
		}(&wg, ch, endpoint)
	}
	wg.Wait()
}
