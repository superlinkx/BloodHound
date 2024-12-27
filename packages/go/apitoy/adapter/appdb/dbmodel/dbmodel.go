package dbmodel

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/specterops/bloodhound/packages/go/apitoy/model"
)

type SQLFilter struct {
	SQLString string
	Params    []any
}

// Basic is a struct which includes the following basic fields: CreatedAt, UpdatedAt, DeletedAt.
type Basic struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"-"`
}

// Unique is a struct is a struct which includes the following basic fields: ID, CreatedAt, UpdatedAt, DeletedAt.
type Unique struct {
	ID uuid.UUID `gorm:"primaryKey"`

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
	ID int32 `gorm:"primaryKey"`

	Basic
}

// BigSerial is a struct that follows the same design principles as Serial but with one exception:
// the ID type is set to int64 to support an ID sequence limit of up to 9223372036854775807.
type BigSerial struct {
	ID int64 `gorm:"primaryKey"`

	Basic
}

func BuildSQLFilter(filters model.Filters) (SQLFilter, error) {
	var (
		result      strings.Builder
		firstFilter = true
		predicate   string
		params      []any
	)

	for name, filterOperations := range filters {
		for _, filter := range filterOperations {
			if !firstFilter {
				result.WriteString(" AND ")
			}

			switch filter.Operator {
			case model.GreaterThan:
				predicate = model.GreaterThanSymbol
			case model.GreaterThanOrEquals:
				predicate = model.GreaterThanOrEqualsSymbol
			case model.LessThan:
				predicate = model.LessThanSymbol
			case model.LessThanOrEquals:
				predicate = model.LessThanOrEqualsSymbol
			case model.Equals:
				predicate = model.EqualsSymbol
			case model.NotEquals:
				predicate = model.NotEqualsSymbol
			case model.ApproximatelyEquals:
				predicate = model.ApproximatelyEqualSymbol
				filter.Value = fmt.Sprintf("%%%s%%", filter.Value)
			default:
				return SQLFilter{}, fmt.Errorf("invalid filter predicate specified")
			}

			result.WriteString(name)
			result.WriteString(" ")
			result.WriteString(predicate)
			result.WriteString(" ?")

			params = append(params, filter.Value)
			firstFilter = false
		}
	}

	return SQLFilter{SQLString: result.String(), Params: params}, nil
}

func BuildSQLSort(sort model.Sort) string {
	var sqlSort = make([]string, 0, len(sort))
	for _, sortItem := range sort {
		var column string
		if sortItem.Direction == model.DescendingSortDirection {
			column = sortItem.Column + " desc"
		} else {
			column = sortItem.Column
		}

		sqlSort = append(sqlSort, column)
	}
	return strings.Join(sqlSort, ",")
}
