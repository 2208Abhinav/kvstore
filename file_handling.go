package kvstore

import (
	"os"
	"strconv"
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
