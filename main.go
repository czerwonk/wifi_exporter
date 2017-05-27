package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/czerwonk/wifi_exporter/unifi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const version string = "0.3.0"

var (
	showVersion   = flag.Bool("version", false, "Print version information.")
	listenAddress = flag.String("web.listen-address", ":9120", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	apiUrl        = flag.String("api.url", "http://unifi", "base URL to the Unifi Controller API")
	apiUser       = flag.String("api.user", "username", "username to access the Unifi Controller API")
	apiPass       = flag.String("api.pass", "password", "password to authorize user")
)

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	startServer()
}

func printVersion() {
	fmt.Println("wifi_exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): Daniel Czerwonk")
	fmt.Println("Metric exporter for wifi controllers")
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
			log.Println(err)
		}
	}
}

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) error {
	c, err := unifi.NewUnifiCollector(*apiUrl, *apiUser, *apiPass)

	if err != nil {
		return err
	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)

	return nil
}
