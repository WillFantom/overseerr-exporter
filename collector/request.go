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

type RequestCollector struct {
	client    *goverseerr.Overseerr
	doGenre   bool
	doCompany bool

	Count *prometheus.Desc
}

type RequestMetricLabel struct {
	MediaType     string
	RequestStatus string
	MediaStatus   string
	UHD           string
	Company       string
	Genre         string
}

func NewRequestCollector(client *goverseerr.Overseerr, doGenre, doCompany bool) *RequestCollector {
	logrus.Traceln("defining request collector")
	specificNamespace := "requests"
	return &RequestCollector{
		client:    client,
		doGenre:   doGenre,
		doCompany: doCompany,

		Count: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "count"),
			"Number of requests on Overseerr",
			[]string{"media_type", "is_4k", "request_status", "media_status", "genre", "company"},
			nil,
		),
	}
}

func (rc *RequestCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- rc.Count
}

func (rc *RequestCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Debug("collecting request data...")
	start := time.Now()

	counts := make(map[RequestMetricLabel]int)
	allRequests := fetchAllRequests(rc.client)

	for _, request := range allRequests {

		genre := "not_collected"
		company := "not_collected"
		switch request.Media.MediaType {
		case goverseerr.MediaTypeMovie:
			if details, err := request.GetMovieDetails(rc.client); err == nil {
				if len(details.Genres) > 0 {
					genre = details.Genres[0].Name
				}
				if len(details.ProductionCompanies) > 0 {
					company = details.ProductionCompanies[0].Name
				}
			}
		case goverseerr.MediaTypeTV:
			if details, err := request.GetTVDetails(rc.client); err == nil {
				if len(details.Genres) > 0 {
					genre = details.Genres[0].Name
				}
				if len(details.Networks) > 0 {
					company = details.Networks[0].Name
				}
			}
		}

		counts[RequestMetricLabel{
			RequestStatus: request.Status.ToString(),
			MediaStatus:   request.Media.Status.ToString(),
			MediaType:     string(request.Media.MediaType),
			UHD:           strconv.FormatBool(request.IsUHD),
			Genre:         genre,
			Company:       company,
		}]++
	}

	for labelSet, count := range counts {
		ch <- prometheus.MustNewConstMetric(
			rc.Count,
			prometheus.GaugeValue,
			float64(count),
			labelSet.MediaType, labelSet.UHD, labelSet.RequestStatus, labelSet.MediaStatus, labelSet.Genre, labelSet.Company,
		)
	}

	elapsed := time.Since(start)
	logrus.WithField("time_elapsed", elapsed).Debugln("request data collected")
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
