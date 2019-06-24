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

package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/infonova/prometheusbeat/config"
	"github.com/infonova/prometheusbeat/prometheus"
)

type Prometheusbeat struct {
	done             chan struct{}
	config           config.Config
	client           beat.Client
	prometheusServer *prometheus.PrometheusServer
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Prometheusbeat{
		done:             make(chan struct{}),
		config:           config,
		prometheusServer: prometheus.NewPrometheusServer(config),
	}
	return bt, nil
}

func (bt *Prometheusbeat) Run(b *beat.Beat) error {
	logp.Info("prometheusbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	prometheusEvents := make(chan common.MapStr, 100000)

	go func(events chan common.MapStr) {
		bt.prometheusServer.Start(events)
	}(prometheusEvents)

	var pevent common.MapStr

	for {
		select {
		case <-bt.done:
			return nil
		case pevent = <-prometheusEvents:
			// See #16
			var ts time.Time
			if pevent["timestamp"] != nil {
				ts = time.Unix(0, pevent["timestamp"].(int64)*1000000)
			} else {
				ts = time.Unix(0, time.Now().UnixNano())
			}

			event := beat.Event{
				Timestamp: ts,
				Fields: common.MapStr{
					"name":   pevent["name"],
					"value":  pevent["value"],
					"labels": pevent["labels"],
					"tags":   pevent["tags"],
				},
			}
			bt.client.Publish(event)
		}
	}
}

func (bt *Prometheusbeat) Stop() {
	bt.client.Close()
	close(bt.done)
	bt.prometheusServer.Shutdown()
}
