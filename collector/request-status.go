package collector

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/willfantom/goverseerr"
)

const (
	requestPageSize int = 50
)

type RequestStatusCollector struct {
	client *goverseerr.Overseerr

	Approved *prometheus.Desc
	Declined *prometheus.Desc
	Pending  *prometheus.Desc
}

type RequestStatusLabel struct {
	MediaType string
	UHD       string
}

func NewRequestStatusCollector(client *goverseerr.Overseerr) *RequestStatusCollector {
	logrus.Traceln("defining request status collector")
	specificNamespace := "request_status"
	return &RequestStatusCollector{
		client: client,

		Approved: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "approved"),
			"Number of requests that are approved",
			[]string{"media_type", "is_4k"},
			nil,
		),
		Declined: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "declined"),
			"Number of requests that are declined",
			[]string{"media_type", "is_4k"},
			nil,
		),
		Pending: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "pending"),
			"Number of requests that are still pending",
			[]string{"media_type", "is_4k"},
			nil,
		),
	}
}

func (rc *RequestStatusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- rc.Approved
	ch <- rc.Declined
	ch <- rc.Pending
}

func (rc *RequestStatusCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Debug("collecting request status data...")
	start := time.Now()

	approved := make(map[RequestStatusLabel]int)
	declined := make(map[RequestStatusLabel]int)
	pending := make(map[RequestStatusLabel]int)

	allRequests := fetchAllRequests(rc.client)

	for _, request := range allRequests {
		switch request.Status {
		case goverseerr.RequestStatusApproved:
			approved[RequestStatusLabel{
				UHD:       strconv.FormatBool(request.IsUHD),
				MediaType: string(request.Media.MediaType),
			}]++
		case goverseerr.RequestStatusDeclined:
			declined[RequestStatusLabel{
				UHD:       strconv.FormatBool(request.IsUHD),
				MediaType: string(request.Media.MediaType),
			}]++
		case goverseerr.RequestStatusPending:
			pending[RequestStatusLabel{
				UHD:       strconv.FormatBool(request.IsUHD),
				MediaType: string(request.Media.MediaType),
			}]++
		}
	}

	for labelSet, count := range approved {
		ch <- prometheus.MustNewConstMetric(
			rc.Approved,
			prometheus.GaugeValue,
			float64(count),
			labelSet.MediaType, labelSet.UHD,
		)
	}
	for labelSet, count := range declined {
		ch <- prometheus.MustNewConstMetric(
			rc.Declined,
			prometheus.GaugeValue,
			float64(count),
			labelSet.MediaType, labelSet.UHD,
		)
	}
	for labelSet, count := range pending {
		ch <- prometheus.MustNewConstMetric(
			rc.Pending,
			prometheus.GaugeValue,
			float64(count),
			labelSet.MediaType, labelSet.UHD,
		)
	}

	elapsed := time.Since(start)
	logrus.WithField("time_elapsed", elapsed).Debugln("request status data collected")
}

func fetchAllRequests(overseerr *goverseerr.Overseerr) []*goverseerr.MediaRequest {
	var allRequests []*goverseerr.MediaRequest
	page := 0
	for {
		logrus.WithField("page", page).Traceln("fetching request list from overseerr")
		requests, pageInfo, err := overseerr.GetRequests(page, requestPageSize, goverseerr.RequestFileterAll, goverseerr.RequestSortAdded)
		if err != nil {
			logrus.WithField("page", page).Errorln("failed to get page of media requests from overseerr")
			return nil
		}
		allRequests = append(allRequests, requests...)
		page++
		if page >= pageInfo.Pages {
			break
		}
	}
	logrus.WithField("total_requests", len(allRequests)).Traceln("fetched all requests from overseerr")
	return allRequests
}
