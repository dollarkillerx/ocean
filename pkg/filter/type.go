package filter

type FilterType string

const (
	FilterBool FilterType = "bool"
	FilterAnd  FilterType = "must"
	FilterOr   FilterType = "should"
	FilterEq   FilterType = "eq"   // =
	FilterNeq  FilterType = "neq"  // !=
	FilterGt   FilterType = "gt"   // >
	FilterEgt  FilterType = "egt"  // >=
	FilterLt   FilterType = "lt"   // <
	FilterElt  FilterType = "elt"  // <=
	FilterLike FilterType = "like" // like
)

func (f FilterType) String() string {
	return string(f)
}

type SortType string

const (
	SortDesc SortType = "desc"
	SortAsc  SortType = "asc"
)
