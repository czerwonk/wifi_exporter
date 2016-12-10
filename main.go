/*
Copyright 2016 Daniel Czerwonk (d.czerwonk@gmail.com)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const version string = "0.1"

var (
	showVersion   = flag.Bool("version", false, "Print version information.")
	listenAddress = flag.String("web.listen-address", ":9120", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	apiUrl        = flag.String("controller.url", "http://unifi", "base URL to the Unifi Controller API")
	apiUser       = flag.String("controller.user", "username", "username to access the Unifi Controller API")
	apiPass       = flag.String("controller.pass", "password", "password to authorize user")
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
	fmt.Println("ubnt_wifi_exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Metric exporter for unifi controller")
}

func startServer() {
	fmt.Printf("Starting unifi exporter (Version: %s)\n", version)
	http.HandleFunc(*metricsPath, handleMetricsRequest)

	fmt.Printf("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	cookie, err := getCookie()

	if err != nil {
		log.Println(err)
		return
	}

	printMetricsForSites(cookie, w)
}

func printMetricsForSites(cookie string, w io.Writer) {
	sites, err := getSites(cookie)

	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%d sites\n", len(sites))

	for _, s := range sites {
		printMetricsForSite(s, cookie, w)
	}
}

func printMetricsForSite(s *site, cookie string, w io.Writer) {
	aps, err := getAccessPoints(s.id, cookie)

	if err != nil {
		log.Println(err)
		return
	}

	for _, ap := range aps {
		printMetricsForAccessPoint(ap, s, w)
	}
}

func printMetricsForAccessPoint(ap *accessPoint, s *site, w io.Writer) {
	fmt.Fprintf(w, "unifi_ap_state{site=\"%s\",ap_name=\"%s\"} %d\n", s.name, ap.name, ap.state)
	fmt.Fprintf(w, "unifi_ap_clients{site=\"%s\",ap_name=\"%s\",radio=\"na\"} %d\n", s.name, ap.name, ap.clientsN)
	fmt.Fprintf(w, "unifi_ap_clients{site=\"%s\",ap_name=\"%s\",radio=\"ng\"} %d\n", s.name, ap.name, ap.clientsG)
}
