package cosmosdb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/azure/cosmosdb"
)

func TestRun(t *testing.T) {
	ctx := context.Background()

	cosmosDBCtr, err := cosmosdb.Run(ctx, "mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:latest")
	testcontainers.CleanupContainer(t, cosmosDBCtr)
	require.NoError(t, err)

	t.Run("get-database", func(t *testing.T) {
	})
}
