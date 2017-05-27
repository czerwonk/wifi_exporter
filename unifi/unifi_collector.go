package unifi

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	apStateDesc       *prometheus.Desc
	apClientsDesc     *prometheus.Desc
	apClientsSsidDesc *prometheus.Desc
)

type UnifiCollector struct {
	Url    string
	Cookie string
}

func init() {
	labels := make([]string, 0)
	labels = append(labels, "site", "ap_name")

	apStateDesc = prometheus.NewDesc("unifi_ap_state", "State of the access point", labels, nil)

	labels = append(labels, "radio")
	apClientsDesc = prometheus.NewDesc("unifi_ap_clients", "Number of users", labels, nil)

	labels = append(labels, "ssid")
	apClientsSsidDesc = prometheus.NewDesc("unifi_ap_ssid_clients", "Number of users for ssid", labels, nil)
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
	for _, s := range sites {
		c.exportForSite(s, ch)
	}
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

	ch <- prometheus.MustNewConstMetric(apClientsDesc, prometheus.GaugeValue, float64(ap.clientsG), append(labelValues, "ng")...)
	ch <- prometheus.MustNewConstMetric(apClientsDesc, prometheus.GaugeValue, float64(ap.clientsN), append(labelValues, "na")...)

	for _, ssid := range ap.ssids {
		ch <- prometheus.MustNewConstMetric(apClientsSsidDesc, prometheus.GaugeValue, float64(ssid.clientsG), append(labelValues, "ng", ssid.name)...)
		ch <- prometheus.MustNewConstMetric(apClientsSsidDesc, prometheus.GaugeValue, float64(ssid.clientsN), append(labelValues, "na", ssid.name)...)
	}
}

// Describe implements Prometheus Collector interface
func (c *UnifiCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- apStateDesc
	ch <- apClientsDesc
	ch <- apClientsSsidDesc
}
