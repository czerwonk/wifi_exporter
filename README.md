# ubnt_wifi_exporter [![Build Status](https://travis-ci.org/czerwonk/ubnt_wifi_exporter.svg)][travis]
Metric exporter for ubnt wireless controllers to use with Prometheus

# Remarks
This is an early version. It uses undocumented API calls so it can break at any  time.

# Install
```
go get github.com/czerwonk/ubnt_wifi_exporter
```

# Use
```
ubnt_wifi_exporter -api.url $apiurl -api.user $user -api.pass $pass
```

# Use with Docker
```
docker run -d -e apiurl="http://unifi" -e user="username" -e pass="secret" czerwonk/ubnt_wifi_exporter
```

# Prometheus
see https://prometheus.io/

# Unfi Controller
Unifi is a registered trademark of Ubiquiti Networks

see https://www.ubnt.com/enterprise/software/

[travis]: https://travis-ci.org/czerwonk/ubnt_wifi_exporter
