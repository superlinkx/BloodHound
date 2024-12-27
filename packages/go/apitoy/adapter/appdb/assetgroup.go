package appdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/specterops/bloodhound/packages/go/apitoy/adapter/appdb/dbmodel"
	"github.com/specterops/bloodhound/packages/go/apitoy/model"
	"github.com/specterops/bloodhound/src/database"
	"gorm.io/gorm"
)

func (s Adapter) GetAllAssetGroups(ctx context.Context, sort model.Sort, filters model.Filters) (model.AssetGroups, error) {
	var sqlSort = dbmodel.BuildSQLSort(sort)

	if sqlFilter, err := dbmodel.BuildSQLFilter(filters); err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrInvalidFilter, err)
	} else if assetGroups, err := getAllAssetGroups(ctx, s.db, sqlSort, sqlFilter); errors.Is(err, database.ErrNotFound) {
		return assetGroups, fmt.Errorf("%w: %v", model.ErrNotFound, err)
	} else if err != nil {
		return assetGroups, fmt.Errorf("%w: %v", model.ErrGenericDatabaseFailure, err)
	} else {
		return assetGroups, nil
	}
}

func getAllAssetGroups(ctx context.Context, db *gorm.DB, order string, filter dbmodel.SQLFilter) (model.AssetGroups, error) {
	var (
		assetGroups dbmodel.AssetGroups
		result      = preload(db, dbmodel.AssetGroupAssociations()).WithContext(ctx)
	)

	if order != "" {
		result = result.Order(order)
	}

	if filter.SQLString != "" {
		result = result.Where(filter.SQLString, filter.Params...)
	}

	if result = result.Find(&assetGroups); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return assetGroups, model.ErrNotFound
	} else if result.Error != nil {
		return assetGroups, fmt.Errorf("%w: %v", model.ErrGenericDatabaseFailure, result.Error)
	}

	for idx := range assetGroups {
		if latestCollection, collectionErr := getLatestAssetGroupCollection(ctx, db, assetGroups[idx].ID); errors.Is(collectionErr, model.ErrNotFound) {
			assetGroups[idx].MemberCount = 0
		} else if collectionErr != nil {
			return assetGroups, fmt.Errorf("%w: get latest collection for asset group %s: %v", model.ErrGenericDatabaseFailure, assetGroups[idx].Name, collectionErr)
		} else {
			assetGroups[idx].MemberCount = len(latestCollection.Entries)
		}
	}
	return assetGroups, nil
}

func getLatestAssetGroupCollection(ctx context.Context, db *gorm.DB, assetGroupID int32) (model.AssetGroupCollection, error) {
	var (
		latestCollection dbmodel.AssetGroupCollection
		result           = preload(db, dbmodel.AssetGroupCollectionAssociations()).
					WithContext(ctx).
					Where("asset_group_id = ?", assetGroupID).
					Order("created_at DESC").
					First(&latestCollection)
	)

	return latestCollection, checkError(result)
}
