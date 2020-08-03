package eventbus

import (
	"testing"
	"time"
)

func TestEventBus(t *testing.T) {
	// Channel 1 for getting info about event 1
	ch1 := make(chan DataEvent, 1)

	// Channel 2 for getting info about event 2
	ch2 := make(chan DataEvent, 1)

	// Subscribe to the events
	Instance.Subscribe("event1", ch1)
	Instance.Subscribe("event2", ch2)

	// Publish events
	go func() {
		for i := 0; i < 5; i++ {
			Instance.Publish("event1", 1)
			time.Sleep(5 * time.Millisecond)
		}
	}()
	go func() {
		for i := 0; i < 2; i++ {
			Instance.Publish("event2", 1)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	totalPublishes := 7

	for i := 0; i < totalPublishes; i++ {
		select {
		case d := <-ch1:
			if d.Topic != "event1" {
				t.Errorf("Expected event %s got event %s", "event1", d.Topic)
			}
		case d := <-ch2:
			if d.Topic != "event2" {
				t.Errorf("Expected event %s got event %s", "event2", d.Topic)
			}
		}
	}
}
