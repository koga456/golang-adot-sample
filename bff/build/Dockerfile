FROM golang:1-alpine as builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build -ldflags '-s -w' ./cmd/main.go

FROM golang:1-alpine
RUN apk update && apk add curl
COPY --from=builder /build/main /opt/app/
ENTRYPOINT ["/opt/app/main"]