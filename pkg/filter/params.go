package filter

type Params struct {
	Param      []Param    `json:"param"`
	FilterType FilterType `json:"filter_type"` // filter type

	// 分页
	From int `json:"from"`
	Size int `json:"size"`

	// 排序
	Sort []FilterSort `json:"sort"`
}

type FilterSort struct {
	Key      string   `json:"key"`
	SortType SortType `json:"sort_type"`
}

type Param struct {
	FilterType FilterType  `json:"filter_type"` // filter type
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`

	Params []Param `json:"params"`
}

type QueryDSL struct {
	// 分页
	From int `json:"from,omitempty"`
	Size int `json:"size,omitempty"`

	// 排序
	Sort map[string]SortType `json:"sort,omitempty"`

	Query map[string]interface{} `json:"query,omitempty"`
}

/**
Bool struct {
	Should []struct {
		Eq struct {
			Key int `json:"key"`
		} `json:"eq"`
	} `json:"should"`
} `json:"bool"`
*/
