package kvstore

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
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

// Init is called by the client to initialize the key value store
func Init(storePath string) (*Store, error) {
	var storeFile *os.File
	var err error

	if _, err = os.Stat(storePath); err == nil {
		// store file already exists
		storeFile, err = os.OpenFile(storePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			return nil, err
		}
	} else if os.IsNotExist(err) {
		// store file doesn't exist
		if len(storePath) == 0 {
			// if no path is provided then by default the store file is created in
			// the parent directory of the project that is using this library with
			// file name as the current unix timestamp
			storePath = fmt.Sprintf("%v.store", time.Now().Unix())
		}
		storeFile, err = os.OpenFile(storePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			return nil, err
		}

		// this acts as a flag to prevent other client from using this key value store if it
		// is already in use. If this value is 0 then the store is not in use otherwise the store
		// is in use if the value is 1
		_, err = storeFile.WriteString("0")
		if err != nil {
			return nil, err
		}
	} else {
		// some unknown error occurred
		return nil, err
	}

	flag, err := getFlag(storeFile)
	if err != nil {
		return nil, err
	} else if flag == 1 {
		return nil, errors.New("store already in use by some other client")
	}

	err = toggleFlag(storeFile)
	if err != nil {
		return nil, err
	}

	store := &Store{StoreFile: storeFile}
	if store.StoreMap, err = readStoreFile(store.StoreFile); err != nil {
		return nil, err
	}

	return store, nil
}

// Close is called by client to close the key value store
func Close(store *Store) error {
	if store.deletesCount > 0 {
		if err := updateStoreFile(store, store.StoreMap); err != nil {
			return err
		}
	}

	if err := toggleFlag(store.StoreFile); err != nil {
		return err
	}

	return store.StoreFile.Close()
}
