package app

import (
	"context"
	"fmt"

	appModel "github.com/specterops/bloodhound/packages/go/apitoy/model"
	"github.com/specterops/bloodhound/src/model"
)

func (s BHApp) GetAllAssetGroups(ctx context.Context, filters appModel.Filters, sort appModel.Sort) (model.AssetGroups, error) {
	if assetGroups, err := s.dbAdapter.GetAllAssetGroups(ctx, sort, filters); err != nil {
		return assetGroups, fmt.Errorf("get all asset groups: %w", err)
	} else {
		return assetGroups, nil
	}
}
