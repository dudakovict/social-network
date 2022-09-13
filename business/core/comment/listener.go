package comment

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/dudakovict/social-network/business/core/comment/db"
	"github.com/dudakovict/social-network/business/sys/nats"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

type Listener struct {
	log   *zap.SugaredLogger
	nats  *nats.NATS
	store db.Store
}

func (l Listener) PostCreated() error {
	l.nats.Subscribe("post-created", "posts", func(m *stan.Msg) {
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
	l.nats.Subscribe("post-updated", "posts", func(m *stan.Msg) {
		buf := bytes.NewReader(m.Data)
		dec := gob.NewDecoder(buf)

		var dbP db.Post

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
	l.nats.Subscribe("post-deleted", "posts", func(m *stan.Msg) {
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
