package app

import (
	"github.com/specterops/bloodhound/packages/go/apitoy/adapter/appdb"
	"github.com/specterops/bloodhound/packages/go/apitoy/adapter/file"
	"github.com/specterops/bloodhound/src/config"
)

// BHApp is the application object, containing all valid methods of the application
type BHApp struct {
	dbAdapter   appdb.Adapter
	fileAdapter file.Adapter
	cfg         config.Configuration
}

// NewBHApp creates a new BHApp instance with injected dependencies
func NewBHApp(dbAdapter appdb.Adapter, fileAdapter file.Adapter, cfg config.Configuration) BHApp {
	return BHApp{
		dbAdapter:   dbAdapter,
		fileAdapter: fileAdapter,
		cfg:         cfg,
	}
}
