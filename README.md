# wifi_exporter 
[![Build Status](https://travis-ci.org/czerwonk/wifi_exporter.svg)][travis]
[![Docker Build Statu](https://img.shields.io/docker/build/czerwonk/wifi_exporter.svg)][dockerbuild]
[![Go Report Card](https://goreportcard.com/badge/github.com/czerwonk/wifi_exporter)][goreportcard]

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

## Unfi Controller
Unifi is a registered trademark of Ubiquiti Networks

see https://www.ubnt.com/enterprise/software/

[travis]: https://travis-ci.org/czerwonk/wifi_exporter
[dockerbuild]: https://hub.docker.com/r/czerwonk/wifi_exporter/builds
[goreportcard]: https://goreportcard.com/report/github.com/czerwonk/wifi_exporter
