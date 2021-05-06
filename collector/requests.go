package collector

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/willfantom/goverseerr"
)

const (
	requestPageSize int = 50
)

type RequestsCollector struct {
	client *goverseerr.Overseerr

	Approved  *prometheus.Desc
	Declined  *prometheus.Desc
	Pending   *prometheus.Desc
	Available *prometheus.Desc
}

func NewRequestsCollector(client *goverseerr.Overseerr) *RequestsCollector {
	logrus.Traceln("🛠	defining requests collector")
	specificNamespace := "requests"
	return &RequestsCollector{
		client: client,

		Approved: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "approved"),
			"Number of requests that are approved",
			[]string{"media_type", "res_type"},
			nil,
		),
		Declined: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "declined"),
			"Number of requests that are declined",
			[]string{"media_type", "res_type"},
			nil,
		),
		Pending: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "pending"),
			"Number of requests that are still pending",
			[]string{"media_type", "res_type"},
			nil,
		),
		Available: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "available"),
			"Number of requests that are available to watch",
			[]string{"media_type", "res_type"},
			nil,
		),
	}
}

func (rc *RequestsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- rc.Approved
	ch <- rc.Declined
	ch <- rc.Pending
	ch <- rc.Available
}

func (rc *RequestsCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Debug("👀	collecting request data...")

	var allRequests []*goverseerr.MediaRequest

	page := 0
	for {
		logrus.WithField("page", page).Traceln("fetching request list from overseerr")
		requests, pageInfo, err := rc.client.GetRequests(page, requestPageSize, goverseerr.RequestFileterAll, goverseerr.RequestSortAdded)
		if err != nil {
			logrus.WithField("page", page).Errorln("failed to get page of media requests from overseerr")
			return
		}
		allRequests = append(allRequests, requests...)
		page++
		if page >= pageInfo.Pages {
			break
		}
	}
	logrus.WithField("total_requests", len(allRequests)).Traceln("fetched all requests from overseerr")

	requestsByMediaType := byMediaType(allRequests)
	for mediaType, requests := range requestsByMediaType {
		logrus.WithField("media_type", mediaType).Traceln("aggregateing data for requests")
		approved := 0
		approvedUHD := 0
		declined := 0
		declinedUHD := 0
		pending := 0
		pendingUHD := 0
		available := 0
		availableUHD := 0
		var wg sync.WaitGroup
		for _, req := range requests {
			wg.Add(1)
			go func(r *goverseerr.MediaRequest) {
				defer wg.Done()
				switch r.Status {
				case goverseerr.RequestStatusApproved:
					approved++
					if r.IsUHD {
						approvedUHD++
					}
				case goverseerr.RequestStatusDeclined:
					declined++
					if r.IsUHD {
						declinedUHD++
					}
				case goverseerr.RequestStatusPending:
					pending++
					if r.IsUHD {
						pendingUHD++
					}
				case goverseerr.RequestStatusAvailable:
					available++
					if r.IsUHD {
						availableUHD++
					}
				}
			}(req)
		}
		wg.Wait()

		ch <- prometheus.MustNewConstMetric(
			rc.Approved,
			prometheus.GaugeValue,
			float64(approved),
			mediaType, "regular",
		)
		ch <- prometheus.MustNewConstMetric(
			rc.Approved,
			prometheus.GaugeValue,
			float64(approvedUHD),
			mediaType, "4k",
		)

		ch <- prometheus.MustNewConstMetric(
			rc.Declined,
			prometheus.GaugeValue,
			float64(declined),
			mediaType, "regular",
		)
		ch <- prometheus.MustNewConstMetric(
			rc.Declined,
			prometheus.GaugeValue,
			float64(declinedUHD),
			mediaType, "4k",
		)

		ch <- prometheus.MustNewConstMetric(
			rc.Pending,
			prometheus.GaugeValue,
			float64(pending),
			mediaType, "regular",
		)
		ch <- prometheus.MustNewConstMetric(
			rc.Pending,
			prometheus.GaugeValue,
			float64(pendingUHD),
			mediaType, "4k",
		)

		ch <- prometheus.MustNewConstMetric(
			rc.Available,
			prometheus.GaugeValue,
			float64(available),
			mediaType, "regular",
		)
		ch <- prometheus.MustNewConstMetric(
			rc.Available,
			prometheus.GaugeValue,
			float64(availableUHD),
			mediaType, "4k",
		)
	}

	logrus.Debugln("✅	request data collected")

}

func byMediaType(requestList []*goverseerr.MediaRequest) map[string][]*goverseerr.MediaRequest {
	requestsMap := make(map[string][]*goverseerr.MediaRequest)
	for _, req := range requestList {
		requestsMap[string(req.Media.MediaType)] = append(requestsMap[string(req.Media.MediaType)], req)
	}
	return requestsMap
}
