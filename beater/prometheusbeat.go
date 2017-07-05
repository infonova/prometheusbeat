package beater

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/infonova/prometheusbeat/config"
	"github.com/infonova/prometheusbeat/prometheus"
)

type Prometheusbeat struct {
	done             chan struct{}
	config           config.Config
	client           publisher.Client
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

	bt.client = b.Publisher.Connect()

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
			pevent["type"] = b.Name
			bt.client.PublishEvent(pevent)
		}
	}
}

func (bt *Prometheusbeat) Stop() {
	bt.client.Close()
	close(bt.done)
	bt.prometheusServer.Shutdown()
}
