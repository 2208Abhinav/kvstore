package kvstore

import (
	"os"
	"testing"
)

var store *Store

func TestInit(t *testing.T) {
	var err error
	store, err = Init("test.store")
	if err != nil {
		t.Errorf("%v\n", err)
	}
}

func TestCreate(t *testing.T) {
	if store == nil {
		t.Errorf("store = nil was provided")
		return
	}

	keyValues := []*KeyValue{
		&KeyValue{
			Key: "abc",
			Value: map[string]interface{}{
				"a": "sasas",
				"b": map[string]interface{}{
					"abhinav": "singh",
				},
			},
			Time: 1000, // valid for 1000 seconds
		},
		&KeyValue{
			Key: "fgdfgsjewhdbdjsgdjhekfhdjkdshfhgfjkdldhdfkghghg",
			Value: map[string]interface{}{
				"asa": 100,
				"vdd": map[string]interface{}{
					"igdfg": "assaas",
				},
			},
			Time: 2000, // valid for 2000 seconds
		},
		&KeyValue{
			Key: "jds",
			Value: map[string]interface{}{
				"ghkgrf": 100000,
				"fgihjf": map[string]interface{}{
					"fsh": "feoapufe",
				},
			},
			Time: 0, // valid forever
		},
		&KeyValue{
			Key: "abc",
			Value: map[string]interface{}{
				"sasfa": "sasasassaas",
				"b": map[string]interface{}{
					"asdfaf": "gfhytrhjrhr",
				},
			},
			Time: 5,
		},
		&KeyValue{
			Key: "",
			Value: map[string]interface{}{
				"sasfa": "sasasassaas",
				"b": map[string]interface{}{
					"asdfaf": "gfhytrhjrhr",
				},
			},
			Time: 5,
		},
	}

	for i, keyValue := range keyValues {
		err := Create(store, keyValue.Key, keyValue.Value, keyValue.Time)

		if i == 0 || i == 2 {
			if err != nil {
				t.Errorf("key: %s, error: %v\n", keyValue.Key, err)
				return
			}
		} else if i == 1 {
			if err == nil || err.Error() != "key size cannot exceed 32 characters" {
				t.Errorf("key: %s, error: %v\n", keyValue.Key, err)
				return
			}
		} else if i == 3 {
			if err == nil || err.Error() != "key already present" {
				t.Errorf("key: %s, error: %v\n", keyValue.Key, err)
				return
			}
		} else if i == 4 {
			if err == nil || err.Error() != "key cannot be empty" {
				t.Errorf("key: %s, error: %v\n", keyValue.Key, err)
				return
			}
		}
	}
}

func TestRead(t *testing.T) {
	if store == nil {
		t.Errorf("store = nil was provided")
		return
	}

	var key string
	var err error

	key = "abc"
	if _, err = Read(store, key); err != nil {
		t.Errorf("key: %s, error: %v\n", key, err)
		return
	}

	key = "zcege"
	if _, err = Read(store, key); err == nil || err.Error() != "key not found" {
		t.Errorf("key: %s, error: %v", key, err)
	}
}

func TestDelete(t *testing.T) {
	if store == nil {
		t.Errorf("store = nil was provided")
		return
	}

	var key string
	var err error

	key = "zfhegf"
	if err = Delete(store, key); err == nil || err.Error() != "key not present" {
		t.Errorf("key: %s, error: %v\n", key, err)
		return
	}

	key = "abc"
	if err = Delete(store, key); err != nil {
		t.Errorf("key: %s, error: %v", key, err)
	}
}

func TestClose(t *testing.T) {
	if store == nil {
		t.Errorf("store = nil was provided")
		return
	}

	err := Close(store)
	if err != nil {
		t.Errorf("%v\n", err)
		return
	}

	// remove temporary test store file
	os.Remove("test.store")
}
