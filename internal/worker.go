package internal

import (
	"log"
	"runtime"

	"github.com/Jeffail/tunny"
)

type Worker struct {
	SubscribePool *tunny.Pool
}

type SubscribeWork struct {
	Event       []byte
	EventSource *EventSource
}

func NewWorker() Worker {
	var worker Worker

	numCPU := runtime.NumCPU()

	worker.SubscribePool = tunny.NewFunc(numCPU, func(work any) any {
		data, ok := work.(*SubscribeWork)
		if !ok {
			log.Println("Worker: invalid work input")

			return nil
		}

		topic := data.EventSource.Topic
		eventData := data.Event
		eventSource := data.EventSource
		event := NewEvent(topic, eventData)

		for i, subscriber := range eventSource.Subscribers {
			err := WriteData(event, subscriber)
			eventSource.Metrics.IncSuccess()
			if err != nil {
				log.Printf("err while sending event to client: %s", err.Error())
				eventSource.Subscribers = append(eventSource.Subscribers[:i], eventSource.Subscribers[i+1:]...)
			}
		}

		return nil
	})

	return worker
}

func NewSubscribeWork(event []byte, eventSource *EventSource) *SubscribeWork {
	return &SubscribeWork{Event: event, EventSource: eventSource}
}
