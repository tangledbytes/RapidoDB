package eventbus

import (
	"reflect"
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

func TestNewDataEvent(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want DataEvent
	}{
		{
			"CREATE A DATA EVENT",
			args{"k1", 1234},
			DataEvent{"k1", 1234},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDataEvent(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDataEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventBus_ChannelMultiplexer(t *testing.T) {
	// Subscribe to 3 events with 0 buffered channels
	muxcd := ChannelMultiplexer(Instance, 0, "ev1", "ev2", "ev3")

	// Publish events
	go func() {
		for i := 0; i < 5; i++ {
			Instance.Publish("ev1", DataEvent{"k1", 1})
			time.Sleep(5 * time.Millisecond)
		}
	}()
	go func() {
		for i := 0; i < 2; i++ {
			Instance.Publish("ev2", DataEvent{"k2", 10})
			time.Sleep(10 * time.Millisecond)
		}
	}()
	go func() {
		for i := 0; i < 3; i++ {
			Instance.Publish("ev3", DataEvent{"k3", 10})
			time.Sleep(10 * time.Millisecond)
		}
	}()

	totalPublishes := 10
	totalk1 := 0
	totalk2 := 0
	totalk3 := 0

	for i := 0; i < totalPublishes; i++ {
		select {
		case d := <-muxcd:
			if d.key != "k1" && d.key != "k2" && d.key != "k3" {
				t.Errorf("Expected key k1, k2 or k3 got key %s", d.key)
			}
			if d.value != 1 && d.value != 10 {
				t.Errorf("Expected value 1 or 10 got value %v", d.value)
			}

			switch d.key {
			case "k1":
				totalk1++
			case "k2":
				totalk2++
			case "k3":
				totalk3++
			}
		}
	}

	if totalk1 != 5 {
		t.Errorf("Expected k1 to be %v got %v", 5, totalk1)
	}
	if totalk2 != 2 {
		t.Errorf("Expected k2 to be %v got %v", 2, totalk2)
	}
	if totalk3 != 3 {
		t.Errorf("Expected k3 to be %v got %v", 3, totalk3)
	}
}
