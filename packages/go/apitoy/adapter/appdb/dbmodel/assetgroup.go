package dbmodel

import (
	"github.com/specterops/bloodhound/packages/go/apitoy/model"
	"github.com/specterops/bloodhound/src/database/types"
)

type AssetGroupSelector struct {
	AssetGroupID   int32  `gorm:"UNIQUE_INDEX:compositeindex"`
	Name           string `gorm:"UNIQUE_INDEX:compositeindex"`
	Selector       string
	SystemSelector bool

	Serial
}

func (s AssetGroupSelector) ConvertToAppModel() model.AssetGroupSelector {
	return model.AssetGroupSelector{
		AssetGroupID:   s.AssetGroupID,
		Name:           s.Name,
		Selector:       s.Selector,
		SystemSelector: s.SystemSelector,
		Serial:         s.Serial.ConvertToAppModel(),
	}
}

type AssetGroupSelectors []AssetGroupSelector

func (s AssetGroupSelectors) ConvertToAppModel() model.AssetGroupSelectors {
	var converted = make(model.AssetGroupSelectors, 0, len(s))

	for _, assetGroupSelector := range s {
		converted = append(converted, assetGroupSelector.ConvertToAppModel())
	}

	return converted
}

// AssetGroupAssociations returns a list of AssetGroup model associations to load eagerly by default with GORM
// Preload(...). Note: this does not include the "Collections" association on-purpose since this collection grows
// over time and may require additional parameters for fetching.
func AssetGroupAssociations() []string {
	return []string{
		"Selectors",
	}
}

type AssetGroup struct {
	Name        string
	Tag         string
	SystemGroup bool
	Selectors   AssetGroupSelectors   `gorm:"constraint:OnDelete:CASCADE;"`
	Collections AssetGroupCollections `gorm:"constraint:OnDelete:CASCADE;"`
	MemberCount int                   `gorm:"-"`

	Serial
}

func (s AssetGroup) ConvertToAppModel() model.AssetGroup {
	return model.AssetGroup{
		Name:        s.Name,
		Tag:         s.Tag,
		SystemGroup: s.SystemGroup,
		Selectors:   s.Selectors.ConvertToAppModel(),
		Collections: s.Collections.ConvertToAppModel(),
		MemberCount: s.MemberCount,
		Serial:      s.Serial.ConvertToAppModel(),
	}
}

type AssetGroups []AssetGroup

func (s AssetGroups) ConvertToAppModel() model.AssetGroups {
	var converted = make(model.AssetGroups, 0, len(s))

	for _, assetGroup := range s {
		converted = append(converted, assetGroup.ConvertToAppModel())
	}

	return converted
}

// AssetGroupCollectionAssociations returns a list of AssetGroupCollection model associations to eagerly by default
// with GORM Preload(...).
func AssetGroupCollectionAssociations() []string {
	return []string{"Entries"}
}

type AssetGroupCollection struct {
	AssetGroupID int32
	Entries      AssetGroupCollectionEntries `gorm:"constraint:OnDelete:CASCADE;"`

	BigSerial
}

func (s AssetGroupCollection) ConvertToAppModel() model.AssetGroupCollection {
	return model.AssetGroupCollection{
		AssetGroupID: s.AssetGroupID,
		Entries:      s.Entries.ConvertToAppModel(),
		BigSerial:    s.BigSerial.ConvertToAppModel(),
	}
}

type AssetGroupCollections []AssetGroupCollection

func (s AssetGroupCollections) ConvertToAppModel() model.AssetGroupCollections {
	var converted = make(model.AssetGroupCollections, 0, len(s))

	for _, assetGroupCollection := range s {
		converted = append(converted, assetGroupCollection.ConvertToAppModel())
	}

	return converted
}

type AssetGroupCollectionEntry struct {
	AssetGroupCollectionID int64
	ObjectID               string
	NodeLabel              string
	Properties             types.JSONUntypedObject

	BigSerial
}

func (s AssetGroupCollectionEntry) ConvertToAppModel() model.AssetGroupCollectionEntry {
	return model.AssetGroupCollectionEntry{
		AssetGroupCollectionID: s.AssetGroupCollectionID,
		ObjectID:               s.ObjectID,
		NodeLabel:              s.NodeLabel,
		Properties:             s.Properties,
		BigSerial:              s.BigSerial.ConvertToAppModel(),
	}
}

type AssetGroupCollectionEntries []AssetGroupCollectionEntry

func (s AssetGroupCollectionEntries) ConvertToAppModel() model.AssetGroupCollectionEntries {
	var converted = make(model.AssetGroupCollectionEntries, 0, len(s))

	for _, assetGroupCollectionEntry := range s {
		converted = append(converted, assetGroupCollectionEntry.ConvertToAppModel())
	}

	return converted
}

type AssetGroupSelectorSpec struct {
	SelectorName   string
	EntityObjectID string
	Action         string
}

type UpdatedAssetGroupSelectors struct {
	Added   AssetGroupSelectors
	Removed AssetGroupSelectors
}

const (
	SelectorSpecActionAdd    = "add"
	SelectorSpecActionRemove = "remove"
	TierZeroAssetGroupName   = "Admin Tier Zero"
	TierZeroAssetGroupTag    = "admin_tier_0"
	OwnedAssetGroupName      = "Owned"
	OwnedAssetGroupTag       = "owned"
)
