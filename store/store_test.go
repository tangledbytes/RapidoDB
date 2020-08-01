package store

import (
	"testing"
)

func BenchmarkStore_Set(b *testing.B) {
	// Create a store
	var store = New(-1)
	for i := 0; i < b.N; i++ {
		store.Set("key"+string(i), i, NeverExpire)
	}
}

func BenchmarkStore_Get(b *testing.B) {
	// Create a store
	var store = New(-1)
	for i := 0; i < b.N; i++ {
		store.Get("key" + string(i))
	}
}

func TestStore(t *testing.T) {
	ts := New(NeverExpire)

	i1, ok := ts.Get("k1")

	if ok || i1 != nil {
		t.Error("Found value for k1 even though it shouldn't exist", i1)
	}

	// Add k2 to store now
	ts.Set("k2", 12345, ts.DefaultExpiry)

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
}
