package comment

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/dudakovict/social-network/business/core/comment/db"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

type Listener struct {
	log    *zap.SugaredLogger
	store  db.Store
	client stan.Conn
}

func (l Listener) PostCreated() error {
	l.client.Subscribe("post-created", func(m *stan.Msg) {
		buf := bytes.NewReader(m.Data)
		dec := gob.NewDecoder(buf)

		var dbP db.Post

		err := dec.Decode(&dbP)
		if err != nil {
			l.log.Infof("decoding: %w", err)
		}

		if err := l.store.CreatePost(context.Background(), dbP); err != nil {
			l.log.Infof("create: %w", err)
		}
	})
	return nil
}

func (l Listener) PostUpdated() error {
	l.client.Subscribe("post-updated", func(m *stan.Msg) {
		buf := bytes.NewReader(m.Data)
		dec := gob.NewDecoder(buf)

		var dbP db.Post

		l.log.Infow("POST-UPDATED", "===========================================")
		err := dec.Decode(&dbP)
		if err != nil {
			l.log.Infof("decoding: %w", err)
		}
		l.log.Infow("postupdated", dbP)
		if err := l.store.UpdatePost(context.Background(), dbP); err != nil {
			l.log.Infof("create: %w", err)
		}
	})
	return nil
}

func (l Listener) PostDeleted() error {
	l.client.Subscribe("post-deleted", func(m *stan.Msg) {
		buf := bytes.NewReader(m.Data)
		dec := gob.NewDecoder(buf)

		var postID string

		err := dec.Decode(&postID)
		if err != nil {
			l.log.Infof("decoding: %w", err)
		}

		if err := l.store.DeletePost(context.Background(), postID); err != nil {
			l.log.Infof("create: %w", err)
		}
	})
	return nil
}
