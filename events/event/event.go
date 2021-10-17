package event

import (
	"log"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

const packageTitle = "event: "

type Event struct {
	nc *nats.Conn

	subject string
}

func New(subject string, nc *nats.Conn) *Event {
	return &Event{
		subject: subject,
		nc:      nc,
	}
}

func (e *Event) Run() {
	go e.run()
}

func (e *Event) run() {
	const funcTitle = packageTitle + "*Event.run"

	ticker := time.NewTicker(time.Millisecond * 100) // nolint:gomnd
	defer ticker.Stop()

	var i int

	for {
		i++

		msg := e.subject + "." + strconv.Itoa(i)
		if err := e.nc.Publish(e.subject, []byte(msg)); err != nil {
			log.Fatalf("%s subject %s message %s", funcTitle, e.subject, msg) // nolint:gocritic
		}

		<-ticker.C
	}
}
