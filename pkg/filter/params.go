package filter

type Params struct {
	Param []Param `json:"param"`
}

type Param struct {
	FilterType FilterType  `json:"filter_type"` // filter type
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
}
