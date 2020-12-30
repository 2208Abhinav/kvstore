package kvstore

import (
	"os"
	"sync"
)

// mutex lock for making the functions thread safe
// this will be used where we are manipulating the store file to prevent the
// corruption of file due to manipulation by many threads at the same time
var lock sync.Mutex

const (
	_KB              = 1024
	_MB              = 1024 * 1024
	_GB              = 1024 * 1024 * 1024
	_DeleteThreshold = 256
)

// KeyValue is the entity stored in the store file
/*
	Format of key value pair in the store file:

	"key": {
		"value": JSONObject (marshalled into string),
		"time": seconds set by user (0 means not set and the key value pair will never expire),
		"validTill": unix timestamp till which the key value pair is valid
	}
*/
type KeyValue struct {
	Key       string
	Value     map[string]interface{}
	Time      int64
	ValidTill int64
}

// Store will contain the file of the store, map of store and
// count of deleted key value pairs
type Store struct {
	StoreFile    *os.File
	StoreMap     *map[string]*KeyValue
	deletesCount int
}
