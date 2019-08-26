# Prometheusbeat

Prometheusbeat is an elastic [beat](https://www.elastic.co/products/beats) that can receive Prometheus metrics via the remote write feature.

Example Prometheusbeat configuration:

```
prometheusbeat:
  listen: ":8080"
  context: "/prometheus"

[...]
```

Example Prometheus configuration:

```
[...]

remote_write:
  url: "http://localhost:8080/prometheus"

[...]
```

Example Prometheusbeat event:

```
{
  "@timestamp": "2019-08-26T11:34:04.253Z",
  "@metadata": {
    "beat": "prometheusbeat",
    "type": "_doc",
    "version": "7.3.1"
  },
  "labels": {
    "instance": "localhost:9090",
    "job": "prometheus"
  },
  "ecs": {
    "version": "1.0.1"
  },
  "host": {
    "containerized": false,
    "hostname": "test",
    "architecture": "x86_64",
    "os": {
      "kernel": "4.14.14-1.el7.elrepo.x86_64",
      "codename": "Core",
      "platform": "centos",
      "version": "7 (Core)",
      "family": "redhat",
      "name": "CentOS Linux"
    },
    "id": "338bc9f83bf343dfa1983fc2bc43bd0f",
    "name": "test"
  },
  "agent": {
    "version": "7.3.1",
    "type": "prometheusbeat",
    "ephemeral_id": "51cadc55-5746-4ad2-a278-b0dcac12942b",
    "hostname": "bhois.station",
    "id": "f9c5188e-a1ea-4fbb-adf6-8494f56c59bf"
  },
  "name": "scrape_series_added",
  "value": 349
}
```

## Getting Started with Prometheusbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init Project
To get running with Prometheusbeat and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Prometheusbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/infonova/prometheusbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Prometheusbeat run the command below. This will generate a binary
in the same directory with the name prometheusbeat.

```
mage build
```


### Run

To run Prometheusbeat with debugging output enabled, run:

```
./prometheusbeat -c prometheusbeat.yml -e -d "*"
```


### Test

To test Prometheusbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/prometheusbeat.template.json and etc/prometheusbeat.asciidoc

```
make update
```


### Cleanup

To clean  Prometheusbeat source code, run the following command:

```
make fmt
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Prometheusbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/infonova
cd ${GOPATH}/github.com/infonova
git clone https://github.com/infonova/prometheusbeat
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make release
```

This will fetch and create all images required for the build process. The whole process to finish can take several minutes.

## Building Docker Image

Prometheusbeat Docker image can be customized and build to run on any environments

```
docker build . -t infonova/prometheusbeat:latest
```

## Running Docker Image

Prometheusbeat Docker image can be run using below command

```
docker run -d  --name prometheusbeat -p <<dockerhost-port>>:8080 -v <<host-config-path>>:/prometheusbeat.yml infonova/prometheusbeat:latest
```
