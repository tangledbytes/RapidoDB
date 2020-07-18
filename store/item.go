package store

import (
	"time"
)

// Item struct encapsualtes how an item
// is represented inside the store
type Item struct {
	// expireAt stores the unix timestamp in NANOSECONDS
	// of the time when a data should be expired
	expireAt int64

	// data can be anything of any type
	data interface{}
}

// newItem returns a new item that can be stored in the database
// newItem takes in the data to be stored as its first parameter which
// can be of any type
//
// newItem takes in the expire duration as it's second arguments which
// must be of type time.Duration
func newItem(data interface{}, expireIn time.Duration) Item {
	// Here int64 is important as NeverExpire is of type int
	// and hence can be casted to int64 and UnixNano return type
	// is int64
	var expiry int64 = NeverExpire

	if expireIn != NeverExpire {
		expiry = time.Now().Add(expireIn).UnixNano()
	}

	return Item{expiry, data}
}

// isExpired returns true if an item is expired
func (item Item) isExpired() bool {
	if item.expireAt == NeverExpire {
		return false
	}

	return item.expireAt < time.Now().UnixNano()
}
