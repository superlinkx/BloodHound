package appdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/specterops/bloodhound/packages/go/apitoy/adapter/appdb/dbmodel"
	appModel "github.com/specterops/bloodhound/packages/go/apitoy/model"
	"github.com/specterops/bloodhound/src/database"
	"github.com/specterops/bloodhound/src/model"
	"gorm.io/gorm"
)

func (s Adapter) GetAllAssetGroups(ctx context.Context, sort appModel.Sort, filters appModel.Filters) (model.AssetGroups, error) {
	var sqlSort = dbmodel.BuildSQLSort(sort)

	if sqlFilter, err := dbmodel.BuildSQLFilter(filters); err != nil {
		return nil, fmt.Errorf("%w: %v", appModel.ErrInvalidFilter, err)
	} else if assetGroups, err := getAllAssetGroups(ctx, s.db, sqlSort, sqlFilter); errors.Is(err, database.ErrNotFound) {
		return assetGroups, fmt.Errorf("%w: %v", appModel.ErrNotFound, err)
	} else if err != nil {
		return assetGroups, fmt.Errorf("%w: %v", appModel.ErrGenericDatabaseFailure, err)
	} else {
		return assetGroups, nil
	}
}

func getAllAssetGroups(ctx context.Context, db *gorm.DB, order string, filter dbmodel.SQLFilter) (model.AssetGroups, error) {
	var (
		assetGroups model.AssetGroups
		result      = preload(db, model.AssetGroupAssociations()).WithContext(ctx)
	)

	if order != "" {
		result = result.Order(order)
	}

	if filter.SQLString != "" {
		result = result.Where(filter.SQLString, filter.Params...)
	}

	if result = result.Find(&assetGroups); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return assetGroups, appModel.ErrNotFound
	} else if result.Error != nil {
		return assetGroups, fmt.Errorf("%w: %v", appModel.ErrGenericDatabaseFailure, result.Error)
	}

	for idx := range assetGroups {
		if latestCollection, collectionErr := getLatestAssetGroupCollection(ctx, db, assetGroups[idx].ID); errors.Is(collectionErr, appModel.ErrNotFound) {
			assetGroups[idx].MemberCount = 0
		} else if collectionErr != nil {
			return assetGroups, fmt.Errorf("%w: get latest collection for asset group %s: %v", appModel.ErrGenericDatabaseFailure, assetGroups[idx].Name, collectionErr)
		} else {
			assetGroups[idx].MemberCount = len(latestCollection.Entries)
		}
	}
	return assetGroups, nil
}

func getLatestAssetGroupCollection(ctx context.Context, db *gorm.DB, assetGroupID int32) (model.AssetGroupCollection, error) {
	var (
		latestCollection model.AssetGroupCollection
		result           = preload(db, model.AssetGroupCollectionAssociations()).
					WithContext(ctx).
					Where("asset_group_id = ?", assetGroupID).
					Order("created_at DESC").
					First(&latestCollection)
	)

	return latestCollection, checkError(result)
}
