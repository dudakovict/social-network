// Package nats provides support for connecting to NATS streaming server.
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
	QueueGroupName string
	AckWait        time.Duration
	Client         stan.Conn
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
		QueueGroupName: "comments-service",
		AckWait:        aw,
		Client:         sc,
	}

	return &n, nil
}

func (n NATS) Subscribe(subject string, cb stan.MsgHandler) error {
	_, err := n.Client.QueueSubscribe(subject, n.QueueGroupName, cb,
		stan.DeliverAllAvailable(),
		stan.SetManualAckMode(),
		stan.AckWait(n.AckWait),
		stan.DurableName(n.QueueGroupName),
	)

	if err != nil {
		return err
	}

	return nil
}
