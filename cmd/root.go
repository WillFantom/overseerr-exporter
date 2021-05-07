package cmd

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/willfantom/goverseerr"
	"github.com/willfantom/overseerr-exporter/collector"
)

// persistent flags
var (
	logLevel           string
	overseerrAddress   string
	overseerrAPIKey    string
	listenAddress      string
	metricsPath        string
	overseerrAPILocale string

	scrapeGenres    bool
	scrapeCompanies bool
)

// instance to use
var overseerr *goverseerr.Overseerr

var RootCmd = &cobra.Command{
	Use:   "overseerr-exporter",
	Short: "Export request metrics from Overseerr",
	Long:  `Export request metrics from an Overseerr instance to a prometheus database`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setupLogger()
		logrus.WithFields(logrus.Fields{
			"command": cmd.Name(),
			"args":    args,
		}).Debugln("running command")
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		setOverseer()
	},
	Run: func(cmd *cobra.Command, args []string) {
		prometheus.MustRegister(prometheus.NewBuildInfoCollector())
		prometheus.MustRegister(collector.NewRequestCollector(overseerr, scrapeGenres, scrapeCompanies))
		prometheus.MustRegister(collector.NewUserCollector(overseerr))

		handler := promhttp.Handler()
		http.Handle(metricsPath, handler)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
		<head><title>Overseerr Exporter</title></head>
		<body>
		<h1>Overseerr Exporter</h1>
		<p><a href="` + metricsPath + `">Metrics</a></p>
		</body>
		</html>`))
		})

		if err := http.ListenAndServe(listenAddress, nil); err != nil {
			logrus.WithField("err msg", err.Error()).Fatalln("🆘	http server failed: exiting")
		}
	},
}

func setupLogger() {
	if level, err := logrus.ParseLevel(logLevel); err != nil {
		logrus.SetLevel(logrus.FatalLevel)
	} else {
		logrus.SetLevel(level)
	}
}

func setOverseer() {
	if o, err := goverseerr.NewKeyAuth(overseerrAddress, nil, overseerrAPILocale, overseerrAPIKey); err != nil {
		logrus.WithField("message", err.Error()).Fatalln("Could not connect to Overseerr")
	} else {
		overseerr = o
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&logLevel, "log", "fatal", "set the log level (fatal, error, info, debug, trace)")

	// overseerr setup
	RootCmd.PersistentFlags().StringVar(&overseerrAddress, "overseerr.address", "", "Address at which Overseerr is hosted.")
	RootCmd.PersistentFlags().StringVar(&overseerrAPIKey, "overseerr.api-key", "", "API key for admin access to the Overseerr instance.")
	RootCmd.PersistentFlags().StringVar(&overseerrAPILocale, "overseerr.locale", "en", "Locale of the Overseerr instance.")
	RootCmd.PersistentFlags().BoolVar(&scrapeGenres, "overseerr.scrape.genres", true, "Scrape genere details from the media requests.")
	RootCmd.PersistentFlags().BoolVar(&scrapeGenres, "overseerr.scrape.companies", true, "Scrape company/network details from the media requests.")
	RootCmd.MarkPersistentFlagRequired("overseerr.address")
	RootCmd.MarkPersistentFlagRequired("overseerr.api-key")

	// setup vars (based on ha proxy exporter)
	RootCmd.PersistentFlags().StringVar(&listenAddress, "web.listen-address", ":9850", "Address to listen on for web interface and telemetry.")
	RootCmd.PersistentFlags().StringVar(&metricsPath, "web.telemetry-path", "/metrics", "Path under which to expose metrics.")
}
