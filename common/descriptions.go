package common

import "github.com/prometheus/client_golang/prometheus"

var (
	AccessPointUpDesc      *prometheus.Desc
	AccessPointClientsDesc *prometheus.Desc
)

func init() {
	labels := []string{"site", "ap_name"}

	AccessPointUpDesc = prometheus.NewDesc("wifi_ap_up", "0 = AP is down, 1 = AP is running", labels, nil)
	AccessPointClientsDesc = prometheus.NewDesc("wifi_ap_clients", "Number of clients on AP", labels, nil)
}
