package events

import (
	"github.com/KnBrBz/natsrequest/events/event"
	"github.com/nats-io/nats.go"
)

type Events struct {
	subjects []string

	es []*event.Event
}

func NewEvents(subjects []string, nc *nats.Conn) *Events {
	es := make([]*event.Event, 0, len(subjects))
	for _, subject := range subjects {
		es = append(es, event.New(subject, nc))
	}

	return &Events{
		subjects: subjects,
		es:       es,
	}
}

func (es *Events) Run() {
	for _, event := range es.es {
		event.Run()
	}
}
