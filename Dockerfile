FROM golang:1.20 as build

WORKDIR /go/src/imap-mailstat-exporter
COPY . .

RUN go mod download
RUN go vet -v
RUN go test -v

RUN CGO_ENABLED=0 go build -o /go/bin/imap-mailstat-exporter

FROM gcr.io/distroless/static

LABEL org.opencontainers.image.source="https://github.com/bt909/imap-mailstat-exporter"

COPY --from=build /go/bin/imap-mailstat-exporter /
CMD ["/imap-mailstat-exporter"]
EXPOSE 8081
