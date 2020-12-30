# kvstore: A file based key-value store

kvstore is a file based key-value store that supports basic CRD (create, read and delete) operations. This data store is meant to be used as a local storage for one single process on one laptop.

## Installation
Latest version of GoLang must be installed on the system.
```sh
$ go get -u github.com/2208Abhinav/kvstore
```

## Usage
**Import**
```Go
import "github.com/2208Abhinav/kvstore"
```

**Initialization without specifying store file path**
- If the store file path is not specified then the store file will be created in the parent directory of the project. The name of the store file will be `<current unix timestamp>.store`.
```Go
store, err := kvstore.Init("")
if err != nil {
	panic(err)
}
```

**Initialization with store file path**
- Pass the path of the store as a string and the library will create the store at the specified location.
```Go
store, err := kvstore.Init("/Users/abhinav/Documents/abcd.store")
if err != nil {
	panic(err)
}
```

**Create a new key-value pair**
- A new key-value pair can be added to the data store by using the `Create` method. The size of the key is capped at `32 chars` and the size of JSON Object is capped at `16KB`.
Errors with appropriate messages are returned if the key is already present, if the size of key or value exceeds the defined limit.
- Also there is a check on the size of the store file. If the store file size exceeds `1GB` then error is returned by the library.
- Every key supports setting an optional Time-To-Live property when it is created. If provided, it will be evaluated as an integer defining the number of seconds the key must be retained in the data store. Once the time has expired, the key will no longer be available for Read or Delete operations.
```Go
var keyValue *kvstore.KeyValue

// key-value pair with Time-To-Live property set
keyValue = &kvstore.KeyValue{
	Key: "abc",
	Value: map[string]interface{}{
		"a": "sasas",
		"b": map[string]interface{}{
			"gghhh": "sgjkjgk",
		},
	},
	Time: 2000, // valid for 2000 seconds
}
if err := kvstore.Create(store, keyValue.Key, keyValue.Value, keyValue.Time); err != nil {
    return err
}

// key-value pair with no Time-To-Live property.
keyValue = &kvstore.KeyValue{
	Key: "dsjkhjds",
	Value: map[string]interface{}{
		"jkhjks": "ghhjjs",
		"c": map[string]interface{}{
			"zxds": "fdedkag",
		},
	},
	Time: 0, // will not expire
}
if err := kvstore.Create(store, keyValue.Key, keyValue.Value, keyValue.Time); err != nil {
    return err
}
```

**Read a key-value pair**
- A `Read` operation on a key can be performed by providing the key, and receiving the value in response, as a JSON object (in GoLang the JSON Object is represented by `map[string]interface{}`)
- An error is returned if the key is expired or not present in the store.
- Read returns the key-value pair of type `*kvstore.KeyValue`
```Go
keyValue, err := kvstore.Read(store, "abcd")
if err != nil {
    return err
}
```

**Delete a key-value pair**
- A `Delete` operation can be performed by providing the key.
- An error is returned if the key is expired or not present in the store.
```Go
if err := kvstore.Delete(store, "def"); err != nil {
    return err
}
```

**Close the key-value store**
- Closing the store is important as it ensures that the state of the store file is handled properly and it can be used by other clients also.
```Go
if err := kvstore.Close(store); err != nil {
	return err
}
```
