package filter

type FilterType string

const (
	FilterAnd  FilterType = "and"
	FilterOr   FilterType = "or"
	FilterEq   FilterType = "eq"   // =
	FilterNeq  FilterType = "neq"  // !=
	FilterGt   FilterType = "gt"   // >
	FilterEgt  FilterType = "egt"  // >=
	FilterLt   FilterType = "lt"   // <
	FilterElt  FilterType = "elt"  // <=
	FilterLike FilterType = "like" // like
)
