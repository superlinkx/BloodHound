package model

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

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

// Basic is a struct which includes the following basic fields: CreatedAt, UpdatedAt, DeletedAt.
type Basic struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// Unique is a struct is a struct which includes the following basic fields: ID, CreatedAt, UpdatedAt, DeletedAt.
type Unique struct {
	ID uuid.UUID

	Basic
}

// Serial is a struct which includes the following basic fields: ID, CreatedAt, UpdatedAt, DeletedAt.
// This was chosen over the default gorm model so that ID retains the bare int type. We do this because
// uint has no meaning with regards to the underlying database storage engine - at least where postgresql is
// concerned. To avoid type gnashing and unexpected pain with sql.NullInt32 the bare int type is a better
// choice all around.
//
// See: https://www.postgresql.org/docs/current/datatype-numeric.html
type Serial struct {
	ID int32

	Basic
}

// BigSerial is a struct that follows the same design principles as Serial but with one exception:
// the ID type is set to int64 to support an ID sequence limit of up to 9223372036854775807.
type BigSerial struct {
	ID int64

	Basic
}

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
