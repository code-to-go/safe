package main

/*
typedef struct Result{
    char* res;
	char* err;
} Result;

typedef struct App {
	void (*feed)(char* name, char* data, int eof);
} App;

#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"

	"github.com/code-to-go/safepool.lib/api"
	"github.com/code-to-go/safepool.lib/pool"
)

func cResult(v any, err error) C.Result {
	var res []byte

	if err != nil {
		return C.Result{nil, C.CString(err.Error())}
	}
	if v == nil {
		return C.Result{nil, nil}
	}

	res, err = json.Marshal(v)
	if err == nil {
		return C.Result{C.CString(string(res)), nil}
	}
	return C.Result{nil, C.CString(err.Error())}
}

func cInput(err error, i *C.char, v any) error {
	if err != nil {
		return err
	}
	data := C.GoString(i)
	return json.Unmarshal([]byte(data), v)
}

//export start
func start(dbPath *C.char) C.Result {
	p := C.GoString(dbPath)
	return cResult(nil, api.Start(p))
}

//export stop
func stop() C.Result {
	return cResult(nil, nil)
}

//export getSelfId
func getSelfId() C.Result {
	return cResult(api.Self.Id(), nil)
}

//export getSelf
func getSelf() C.Result {
	return cResult(api.Self, nil)
}

//export getPoolList
func getPoolList() C.Result {
	return cResult(pool.List(), nil)
}

//export createPool
func createPool(config *C.char, apps *C.char) C.Result {
	var c pool.Config
	var apps_ []string

	err := cInput(nil, config, &c)
	err = cInput(err, apps, &apps_)
	if err != nil {
		return cResult(nil, err)
	}

	err = api.CreatePool(c, apps_)
	return cResult(nil, err)
}

//export addPool
func addPool(token *C.char) C.Result {
	err := api.AddPool(C.GoString(token))

	return cResult(nil, err)
}
