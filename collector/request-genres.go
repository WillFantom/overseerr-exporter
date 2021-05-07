package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/willfantom/goverseerr"
)

type RequestGenreCollector struct {
	client *goverseerr.Overseerr

	RequestCount *prometheus.Desc
}

type RequestGenreLabel struct {
	MediaType string
	Genre     string
}

func NewRequestGenereCollector(client *goverseerr.Overseerr) *RequestGenreCollector {
	logrus.Traceln("defining request genre collector")
	specificNamespace := "request_genre"
	return &RequestGenreCollector{
		client: client,

		RequestCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "count"),
			"Number of requests for a given genre",
			[]string{"genre", "media_type"},
			nil,
		),
	}
}

func (rc *RequestGenreCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- rc.RequestCount
}

func (rc *RequestGenreCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Debug("collecting request genre data...")
	start := time.Now()

	genreCounts := make(map[RequestGenreLabel]int)

	allRequests := fetchAllRequests(rc.client)

	for _, request := range allRequests {
		switch request.Media.MediaType {
		case goverseerr.MediaTypeMovie:
			if details, err := request.GetMovieDetails(rc.client); err == nil {
				for _, g := range details.Genres {
					genreCounts[RequestGenreLabel{
						MediaType: string(request.Media.MediaType),
						Genre:     g.Name,
					}]++
				}
			}
		case goverseerr.MediaTypeTV:
			if details, err := request.GetTVDetails(rc.client); err == nil {
				for _, g := range details.Genres {
					genreCounts[RequestGenreLabel{
						MediaType: string(request.Media.MediaType),
						Genre:     g.Name,
					}]++
				}
			}
		}
	}

	for labelSet, count := range genreCounts {
		ch <- prometheus.MustNewConstMetric(
			rc.RequestCount,
			prometheus.GaugeValue,
			float64(count),
			labelSet.Genre, labelSet.MediaType,
		)
	}

	elapsed := time.Since(start)
	logrus.WithField("time_elapsed", elapsed).Debugln("request genre data collected")
}
