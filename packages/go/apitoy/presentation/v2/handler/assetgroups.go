package handler

import (
	"net/http"

	"github.com/specterops/bloodhound/packages/go/apitoy/presentation/common"
	"github.com/specterops/bloodhound/packages/go/apitoy/presentation/v2/view"
	"github.com/specterops/bloodhound/src/api"
)

func (s Handler) ListAssetGroups(response http.ResponseWriter, request *http.Request) {
	var (
		sortByColumns = request.URL.Query()[api.QueryParameterSortBy]
	)

	if sort, err := view.ValidSort(sortByColumns, view.AssetGroupsSortableColumns()); err != nil {
		common.WriteErrorResponse(request.Context(), request, response, err)
	} else if filters, err := view.GetValidFiltersFromQuery(request.URL.Query(), view.AssetGroupsFilters()); err != nil {
		common.WriteErrorResponse(request.Context(), request, response, err)
	} else if assetGroups, err := s.bhApp.GetAllAssetGroups(request.Context(), filters, sort); err != nil {
		common.WriteErrorResponse(request.Context(), request, response, err)
	} else {
		api.WriteBasicResponse(request.Context(), view.GenerateListAssetGroupsResponse(assetGroups), http.StatusOK, response)
	}
}
