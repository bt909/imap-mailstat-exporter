FROM golang:1.20 as build

WORKDIR /go/src/imap-mailstat-exporter
COPY . .

RUN go mod download
RUN go vet -v cmd/imap-mailstat-exporter/main.go
RUN go test -v cmd/imap-mailstat-exporter/main.go

RUN CGO_ENABLED=0 go build -o /go/bin/imap-mailstat-exporter cmd/imap-mailstat-exporter/main.go

FROM gcr.io/distroless/static

LABEL org.opencontainers.image.source="https://github.com/bt909/imap-mailstat-exporter"

COPY --from=build /go/bin/imap-mailstat-exporter /
CMD ["/imap-mailstat-exporter"]
EXPOSE 8081
