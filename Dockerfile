FROM golang:1.11 as builder

RUN apt update &&\
    apt install -y pkg-config libsystemd-dev git gcc curl

COPY . /go/src/github.com/infonova/prometheusbeat

WORKDIR /go/src/github.com/infonova/prometheusbeat

RUN go test -race . ./beater &&\
    go build -ldflags '-s -w' -o /prometheusbeat


FROM debian:stretch-slim

RUN groupadd prometheus &&\
    useradd prometheusbeat &&\
    mkdir /data

COPY --from=builder /prometheusbeat /
COPY prometheusbeat.yml /
RUN chown prometheusbeat. /prometheusbeat /prometheusbeat.yml /data

USER prometheusbeat

ENTRYPOINT ["/prometheusbeat"]

CMD ["-e", "-c", "prometheusbeat.yml"]
