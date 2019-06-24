// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package prometheus

import (
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"

	version "github.com/hashicorp/go-version"
	"github.com/infonova/prometheusbeat/config"
	"github.com/prometheus/prometheus/prompb"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
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

func removeSpecialCharacters(input string, characters string) string {
	filter := func(r rune) rune {
		if strings.IndexRune(characters, r) < 0 {
			return r
		}
		return -1
	}

	trimmedInput := strings.TrimLeft(input, "_")
	return strings.Map(filter, trimmedInput)
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
			// Move timeseries name to root level
			if l.Name == "__name__" {
				event["name"] = l.Value
			} else {
				// Remove special characters
				fieldName := removeSpecialCharacters(l.Name, ", ")
				labels[fieldName] = l.Value
			}
		}
		event["labels"] = labels

		for _, s := range ts.Samples {
			if math.IsNaN(s.Value) {
				event["tags"] = []string{"nan"}
			} else if math.IsInf(s.Value, 0) {
				event["tags"] = []string{"inf"}
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
