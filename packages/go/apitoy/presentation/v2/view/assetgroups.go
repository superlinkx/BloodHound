package view

import (
	"errors"

	appModel "github.com/specterops/bloodhound/packages/go/apitoy/model"
	"github.com/specterops/bloodhound/slicesext"
	"github.com/specterops/bloodhound/src/model"
)

var ErrNotSortableOnColumn = errors.New("not sortable on column")

type AssetGroupSelector struct {
	AssetGroupID   int32  `json:"asset_group_id"`
	Name           string `json:"name"`
	Selector       string `json:"selector"`
	SystemSelector bool   `json:"system_selector"`

	Serial
}

type AssetGroupSelectors []AssetGroupSelector

type AssetGroup struct {
	Name        string                `json:"name"`
	Tag         string                `json:"tag"`
	SystemGroup bool                  `json:"system_group"`
	Selectors   AssetGroupSelectors   `json:"selectors"`
	Collections AssetGroupCollections `json:"-"`
	MemberCount int                   `json:"member_count"`

	Serial
}

type AssetGroups []AssetGroup

type AssetGroupCollection struct {
	AssetGroupID int32                       `json:"-"`
	Entries      AssetGroupCollectionEntries `json:"entries"`

	BigSerial
}

type AssetGroupCollections []AssetGroupCollection

type AssetGroupCollectionEntry struct {
	AssetGroupCollectionID int64          `json:"-"`
	ObjectID               string         `json:"object_id"`
	NodeLabel              string         `json:"node_label"`
	Properties             map[string]any `json:"properties"`

	BigSerial
}

type AssetGroupCollectionEntries []AssetGroupCollectionEntry

// ListAssetGroupsResponse holds the data returned to a list asset groups request
type ListAssetGroupsResponse struct {
	AssetGroups AssetGroups `json:"asset_groups"`
}

type AssetGroupCollectionsResponse struct {
	Data []any `json:"data"`
}

func AssetGroupsSortableColumns() []string {
	return []string{"name", "tag", "member_count"}
}

func AssetGroupsFilters() appModel.ValidFilters {
	return appModel.ValidFilters{
		"name":         {appModel.Equals, appModel.NotEquals},
		"tag":          {appModel.Equals, appModel.NotEquals},
		"system_group": {appModel.Equals, appModel.NotEquals},
		"member_count": {appModel.Equals, appModel.GreaterThan, appModel.GreaterThanOrEquals, appModel.LessThan, appModel.LessThanOrEquals, appModel.NotEquals},
		"id":           {appModel.Equals, appModel.GreaterThan, appModel.GreaterThanOrEquals, appModel.LessThan, appModel.LessThanOrEquals, appModel.NotEquals},
		"created_at":   {appModel.Equals, appModel.GreaterThan, appModel.GreaterThanOrEquals, appModel.LessThan, appModel.LessThanOrEquals, appModel.NotEquals},
		"updated_at":   {appModel.Equals, appModel.GreaterThan, appModel.GreaterThanOrEquals, appModel.LessThan, appModel.LessThanOrEquals, appModel.NotEquals},
		"deleted_at":   {appModel.Equals, appModel.GreaterThan, appModel.GreaterThanOrEquals, appModel.LessThan, appModel.LessThanOrEquals, appModel.NotEquals},
	}
}

func GenerateListAssetGroupsResponse(assetGroups model.AssetGroups) ListAssetGroupsResponse {
	return ListAssetGroupsResponse{
		AssetGroups: slicesext.Map(assetGroups, convertToAssetGroup),
	}
}

func convertToAssetGroup(assetGroup model.AssetGroup) AssetGroup {
	return AssetGroup{
		Name:        assetGroup.Name,
		Tag:         assetGroup.Tag,
		SystemGroup: assetGroup.SystemGroup,
		Selectors:   slicesext.Map(assetGroup.Selectors, convertToAssetGroupSelector),
		Collections: slicesext.Map(assetGroup.Collections, convertToAssetGroupCollection),
		MemberCount: assetGroup.MemberCount,
		Serial:      convertToSerial(assetGroup.Serial),
	}
}

func convertToAssetGroupSelector(selector model.AssetGroupSelector) AssetGroupSelector {
	return AssetGroupSelector{
		AssetGroupID:   selector.AssetGroupID,
		Name:           selector.Name,
		Selector:       selector.Selector,
		SystemSelector: selector.SystemSelector,
		Serial:         convertToSerial(selector.Serial),
	}
}

func convertToAssetGroupCollection(collection model.AssetGroupCollection) AssetGroupCollection {
	return AssetGroupCollection{
		AssetGroupID: collection.AssetGroupID,
		Entries:      slicesext.Map(collection.Entries, convertToAssetGroupCollectionEntry),
		BigSerial:    convertToBigSerial(collection.BigSerial),
	}
}

func convertToAssetGroupCollectionEntry(entry model.AssetGroupCollectionEntry) AssetGroupCollectionEntry {
	return AssetGroupCollectionEntry{
		AssetGroupCollectionID: entry.AssetGroupCollectionID,
		ObjectID:               entry.ObjectID,
		NodeLabel:              entry.NodeLabel,
		Properties:             entry.Properties,
		BigSerial:              convertToBigSerial(entry.BigSerial),
	}
}
