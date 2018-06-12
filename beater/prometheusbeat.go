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
			event := beat.Event{
				Timestamp: time.Unix(0, pevent["timestamp"].(int64)*1000000),
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
