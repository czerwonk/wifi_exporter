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
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const version string = "0.2"

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
	fmt.Println("ubnt_wifi_exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Metric exporter for unifi controller")
}

func startServer() {
	fmt.Printf("Starting ubnt wifi exporter (Version: %s)\n", version)
	http.HandleFunc(*metricsPath, errorHandler(handleMetricsRequest))

	fmt.Printf("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func errorHandler(f func(io.Writer, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		wr := bufio.NewWriter(&buf)
		err := f(wr, r)
		wr.Flush()

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_, err = w.Write(buf.Bytes())

		if err != nil {
			log.Println(err)
		}
	}
}

func handleMetricsRequest(w io.Writer, r *http.Request) error {
	cookie, err := getCookie()

	if err != nil {
		return err
	}

	return printMetricsForSites(cookie, w)
}

func printMetricsForSites(cookie string, w io.Writer) error {
	sites, err := getSites(cookie)

	if err != nil {
		return err
	}

	log.Printf("%d sites\n", len(sites))

	for _, s := range sites {
		if err = printMetricsForSite(s, cookie, w); err != nil {
			return err
		}
	}

	return nil
}

func printMetricsForSite(s *site, cookie string, w io.Writer) error {
	aps, err := getAccessPoints(s.id, cookie)

	if err != nil {
		return err
	}

	for _, ap := range aps {
		printMetricsForAccessPoint(ap, s, w)
	}

	return nil
}

func printMetricsForAccessPoint(ap *accessPoint, s *site, w io.Writer) {
	fmt.Fprintf(w, "unifi_ap_state{site=\"%s\",ap_name=\"%s\"} %d\n", s.name, ap.name, ap.state)
	fmt.Fprintf(w, "unifi_ap_clients{site=\"%s\",ap_name=\"%s\",radio=\"na\"} %d\n", s.name, ap.name, ap.clientsN)
	fmt.Fprintf(w, "unifi_ap_clients{site=\"%s\",ap_name=\"%s\",radio=\"ng\"} %d\n", s.name, ap.name, ap.clientsG)
}