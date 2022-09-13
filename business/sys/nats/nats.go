// Package nats provides support for connecting to NATS server.
package nats

import (
	"time"

	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
)

type Config struct {
	ClusterID string
	ClientID  string
	Host      string
}

type NATS struct {
	AckWait time.Duration
	Client  stan.Conn
}

func Connect(cfg Config) (*NATS, error) {
	nc, err := nats.Connect(cfg.Host)
	if err != nil {
		return nil, err
	}

	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsConn(nc))
	if err != nil {
		return nil, err
	}

	aw, _ := time.ParseDuration("60s")

	n := NATS{
		AckWait: aw,
		Client:  sc,
	}

	return &n, nil
}

func (n NATS) Subscribe(subject string, queueGroupName string, cb stan.MsgHandler) error {
	_, err := n.Client.QueueSubscribe(subject, queueGroupName, cb,
		stan.DeliverAllAvailable(),
		stan.SetManualAckMode(),
		stan.AckWait(n.AckWait),
		stan.DurableName(queueGroupName),
	)

	if err != nil {
		return err
	}

	return nil
}
