package dbmodel

import (
	"fmt"
	"strings"

	"github.com/specterops/bloodhound/packages/go/apitoy/model"
)

type SQLFilter struct {
	SQLString string
	Params    []any
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
