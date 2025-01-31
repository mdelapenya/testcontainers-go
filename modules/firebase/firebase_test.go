package firebase_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/firebase"
)

func TestFirebase(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Minute))
	defer cancel()

	ctr, err := firebase.Run(ctx, "ghcr.io/u-health/docker-firebase-emulator:13.29.2",
		firebase.WithRoot(filepath.Join(".", "firebase")),
	)
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	// perform assertions
	// Ports are linked to the example config in firebase/firebase.json

	firestoreUrl, err := ctr.ConnectionString(ctx, "8080/tcp")
	assert.NoError(t, err)
	assert.NotEmpty(t, firestoreUrl)
	_ = os.Setenv("FIRESTORE_EMULATOR_HOST", firestoreUrl)
	c, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	require.NoError(t, err)
	defer c.Close()
	w, err := c.Collection("example").Doc("one").Set(ctx, map[string]interface{}{
		"foo": "bar",
	})
	require.NotNil(t, w)
	require.NoError(t, err)
	snap, err := c.Collection("example").Doc("one").Get(ctx)
	require.NoError(t, err)
	var out map[string]string
	err = snap.DataTo(&out)
	require.NoError(t, err)
	assert.Equal(t, "bar", out["foo"])
}

func TestFirebaseBadDirectory(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Minute))
	defer cancel()

	ctr, err := firebase.Run(ctx, "ghcr.io/u-health/docker-firebase-emulator:13.29.2",
		firebase.WithRoot(filepath.Join(".", "failure")),
	)
	// In this case, the file gets copied over at /srv/failure (instead of /srv/firebase)
	// and this stops working.
	// What would be a solution here? Previously I just added an assertion that the root must
	// end in "/firebase"... I could do the same.
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)
}

func TestFirebaseRequiresRoot(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Minute))
	defer cancel()

	ctr, err := firebase.Run(ctx, "ghcr.io/u-health/docker-firebase-emulator:13.29.2")
	testcontainers.CleanupContainer(t, ctr)
	require.Error(t, err)
	require.ErrorContains(t, err, "unable to boot without configuration root")
}
