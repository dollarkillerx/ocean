package query

import (
	"encoding/json"
	"errors"

	"github.com/dollarkillerx/ocean/pkg/filter"
)

func ParseQuery(input []byte) (params *filter.Params, err error) {
	var dsl filter.QueryDSL
	err = json.Unmarshal(input, &dsl)
	if err != nil {
		return nil, err
	}

	var result filter.Params

	boolMap, ex := dsl.Query["bool"]
	if ex {
		m, ok := boolMap.(map[string]interface{})
		if !ok {
			return nil, errors.New("parse error")
		}
		andList, andEx := m[filter.FilterAnd.String()]
		orList, orEx := m[filter.FilterOr.String()]
		if !andEx && !orEx {
			return nil, errors.New("parse error")
		}

		if andList != nil {
			aList, ok := andList.([]interface{})
			if !ok {
				return nil, errors.New("parse error")
			}
			for _, v := range aList {
				v, ok := v.(map[string]interface{})
				if !ok {
					return nil, errors.New("parse error")
				}

				param, err := parseParam(v)
				if err != nil {
					return nil, err
				}

				result.Param = append(result.Param, *param)
			}
		}
		if orList != nil {
			oList, ok := orList.([]interface{})
			if !ok {
				return nil, errors.New("parse error")
			}
			for _, v := range oList {
				v, ok := v.(map[string]interface{})
				if !ok {
					return nil, errors.New("parse error")
				}

				param, err := parseParam(v)
				if err != nil {
					return nil, err
				}

				result.Param = append(result.Param, *param)
			}
		}
	}

	result.From = dsl.From
	result.Size = dsl.Size
	for k, v := range dsl.Sort {
		result.Sort = append(result.Sort, filter.FilterSort{
			Key:      k,
			SortType: v,
		})
	}

	return &result, nil
}

func parseParam(r map[string]interface{}) (*filter.Param, error) {
	if len(r) != 1 {
		return nil, errors.New("parse error")
	}

	for k, v := range r {
		switch k {
		case filter.FilterEq.String():
		case filter.FilterGt.String():
		case filter.FilterLike.String():
		case filter.FilterEgt.String():
		case filter.FilterLt.String():
		case filter.FilterNeq.String():
		case filter.FilterElt.String():
		case filter.FilterBool.String():
			mc, ok := v.(map[string]interface{})
			if !ok {
				return nil, errors.New("parse error")
			}

			andList, andEx := mc[filter.FilterAnd.String()]
			orList, orEx := mc[filter.FilterOr.String()]
			if !andEx && !orEx {
				return nil, errors.New("parse error")
			}

			var pa = filter.Param{
				FilterType: filter.FilterOr,
				Params:     []filter.Param{},
			}

			var aPa = filter.Param{
				FilterType: filter.FilterAnd,
			}
			var oPa = filter.Param{
				FilterType: filter.FilterOr,
			}

			if andList != nil {
				aList, ok := andList.([]interface{})
				if !ok {
					return nil, errors.New("parse error")
				}

				for _, v2 := range aList {
					vk, ok := v2.(map[string]interface{})
					if !ok {
						return nil, errors.New("parse error")
					}

					param, err := parseParam(vk)
					if err != nil {
						return nil, err
					}
					aPa.Params = append(aPa.Params, *param)
				}
			}
			if orList != nil {
				oList, ok := orList.([]interface{})
				if !ok {
					return nil, errors.New("parse error")
				}

				for _, v2 := range oList {
					vk, ok := v2.(map[string]interface{})
					if !ok {
						return nil, errors.New("parse error")
					}

					param, err := parseParam(vk)
					if err != nil {
						return nil, err
					}
					oPa.Params = append(oPa.Params, *param)
				}
			}

			if len(aPa.Params) != 0 {
				pa.Params = append(pa.Params, aPa)
			}
			if len(oPa.Params) != 0 {
				pa.Params = append(pa.Params, oPa)
			}
			return &pa, nil
		default:
			return nil, errors.New("parse error")
		}

		m, ok := v.(map[string]interface{})
		if !ok {
			return nil, errors.New("parse error")
		}
		if len(m) != 1 {
			return nil, errors.New("parse error")
		}

		for k1, v1 := range m {
			return &filter.Param{
				FilterType: filter.FilterType(k),
				Key:        k1,
				Value:      v1,
			}, nil
		}
	}

	return nil, errors.New("parse error")
}
