FROM golang:1.12 as builder

WORKDIR /go/src/github.com/CarsonBull/mobileCICD/scheduler
COPY . .


RUN go get -u -v github.com/swaggo/swag/cmd/swag && \
    go install github.com/swaggo/swag/cmd/swag && \
    swag init

RUN go get -d -v ./...
RUN go install -v ./...

# CGO_ENABLED=0 has to do with c linkings so we can run this anywhere
# GOOS is setting the os?
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler

FROM alpine:latest as certs
RUN apk --update add ca-certificates

# TODO: Add a test layer in here

FROM scratch

ARG BASE_VERSION=0.0.1
ENV VERSION=${BASE_VERSION}

WORKDIR /be

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/CarsonBull/mobileCICD/scheduler/scheduler .

ENTRYPOINT ["./scheduler"]