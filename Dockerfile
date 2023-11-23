#FROM registry.access.redhat.com/ubi8/go-toolset:1.20.10-3
FROM golang:1.21.3
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o /descaler
CMD [ "/descaler" ]