//go:build serial_integration
// +build serial_integration

package analysis_test

import (
	"context"
	"github.com/specterops/bloodhound/analysis"
	adAnalysis "github.com/specterops/bloodhound/analysis/ad"
	azureAnalysis "github.com/specterops/bloodhound/analysis/azure"
	"github.com/specterops/bloodhound/dawgs/graph"
	"github.com/specterops/bloodhound/dawgs/query"
	"github.com/specterops/bloodhound/graphschema"
	"github.com/specterops/bloodhound/graphschema/ad"
	"github.com/specterops/bloodhound/graphschema/azure"
	"github.com/specterops/bloodhound/src/test/integration"
	"github.com/stretchr/testify/require"
	"testing"
)

// This is a test to validate when we have a situ such
// There exists an AD user and an Azure user that represent the same principal (the same user identity)
//
// This connection is made by correlating properties that are inserted when data from Active Directory or Azure is ingested
// into the system. These properties are referenced in the function in bhce/packages/go/analysis/hybrid/hybrid.go - hasOnPremUser(...)
// and then mapped to AD users for creation of the SyncedToEntraUser and SyncedToADUser edges.
//
// Hybrid post-processing is driven by https://learn.microsoft.com/en-us/azure/architecture/reference-architectures/identity/azure-ad - current
// limitations of the implementation in MS means that the relationship between User and AZUser is 1:* where a AZUser may only be connected to
// one AD principal.
func TestDeleteTransitEdges(t *testing.T) {
	var (
		// This creates a new live integration test context with the graph database
		// This call will load whatever BHE configuration the environment variable `INTEGRATION_CONFIG_PATH` points to.
		textCtx = integration.NewGraphTestContext(t, graphschema.DefaultGraphSchema())

		// For this test we need to validate BED-4954 - this requires, at minimum, an AD user and an Entra (Azure) user. The lines below
		// will utilize the test context to put the data directly into the graph.

		// AD user first
		adUser = textCtx.NewNode(graph.AsProperties(map[string]any{
			"name":     "ad_user",
			"objectid": "1234",
		}), ad.Entity, ad.User)

		// Azure user second
		azureUser = textCtx.NewNode(graph.AsProperties(map[string]any{
			"name":     "azure_user",
			"objectid": "4321",
		}), azure.Entity, azure.User)
	)

	// In order to validate that DeleteTransitEdges and the updated PostProcessedRelationships for both AD and Azure are correct, we need to simulate
	// the completion of post-processing in: lib/go/analysis/azure/post.go
	//
	// The specific function that is responsible for creating the edges below can be found in bhce/packages/go/analysis/hybrid/hybrid.go - PostHybrid(...)
	//
	// Here, we are choosing to create these edges such that the data describes what we would expect to see after a successful execution of the logic
	// in lib/go/analysis/azure/post.go.
	textCtx.NewRelationship(adUser, azureUser, ad.SyncedToEntraUser)
	textCtx.NewRelationship(azureUser, adUser, azure.SyncedToADUser)

	// The way post-processing operates is that all edges created during post-processing are deleted before each analysis run. This helps keep the graph consistent
	// where certain graph conditions (edges, node properties, etc.) that once existed were removed or modified due to the user's environment changing.

	// This first run removes all Azure post-processed relationships - expected outcome is that SyncedToADUser is removed at this stage
	_, err := analysis.DeleteTransitEdges(context.Background(), textCtx.Graph.Database, graph.Kinds{ad.Entity, azure.Entity}, azureAnalysis.PostProcessedRelationships()...)

	// Deleting transit edges must not return an error
	require.Nil(t, err)

	err = textCtx.Graph.Database.ReadTransaction(context.Background(), func(tx graph.Transaction) error {
		numEdges, err := tx.Relationships().Filter(query.Kind(query.Relationship(), azure.SyncedToADUser)).Count()

		// This must be true which would mean that the above created SyncedToADUser was correctly deleted by the DeleteTransitEdges call
		require.Equal(t, int64(0), numEdges)
		return err
	})

	// The DB must not return any errors
	require.Nil(t, err)

	// This first run removes all AD post-processed relationships - expected outcome is that SyncedToEntraUser is removed at this stage
	_, err = analysis.DeleteTransitEdges(context.Background(), textCtx.Graph.Database, graph.Kinds{ad.Entity, azure.Entity}, adAnalysis.PostProcessedRelationships()...)
	// Deleting transit edges must not return an error
	require.Nil(t, err)

	err = textCtx.Graph.Database.ReadTransaction(context.Background(), func(tx graph.Transaction) error {
		numEdges, err := tx.Relationships().Filter(query.Kind(query.Relationship(), ad.SyncedToEntraUser)).Count()

		// This must be true which would mean that the above created SyncedToADUser was correctly deleted by the DeleteTransitEdges call
		require.Equal(t, int64(0), numEdges)
		return err
	})

	// The DB must not return any errors
	require.Nil(t, err)
}