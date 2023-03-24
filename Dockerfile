ARG GO_IMAGE=golang:alpine3.17
ARG ALPINE_IMAGE=alpine:3.17.2

FROM ${GO_IMAGE} AS builder
WORKDIR /app
COPY go.mod main.go ./
RUN go build

FROM ${ALPINE_IMAGE}
WORKDIR /app
COPY --from=builder /app/udp-echo-test .

ENTRYPOINT ["./udp-echo-test"]

EXPOSE 53/tcp
EXPOSE 53/udp