package app

import (
	"context"
	"errors"
	"fmt"
	"strings"

	appModel "github.com/specterops/bloodhound/packages/go/apitoy/model"
	"github.com/specterops/bloodhound/src/database"
	"github.com/specterops/bloodhound/src/model"
)

func (s BHApp) GetAllAssetGroups(ctx context.Context, filters appModel.Filters, sort appModel.Sort) (model.AssetGroups, error) {
	var sqlSort = BuildSQLSort(sort)

	if sqlFilter, err := BuildSQLFilter(filters); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGeneralApplication, err)
	} else if assetGroups, err := s.db.GetAllAssetGroups(ctx, strings.Join(sqlSort, ", "), sqlFilter); errors.Is(err, database.ErrNotFound) {
		return assetGroups, fmt.Errorf("%w: %v", ErrNotFound, err)
	} else if err != nil {
		return assetGroups, fmt.Errorf("%w: %v", ErrGenericDatabase, err)
	} else {
		return assetGroups, nil
	}
}
