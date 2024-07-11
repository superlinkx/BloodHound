package model

import "fmt"

type FilterOperator string

const (
	GreaterThan         FilterOperator = "gt"
	GreaterThanOrEquals FilterOperator = "gte"
	LessThan            FilterOperator = "lt"
	LessThanOrEquals    FilterOperator = "lte"
	Equals              FilterOperator = "eq"
	NotEquals           FilterOperator = "neq"
	ApproximatelyEquals FilterOperator = "~eq"

	GreaterThanSymbol         string = ">"
	GreaterThanOrEqualsSymbol string = ">="
	LessThanSymbol            string = "<"
	LessThanOrEqualsSymbol    string = "<="
	EqualsSymbol              string = "="
	NotEqualsSymbol           string = "<>"
	ApproximatelyEqualSymbol  string = "ILIKE"

	TrueString     = "true"
	FalseString    = "false"
	IdString       = "id"
	ObjectIdString = "objectid"
)

func ParseFilterOperator(raw string) (FilterOperator, error) {
	switch FilterOperator(raw) {
	case GreaterThan:
		return GreaterThan, nil

	case GreaterThanOrEquals:
		return GreaterThanOrEquals, nil

	case LessThan:
		return LessThan, nil

	case LessThanOrEquals:
		return LessThanOrEquals, nil

	case Equals:
		return Equals, nil

	case NotEquals:
		return NotEquals, nil

	case ApproximatelyEquals:
		return ApproximatelyEquals, nil

	default:
		return "", fmt.Errorf("unknown query parameter filter predicate: %s", raw)
	}
}

type Filter struct {
	Operator FilterOperator
	Value    string
}

type Filters map[string][]Filter

type ValidFilters map[string][]FilterOperator

type SortDirection int

const (
	InvalidSortDirection SortDirection = iota
	AscendingSortDirection
	DescendingSortDirection
)

type SortItem struct {
	Direction SortDirection
	Column    string
}

type Sort []SortItem
