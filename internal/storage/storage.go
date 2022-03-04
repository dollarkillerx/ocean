package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/dollarkillerx/ocean/internal/utils"
	"github.com/dollarkillerx/ocean/pkg/enum"
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

	schema["ocean_id"] = enum.SchemaString
	s.Schema[index] = &schema
	_, ex := s.ListData[index]
	if !ex {
		s.ListData[index] = &[]map[string]interface{}{}
	}

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
	schema := s.getSchema(index)
	if schema == nil {
		return 0, fmt.Errorf("not found index: %s", index)
	}

	lock := s.DataRWLock.Lock(index)
	defer lock.Unlock()

	da := s.getListData(index)
	if da == nil {
		return 0, fmt.Errorf("not found index: %s", index)
	}

	for _, v := range datas {
		v["ocean_id"] = utils.GenerateID()
		for k1, v2 := range v {
			schemaType, ok := (*schema)[k1]
			if !ok {
				return 0, fmt.Errorf("illegal field: %s", k1)
			}
			switch schemaType {
			case enum.SchemaInt64, enum.SchemaTimestamp:
				{
					i, ok := v2.(int)
					if ok {
						v[k1] = int64(i)
					}
				}
				{
					i, ok := v2.(int32)
					if ok {
						v[k1] = int64(i)
					}
				}
			case enum.SchemaFloat64:
				{
					i, ok := v2.(float32)
					if ok {
						v[k1] = float64(i)
					}
				}
				{
					i, ok := v2.(int64)
					if ok {
						v[k1] = float64(i)
					}
				}
				{
					i, ok := v2.(int)
					if ok {
						v[k1] = float64(i)
					}
				}
				{
					i, ok := v2.(int32)
					if ok {
						v[k1] = float64(i)
					}
				}
			}
		}
		*da = append(*da, v)
	}

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

func (s *Storage) searchData(index string, filterParams filter.Params) (result []interface{}, err error) {
	defer func() {
		if er := recover(); er != nil {
			log.Println(er)
			err = errors.New("internal filter sql key is nil")
			return
		}
	}()

	marshal, err := json.Marshal(filterParams)
	fmt.Println(string(marshal))

	schema := s.getSchema(index)
	if schema == nil {
		return nil, fmt.Errorf("not found: %s", index)
	}

	lock := s.DataRWLock.Lock(index)
	defer lock.Unlock()

	da := s.getListData(index)
	if da == nil {
		return nil, fmt.Errorf("not found index: %s", index)
	}

	if len(filterParams.Param) == 0 {
		for _, v := range *da {
			result = append(result, v)
		}

		return result, nil
	}

	vrs := map[string]andStruct{}
	for _, v := range filterParams.Param {
		switch v.FilterType {
		case filter.FilterOr, filter.FilterAnd, filter.FilterLt, filter.FilterEq, filter.FilterNeq, filter.FilterGt, filter.FilterEgt, filter.FilterElt, filter.FilterLike:
			data, err := s.searchBaseData(v, *schema, *da)
			if err != nil {
				return nil, err
			}
			for iv := range data {
				m, ok := data[iv].(map[string]interface{})
				if !ok {
					continue
				}
				pOceanID, ok := m["ocean_id"]
				if !ok {
					continue
				}
				oceanID, ok := pOceanID.(string)
				if !ok {
					continue
				}
				pdata, ex := vrs[oceanID]
				if !ex {
					pdata = andStruct{
						Data:  data[iv],
						Count: 0,
					}

					vrs[oceanID] = pdata
				}
				pdata.Count += 1

				vrs[oceanID] = pdata
			}
		default:
			return nil, fmt.Errorf("v1 illegal parameter: %s", v.FilterType)
		}
	}

	switch filterParams.FilterType {
	case filter.FilterAnd:
		for _, vb := range vrs {
			if vb.Count == len(filterParams.Param) {
				result = append(result, vb.Data)
			}
		}
	case filter.FilterOr:
		for _, vb := range vrs {
			result = append(result, vb.Data)
		}
	default:
		return nil, fmt.Errorf("v2 illegal parameter: %s", filterParams.FilterType)
	}

	// order by

	// limit offset

	return result, nil
}

