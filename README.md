# wifi_exporter 
[![Build Status](https://travis-ci.org/czerwonk/wifi_exporter.svg)](https://travis-ci.org/czerwonk/wifi_exporter)
[![Docker Build Statu](https://img.shields.io/docker/build/czerwonk/wifi_exporter.svg)](https://hub.docker.com/r/czerwonk/wifi_exporter/builds)
[![Go Report Card](https://goreportcard.com/badge/github.com/czerwonk/wifi_exporter)](https://goreportcard.com/report/github.com/czerwonk/wifi_exporter)

Metric exporter for wireless controllers to use with Prometheus

## Remarks
This is an early version. It uses undocumented API calls so it can break at any  time.

## Install
```
go get -u github.com/czerwonk/wifi_exporter
```

## Use
```
wifi_exporter -config.path config.yml
```

## Use with Docker
```
docker run -d -v config.yml:/etc/wifi_exporter.yml
```

## License
(c) Daniel Czerwonk, 2016. Licensed under [MIT](LICENSE) license.

## Prometheus
see https://prometheus.io/

## Unifi Controller
Unifi is a registered trademark of Ubiquiti Networks

see https://www.ubnt.com/enterprise/software/
