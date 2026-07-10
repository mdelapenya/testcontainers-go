package eventstoredb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/eventstoredb"
)

func TestEventStoreDB(t *testing.T) {
	ctx := context.Background()

	ctr, err := eventstoredb.Run(ctx, "eventstore/eventstore:24.2")
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	t.Run("ConnectionString_insecure_default", func(t *testing.T) {
		connStr, err := ctr.ConnectionString(ctx)
		require.NoError(t, err)
		require.Contains(t, connStr, "esdb://")
		require.Contains(t, connStr, "tls=false")
	})
}

func TestEventStoreDB_WithInsecure(t *testing.T) {
	ctx := context.Background()

	ctr, err := eventstoredb.Run(ctx, "eventstore/eventstore:24.2",
		eventstoredb.WithInsecure(),
	)
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	connStr, err := ctr.ConnectionString(ctx)
	require.NoError(t, err)
	require.Contains(t, connStr, "esdb://")
	require.Contains(t, connStr, "tls=false")
}
