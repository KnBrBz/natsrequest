// nolint:gomnd
package main

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/KnBrBz/natsrequest/cnst"
	"github.com/KnBrBz/natsrequest/destination"
	"github.com/KnBrBz/natsrequest/destination/answer"
	"github.com/KnBrBz/natsrequest/events"
	"github.com/KnBrBz/natsrequest/source"
	"github.com/KnBrBz/natsrequest/source/request"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

func natsConn() *nats.Conn {
	url := os.Getenv("NATS_URL")
	user := os.Getenv("NATS_USER")
	password := os.Getenv("NATS_PASSWORD")

	nc, err := nats.Connect(url, nats.UserInfo(user, password))
	if err != nil {
		log.Fatal(err)
	}

	return nc
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	ncSrc := natsConn()
	defer ncSrc.Close()

	ncDst := natsConn()
	defer ncDst.Close()

	necDst, err := nats.NewEncodedConn(ncDst, nats.DEFAULT_ENCODER)
	if err != nil {
		log.Fatal(err) // nolint:gocritic
	}
	defer necDst.Close()

	ncEvent := natsConn()
	defer ncEvent.Close()

	subjects := makeSubjects(100)

	events := events.NewEvents(subjects, ncEvent)
	events.Run()

	var (
		wg       sync.WaitGroup
		sourceID int
	)

	sources := make([]*source.Source, 0, 1000)
	destinations := make([]*destination.Destination, 0, 1000)

	for i := 0; i <= 10; i++ {
		log.Println("Iteration", i, "start")

		for k := 0; k < 200; k++ {
			sourceID++

			destination := destination.New(answer.New(necDst))
			destination.Run()
			destination.Subscribe(cnst.SourceTitle + strconv.Itoa(sourceID))
			destination.SubscribeToEvents(subjects)
			destinations = append(destinations, destination) // nolint:staticcheck

			source := source.New(sourceID, request.New(ncSrc))
			source.Run()
			sources = append(sources, source)
		}

		for _, source := range sources {
			wg.Add(1)
			source.Act(&wg)
		}

		wg.Wait()
		log.Println("Iteration", i, "end")
	}
}

func makeSubjects(subjectCount int) []string {
	subjects := make([]string, subjectCount)
	for i := 0; i < subjectCount; i++ {
		subjects[i] = "event." + strconv.Itoa(i)
	}

	return subjects
}
