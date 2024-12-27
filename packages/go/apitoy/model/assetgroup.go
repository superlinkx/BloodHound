package model

type AssetGroupSelector struct {
	AssetGroupID   int32
	Name           string
	Selector       string
	SystemSelector bool

	Serial
}

type AssetGroupSelectors []AssetGroupSelector

type AssetGroup struct {
	Name        string
	Tag         string
	SystemGroup bool
	Selectors   AssetGroupSelectors
	Collections AssetGroupCollections
	MemberCount int

	Serial
}

type AssetGroups []AssetGroup

type AssetGroupCollection struct {
	AssetGroupID int32
	Entries      AssetGroupCollectionEntries

	BigSerial
}

type AssetGroupCollections []AssetGroupCollection

type AssetGroupCollectionEntry struct {
	AssetGroupCollectionID int64
	ObjectID               string
	NodeLabel              string
	Properties             map[string]any

	BigSerial
}

type AssetGroupCollectionEntries []AssetGroupCollectionEntry
