FROM golang:1-alpine as builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build -ldflags '-s -w' ./cmd/main.go

FROM golang:1-alpine as download

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.24 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /build/main /opt/app/
COPY --from=download /bin/grpc_health_probe /bin/grpc_health_probe
ENTRYPOINT ["/opt/app/main"]