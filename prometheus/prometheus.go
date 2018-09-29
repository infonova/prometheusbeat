package prometheus

import (
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	version "github.com/hashicorp/go-version"
	"github.com/infonova/prometheusbeat/config"
	"github.com/prometheus/prometheus/prompb"
)

type PrometheusServer struct {
	config           config.Config
	prometheusEvents chan common.MapStr
}

func NewPrometheusServer(cfg config.Config) *PrometheusServer {
	promSrv := &PrometheusServer{
		config: cfg,
	}

	return promSrv
}

func (promSrv *PrometheusServer) Start(events chan common.MapStr) {
	promSrv.prometheusEvents = events

	http.HandleFunc(promSrv.config.Context, promSrv.handlePrometheus)
	log.Fatal(http.ListenAndServe(promSrv.config.Listen, nil))
}

func (promSrv *PrometheusServer) handlePrometheus(w http.ResponseWriter, r *http.Request) {

	v := r.Header.Get("X-Prometheus-Remote-Write-Version")
	//No header indicates old prometheus version
	if len(v) == 0 {
		v = "0.0.0"
	}

	baseVer, _ := version.NewVersion("0.1.0")
	reqVer, err := version.NewVersion(v)
	if err != nil {
		logp.Err(strings.Join([]string{"wrong prometheus remote write version:", v}, " "))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var reqBuf []byte
	if reqVer.LessThan(baseVer) {
		reqBuf, err = ioutil.ReadAll(snappy.NewReader(r.Body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		compressed, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reqBuf, err = snappy.Decode(nil, compressed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var req prompb.WriteRequest
	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, ts := range req.Timeseries {
		event := map[string]interface{}{}
		labels := map[string]interface{}{}
		for _, l := range ts.Labels {
			fieldName := strings.Trim(l.Name, "_")
			if fieldName == "name" && promSrv.config.NameUnderRoot {
				event[fieldName] = l.Value
			} else {
				labels[fieldName] = l.Value
			}
		}
		event["labels"] = labels

		for _, s := range ts.Samples {
			if math.IsNaN(s.Value) {
				event["tags"] = []string{"nan"}
			} else {
				event["value"] = s.Value
			}
			event["timestamp"] = s.Timestamp
		}

		promSrv.prometheusEvents <- event
	}
}

func (promSrv *PrometheusServer) Shutdown() {
	close(promSrv.prometheusEvents)
}
