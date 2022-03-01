package storage

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/dollarkillerx/ocean/internal/utils"
)

func TestSchema(t *testing.T) {
	//var schema map[string]*models.Schema
	//m := schema["asdas"]
	//if m == nil {
	//	fmt.Println("ok")
	//}

	var ListData = map[string]*[]map[string]interface{}{}

	ListData["ppx"] = &[]map[string]interface{}{}

	bs(ListData["ppx"])

	marshal, err := json.Marshal(ListData)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(marshal))
}

func bs(r *[]map[string]interface{}) {
	*r = append(*r, map[string]interface{}{
		"asd":  "sdsd",
		"asdb": "sdsd",
	})
}

func TestBBC(t *testing.T) {
	lock := utils.NewRWLock()

	go func() {
		l := lock.Lock("aaa")
		defer l.Unlock()

		fmt.Println("lock")
		time.Sleep(time.Second)
	}()
	time.Sleep(time.Millisecond * 20)
	rLock := lock.RLock("aaa")
	defer rLock.Unlock()

	fmt.Println("rw unlock")
}
