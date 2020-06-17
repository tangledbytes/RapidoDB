package db

import "time"

// Item struct represents the Item being stored in the database
type Item struct {
	expireAt int64
	data     interface{}
}

// NewItem creates a new item and returns it
func NewItem(data interface{}, expireIn time.Duration) Item {
	var expiry int64 = NeverExpire
	if expireIn != NeverExpire {
		expiry = time.Now().Add(expireIn).UnixNano()
	}

	return Item{
		expireAt: expiry,
		data:     data,
	}
}

// IsExpired checks if a item is expired
func (item Item) IsExpired() bool {
	if item.expireAt == NeverExpire {
		return false
	}

	return item.expireAt < time.Now().Unix()
}
