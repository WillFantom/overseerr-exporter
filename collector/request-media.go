package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/willfantom/goverseerr"
)

type RequestMediaStatusCollector struct {
	client *goverseerr.Overseerr

	Available     *prometheus.Desc
	PartAvailable *prometheus.Desc
	Processing    *prometheus.Desc
	Pending       *prometheus.Desc
	Unknown       *prometheus.Desc
}

func NewRequestMediaStatusCollector(client *goverseerr.Overseerr) *RequestMediaStatusCollector {
	logrus.Traceln("defining request media status collector")
	specificNamespace := "request_media_status"
	return &RequestMediaStatusCollector{
		client: client,

		Available: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "available"),
			"Number of requests where the media is available to watch",
			nil,
			nil,
		),
		PartAvailable: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "part_available"),
			"Number of requests where the media is partially available to watch",
			nil,
			nil,
		),
		Processing: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "processing"),
			"Number of requests where the media is currently processing",
			nil,
			nil,
		),
		Pending: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "pending"),
			"Number of requests where the media is currently pending",
			nil,
			nil,
		),
		Unknown: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "unknown"),
			"Number of requests where the media status is unknown",
			nil,
			nil,
		),
	}
}

func (rc *RequestMediaStatusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- rc.Available
	ch <- rc.PartAvailable
	ch <- rc.Processing
	ch <- rc.Pending
	ch <- rc.Unknown
}

func (rc *RequestMediaStatusCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Debug("collecting request media status data...")
	start := time.Now()

	available := 0
	partAvailable := 0
	processing := 0
	pending := 0
	unknown := 0

	allRequests := fetchAllRequests(rc.client)

	for _, request := range allRequests {
		switch request.Media.Status {
		case goverseerr.MediaStatusAvailable:
			available++
		case goverseerr.MediaStatusPartial:
			partAvailable++
		case goverseerr.MediaStatusProcessing:
			processing++
		case goverseerr.MediaStatusPending:
			pending++
		case goverseerr.MediaStatusUnknown:
			unknown++
		}
	}
	ch <- prometheus.MustNewConstMetric(
		rc.Available,
		prometheus.GaugeValue,
		float64(available),
	)
	ch <- prometheus.MustNewConstMetric(
		rc.PartAvailable,
		prometheus.GaugeValue,
		float64(partAvailable),
	)
	ch <- prometheus.MustNewConstMetric(
		rc.Processing,
		prometheus.GaugeValue,
		float64(processing),
	)
	ch <- prometheus.MustNewConstMetric(
		rc.Pending,
		prometheus.GaugeValue,
		float64(pending),
	)
	ch <- prometheus.MustNewConstMetric(
		rc.Unknown,
		prometheus.GaugeValue,
		float64(unknown),
	)

	elapsed := time.Since(start)
	logrus.WithField("time_elapsed", elapsed).Debugln("request media status data collected")
}
