package app

import (
	"context"
	"fmt"

	appModel "github.com/specterops/bloodhound/packages/go/apitoy/model"
)

func (s BHApp) GetAllAssetGroups(ctx context.Context, filters appModel.Filters, sort appModel.Sort) (appModel.AssetGroups, error) {
	if assetGroups, err := s.dbAdapter.GetAllAssetGroups(ctx, sort, filters); err != nil {
		return assetGroups, fmt.Errorf("get all asset groups: %w", err)
	} else {
		return assetGroups, nil
	}
}
