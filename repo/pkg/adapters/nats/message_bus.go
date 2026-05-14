package nats

import (
	"context"
	"fmt"

	"github.com/kubercloud/ani/pkg/ports"
	natsgo "github.com/nats-io/nats.go"
)

type MessageBus struct {
	js natsgo.JetStreamContext
}

var _ ports.MessageBus = (*MessageBus)(nil)

func NewMessageBus(js natsgo.JetStreamContext) *MessageBus {
	return &MessageBus{js: js}
}

func (b *MessageBus) Publish(ctx context.Context, event ports.EventEnvelope, opts ports.PublishOptions) error {
	if opts.Subject == "" {
		return fmt.Errorf("message bus publish: subject required")
	}
	pubOpts := []natsgo.PubOpt{natsgo.Context(ctx)}
	if opts.Key != "" {
		pubOpts = append(pubOpts, natsgo.MsgId(opts.Key))
	}
	_, err := b.js.Publish(opts.Subject, event.Payload, pubOpts...)
	if err != nil {
		return fmt.Errorf("message bus publish: %w", err)
	}
	return nil
}

func (b *MessageBus) Subscribe(ctx context.Context, opts ports.SubscribeOptions, handler ports.MessageHandler) (ports.Subscription, error) {
	if opts.Subject == "" {
		return nil, fmt.Errorf("message bus subscribe: subject required")
	}
	subOpts := []natsgo.SubOpt{natsgo.ManualAck()}
	if opts.Consumer != "" {
		subOpts = append(subOpts, natsgo.Durable(opts.Consumer))
	}
	if opts.MaxInflight > 0 {
		subOpts = append(subOpts, natsgo.MaxAckPending(opts.MaxInflight))
	}
	handlerFunc := func(msg *natsgo.Msg) {
		_ = handler(ctx, message{msg: msg})
	}
	var (
		sub *natsgo.Subscription
		err error
	)
	if opts.Queue != "" {
		sub, err = b.js.QueueSubscribe(opts.Subject, opts.Queue, handlerFunc, subOpts...)
	} else {
		sub, err = b.js.Subscribe(opts.Subject, handlerFunc, subOpts...)
	}
	if err != nil {
		return nil, fmt.Errorf("message bus subscribe: %w", err)
	}
	return subscription{sub: sub}, nil
}

type message struct {
	msg *natsgo.Msg
}

func (m message) Subject() string {
	return m.msg.Subject
}

func (m message) Data() []byte {
	return m.msg.Data
}

func (m message) Ack(context.Context) error {
	return m.msg.Ack()
}

func (m message) Nack(context.Context) error {
	return m.msg.Nak()
}

type subscription struct {
	sub *natsgo.Subscription
}

func (s subscription) Drain(context.Context) error {
	return s.sub.Drain()
}
