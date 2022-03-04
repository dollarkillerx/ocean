package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/dollarkillerx/ocean/internal/utils"
	"github.com/dollarkillerx/ocean/pkg/enum"
	"github.com/dollarkillerx/ocean/pkg/filter"
	"github.com/dollarkillerx/ocean/pkg/models"
)

func TestInsertAndSelect(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	storage := New()

	idx := "storage_v1"

	// create index

	err := storage.CreateIndex(idx, models.Schema{
		"name":        enum.SchemaString,
		"age":         enum.SchemaInt64,
		"money":       enum.SchemaFloat64,
		"create_time": enum.SchemaTimestamp,
	})
	if err != nil {
		panic(err)
	}

	// insert

	for i := 0; i < 100; i++ {
		_, err := storage.InsertDatas(idx, []map[string]interface{}{
			{
				"name":        fmt.Sprintf("wamg: %d", i),
				"age":         i,
				"money":       float64(i) + 100 + 0.1,
				"create_time": time.Now().Unix(),
			},
		})
		if err != nil {
			log.Fatalln(err)
		}
		//fmt.Println("ins: ", datas)
	}

	_, err = storage.InsertDatas(idx, []map[string]interface{}{
		{
			"name":        fmt.Sprintf("wamg: %d", 9),
			"age":         80,
			"create_time": time.Now().Unix(),
		},
	})

	data, err := storage.searchData(idx, filter.Params{})
	if err != nil {
		panic(err)
	}

	//log.Println(len(data))
	//log.Println(data)

	data, err = storage.searchData(idx, filter.Params{
		FilterType: filter.FilterAnd,
		Param: []filter.Param{
			{
				FilterType: filter.FilterGt,
				Key:        "age",
				Value:      30,
			},
			{
				FilterType: filter.FilterLike,
				Key:        "name",
				Value:      "wamg",
			},
			{
				FilterType: filter.FilterAnd,
				Params: []filter.Param{
					{
						FilterType: filter.FilterGt,
						Key:        "age",
						Value:      20,
					},
					{
						FilterType: filter.FilterGt,
						Key:        "money",
						Value:      60,
					},
					{
						FilterType: filter.FilterAnd,
						Params: []filter.Param{
							{
								FilterType: filter.FilterGt,
								Key:        "money",
								Value:      180,
							},
						},
					},
				},
			},
		},

		Sort: []filter.FilterSort{
			{
				Key:      "age",
				SortType: filter.SortDesc,
			},
		},

		From: 0,
		Size: 10,
	})
	if err != nil {
		panic(err)
	}

	log.Println(len(data))
	log.Println(data)
}

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
