package handler

import "github.com/specterops/bloodhound/packages/go/apitoy/app"

// Handler stores dependencies of all handlers (currently just the BHEApp interface)
type Handler struct {
	bhApp app.BHApp
}

// NewHandler initializes a Handlers struct with a BHEApp
func NewHandler(bhApp app.BHApp) Handler {
	return Handler{
		bhApp: bhApp,
	}
}
