package storage

import (
	"github.com/dollarkillerx/ocean/pkg/filter"
	"github.com/dollarkillerx/ocean/pkg/models"
)

type Interface interface {

	// index scheam
	ExIndex(index string) error
	CreateIndex(index string, schema models.Schema) error
	DelIndex(index string) error
	UpdateIndex(index string, schema models.Schema) error

	// CURD

	// InsertDatas ..
	InsertDatas(index string, datas []interface{}) (count int, err error)
	// UpdateData ..
	UpdateData(index string, filter filter.Params, update map[string]interface{}) (err error)
	// DelData ..
	DelData(index string, filter filter.Params) (count int, err error)
	// SearchData ..
	SearchData(index string, filter filter.Params) ([]interface{}, error)

	// raft

}
