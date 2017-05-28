package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/czerwonk/wifi_exporter/configuration"
	"github.com/czerwonk/wifi_exporter/unifi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

const version string = "0.4.0"

var (
	showVersion   = flag.Bool("version", false, "Print version information.")
	listenAddress = flag.String("web.listen-address", ":9120", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	configPath    = flag.String("config.path", "config.yml", "Path to config file")
	config        *configuration.Config
)

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	err := loadConfig()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	startServer()
}

func printVersion() {
	fmt.Println("wifi_exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): Daniel Czerwonk")
	fmt.Println("Metric exporter for wifi controllers")
}

func loadConfig() error {
	var err error
	config, err = configuration.Load(*configPath)

	if err != nil {
		return err
	}

	return nil
}

func startServer() {
	fmt.Printf("Starting wifi exporter (Version: %s)\n", version)
	http.HandleFunc(*metricsPath, errorHandler(handleMetricsRequest))

	fmt.Printf("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)

		if err != nil {
			log.Errorln(err)
		}
	}
}

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) error {
	reg := prometheus.NewRegistry()
	addUnifiCollectors(reg)
	addRuckusCollectors(reg)

	h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		ErrorLog:      log.NewErrorLogger(),
		ErrorHandling: promhttp.ContinueOnError})
	h.ServeHTTP(w, r)

	return nil
}

func addUnifiCollectors(reg *prometheus.Registry) {
	for _, u := range config.Unifi {
		c, err := unifi.NewUnifiCollector(u.ApiUrl, u.ApiUser, u.ApiPass)

		if err != nil {
			log.Errorln(err)
		} else {
			reg.MustRegister(c)
		}
	}
}

func addRuckusCollectors(reg *prometheus.Registry) {

}