func (s *Storage) searchBaseData(fil filter.Param, schema models.Schema, da []map[string]interface{}) (result []interface{}, err error) {
	defer func() {
		if er := recover(); er != nil {
			log.Println(er)
			err = errors.New("internal filter sql key is nil")
			return
		}
	}()

	if fil.Key == "" && !(fil.FilterType == filter.FilterAnd || fil.FilterType == filter.FilterOr) {
		return nil, errors.New("filter sql key is nil")
	}
	if fil.Value == nil && !(fil.FilterType == filter.FilterEq || fil.FilterType == filter.FilterNeq || fil.FilterType == filter.FilterAnd || fil.FilterType == filter.FilterOr) {
		return nil, fmt.Errorf("filter sql val is nil KEY: %s", fil.Key)
	}

	switch fil.FilterType {
	case filter.FilterAnd:
		vrs := map[string]andStruct{}
		for vi := range fil.Params {
			data, err := s.searchBaseData(fil.Params[vi], schema, da)
			if err != nil {
				return nil, err
			}

			for iv := range data {
				m, ok := data[iv].(map[string]interface{})
				if !ok {
					continue
				}
				pOceanID, ok := m["ocean_id"]
				if !ok {
					continue
				}
				oceanID, ok := pOceanID.(string)
				if !ok {
					continue
				}
				dt, ex := vrs[oceanID]
				if !ex {
					dt = andStruct{
						Data:  data[iv],
						Count: 0,
					}
				}
				dt.Count += 1
				vrs[oceanID] = dt
			}
		}
		// 求和
		for _, vb := range vrs {
			if vb.Count == len(fil.Params) {
				result = append(result, vb.Data)
			}
		}

		return result, nil
	case filter.FilterOr:
		vrs := map[string]andStruct{}

		for vi := range fil.Params {
			data, err := s.searchBaseData(fil.Params[vi], schema, da)
			if err != nil {
				return nil, err
			}

			for iv := range data {
				m, ok := data[iv].(map[string]interface{})
				if !ok {
					continue
				}
				pOceanID, ok := m["ocean_id"]
				if !ok {
					continue
				}
				oceanID, ok := pOceanID.(string)
				dt, ex := vrs[oceanID]
				if !ex {
					dt = andStruct{
						Data:  data[iv],
						Count: 0,
					}
				}
				dt.Count += 1
				vrs[oceanID] = dt
			}
		}

		for _, vb := range vrs {
			result = append(result, vb.Data)
		}

		return result, nil
	}

	for _, v := range da {
		schemaType, ex := schema[fil.Key]
		if !ex {
			return nil, fmt.Errorf("nonexistent field: %s", fil.Key)
		}

		kVal, ex := v[fil.Key]

		switch fil.FilterType {
		case filter.FilterEq:
			if !ex {
				if fil.Value == nil {
					result = append(result, v)
				}
			}

			switch schemaType {
			case enum.SchemaString:
				i, ok := kVal.(string)
				i2, i2Ok := fil.Value.(string)
				if ok && i2Ok {
					if i == i2 {
						result = append(result, v)
					}
				}
			case enum.SchemaFloat64:
				i, ok := kVal.(float64)
				i2, i2Ok := pFloat64(fil.Value)
				if ok && i2Ok {
					if i == i2 {
						result = append(result, v)
					}
				}
			case enum.SchemaInt64, enum.SchemaTimestamp:
				{
					i, ok := kVal.(int64)
					i2, i2Ok := pInt64(fil.Value)
					if ok && i2Ok {
						if i == i2 {
							result = append(result, v)
						}
					}
				}
			case enum.SchemaBool:
				i, ok := kVal.(bool)
				i2, i2Ok := fil.Value.(bool)
				if ok && i2Ok {
					if i == i2 {
						result = append(result, v)
					}
				}
			}
		case filter.FilterNeq:
			if !ex {
				continue
			}
			if fil.Value == nil {
				result = append(result, v)
			}

			switch schemaType {
			case enum.SchemaString:
				i, ok := kVal.(string)
				i2, i2Ok := fil.Value.(string)
				if ok && i2Ok {
					if i != i2 {
						result = append(result, v)
					}
				}
			case enum.SchemaFloat64:
				i, ok := kVal.(float64)
				i2, i2Ok := pFloat64(fil.Value)
				if ok && i2Ok {
					if i != i2 {
						result = append(result, v)
					}
				}
			case enum.SchemaInt64, enum.SchemaTimestamp:
				{
					i, ok := kVal.(int64)
					i2, i2Ok := pInt64(fil.Value)
					if ok && i2Ok {
						if i != i2 {
							result = append(result, v)
						}
					}
				}
			case enum.SchemaBool:
				i, ok := kVal.(bool)
				i2, i2Ok := fil.Value.(bool)
				if ok && i2Ok {
					if i != i2 {
						result = append(result, v)
					}
				}
			}
		case filter.FilterLt:
			if !ex {
				continue
			}

			switch schemaType {
			case enum.SchemaInt64, enum.SchemaTimestamp:
				{
					i, ok := kVal.(int64)
					i2, i2Ok := pInt64(fil.Value)
					if ok && i2Ok {
						if i < i2 {
							result = append(result, v)
						}
					}
				}
			case enum.SchemaFloat64:
				i, ok := kVal.(float64)
				i2, i2Ok := pFloat64(fil.Value)
				if ok && i2Ok {
					if i < i2 {
						result = append(result, v)
					}
				}
			default:
				return nil, fmt.Errorf("wrong type: %s", fil.Key)
			}
		case filter.FilterGt:
			if !ex {
				continue
			}
			switch schemaType {
			case enum.SchemaInt64, enum.SchemaTimestamp:
				{
					i, ok := kVal.(int64)
					i2, i2Ok := pInt64(fil.Value)
					if ok && i2Ok {
						if i > i2 {
							result = append(result, v)
						}
					}
				}
			case enum.SchemaFloat64:
				i, ok := kVal.(float64)
				i2, i2Ok := pFloat64(fil.Value)
				if ok && i2Ok {
					if i > i2 {
						result = append(result, v)
					}
				}
			default:
				return nil, fmt.Errorf("wrong type: %s", fil.Key)
			}
		case filter.FilterEgt:
			if !ex {
				continue
			}

			switch schemaType {
			case enum.SchemaInt64, enum.SchemaTimestamp:
				{
					i, ok := kVal.(int64)
					i2, i2Ok := pInt64(fil.Value)
					if ok && i2Ok {
						if i >= i2 {
							result = append(result, v)
						}
					}
				}
			case enum.SchemaFloat64:
				i, ok := kVal.(float64)
				i2, i2Ok := pFloat64(fil.Value)
				if ok && i2Ok {
					if i >= i2 {
						result = append(result, v)
					}
				}
			default:
				return nil, fmt.Errorf("wrong type: %s", fil.Key)
			}
		case filter.FilterElt:
			if !ex {
				continue
			}

			switch schemaType {
			case enum.SchemaInt64, enum.SchemaTimestamp:
				{
					i, ok := kVal.(int64)
					i2, i2Ok := pInt64(fil.Value)
					if ok && i2Ok {
						if i <= i2 {
							result = append(result, v)
						}
					}
				}
			case enum.SchemaFloat64:
				i, ok := kVal.(float64)
				i2, i2Ok := pFloat64(fil.Value)
				if ok && i2Ok {
					if i <= i2 {
						result = append(result, v)
					}
				}
			default:
				return nil, fmt.Errorf("wrong type: %s", fil.Key)
			}
		case filter.FilterLike:
			if !ex {
				continue
			}

			switch schemaType {
			case enum.SchemaString:
				i, ok := kVal.(string)
				i2, i2Ok := fil.Value.(string)
				if ok && i2Ok {
					if strings.Contains(i, i2) {
						result = append(result, v)
					}
				}
			default:
				return nil, fmt.Errorf("wrong type: %s", fil.Key)
			}
		default:
			return nil, fmt.Errorf("v3 illegal parameter: %s", fil.FilterType)
		}
	}

	return result, nil
}

func pInt64(i interface{}) (int64, bool) {
	switch v := i.(type) {
	case int64:
		return v, true
	case int32:
		return int64(v), true
	case int:
		return int64(v), true
	}

	return 0, false
}

func pFloat64(i interface{}) (float64, bool) {
	switch v := i.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	case int:
		return float64(v), true
	}

	return 0, false
}

type andStruct struct {
	Data  interface{}
	Count int
}
