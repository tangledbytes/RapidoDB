package store

import (
	"strconv"
	"testing"
	"time"
)

func BenchmarkStore_Set(b *testing.B) {
	// Create a store
	var store = New(-1)
	for i := 0; i < b.N; i++ {
		store.Set("key"+strconv.Itoa(i), i, NeverExpire)
	}
}

func BenchmarkStore_Get(b *testing.B) {
	// Create a store
	var store = New(-1)
	for i := 0; i < b.N; i++ {
		store.Get("key" + strconv.Itoa(i))
	}
}

func TestStore(t *testing.T) {
	ts := New(NeverExpire)

	i1, ok := ts.Get("k1")

	if ok || i1 != nil {
		t.Error("Found value for k1 even though it shouldn't exist", i1)
	}

	// Add k2 to store now
	ts.Set("k2", 12345, ts.DefaultExpiry())

	i2, ok := ts.Get("k2")

	if !ok || i2 == nil {
		t.Error("Didn't found any value for k2 even though it was added")
	}

	// Delete k1 from store which doesn't exist
	i3, ok := ts.Delete("k1")

	if ok || i3 != nil {
		t.Error("Deleted a key even though it shouldn't exist in the store", i3)
	}

	// Delete k2 from the store which does exists in the store
	i4, ok := ts.Delete("k2")

	if !ok || i4 == nil {
		t.Error("Didn't delete k2 even though it is present in the store")
	}

	// Find k2 in the store
	i5, ok := ts.Get("k2")

	if ok || i5 != nil {
		t.Error("k2 shouldn't exist in the store once the item is deleted", i5)
	}

	// Add multiple keys
	ts.Set("k1", 123, ts.DefaultExpiry())
	ts.Set("k2", 1234, ts.DefaultExpiry())
	ts.Set("k3", "Hello World", ts.DefaultExpiry())
	ts.Set("k4", 345.0983, ts.DefaultExpiry())

	//  Wipe the entire store
	ts.Wipe()

	// Check if any of the above values exists
	for i := 0; i < 5; i++ {
		key := "k" + strconv.Itoa(i+1)

		item, ok := ts.Get(key)

		if ok || item != nil {
			t.Error("Key shouldn't exist after wiping the store", key)
		}
	}
}

func TestStoreJanitor(t *testing.T) {
	// Create a store without using the new method
	// to pass in a custom janitor interval
	store := &Store{
		defaultExpiry: NeverExpire,
		data:          make(map[string]Item),
		janitor:       newJanitor(1 * time.Millisecond),
	}

	// Setup the janitor
	setupJanitor(store)

	// Add items to the store
	store.Set("k1", 123, 10*time.Millisecond)
	store.Set("k2", "Hello World", 1*time.Millisecond)

	// Check if the items exists in the store without
	// using the Get method defined on the store as
	// that method will never return an item that has expired
	// even if that item exists in the store
	store.RLock()
	v1, ok := store.data["k1"]
	store.RUnlock()

	if !ok {
		t.Error("Expected key to exist", v1)
	}

	// Sleep for 5 millisecond
	time.Sleep(5 * time.Millisecond)

	store.RLock()
	v2, ok := store.data["k2"]
	store.RUnlock()

	if ok {
		t.Error("Item exists in the store even after expiring", v2)
	}
}
