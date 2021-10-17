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

	nec *nats.EncodedConn

	subs []*nats.Subscription
}

func New(nec *nats.EncodedConn) *Answer {
	return &Answer{
		subs:   make([]*nats.Subscription, 0),
		nec:    nec,
		events: make(chan []byte, 1000),
	}
}

func (ans *Answer) Events() <-chan []byte {
	return ans.events
}

func (ans *Answer) SubscribeEvents(subjects []string) {
	const funcTitle = packageTitle + ""

	conn := ans.nec

	for _, subject := range subjects {
		sub, err := conn.BindRecvChan(subject, ans.events)
		if err != nil {
			log.Fatal(errors.Wrapf(err, "%s subject", funcTitle))
		}

		ans.subs = append(ans.subs, sub)
	}
}

func (ans *Answer) Subscribe(request string, eventHandler func(subj, reply string, msg []byte) interface{}) error {
	const funcTitle = packageTitle + "*Answer.Subscribe"

	conn := ans.nec

	sub, err := conn.Subscribe(request, func(subj, reply string, msg []byte) {
		// log.Printf("Incoming Request: subj `%s` reply `%s` msg `%s`\n", subj, reply, msg)
		if v := eventHandler(subj, reply, msg); v != nil {
			if pubErr := conn.Publish(reply, v); pubErr != nil {
				log.Fatal(errors.Wrap(pubErr, funcTitle))
			}
		}
		// log.Printf("Complete Request: subj `%s` reply `%s`\n", subj, reply)
	})
	if err != nil {
		return errors.Wrap(err, funcTitle)
	}

	ans.subs = append(ans.subs, sub)

	return nil
}

func (ans *Answer) Publish(reply string, data []byte) {
	const funcTitle = packageTitle + "*Answer.Publish"

	conn := ans.nec
	if pubErr := conn.Publish(reply, data); pubErr != nil {
		log.Fatalf("%s reply %s data %s", funcTitle, reply, data)
	}
}
