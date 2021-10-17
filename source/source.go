// nolint:gomnd
package source

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/KnBrBz/natsrequest/cnst"
	"github.com/KnBrBz/natsrequest/source/request"
)

type Source struct {
	req *request.Request

	act chan *sync.WaitGroup

	id int
}

func New(id int, req *request.Request) *Source {
	return &Source{
		id:  id,
		req: req,
		act: make(chan *sync.WaitGroup, 1),
	}
}

func (src *Source) Run() {
	go src.run()
}

func (src *Source) Act(wg *sync.WaitGroup) {
	src.act <- wg
}

func (src *Source) run() {
	for wg := range src.act {
		src.Publish(wg)
	}
}

func (src *Source) Publish(wg *sync.WaitGroup) {
	defer wg.Done()

	subject := cnst.SourceTitle + strconv.Itoa(src.id)
	timeout := time.Second * 1

	for i := 0; i < 200; i++ {
		start := time.Now()

		answer, err := src.req.Publish(subject, []byte(subject+"."+strconv.Itoa(i)), timeout)
		if err != nil {
			log.Fatal(subject, ".", i, err) // nolint:gocritic
		}

		log.Printf("%s answer `%s` time %s", subject, answer, time.Since(start))
		time.Sleep(time.Millisecond * 500)
	}
}
