package query

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestParseQuery(t *testing.T) {
	r := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"eq": map[string]interface{}{
							"key": 16,
						},
					},
					{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{
									"neq": map[string]interface{}{
										"k2": "66",
									},
								},
							},
							"should": []map[string]interface{}{
								{
									"neq": map[string]interface{}{
										"k3": "33",
									},
								},
							},
						},
					},
				},
			},
		},
		"sort": map[string]interface{}{
			"key": "desc",
		},
		"from": 0,
		"size": 10,
	}

	marshal, err := json.Marshal(r)
	if err != nil {
		log.Fatalln(err)
	}

	query, err := ParseQuery(marshal)
	if err != nil {
		panic(err)
	}

	fmt.Println(query)

	bytes, err := json.Marshal(query)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}

func TestParseQuery2(t *testing.T) {
	r := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"eq": map[string]interface{}{
							"key": 16,
						},
					},
					{
						"neq": map[string]interface{}{
							"key": 16,
						},
					},
				},
			},
		},
		"sort": map[string]interface{}{
			"key": "desc",
		},
		"from": 0,
		"size": 10,
	}

	marshal, err := json.Marshal(r)
	if err != nil {
		log.Fatalln(err)
	}

	query, err := ParseQuery(marshal)
	if err != nil {
		panic(err)
	}

	fmt.Println(query)

	bytes, err := json.Marshal(query)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}
