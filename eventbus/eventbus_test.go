package eventbus

import (
	"testing"
	"time"
)

func TestEventBus(t *testing.T) {
	// Subscribe to the events
	ch1 := Instance.Subscribe("event1", 1)
	ch2 := Instance.Subscribe("event2", 1)

	// Publish events
	go func() {
		for i := 0; i < 5; i++ {
			Instance.Publish("event1", 1)
			time.Sleep(5 * time.Millisecond)
		}
	}()
	go func() {
		for i := 0; i < 2; i++ {
			Instance.Publish("event2", 10)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	totalPublishes := 7

	for i := 0; i < totalPublishes; i++ {
		select {
		case d := <-ch1:
			if d != 1 {
				t.Errorf("Expected value %s got value %s", "event1", d)
			}
		case d := <-ch2:
			if d != 10 {
				t.Errorf("Expected value %s got value %s", "event2", d)
			}
		}
	}
}
