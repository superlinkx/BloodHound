package app

import (
	"fmt"
	"strings"

	appModel "github.com/specterops/bloodhound/packages/go/apitoy/model"
	"github.com/specterops/bloodhound/src/config"
	"github.com/specterops/bloodhound/src/database"
	"github.com/specterops/bloodhound/src/model"
)

// BHApp is the application object, containing all valid methods of the application
type BHApp struct {
	// adapter adapter.PostgresAdapter
	db  database.Database
	cfg config.Configuration
}

// NewBHApp creates a new BHApp instance with injected dependencies
func NewBHApp(db database.Database, cfg config.Configuration) BHApp {
	return BHApp{
		db:  db,
		cfg: cfg,
	}
}

func BuildSQLFilter(filters appModel.Filters) (model.SQLFilter, error) {
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
			case appModel.GreaterThan:
				predicate = appModel.GreaterThanSymbol
			case appModel.GreaterThanOrEquals:
				predicate = appModel.GreaterThanOrEqualsSymbol
			case appModel.LessThan:
				predicate = appModel.LessThanSymbol
			case appModel.LessThanOrEquals:
				predicate = appModel.LessThanOrEqualsSymbol
			case appModel.Equals:
				predicate = appModel.EqualsSymbol
			case appModel.NotEquals:
				predicate = appModel.NotEqualsSymbol
			case appModel.ApproximatelyEquals:
				predicate = appModel.ApproximatelyEqualSymbol
				filter.Value = fmt.Sprintf("%%%s%%", filter.Value)
			default:
				return model.SQLFilter{}, fmt.Errorf("invalid filter predicate specified")
			}

			result.WriteString(name)
			result.WriteString(" ")
			result.WriteString(predicate)
			result.WriteString(" ?")

			params = append(params, filter.Value)
			firstFilter = false
		}
	}

	return model.SQLFilter{SQLString: result.String(), Params: params}, nil
}

func BuildSQLSort(sort appModel.Sort) []string {
	var sqlSort = make([]string, 0, len(sort))
	for _, sortItem := range sort {
		var column string
		if sortItem.Direction == appModel.DescendingSortDirection {
			column = sortItem.Column + " desc"
		} else {
			column = sortItem.Column
		}

		sqlSort = append(sqlSort, column)
	}
	return sqlSort
}
