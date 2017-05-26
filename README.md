# wifi_exporter [![Build Status](https://travis-ci.org/czerwonk/wifi_exporter.svg)][travis]
Metric exporter for wireless controllers to use with Prometheus

# Remarks
This is an early version. It uses undocumented API calls so it can break at any  time.

# Install
```
go get -u github.com/czerwonk/wifi_exporter
```

# Use
```
wifi_exporter -api.url $apiurl -api.user $user -api.pass $pass
```

# Use with Docker
```
docker run -d -e apiurl="http://unifi" -e user="username" -e pass="secret" czerwonk/wifi_exporter
```

## License
(c) Daniel Czerwonk, 2016. Licensed under [MIT](LICENSE) license.

# Prometheus
see https://prometheus.io/

# Unfi Controller
Unifi is a registered trademark of Ubiquiti Networks

see https://www.ubnt.com/enterprise/software/

[travis]: https://travis-ci.org/czerwonk/wifi_exporter
