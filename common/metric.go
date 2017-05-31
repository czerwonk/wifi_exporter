package common

import "github.com/prometheus/client_golang/prometheus"

func MustNewMetricForUp(site string, name string, status int) prometheus.Metric {
	l := []string{site, name}

	return prometheus.MustNewConstMetric(AccessPointUpDesc, prometheus.GaugeValue, float64(status), l...)
}

func MustNewMetricForClients(site string, name string, clients int) prometheus.Metric {
	l := []string{site, name}

	return prometheus.MustNewConstMetric(AccessPointClientsDesc, prometheus.GaugeValue, float64(clients), l...)
}
