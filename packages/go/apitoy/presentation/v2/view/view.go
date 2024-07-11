package view

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"time"

	"github.com/gofrs/uuid"
	appModel "github.com/specterops/bloodhound/packages/go/apitoy/model"
	"github.com/specterops/bloodhound/src/model"
)

var filterRegex = regexp.MustCompile(`([~\w]+):([\w\--_ ]+)`)

// Basic is a struct which includes the following basic fields: CreatedAt, UpdatedAt, DeletedAt.
type Basic struct {
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

// Unique is a struct is a struct which includes the following basic fields: ID, CreatedAt, UpdatedAt, DeletedAt.
type Unique struct {
	ID uuid.UUID `json:"id"`

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
	ID int32 `json:"id"`

	Basic
}

// BigSerial is a struct that follows the same design principles as Serial but with one exception:
// the ID type is set to int64 to support an ID sequence limit of up to 9223372036854775807.
type BigSerial struct {
	ID int64 `json:"id"`

	Basic
}

// PaginatedResponse has been DEPRECATED as part of V1. Please use api.ResponseWrapper instead
type PaginatedResponse struct {
	Count int `json:"count"`
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
	Data  any `json:"data"`
}

func ValidSort(requestColumns []string, sortableColumns []string) (appModel.Sort, error) {
	var sort = make(appModel.Sort, 0, len(requestColumns))

	for _, column := range requestColumns {
		var sortItem appModel.SortItem
		if string(column[0]) == "-" {
			sortItem.Direction = appModel.DescendingSortDirection
			sortItem.Column = column[1:]
		} else {
			sortItem.Direction = appModel.AscendingSortDirection
			sortItem.Column = column
		}

		if !slices.Contains(sortableColumns, sortItem.Column) {
			return sort, fmt.Errorf("%w: %s", ErrNotSortableOnColumn, sortItem.Column)
		}

		sort = append(sort, sortItem)
	}

	return sort, nil
}

func GetValidFiltersFromQuery(queryParams map[string][]string, validFilters appModel.ValidFilters) (appModel.Filters, error) {
	filters := make(appModel.Filters)

	for name, values := range queryParams {
		// ignore pagination query params
		if slices.Contains(AllPaginationQueryParameters(), name) {
			continue
		}

		if slices.Contains(IgnoreFilters(), name) {
			continue
		}

		if validPredicates, ok := validFilters[name]; !ok {
			return filters, errors.New("invalid filter")
		} else {
			for _, value := range values {
				if subgroups := filterRegex.FindStringSubmatch(value); len(subgroups) > 0 {
					if filterPredicate, err := appModel.ParseFilterOperator(subgroups[1]); err != nil {
						return filters, err
					} else if !slices.Contains(validPredicates, filterPredicate) {
						return filters, errors.New("invalid filter predicate")
					} else {
						if _, ok := filters[name]; !ok {
							filters[name] = make([]appModel.Filter, 0, 4)
						}

						filters[name] = append(filters[name], appModel.Filter{
							Operator: filterPredicate,
							Value:    subgroups[2],
						})
					}
				}
			}
		}
	}

	return filters, nil
}

var ErrNotFiltered = errors.New("parameter value is not filtered")

const (
	PaginationQueryParameterBefore = "before"
	PaginationQueryParameterAfter  = "after"
	PaginationQueryParameterLimit  = "limit"
	PaginationQueryParameterOffset = "offset"
	PaginationQueryParameterSkip   = "skip"
	PaginationQueryParameterSortBy = "sort_by"
)

func AllPaginationQueryParameters() []string {
	return []string{
		PaginationQueryParameterAfter,
		PaginationQueryParameterLimit,
		PaginationQueryParameterBefore,
		PaginationQueryParameterOffset,
		PaginationQueryParameterSkip,
		PaginationQueryParameterSortBy}
}

func IgnoreFilters() []string {
	return []string{
		"scope",
	}
}

func convertToSerial(serial model.Serial) Serial {
	return Serial{
		ID:    serial.ID,
		Basic: Basic(serial.Basic),
	}
}

func convertToBigSerial(bigSerial model.BigSerial) BigSerial {
	return BigSerial{
		ID:    bigSerial.ID,
		Basic: Basic(bigSerial.Basic),
	}
}
