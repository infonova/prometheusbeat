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
  "@timestamp": "2018-09-28T11:44:07.006Z",
  "@metadata": {
    "beat": "prometheusbeat",
    "type": "doc",
    "version": "6.4.1"
  },
  "labels": {
    "le": "10",
    "name": "prometheus_tsdb_tombstone_cleanup_seconds_bucket",
    "instance": "localhost:9090",
    "job": "prometheus"
  },
  "value": 0,
  "host": {
    "name": "example.com"
  },
  "beat": {
    "version": "6.4.1",
    "name": "prometheusbeat",
    "hostname": "example.com"
  }
}
```

Ensure that this folder is at the following location:
`${GOPATH}/github.com/infonova/prometheusbeat`

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
make
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

To clean  Prometheusbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Prometheusbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/infonova/prometheusbeat
cd ${GOPATH}/github.com/infonova/prometheusbeat
git clone https://github.com/infonova/prometheusbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The whole process to finish can take several minutes.
