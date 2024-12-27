package dbmodel

import (
	"github.com/specterops/bloodhound/src/database/types"
)

type AssetGroupSelector struct {
	AssetGroupID   int32  `gorm:"UNIQUE_INDEX:compositeindex"`
	Name           string `gorm:"UNIQUE_INDEX:compositeindex"`
	Selector       string
	SystemSelector bool

	Serial
}

type AssetGroupSelectors []AssetGroupSelector

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

type AssetGroups []AssetGroup

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

type AssetGroupCollections []AssetGroupCollection

type AssetGroupCollectionEntry struct {
	AssetGroupCollectionID int64
	ObjectID               string
	NodeLabel              string
	Properties             types.JSONUntypedObject

	BigSerial
}

type AssetGroupCollectionEntries []AssetGroupCollectionEntry

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
