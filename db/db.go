package db

import "time"

type DB struct {
}

type Item struct {
	expireIn time.Duration
	data     interface{}
}
