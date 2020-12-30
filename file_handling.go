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
