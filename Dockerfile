FROM golang:1.18 as build

WORKDIR /go/src/imap-mailstat-exporter
COPY . .

RUN go mod download
RUN go vet -v
RUN go test -v

RUN CGO_ENABLED=0 go build -o /go/bin/imap-mailstat-exporter

FROM gcr.io/distroless/static

COPY --from=build /go/bin/imap-mailstat-exporter /
CMD ["/imap-mailstat-exporter"]
EXPOSE 8081
