package request

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

const packageTitle = "request: "

type Request struct {
	nc *nats.Conn
}

func New(nc *nats.Conn) *Request {
	return &Request{
		nc: nc,
	}
}

func (req *Request) Publish(subject string, data []byte, timeout time.Duration) (answer []byte, err error) {
	const funcTitle = packageTitle + "*Request.Publish"

	inbox := nats.NewInbox()

	sub, err := req.nc.SubscribeSync(inbox)
	if err != nil {
		err = errors.Wrap(err, funcTitle)
		return
	}

	defer func() {
		if errUnsub := sub.Unsubscribe(); errUnsub != nil {
			log.Printf("[ERR] %v\n", errors.Wrapf(errUnsub, "%s: unsubscribe", funcTitle))
		}
	}()

	if err = req.nc.PublishRequest(subject, inbox, data); err != nil {
		err = errors.Wrap(err, funcTitle)
		return
	}

	if err = req.nc.FlushTimeout(timeout); err != nil {
		err = errors.Wrap(err, funcTitle)
		return
	}

	msg, err := sub.NextMsg(timeout)
	if err != nil {
		err = errors.Wrap(err, funcTitle)
		return
	}

	answer = msg.Data

	return
}
