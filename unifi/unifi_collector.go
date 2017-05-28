package unifi

import (
	"sync"

	"github.com/czerwonk/wifi_exporter/common"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	apStateDesc   *prometheus.Desc
	apClientsDesc *prometheus.Desc
)

type UnifiCollector struct {
	Url    string
	Cookie string
}

func init() {
	labels := make([]string, 0)
	labels = append(labels, "site", "ap_name")

	apStateDesc = prometheus.NewDesc("wifi_unifi_ap_state", "State of the access point", labels, nil)

	labels = append(labels, "radio", "ssid")
	apClientsDesc = prometheus.NewDesc("wifi_unifi_ap_clients", "Number of clients on AP", labels, nil)
}

// NewUnifiCollector create a collector to get metrics from unifi controller
func NewUnifiCollector(apiUrl, apiUser, apiPass string) (*UnifiCollector, error) {
	cookie, err := getCookie(apiUrl, apiUser, apiPass)

	if err != nil {
		return nil, err
	}

	return &UnifiCollector{Cookie: cookie, Url: apiUrl}, nil
}

// Collect implements Prometheus Collector interface
func (c *UnifiCollector) Collect(ch chan<- prometheus.Metric) {
	sites, err := getSites(c.Cookie, c.Url)

	if err != nil {
		return
	}

	c.exportForSites(sites, ch)
}

func (c *UnifiCollector) exportForSites(sites []*site, ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	wg.Add(len(sites))

	for _, s := range sites {
		go func(site *site) {
			defer wg.Done()
			c.exportForSite(site, ch)
		}(s)
	}

	wg.Wait()
}

func (c *UnifiCollector) exportForSite(s *site, ch chan<- prometheus.Metric) {
	aps, err := getAccessPoints(s.id, c.Cookie, c.Url)

	if err != nil {
		return
	}

	for _, ap := range aps {
		c.exportForAccessPoint(s, ap, ch)
	}
}

func (c *UnifiCollector) exportForAccessPoint(s *site, ap *accessPoint, ch chan<- prometheus.Metric) {
	labelValues := make([]string, 0)

	name := ap.name
	if len(ap.name) == 0 {
		name = ap.mac
	}

	labelValues = append(labelValues, s.name, name)
	ch <- prometheus.MustNewConstMetric(apStateDesc, prometheus.GaugeValue, float64(ap.state), labelValues...)

	up := 0
	if ap.state == 1 {
		up = 1
	}
	ch <- common.MustNewMetricForUp(s.name, name, up)

	clients := 0
	for _, ssid := range ap.ssids {
		ch <- prometheus.MustNewConstMetric(apClientsDesc, prometheus.GaugeValue, float64(ssid.clientsG), append(labelValues, "ng", ssid.name)...)
		ch <- prometheus.MustNewConstMetric(apClientsDesc, prometheus.GaugeValue, float64(ssid.clientsN), append(labelValues, "na", ssid.name)...)

		clients += ssid.clientsG + ssid.clientsN
	}

	ch <- common.MustNewMetricForClients(s.name, name, clients)
}

// Describe implements Prometheus Collector interface
func (c *UnifiCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- apStateDesc
	ch <- apClientsDesc
	ch <- common.AccessPointUpDesc
	ch <- common.AccessPointClientsDesc
}
