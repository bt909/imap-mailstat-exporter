FROM golang:1.24 as build

WORKDIR /go/src/imap-mailstat-exporter
COPY . .

RUN go mod download
RUN go vet -v ./...
RUN go test -v ./...

RUN CGO_ENABLED=0 go build -o /go/bin/imap-mailstat-exporter cmd/imap-mailstat-exporter/main.go

FROM gcr.io/distroless/static:nonroot

LABEL org.opencontainers.image.source="https://github.com/bt909/imap-mailstat-exporter"
LABEL org.opencontainers.image.author="Thomas Belián <thomas.belian@bt909.de>"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.documentation="https://github.com/bt909/imap-mailstat-exporter"
LABEL org.opencontainers.image.description="Prometheus Exporter to provide some metrics about your IMAP mailboxes"
LABEL org.opencontainers.image.vendor="Thomas Belián"
LABEL org.opencontainers.image.title="IMAP Mailstat Exporter"
LABEL org.opencontainers.image.url="https://github.com/bt909/imap-mailstat-exporter"

COPY --from=build /go/bin/imap-mailstat-exporter /
ENTRYPOINT ["/imap-mailstat-exporter"]
EXPOSE 8081
