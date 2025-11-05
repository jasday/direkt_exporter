
FROM golang:1.25-alpine AS builder

ENV GOOS=linux GOARCH=amd64

WORKDIR /direkt
COPY go.mod go.sum ./

RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -trimpath -o exporter


FROM alpine:3.20
RUN adduser -D exporter
USER exporter
COPY --from=builder /direkt/exporter /usr/local/bin/exporter

EXPOSE 9110
ENTRYPOINT ["exporter"]
