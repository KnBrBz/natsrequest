// nolint:gomnd
package destination

import (
	"log"
	"time"

	"github.com/KnBrBz/natsrequest/destination/answer"
	"github.com/pkg/errors"
)

const packageTitle = "destination: "

type Destination struct {
	ans *answer.Answer
	msg chan Msg
}

type Msg struct {
	reply string
	data  []byte
}

func New(ans *answer.Answer) *Destination {
	return &Destination{
		ans: ans,
		msg: make(chan Msg, 1),
	}
}

func (dst *Destination) Run() {
	go dst.run()
	go dst.runEvents()
}

func (dst *Destination) run() {
	for msg := range dst.msg {
		time.Sleep(time.Millisecond * 100)
		dst.ans.Publish(msg.reply, append(msg.data, []byte(" response")...))
	}
}

func (dst *Destination) runEvents() {
	events := dst.ans.Events()
	for event := range events {
		processEvent(event)
	}
}

func processEvent(event []byte) {
	// log.Printf("%s", event)
}

func (dst *Destination) SubscribeToEvents(subjects []string) {
	dst.ans.SubscribeEvents(subjects)
}

func (dst *Destination) Subscribe(request string) {
	const funcTitle = packageTitle + "*Destination.Subscribe"

	eventHandler := func(subj, reply string, msg []byte) interface{} {
		dst.msg <- Msg{
			reply: reply,
			data:  msg,
		}

		return nil
	}

	if err := dst.ans.Subscribe(request, eventHandler); err != nil {
		log.Fatal(errors.Wrap(err, funcTitle))
	}
}
