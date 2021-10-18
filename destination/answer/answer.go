// nolint:gomnd
package answer

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

const packageTitle = "answer: "

type Answer struct {
	events chan []byte

	nc *nats.Conn

	subs []*nats.Subscription
}

func New(nc *nats.Conn) *Answer {
	return &Answer{
		subs:   make([]*nats.Subscription, 0),
		nc:     nc,
		events: make(chan []byte, 1000),
	}
}

func (ans *Answer) Events() <-chan []byte {
	return ans.events
}

func (ans *Answer) SubscribeEvents(subjects []string) {
	const funcTitle = packageTitle + ""

	conn := ans.nc

	for _, subject := range subjects {
		sub, err := conn.Subscribe(subject, func(m *nats.Msg) {
			ans.events <- m.Data
		})
		if err != nil {
			log.Fatal(errors.Wrapf(err, "%s subject", funcTitle))
		}

		ans.subs = append(ans.subs, sub)
	}
}

func (ans *Answer) Subscribe(request string, eventHandler func(subj, reply string, msg []byte) []byte) error {
	const funcTitle = packageTitle + "*Answer.Subscribe"

	conn := ans.nc

	sub, err := conn.Subscribe(request, func(m *nats.Msg) {
		if v := eventHandler(m.Subject, m.Reply, m.Data); v != nil {
			if pubErr := m.Respond(v); pubErr != nil {
				log.Fatal(errors.Wrap(pubErr, funcTitle))
			}
		}
	})

	if err != nil {
		return errors.Wrap(err, funcTitle)
	}

	ans.subs = append(ans.subs, sub)

	return nil
}

func (ans *Answer) Publish(reply string, data []byte) {
	const funcTitle = packageTitle + "*Answer.Publish"

	conn := ans.nc
	if pubErr := conn.Publish(reply, data); pubErr != nil {
		log.Fatalf("%s reply %s data %s", funcTitle, reply, data)
	}
}
