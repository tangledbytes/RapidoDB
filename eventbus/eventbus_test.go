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
			Instance.Publish("event1", DataEvent{"k1", 1})
			time.Sleep(5 * time.Millisecond)
		}
	}()
	go func() {
		for i := 0; i < 2; i++ {
			Instance.Publish("event2", DataEvent{"k2", 10})
			time.Sleep(10 * time.Millisecond)
		}
	}()

	totalPublishes := 7

	for i := 0; i < totalPublishes; i++ {
		select {
		case d := <-ch1:
			if d.key != "k1" {
				t.Errorf("Expected key %s got key %s", "k1", d.key)
			}
			if d.value != 1 {
				t.Errorf("Expected value %v got value %v", 1, d.value)
			}
		case d := <-ch2:
			if d.key != "k2" {
				t.Errorf("Expected key %s got key %s", "k2", d.key)
			}
			if d.value != 10 {
				t.Errorf("Expected value %v got value %v", 10, d.value)
			}
		}
	}
}
