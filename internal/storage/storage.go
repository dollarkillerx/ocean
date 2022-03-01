package storage

import (
	"errors"
	"sync"

	"github.com/dollarkillerx/ocean/internal/utils"
	"github.com/dollarkillerx/ocean/pkg/filter"
	"github.com/dollarkillerx/ocean/pkg/models"
)

type Storage struct {
	Schema       map[string]*models.Schema            // schema
	ListData     map[string]*[]map[string]interface{} // 具体数据
	GlobalLock   sync.RWMutex                         // 全局锁 用于raft 同步
	SchemaRWLock *utils.RWLock                        // schema 局部锁
	DataRWLock   *utils.RWLock                        // 数据局部锁
}

func New() *Storage {
	return &Storage{
		Schema:       map[string]*models.Schema{},
		ListData:     map[string]*[]map[string]interface{}{},
		SchemaRWLock: utils.NewRWLock(),
		DataRWLock:   utils.NewRWLock(),
	}
}

func (s *Storage) getSchema(index string) *models.Schema {
	s.GlobalLock.Lock()
	defer s.GlobalLock.Unlock()

	return s.Schema[index]
}

func (s *Storage) getListData(index string) *[]map[string]interface{} {
	s.GlobalLock.Lock()
	defer s.GlobalLock.Unlock()

	return s.ListData[index]
}

func (s *Storage) ExIndex(index string) bool {
	schema := s.getSchema(index)
	if schema == nil {
		return false
	}

	return true
}

func (s *Storage) CreateIndex(index string, schema models.Schema) error {
	if len(schema) == 0 {
		return errors.New("invalid schema")
	}

	s.GlobalLock.Lock()
	defer s.GlobalLock.Unlock()

	s.Schema[index] = &schema

	return nil
}

func (s *Storage) DelIndex(index string) error {
	s.GlobalLock.Lock()
	defer s.GlobalLock.Unlock()

	delete(s.Schema, index)

	lock := s.DataRWLock.Lock(index)
	defer lock.Unlock()

	delete(s.ListData, index)
	return nil
}

func (s *Storage) UpdateIndex(index string, schema models.Schema) error {
	if len(schema) == 0 {
		return errors.New("invalid schema")
	}

	s.GlobalLock.Lock()
	defer s.GlobalLock.Unlock()

	s.Schema[index] = &schema

	return nil
}

func (s *Storage) InsertDatas(index string, datas []map[string]interface{}) (count int, err error) {
	lock := s.DataRWLock.Lock(index)
	defer lock.Unlock()

	da := s.getListData(index)
	*da = append(*da, datas...)

	return len(datas), nil
}

func (s *Storage) UpdateData(index string, filterParams filter.Params, update map[string]interface{}) (err error) {
	for _, v := range filterParams.Param {
		switch v.FilterType {
		case filter.FilterOr:
		case filter.FilterLt:
		case filter.FilterBool:
		case filter.FilterAnd:
		case filter.FilterEq:
		case filter.FilterNeq:
		case filter.FilterGt:
		case filter.FilterEgt:
		case filter.FilterElt:
		case filter.FilterLike:
		}
	}
	//TODO implement me
	panic("implement me")
}

func (s *Storage) DelData(index string, filter filter.Params) (count int, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) SearchData(index string, filter filter.Params) ([]interface{}, error) {
	//TODO implement me
	panic("implement me")
}
