package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/willfantom/goverseerr"
)

const (
	userPageSize int = 50
)

type UserCollector struct {
	client *goverseerr.Overseerr

	Requests *prometheus.Desc
}

func NewUserCollector(client *goverseerr.Overseerr) *UserCollector {
	logrus.Traceln("defining user collector")
	specificNamespace := "user"
	return &UserCollector{
		client: client,

		Requests: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, specificNamespace, "requests"),
			"Number of requests made by a user",
			[]string{"email"},
			nil,
		),
	}
}

func (rc *UserCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- rc.Requests
}

func (rc *UserCollector) Collect(ch chan<- prometheus.Metric) {
	logrus.Debugln("collecting user data...")
	start := time.Now()

	var allUsers []*goverseerr.User

	page := 0
	for {
		logrus.WithField("page", page).Traceln("fetching user list from overseerr")
		users, pageInfo, err := rc.client.GetAllUsers(userPageSize, page)
		if err != nil {
			logrus.WithField("page", page).Errorln("failed to get page of users from overseerr")
			return
		}
		allUsers = append(allUsers, users...)
		page++
		if page >= (pageInfo.Pages - 1) {
			break
		}
	}
	logrus.WithField("total_users", len(allUsers)).Traceln("fetched all users from overseerr")

	for _, user := range allUsers {
		ch <- prometheus.MustNewConstMetric(
			rc.Requests,
			prometheus.GaugeValue,
			float64(user.RequestCount),
			user.Email,
		)
	}

	elapsed := time.Since(start)
	logrus.WithField("time_elapsed", elapsed).Debugln("user data collected")
}
