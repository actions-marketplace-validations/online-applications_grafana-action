FROM golang:alpine as build-env
# All these steps will be cached
RUN mkdir /build
WORKDIR /build
# COPY the source code as the last step
RUN apk --no-cache add ca-certificates
COPY grafana.go grafana.go
COPY go.mod go.mod

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/addGrafanaAnnotations

FROM scratch
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /go/bin/addGrafanaAnnotations /go/bin/addGrafanaAnnotations
ENTRYPOINT ["/go/bin/addGrafanaAnnotations"]