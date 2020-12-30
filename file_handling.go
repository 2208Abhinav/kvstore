package kvstore

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// getFlag is used to determine the state of the store file. If it is in use by
// other client or not
// 0 -> not in use by any other client
// 1 -> is in use by some other client
func getFlag(storeFile *os.File) (int, error) {
	var err error
	flag := make([]byte, 1)
	_, err = storeFile.ReadAt(flag, 0)

	if err != nil {
		return -1, err
	}

	return strconv.Atoi(string(flag))
}

// toggleFlag is used to toggle the state of the store file. The state is used
// for determining whether the file is in use by other client or not
func toggleFlag(storeFile *os.File) error {
	flag, err := getFlag(storeFile)
	if err != nil {
		return err
	}

	if flag == 0 {
		// change the flag to 1 because this file is now used by the current client
		_, err = storeFile.WriteAt([]byte{'1'}, 0)
	} else {
		// change the flag to 0 because this file is not used by any client
		_, err = storeFile.WriteAt([]byte{'0'}, 0)
	}

	return err
}

// readStoreFile will read all the content in given store file and will
// unmarshal the content in the usable JSON Object as represented by
// map[string]interface{} in golang
func readStoreFile(fi *os.File) (*map[string]*KeyValue, error) {
	content := []byte{}
	content = append(content, byte('{'))
	// make a read buffer
	r := bufio.NewReader(fi)
	buf := make([]byte, 1024)

	flag := ""
	for {
		// read a chunk of key value pair
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}

		i := 0
		for i < n {
			if len(flag) == 0 && i == 0 {
				flag = string(buf[i])
			} else {
				content = append(content, buf[i])
			}
			i++
		}

		if n == 0 {
			break
		}
	}
	// remove the extra invalid ","
	if len(content) > 1 {
		content = content[:len(content)-1]
	}
	content = append(content, byte('}'))

	var result map[string]*KeyValue
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, err
	}

	for key, keyValue := range result {
		keyValue.Key = key
	}

	return &result, nil
}

// writeToStoreFile will be used to marshal the given key value pair
// and write to the store file for persistence
func writeToStoreFile(storeFile *os.File, keyValue *KeyValue) error {
	info, err := storeFile.Stat()
	if err != nil {
		return err
	}
	// according to the requirement the store file size cannot exceed 1GB
	if info.Size() > _GB {
		return errors.New("size of store file cannot exceed 1GB")
	}

	// set the unix time stamp till which the key is valid
	// otherwise the 0 timestamp means the key value pair is valid forever
	if keyValue.Time != 0 {
		keyValue.ValidTill = time.Now().Unix() + keyValue.Time
	}

	valueStrBytes, err := json.Marshal(keyValue.Value)
	if err != nil {
		return err
	}

	valueStr := string(valueStrBytes)
	// 1 character takes 1 byte in file
	if len(valueStr) > 16*_KB {
		return errors.New("size of value cannot exceed 16KB")
	}

	_, err = storeFile.WriteString(
		fmt.Sprintf("\"%v\":{\"value\":%v,\"time\":%v,\"validTill\":%v},",
			keyValue.Key, valueStr, keyValue.Time, keyValue.ValidTill))

	return err
}

// updateStoreFile will update the store file to a new state after removing the
// deleted keys. This function works by first removing the old store file. Then
// a new store file is created and the remaning key value pairs are written to
// that new store file. And then the old store file is replaced by the new file
func updateStoreFile(store *Store, storeMap *map[string]*KeyValue) error {
	storeFile := store.StoreFile
	storeFilePath, err := filepath.Abs(storeFile.Name())
	if err != nil {
		return err
	}

	if err = os.Remove(storeFile.Name()); err != nil {
		return err
	}

	newStoreFile, err := os.OpenFile(storeFilePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	if _, err = newStoreFile.WriteAt([]byte{'1'}, 0); err != nil {
		return err
	}

	for _, keyValue := range *storeMap {
		writeToStoreFile(newStoreFile, keyValue)
	}

	store.StoreFile = newStoreFile

	return nil
}
