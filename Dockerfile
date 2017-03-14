FROM golang

RUN apt-get install -y git

# Install application
RUN go get github.com/czerwonk/wifi_exporter

# Run the application and expose the port
CMD wifi_exporter -api.url $apiurl -api.user $user -api.pass $pass
EXPOSE 9120