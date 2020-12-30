package kvstore

import (
	"bufio"
	"encoding/json"
	"io"
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
